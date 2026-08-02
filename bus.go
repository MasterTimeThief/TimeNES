package main

var AddressBus, ppuAddressBus uint16
var cpuOpenBus, ppuBus byte
var OAMBusAddress byte

var ppuBusDecay [8]int
var ppuBusDecayConstant int = 1786830

// Read from Address, and return that byte
func Read(Address uint16) byte {
	if Address < 0x2000 {
		//Read from RAM (Accounting for RAM Mirroring)
		return RAM[Address&0x7FF]
	} else if Address < 0x4000 {
		//Reading a PPU Register
		Address &= 0x2007
		switch Address {
		case 0x2000: //PPUCTRL
		case 0x2001: //PPUMASK
		case 0x2002: //PPUSTATUS
			ppuBus &= 0x1F
			ppuBus |= byte(ternary(ppuStatus_VBlank, 0x80, 0x00))
			ppuBus |= byte(ternary(ppuStatus_SpriteZeroHit, 0x40, 0x00))
			ppuBus |= byte(ternary(ppuStatus_Overflow, 0x20, 0x00))
			ppuStatus_VBlank = false
			WriteLatch = false
			UpdatePPUBus2002(ppuBus)
		case 0x2003: //OAM ADDR
		case 0x2004: //OAMDATA
			if ppuScanline < 240 && (ppuMask_RenderBG || ppuMask_RenderSprites) {
				//Return buffer?
				if ppuDot == 0 || ppuDot > 320 {
					ppuBus = SecondaryOAM[0]
				} else if ppuDot > 0 && ppuDot <= 64 {
					ppuBus = 0xFF
				} else {
					return ppuBus
				}
			} else {
				ppuBus = OAM[OAMBusAddress]
			}
			UpdatePPUBus(ppuBus)
		case 0x2005: //PPUSCROLL
		case 0x2006: //PPUADDR
		case 0x2007: //PPUDATA
			if (VRAMAddress & 0x3FFF) > 0x3F00 {
				ppuAddressBus = VRAMAddress
				//Palette data
				data := ReadPPU()
				ppuBus = (ppuBus & 0xC0) | (data & byte(ternary(ppuMask_Greyscale, 0x30, 0x3F)))
				//PPUReadBuffer = ReadPPU((VRAMAddress & 0x2F00) | (VRAMAddress & 0xFF))
				//PPUReadBuffer = ppuBus
				UpdatePPUBus2007Palette(ppuBus)
			} else {
				ppuBus = PPUReadBuffer
				PPUReadBuffer = ReadPPU()
				UpdatePPUBus(ppuBus)
			}

			VRAMAddress += ternary(ppuCtrl_VRAMInc32Mode, 0x20, 0x01)
			VRAMAddress &= 0x3FFF
		default:
			return ppuBus
		}
		return ppuBus
	} else if Address == 0x4015 { //APU Status
		apuStatus := byte(0)
		apuStatus |= byte(ternary(apuDMCInterrupt, 0x80, 0x00))                                                    //DMC Interrupt
		apuStatus |= byte(ternary(apuFrameInterrupt, 0x40, 0x00))                                                  //Frame Interrupt
		apuStatus |= byte(ternary(apuDMC.Length > 0, 0x10, 0x00))                                                  //DMC Active
		apuStatus |= byte(ternary(!(apuNoise.LengthCounter == 0 /*|| apuNoise.LengthCounterHalt*/), 0x08, 0x00))   //Noise Active
		apuStatus |= byte(ternary(!(apuTriangle.LengthCounter == 0 /*|| apuTriangle.Control*/), 0x04, 0x00))       //Triangle Active
		apuStatus |= byte(ternary(!(apuPulse2.LengthCounter == 0 /*|| apuPulse2.LengthCounterHalt*/), 0x02, 0x00)) //Pulse 2 Active
		apuStatus |= byte(ternary(!(apuPulse1.LengthCounter == 0 /*|| apuPulse1.LengthCounterHalt*/), 0x01, 0x00)) //Pulse 1 Active
		apuFrameInterrupt = false
		return apuStatus
	} else if Address == 0x4016 { //Controller 1
		cBit := byte((Controller1ShiftRegister & 0x80) >> 7)
		Controller1ShiftRegister <<= 1
		return cBit
	} else if Address == 0x4017 { //Controller 2
		cBit := byte((Controller2ShiftRegister & 0x80) >> 7)
		Controller2ShiftRegister <<= 1
		return cBit
	} else if Address >= 0x8000 {
		//Read from ROM
		return ROM[(Address-0x8000)&((uint16(Header[4])*0x4000)-1)]
		//
	} else {
		outsideCodeRead = Address
	}
	return 0
}

// Write the Value into the Address given (PPU may have extra steps)
func Write(Address uint16, Value byte) {
	if Address < 0x2000 {
		RAM[Address&0x7FF] = Value
	} else if Address < 0x4000 {
		//Write to PPU Register
		UpdatePPUBus(Value)
		Address &= 0x2007
		switch Address {
		case 0x2000: //PPUCTRL
			ppuCtrl_NametableSelect = Value & 0x03
			ppuCtrl_VRAMInc32Mode = (Value & 0x04) != 0
			ppuCtrl_SpritePatternTable = (Value & 0x08) != 0
			ppuCtrl_BGPatternTable = (Value & 0x10) != 0
			ppuCtrl_Use8x16Sprites = (Value & 0x20) != 0
			ppuCtrl_EnableNMI = (Value & 0x80) != 0

			TransferAddress = (uint16(ppuCtrl_NametableSelect) << 10) | (uint16(TransferAddress) & 0x73FF)
		case 0x2001: //PPUMASK
			ppuMask_Greyscale = (Value & 0x01) != 0
			ppuMask_8pxMaskBG = (Value & 0x02) != 0
			ppuMask_8pxMaskSprites = (Value & 0x04) != 0
			ppuMask_RenderBG = (Value & 0x08) != 0
			ppuMask_RenderSprites = (Value & 0x10) != 0
			//NTSC scanline stuff
			ppuMask_EmphasisRed = (Value & 0x20) != 0
			ppuMask_EmphasisGreen = (Value & 0x40) != 0
			ppuMask_EmphasisBlue = (Value & 0x80) != 0
		case 0x2002: //PPUSTATUS
		case 0x2003: //OAMADDR
			OAMBusAddress = Value
		case 0x2004: //OAMDATA
			if ((ppuScanline >= 240 && ppuScanline < 261) && (ppuMask_RenderBG || ppuMask_RenderSprites)) || (!ppuMask_RenderBG && !ppuMask_RenderSprites) {
				if (OAMBusAddress & 3) == 2 {
					Value &= 0xE3
				}
				OAM[OAMBusAddress] = Value
				OAMBusAddress++
			} else {
				OAMBusAddress += 4
				OAMBusAddress &= 0xFC
			}
		case 0x2005: //PPUSCROLL
			if !WriteLatch {
				ppuScrollFineX = byte(Value & 7)
				TempVRAMAddress = uint16((TempVRAMAddress & 0b0111111111100000) | uint16(Value>>3))
			} else {
				TransferAddress = ((TempVRAMAddress & 0b0000110000011111) | uint16(uint16(Value&0xF8)<<2) | uint16(uint16(Value&7)<<12) | (uint16(ppuCtrl_NametableSelect&1) << 10))
			}
			WriteLatch = !WriteLatch
		case 0x2006: //PPUADDR
			if !WriteLatch {
				//First write sets the high byte
				TempVRAMAddress = (uint16(Value&0x7F) << 8)
				//The actual VRAMAddress isn't changed until the 2nd write
			} else {
				//Second write sets the low byte
				TransferAddress = (TempVRAMAddress | uint16(Value))
				VRAMAddress = TransferAddress /*& 0x3FFF*/
			}
			WriteLatch = !WriteLatch
		case 0x2007: //PPUDATA
			if VRAMAddress < 0x2000 {
				//Write to pattern table. (If the cartridge supports it)
				if Header[5] == 0 {
					CHRROM[VRAMAddress] = Value
				}
				//else, nothing happens because it's CHR-ROM
			} else if VRAMAddress < 0x3F00 {
				//Write to the Nametables
				if (Header[6] & 1) == 0 {
					// Horizontal Mirroring
					VRAM[int(VRAMAddress&0x3FF)|int(VRAMAddress&0x800)>>1] = Value
				} else {
					//Vertical Mirroring
					index := int(VRAMAddress & 0x7FF)
					VRAM[index] = Value
				}
			} else {
				//Write to Palette RAM
				if (VRAMAddress & 3) == 0 {
					PaletteRAM[VRAMAddress&0x0F] = Value
				} else {
					PaletteRAM[VRAMAddress&0x1F] = Value
				}
			}

			VRAMAddress += ternary(ppuCtrl_VRAMInc32Mode, 0x20, 0x01)
			VRAMAddress &= 0x3FFF
		}

	} else if Address < 0x4018 { //APU and I/O

		switch Address {
		// Pulse 1
		case 0x4000:
			apuPulse1.Duty = (Value & 0xC0) >> 6
			apuPulse1.LengthCounterHalt = ((Value & 0x20) >> 5) != 0
			apuPulse1.ConstantVolume = ((Value & 0x10) >> 4) != 0
			apuPulse1.Envelope = (Value & 0xF)
		case 0x4001:
			apuPulse1.SweepUnitEnabled = ((Value & 0x80) >> 7) != 0
			apuPulse1.Period = (Value & 0x70) >> 4
			apuPulse1.Negate = ((Value & 0x8) >> 3) != 0
			apuPulse1.ShiftRegister = (Value & 0x06)
		case 0x4002:
			apuPulse1.Timer = uint16(Value)
		case 0x4003:
			apuPulse1.Timer |= (uint16(Value&0x7) << 8)
			if apuPulse1.Enabled {
				apuPulse1.LengthCounter = LengthCounterLoad(Value >> 3)

			}

		// Pulse 2
		case 0x4004:
			apuPulse2.Duty = (Value & 0xC0)
			apuPulse2.LengthCounterHalt = ((Value & 0x20) >> 5) != 0
			apuPulse2.ConstantVolume = ((Value & 0x10) >> 4) != 0
			apuPulse2.Envelope = (Value & 0xF)
		case 0x4005:
			apuPulse2.SweepUnitEnabled = ((Value & 0x80) >> 7) != 0
			apuPulse2.Period = (Value & 0x70) >> 4
			apuPulse2.Negate = ((Value & 0x8) >> 3) != 0
			apuPulse2.ShiftRegister = (Value & 0x06)
		case 0x4006:
			apuPulse2.Timer = uint16(Value)
		case 0x4007:
			apuPulse2.Timer |= (uint16(Value&0x7) << 8)
			if apuPulse2.Enabled {
				apuPulse2.LengthCounter = LengthCounterLoad(Value >> 3)
			}

		// Triangle
		case 0x4008:
			apuTriangle.Control = ((Value & 0x80) >> 7) != 0
			apuTriangle.LinearCounter = (Value & 0x7F)
		case 0x4009: //Unused
		case 0x400A:
			apuTriangle.Timer = uint16(Value)
		case 0x400B:
			apuTriangle.Timer |= (uint16(Value&0x7) << 8)
			if apuTriangle.Enabled {
				apuTriangle.LengthCounter = LengthCounterLoad(Value >> 3)
			}

		// Noise
		case 0x400C:
			apuNoise.LengthCounterHalt = ((Value & 0x20) >> 5) != 0
			apuNoise.ConstantVolume = ((Value & 0x10) >> 4) != 0
			apuNoise.Envelope = (Value & 0xF)
		case 0x400D: //Unused
		case 0x400E:
			apuNoise.Mode = ((Value & 0x80) >> 7) != 0
			apuNoise.Period = (Value & 0xF)
		case 0x400F:
			if apuNoise.Enabled {
				apuNoise.LengthCounter = LengthCounterLoad(Value >> 3)
			}
		// DMC
		case 0x4010:
			apuDMC.IRQEnable = ((Value & 0x80) >> 7) != 0
			apuDMC.Loop = ((Value & 0x40) >> 6) != 0
			apuDMC.Frequency = (Value & 0xF)
		case 0x4011:
			apuDMC.LoadCounter = (Value & 0x7F)
		case 0x4012:
			apuDMC.Address = Value
		case 0x4013:
			apuDMC.Length = int((Value << 4) + 1)

		case 0x4014: //OAMDMA
			for i := 0; i < 256; i++ {
				OAM[OAMBusAddress] = Read((uint16(Value) << 8) + uint16(i))
				OAMBusAddress++
			}
		case 0x4015: //APU Status
			apuDMC.Length = int((Value & 0x10) >> 4)
			apuNoise.LengthCounter &= (0xFF * ((Value & 0x08) >> 3))
			apuTriangle.LengthCounter &= (0xFF * ((Value & 0x04) >> 2))
			apuPulse2.LengthCounter &= (0xFF * ((Value & 0x02) >> 1))
			apuPulse1.LengthCounter &= (0xFF * (Value & 0x01))

			apuNoise.Enabled = (Value & 0x08) != 0
			apuTriangle.Enabled = (Value & 0x04) != 0
			apuPulse2.Enabled = (Value & 0x02) != 0
			apuPulse1.Enabled = (Value & 0x01) != 0
			apuDMCInterrupt = false
		case 0x4016: //Controller Input
			Controller1ShiftRegister = uint16(Controller1)
			Controller2ShiftRegister = uint16(Controller2)
		case 0x4017: //APU Frame Counter control
			//modeFlagPrev := apuFrameCounterMode
			apuFrameCounterMode = ((Value & 0x80) >> 7) != 0
			apuInhibitIRQ = ((Value & 0x40) >> 6) != 0
			if apuInhibitIRQ {
				apuFrameInterrupt = false
			}
			if /*!modeFlagPrev &&*/ apuFrameCounterMode {
				ClockFrameCounterQuarterFrame()
				ClockFrameCounterHalfFrame()
			}
			apu4017ResetTimer = int(ternary(apuDMAGetCycle, 4, 3))

		}

		/*} else if Address == 0x4014 { //OAM
			for i := 0; i < 256; i++ {
				OAM[i] = Read((uint16(Value) << 8) + uint16(i))
			}
		} else if Address == 0x4016 { //Controller Input
			Controller1ShiftRegister = uint16(Controller1)
			Controller2ShiftRegister = uint16(Controller2)*/
	} else if Address < 0x401B {
		//Audio Processing Unit stuff
		//$4000 - $4017 is APU and I/O registers
		//$4018 - $401F is APU and I/O functions that are normally disabled

	} else if Address < 0x8000 {
		CartRAM[Address&0x1FFF] = Value
	} else {
		outsideCodeWrite = Address
	}
}

func ReadPPU( /*Address uint16*/ ) byte {
	if ppuAddressBus < 0x2000 {
		//Read from pattern table.
		return CHRROM[ppuAddressBus]
		//else, nothing happens
	} else if ppuAddressBus < 0x3F00 {
		//Read from the Nametables
		if (Header[6] & 1) == 0 {
			// Horizontal Mirroring
			return VRAM[int(ppuAddressBus&0x3FF)|(int(ppuAddressBus&0x800)>>1)]
		} else {
			//Vertical Mirroring
			return VRAM[int(ppuAddressBus&0x7FF)]
		}
	} else {
		//Read from Palette RAM
		if (ppuAddressBus & 3) == 0 {
			return PaletteRAM[ppuAddressBus&0x0F]
		} else {
			return PaletteRAM[ppuAddressBus&0x1F]
		}
	}
}

func BuildAddress(Value_Low, Value_High byte) uint16 {
	//b := []byte{Value_Low, Value_High}
	//return binary.LittleEndian.Uint16(b[0:])

	return (uint16(Value_High)<<8 | uint16(Value_Low))

}

func ReadFromPC() byte {
	Value := Read(ProgramCounter)
	if LoggingCPU {
		operands = append(operands, Value)
	}
	ProgramCounter++
	return Value
}

func UpdatePPUBus(Value byte) byte {
	ppuBus = Value
	//PPU decay buffer code here
	for i := 0; i < 8; i++ {
		ppuBusDecay[i] = ppuBusDecayConstant
	}
	//???
	return ppuBus
}

func UpdatePPUBus2002(Value byte) {
	ppuBus = Value
	//Only update the decay constant on the top 3 bits
	for i := 5; i < 8; i++ {
		ppuBusDecay[i] = ppuBusDecayConstant
	}
}

func UpdatePPUBus2007Palette(Value byte) {
	ppuBus = Value
	//Only update the decay constant on the bottom 6 bits
	for i := 0; i < 6; i++ {
		ppuBusDecay[i] = ppuBusDecayConstant
	}
}

func DecayPPUDataBus() {
	DecayBitmask := [8]byte{0xFE, 0xFD, 0xFB, 0xF7, 0xEF, 0xDF, 0xBF, 0x7F}

	for i := range ppuBusDecay {
		if ppuBusDecay[i] > 0 {
			ppuBusDecay[i]--
			if ppuBusDecay[i] == 0 {
				ppuBus &= DecayBitmask[i]
			}
		}
	}
}

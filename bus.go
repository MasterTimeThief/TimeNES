package main

import "mtt/timenes/mappers"

var AddressBus, ppuAddressBus uint16
var cpuOpenBus, ppuBus byte
var OAMBusAddress byte

var ppuBusDecay [8]int
var ppuBusDecayConstant int = 1786830

// Read from Address, and return that byte
func Read(Address uint16) byte {
	var returnValue byte
	if Address < 0x2000 {
		//Read from RAM (Accounting for RAM Mirroring)
		returnValue = RAM[Address&0x7FF]
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
					returnValue = ppuBus
				}
			} else {
				ppuBus = OAM[OAMBusAddress]
			}
			UpdatePPUBus(ppuBus)
		case 0x2005: //PPUSCROLL
		case 0x2006: //PPUADDR
		case 0x2007: //PPUDATA
			ppuAddressBus = VRAMAddress
			if (VRAMAddress & 0x3FFF) > 0x3F00 {
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
			returnValue = ppuBus
		}
		returnValue = ppuBus
	} else if Address == 0x4015 { //APU Status
		apuStatus := byte(0)
		apuStatus |= byte(ternary(apuDMCInterrupt, 0x80, 0x00))                                                                     //DMC Interrupt
		apuStatus |= byte(ternary(apuFrameInterrupt, 0x40, 0x00))                                                                   //Frame Interrupt
		apuStatus |= byte(ternary(apuDMC.BytesRemaining > 0, 0x10, 0x00))                                                           //DMC Active
		apuStatus |= byte(ternary(!(apuNoise.LengthCounter.Counter == 0 /*|| apuNoise.LengthCounter.HaltFlag*/), 0x08, 0x00))       //Noise Active
		apuStatus |= byte(ternary(!(apuTriangle.LengthCounter.Counter == 0 /*|| apuTriangle.LengthCounter.HaltFlag*/), 0x04, 0x00)) //Triangle Active
		apuStatus |= byte(ternary(!(apuPulse2.LengthCounter.Counter == 0 /*|| apuPulse2.LengthCounter.HaltFlag*/), 0x02, 0x00))     //Pulse 2 Active
		apuStatus |= byte(ternary(!(apuPulse1.LengthCounter.Counter == 0 /*|| apuPulse1.LengthCounter.HaltFlag*/), 0x01, 0x00))     //Pulse 1 Active
		apuFrameInterrupt = false
		returnValue = apuStatus
	} else if Address == 0x4016 { //Controller 1
		cBit := byte((Controller1ShiftRegister & 0x80) >> 7)
		Controller1ShiftRegister <<= 1
		returnValue = cBit
	} else if Address == 0x4017 { //Controller 2
		cBit := byte((Controller2ShiftRegister & 0x80) >> 7)
		Controller2ShiftRegister <<= 1
		returnValue = cBit

	} else if Address < 0x7FFF {
		//Could be PRG-RAM on the cartridge
		switch MapperChipID {
		case 1: //MMC1
			returnValue = mappers.MMC1_PRGRAM[Address&0x1FFF]
		default:
			//outsideCodeRead = Address
		}
	} else if Address >= 0x8000 {
		//Read from ROM
		switch MapperChipID {
		case 1: //MMC1
			returnValue = ROM[mappers.MMC1_FetchCPUAddress(Address, PRGROM_Size)]
		//case 3: //CNROM
		//case 4: //MMC3
		default:
			returnValue = ROM[(Address-0x8000)&uint16(PRGROM_Size-1)]
			//outsideCodeRead = Address
		}
	}
	//MasterClockTick("READ")
	return returnValue
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
				TransferAddress = uint16((TransferAddress & 0b0111111111100000) | uint16(Value>>3))
			} else {
				TransferAddress = ((TransferAddress & 0b0000110000011111) | uint16(uint16(Value&0xF8)<<2) | uint16(uint16(Value&7)<<12) /*| (uint16(ppuCtrl_NametableSelect&1) << 10)*/)
			}
			WriteLatch = !WriteLatch
		case 0x2006: //PPUADDR
			if !WriteLatch {
				//First write sets the high byte
				//TempVRAMAddress = (uint16(Value&0x7F) << 8)
				TransferAddress = uint16((TransferAddress & 0xFF) | uint16(Value&0x3F)<<8)
				//The actual VRAMAddress isn't changed until the 2nd write
			} else {
				//Second write sets the low byte
				TransferAddress = ((TransferAddress & 0xFF00) | uint16(Value))
				VRAMAddress = TransferAddress /*& 0x3FFF*/
			}
			WriteLatch = !WriteLatch
		case 0x2007: //PPUDATA
			WritePPU(Value)

			VRAMAddress += ternary(ppuCtrl_VRAMInc32Mode, 0x20, 0x01)
			VRAMAddress &= 0x3FFF
		}

	} else if Address < 0x4018 { //APU and I/O

		switch Address {
		// Pulse 1
		case 0x4000:
			apuPulse1.Duty = (Value & 0xC0) >> 6
			apuPulse1.LengthCounter.HaltFlag = ((Value & 0x20) >> 5) != 0
			apuPulse1.Envelope.ConstantVolume = ((Value & 0x10) >> 4) != 0
			apuPulse1.Envelope.Volume = (Value & 0xF)
			apuPulse1.LoopFlag = apuPulse1.HaltFlag
		case 0x4001:
			apuPulse1.Sweep.Enabled = ((Value & 0x80) >> 7) != 0
			apuPulse1.Sweep.Period = uint16(Value&0x70) >> 4
			apuPulse1.Sweep.Negate = ((Value & 0x8) >> 3) != 0
			apuPulse1.Shift = (Value & 0x06)
		case 0x4002:
			apuPulse1.TimerReloadValue = SetTimerLow(Value, apuPulse1.TimerReloadValue)
			apuPulse1.Timer = apuPulse1.TimerReloadValue
		case 0x4003:
			apuPulse1.TimerReloadValue = SetTimerHi(Value, apuPulse1.TimerReloadValue)
			apuPulse1.Timer = apuPulse1.TimerReloadValue
			if apuPulse1.Enabled {
				apuPulse1.LengthCounter.Counter = LengthCounterLoad(Value >> 3)
			}

		// Pulse 2
		case 0x4004:
			apuPulse2.Duty = (Value & 0xC0) >> 6
			apuPulse2.LengthCounter.HaltFlag = ((Value & 0x20) >> 5) != 0
			apuPulse2.Envelope.ConstantVolume = ((Value & 0x10) >> 4) != 0
			apuPulse2.Envelope.Volume = (Value & 0xF)
			apuPulse2.LoopFlag = apuPulse2.HaltFlag
		case 0x4005:
			apuPulse2.Sweep.Enabled = ((Value & 0x80) >> 7) != 0
			apuPulse2.Sweep.Period = uint16(Value&0x70) >> 4
			apuPulse2.Sweep.Negate = ((Value & 0x8) >> 3) != 0
			apuPulse2.Shift = (Value & 0x06)
		case 0x4006:
			apuPulse2.TimerReloadValue = SetTimerLow(Value, apuPulse2.TimerReloadValue)
			apuPulse2.Timer = apuPulse2.TimerReloadValue
		case 0x4007:
			apuPulse2.TimerReloadValue = SetTimerHi(Value, apuPulse2.TimerReloadValue)
			apuPulse2.Timer = apuPulse2.TimerReloadValue
			if apuPulse2.Enabled {
				apuPulse2.LengthCounter.Counter = LengthCounterLoad(Value >> 3)
			}

		// Triangle
		case 0x4008:
			apuTriangle.LengthCounter.HaltFlag = ((Value & 0x80) >> 7) != 0
			apuTriangle.LinearCounter.Counter = (Value & 0x7F)
		case 0x4009: //Unused
		case 0x400A:
			//apuTriangle.Timer = uint16(Value)
			apuTriangle.TimerReloadValue = SetTimerLow(Value, apuTriangle.TimerReloadValue)
			apuTriangle.Timer = apuTriangle.TimerReloadValue
		case 0x400B:
			//apuTriangle.Timer |= (uint16(Value&0x7) << 8)
			apuTriangle.TimerReloadValue = SetTimerHi(Value, apuTriangle.TimerReloadValue)
			apuTriangle.Timer = apuTriangle.TimerReloadValue
			if apuTriangle.Enabled {
				apuTriangle.LengthCounter.Counter = LengthCounterLoad(Value >> 3)
			}
			apuTriangle.LinearCounter.ReloadFlag = true

		// Noise
		case 0x400C:
			apuNoise.LengthCounter.HaltFlag = ((Value & 0x20) >> 5) != 0
			apuNoise.Envelope.ConstantVolume = ((Value & 0x10) >> 4) != 0
			apuNoise.Envelope.Volume = (Value & 0xF)
		case 0x400D: //Unused
		case 0x400E:
			apuNoise.Mode = ((Value & 0x80) >> 7) != 0
			apuNoise.SetNoiseTimer(Value & 0xF)
		case 0x400F:
			if apuNoise.Enabled {
				apuNoise.LengthCounter.Counter = LengthCounterLoad(Value >> 3)
			}

		// DMC
		case 0x4010:
			apuDMC.IRQEnable = ((Value & 0x80) >> 7) != 0
			apuDMC.Loop = ((Value & 0x40) >> 6) != 0
			apuDMC.SampleRate = apuDMCSampleRateLUT[Value&0xF]
			if !apuDMC.IRQEnable {
				apuDMCInterrupt = false
				IRQLevelDetector = false
			}
		case 0x4011:
			apuDMC.Output = (Value & 0x7F)
		case 0x4012:
			apuDMC.SampleAddress = (0xC000 | (uint16(Value) << 6))
		case 0x4013:
			apuDMC.SampleLength = ((uint16(Value) << 4) | 1)

		case 0x4014: //OAMDMA
			for i := 0; i < 256; i++ {
				OAM[OAMBusAddress] = Read((uint16(Value) << 8) + uint16(i))
				OAMBusAddress++
			}
		case 0x4015: //APU Status
			//apuDMC.BytesRemaining = int((Value & 0x10) >> 4)
			apuNoise.LengthCounter.Counter &= (0xFF * ((Value & 0x08) >> 3))
			apuTriangle.LengthCounter.Counter &= (0xFF * ((Value & 0x04) >> 2))
			apuPulse2.LengthCounter.Counter &= (0xFF * ((Value & 0x02) >> 1))
			apuPulse1.LengthCounter.Counter &= (0xFF * (Value & 0x01))

			apuDMC.Enabled = (Value & 0x10) != 0
			apuNoise.Enabled = (Value & 0x08) != 0
			apuTriangle.Enabled = (Value & 0x04) != 0
			apuPulse2.Enabled = (Value & 0x02) != 0
			apuPulse1.Enabled = (Value & 0x01) != 0

			apuDMCInterrupt = false
			IRQLevelDetector = false

		case 0x4016: //Controller Input
			Controller1ShiftRegister = uint16(Controller1)
			Controller2ShiftRegister = uint16(Controller2)
		case 0x4017: //APU Frame Counter control
			//modeFlagPrev := apuFrameCounterMode
			apuFrameCounterMode = ((Value & 0x80) >> 7) != 0
			apuInhibitIRQ = ((Value & 0x40) >> 6) != 0
			if apuInhibitIRQ {
				apuFrameInterrupt = false
				IRQLevelDetector = false
			}
			if /*!modeFlagPrev &&*/ apuFrameCounterMode {
				ClockFrameCounterQuarterFrame()
				ClockFrameCounterHalfFrame()
			}
			apu4017ResetTimer = int(ternary(apuDMAGetCycle, 4, 3))

		}
	} else if Address < 0x401B {
		//Audio Processing Unit stuff
		//$4000 - $4017 is APU and I/O registers
		//$4018 - $401F is APU and I/O functions that are normally disabled

		//} else if Address < 0x8000 {
		//	CartRAM[Address&0x1FFF] = Value
	} else {
		//Check what mapper chip we're using
		switch MapperChipID {
		case 1: //MMC1
			mappers.MMC1_Write(Value, Address, CPU_TotalCycles)
		case 3: //CNROM
			prgValue := Read(Address)
			mappers.CNROM_Write(Value&prgValue, Address)
		//case 4: //MMC3
		default:
			outsideCodeWrite = Address
		}
	}
	//MasterClockTick("WRITE")
}

func ReadPPU( /*Address uint16*/ ) byte {
	if ppuAddressBus < 0x2000 {
		//Read from pattern table.
		switch MapperChipID {
		case 1: //MMC1
			return CHRROM[mappers.MMC1_FetchPPUAddress(ppuAddressBus, CHRROM_Size)]
		case 3: //CNROM
			return CHRROM[mappers.CNROM_ReadAddress(ppuAddressBus)]
		//case 4: //MMC3
		//	return 0
		default:
			return CHRROM[ppuAddressBus]
		}

		//else, nothing happens
	} else if ppuAddressBus < 0x3F00 {
		//Read from the Nametables
		switch MapperChipID {
		case 1: //MMC1
			return CartVRAM[mappers.MMC1_FetchNametable(ppuAddressBus)]
		//case 3: //CNROM
		//case 4: //MMC3
		default:
			if !AltNametableLayout {
				if IsNametableHorizontal {
					// Horizontal Mirroring
					return VRAM[int(ppuAddressBus&0x3FF)|(int(ppuAddressBus&0x800)>>1)]
				} else {
					//Vertical Mirroring
					return VRAM[int(ppuAddressBus&0x7FF)]
				}
			}
		}
		return 0

	} else {
		//Read from Palette RAM
		if (ppuAddressBus & 3) == 0 {
			return PaletteRAM[ppuAddressBus&0x0F]
		} else {
			return PaletteRAM[ppuAddressBus&0x1F]
		}
	}
}

func WritePPU(Value byte) {
	if VRAMAddress < 0x2000 {
		//Write to pattern table. (If the cartridge supports it)
		if CHRROM_Size == 0 {
			switch MapperChipID {
			case 1: //MMC1
				//if AltNametableLayout {
				//	mappers.MMC1_VRAM[mappers.MMC1_FetchPPUAddress(VRAMAddress, CHRROM_Size)] = Value
				//} else {
				CHRROM[mappers.MMC1_FetchPPUAddress(VRAMAddress, CHRROM_Size)] = Value
				//}
			case 3: //CNROM
				CHRROM[mappers.CNROM_ReadAddress(VRAMAddress)] = Value
			//case 4: //MMC3
			default:
				CHRROM[VRAMAddress] = Value
			}
		}
		//else, nothing happens because it's CHR-ROM
	} else if VRAMAddress < 0x3F00 {
		//Write to the Nametables

		switch MapperChipID {
		case 1: //MMC1
			switch mappers.MMC1_GetNametableArrangement() {
			case 0: //Screen A only
				CartVRAM[int(VRAMAddress&0x3FF)] = Value
			case 1: //Screen B Only
				CartVRAM[int(VRAMAddress&0x3FF)+0x400] = Value
			case 2: //Horizontal
				CartVRAM[int(VRAMAddress&0x7FF)] = Value
			case 3: //Vertical
				CartVRAM[int(VRAMAddress&0x3FF)|int(VRAMAddress&0x800)>>1] = Value
			}
			//CartVRAM[VRAMAddress&0xFFF] = Value
		//case 3: //CNROM
		//case 4: //MMC3
		default:
			if !AltNametableLayout {
				if IsNametableHorizontal {
					// Horizontal Mirroring
					VRAM[int(VRAMAddress&0x3FF)|int(VRAMAddress&0x800)>>1] = Value
				} else {
					//Vertical Mirroring
					index := int(VRAMAddress & 0x7FF)
					VRAM[index] = Value
				}
			}
		}

	} else {
		//Write to Palette RAM
		if (VRAMAddress & 3) == 0 {
			PaletteRAM[VRAMAddress&0x0F] = Value
		} else {
			PaletteRAM[VRAMAddress&0x1F] = Value
		}
	}
}

func BuildAddress(Value_Low, Value_High byte) uint16 {
	//b := []byte{Value_Low, Value_High}
	//return binary.LittleEndian.Uint16(b[0:])

	AddressBus = (uint16(Value_High)<<8 | uint16(Value_Low))
	return AddressBus
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

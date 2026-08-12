package nes

import (
	"mtt/timenes/common"
	"mtt/timenes/nes/apu"
	"mtt/timenes/nes/cartridge"
	"mtt/timenes/nes/cartridge/mappers"
	"mtt/timenes/nes/input"
)

var ppuAddressBus uint16
var cpuOpenBus, ppuBus byte
var OAMBusAddress byte

var ppuBusDecay [8]int
var ppuBusDecayConstant int = 1786830

// Read from Address, and return that byte
func Read(Address uint16) byte {
	var returnValue byte
	if Address < 0x2000 {
		//Read from RAM (Accounting for RAM Mirroring)
		returnValue = cartridge.RAM[Address&0x7FF]
	} else if Address < 0x4000 {
		//Reading a PPU Register
		Address &= 0x2007
		switch Address {
		case 0x2000: //PPUCTRL
		case 0x2001: //PPUMASK
		case 0x2002: //PPUSTATUS
			ppuBus &= 0x1F
			ppuBus |= byte(common.Ternary(ppuStatus_VBlank, 0x80, 0x00))
			ppuBus |= byte(common.Ternary(ppuStatus_SpriteZeroHit, 0x40, 0x00))
			ppuBus |= byte(common.Ternary(ppuStatus_Overflow, 0x20, 0x00))
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
				ppuBus = (ppuBus & 0xC0) | (data & byte(common.Ternary(ppuMask_Greyscale, 0x30, 0x3F)))
				//PPUReadBuffer = ReadPPU((VRAMAddress & 0x2F00) | (VRAMAddress & 0xFF))
				//PPUReadBuffer = ppuBus
				UpdatePPUBus2007Palette(ppuBus)
			} else {
				ppuBus = PPUReadBuffer
				PPUReadBuffer = ReadPPU()
				UpdatePPUBus(ppuBus)
			}

			VRAMAddress += common.Ternary(ppuCtrl_VRAMInc32Mode, 0x20, 0x01)
			VRAMAddress &= 0x3FFF
		default:
			returnValue = ppuBus
		}
		returnValue = ppuBus
	} else if Address == 0x4015 { //APU Status
		apuStatus := byte(0)
		apuStatus |= byte(common.Ternary(apu.APUDMCInterrupt, 0x80, 0x00))                                                                  //DMC Interrupt
		apuStatus |= byte(common.Ternary(apu.APUFrameInterrupt, 0x40, 0x00))                                                                //Frame Interrupt
		apuStatus |= byte(common.Ternary(apu.DMC.BytesRemaining > 0, 0x10, 0x00))                                                           //DMC Active
		apuStatus |= byte(common.Ternary(!(apu.Noise.LengthCounter.Counter == 0 /*|| apuNoise.LengthCounter.HaltFlag*/), 0x08, 0x00))       //Noise Active
		apuStatus |= byte(common.Ternary(!(apu.Triangle.LengthCounter.Counter == 0 /*|| apuTriangle.LengthCounter.HaltFlag*/), 0x04, 0x00)) //Triangle Active
		apuStatus |= byte(common.Ternary(!(apu.Pulse2.LengthCounter.Counter == 0 /*|| apuPulse2.LengthCounter.HaltFlag*/), 0x02, 0x00))     //Pulse 2 Active
		apuStatus |= byte(common.Ternary(!(apu.Pulse1.LengthCounter.Counter == 0 /*|| apuPulse1.LengthCounter.HaltFlag*/), 0x01, 0x00))     //Pulse 1 Active
		apu.APUFrameInterrupt = false
		returnValue = apuStatus
	} else if Address == 0x4016 { //Controller 1
		cBit := byte((input.Controller1ShiftRegister & 0x80) >> 7)
		input.Controller1ShiftRegister <<= 1
		returnValue = cBit
	} else if Address == 0x4017 { //Controller 2
		cBit := byte((input.Controller2ShiftRegister & 0x80) >> 7)
		input.Controller2ShiftRegister <<= 1
		returnValue = cBit

	} else if Address < 0x7FFF {
		//Could be PRG-RAM on the cartridge
		switch cartridge.MapperChipID {
		case 1: //MMC1
			returnValue = mappers.MMC1_PRGRAM[Address&0x1FFF]
		default:
			//outsideCodeRead = Address
		}
	} else if Address >= 0x8000 {
		//Read from ROM
		switch cartridge.MapperChipID {
		case 1: //MMC1
			returnValue = cartridge.PRGROM[mappers.MMC1_FetchCPUAddress(Address, cartridge.PRGROM_Size)]
		//case 3: //CNROM
		//case 4: //MMC3
		default:
			returnValue = cartridge.PRGROM[(Address-0x8000)&uint16(cartridge.PRGROM_Size-1)]
			//outsideCodeRead = Address
		}
	}
	//MasterClockTick("READ")
	return returnValue
}

// Write the Value into the Address given (PPU may have extra steps)
func Write(Address uint16, Value byte) {
	if Address < 0x2000 {
		cartridge.RAM[Address&0x7FF] = Value
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

			VRAMAddress += common.Ternary(ppuCtrl_VRAMInc32Mode, 0x20, 0x01)
			VRAMAddress &= 0x3FFF
		}

	} else if Address < 0x4018 { //APU and I/O

		switch Address {
		// Pulse 1
		case 0x4000:
			apu.Pulse1.Duty = (Value & 0xC0) >> 6
			apu.Pulse1.LengthCounter.HaltFlag = ((Value & 0x20) >> 5) != 0
			apu.Pulse1.Envelope.ConstantVolume = ((Value & 0x10) >> 4) != 0
			apu.Pulse1.Envelope.Volume = (Value & 0xF)
			apu.Pulse1.LoopFlag = apu.Pulse1.HaltFlag
		case 0x4001:
			apu.Pulse1.Sweep.Enabled = ((Value & 0x80) >> 7) != 0
			apu.Pulse1.Sweep.Period = uint16(Value&0x70) >> 4
			apu.Pulse1.Sweep.Negate = ((Value & 0x8) >> 3) != 0
			apu.Pulse1.Shift = (Value & 0x07)
		case 0x4002:
			apu.Pulse1.TimerReloadValue = apu.SetTimerLow(Value, apu.Pulse1.TimerReloadValue)
			apu.Pulse1.Timer = apu.Pulse1.TimerReloadValue
		case 0x4003:
			apu.Pulse1.TimerReloadValue = apu.SetTimerHi(Value, apu.Pulse1.TimerReloadValue)
			apu.Pulse1.Timer = apu.Pulse1.TimerReloadValue
			if apu.Pulse1.Enabled {
				apu.Pulse1.LengthCounter.Counter = apu.LengthCounterLoad(Value >> 3)
			}

		// Pulse 2
		case 0x4004:
			apu.Pulse2.Duty = (Value & 0xC0) >> 6
			apu.Pulse2.LengthCounter.HaltFlag = ((Value & 0x20) >> 5) != 0
			apu.Pulse2.Envelope.ConstantVolume = ((Value & 0x10) >> 4) != 0
			apu.Pulse2.Envelope.Volume = (Value & 0xF)
			apu.Pulse2.LoopFlag = apu.Pulse2.HaltFlag
		case 0x4005:
			apu.Pulse2.Sweep.Enabled = ((Value & 0x80) >> 7) != 0
			apu.Pulse2.Sweep.Period = uint16(Value&0x70) >> 4
			apu.Pulse2.Sweep.Negate = ((Value & 0x8) >> 3) != 0
			apu.Pulse2.Shift = (Value & 0x07)
		case 0x4006:
			apu.Pulse2.TimerReloadValue = apu.SetTimerLow(Value, apu.Pulse2.TimerReloadValue)
			apu.Pulse2.Timer = apu.Pulse2.TimerReloadValue
		case 0x4007:
			apu.Pulse2.TimerReloadValue = apu.SetTimerHi(Value, apu.Pulse2.TimerReloadValue)
			apu.Pulse2.Timer = apu.Pulse2.TimerReloadValue
			if apu.Pulse2.Enabled {
				apu.Pulse2.LengthCounter.Counter = apu.LengthCounterLoad(Value >> 3)
			}

		// Triangle
		case 0x4008:
			apu.Triangle.LengthCounter.HaltFlag = ((Value & 0x80) >> 7) != 0
			if apu.Triangle.Enabled {
				apu.Triangle.LinearCounter.Counter = (Value & 0x7F)
			}
		case 0x4009: //Unused
		case 0x400A:
			//apuTriangle.Timer = uint16(Value)
			apu.Triangle.TimerReloadValue = apu.SetTimerLow(Value, apu.Triangle.TimerReloadValue)
			apu.Triangle.Timer = apu.Triangle.TimerReloadValue
		case 0x400B:
			//apuTriangle.Timer |= (uint16(Value&0x7) << 8)
			apu.Triangle.TimerReloadValue = apu.SetTimerHi(Value, apu.Triangle.TimerReloadValue)
			apu.Triangle.Timer = apu.Triangle.TimerReloadValue
			if apu.Triangle.Enabled {
				apu.Triangle.LengthCounter.Counter = apu.LengthCounterLoad(Value >> 3)
			}
			apu.Triangle.LinearCounter.ReloadFlag = true

		// Noise
		case 0x400C:
			apu.Noise.LengthCounter.HaltFlag = ((Value & 0x20) >> 5) != 0
			apu.Noise.Envelope.ConstantVolume = ((Value & 0x10) >> 4) != 0
			apu.Noise.Envelope.Volume = (Value & 0xF)
		case 0x400D: //Unused
		case 0x400E:
			apu.Noise.Mode = ((Value & 0x80) >> 7) != 0
			apu.Noise.SetNoiseTimer(Value & 0xF)
		case 0x400F:
			if apu.Noise.Enabled {
				apu.Noise.LengthCounter.Counter = apu.LengthCounterLoad(Value >> 3)
			}

		// DMC
		case 0x4010:
			apu.DMC.IRQEnable = ((Value & 0x80) >> 7) != 0
			apu.DMC.Loop = ((Value & 0x40) >> 6) != 0
			apu.DMC.SampleRate = apu.APUDMCSampleRateLUT[Value&0xF]
			if !apu.DMC.IRQEnable {
				apu.APUDMCInterrupt = false
				apu.IRQLevelDetector = false
			}
		case 0x4011:
			apu.DMC.Output = (Value & 0x7F)
		case 0x4012:
			apu.DMC.SampleAddress = (0xC000 | (uint16(Value) << 6))
			apu.DMC.CurrentAddress = apu.DMC.SampleAddress
			PrebufferDMCSamples()
		case 0x4013:
			apu.DMC.SampleLength = ((uint16(Value) << 4) | 1)
			apu.DMC.BytesRemaining = apu.DMC.SampleLength
			PrebufferDMCSamples()

		case 0x4014: //OAMDMA
			for i := 0; i < 256; i++ {
				OAM[OAMBusAddress] = Read((uint16(Value) << 8) + uint16(i))
				OAMBusAddress++
			}
		case 0x4015: //APU Status
			//apuDMC.BytesRemaining = int((Value & 0x10) >> 4)
			apu.Noise.LengthCounter.Counter &= (0xFF * ((Value & 0x08) >> 3))
			apu.Triangle.LengthCounter.Counter &= (0xFF * ((Value & 0x04) >> 2))
			apu.Triangle.LinearCounter.Counter &= (0xFF * ((Value & 0x04) >> 2))
			apu.Pulse2.LengthCounter.Counter &= (0xFF * ((Value & 0x02) >> 1))
			apu.Pulse1.LengthCounter.Counter &= (0xFF * (Value & 0x01))

			apu.DMC.Enabled = (Value & 0x10) != 0
			apu.Noise.Enabled = (Value & 0x08) != 0
			apu.Triangle.Enabled = (Value & 0x04) != 0
			apu.Pulse2.Enabled = (Value & 0x02) != 0
			apu.Pulse1.Enabled = (Value & 0x01) != 0

			apu.APUDMCInterrupt = false
			apu.IRQLevelDetector = false

		case 0x4016: //Controller Input
			input.UpdateControllers()
		case 0x4017: //APU Frame Counter control
			//modeFlagPrev := apuFrameCounterMode
			apu.APUFrameCounterMode = ((Value & 0x80) >> 7) != 0
			apu.APUInhibitIRQ = ((Value & 0x40) >> 6) != 0
			if apu.APUInhibitIRQ {
				apu.APUFrameInterrupt = false
				apu.IRQLevelDetector = false
			}
			if /*!modeFlagPrev &&*/ apu.APUFrameCounterMode {
				apu.ClockFrameCounterQuarterFrame()
				apu.ClockFrameCounterHalfFrame()
			}
			apu.Set4017ResetTimer()

		}
	} else if Address < 0x401B {
		//Audio Processing Unit stuff
		//$4000 - $4017 is APU and I/O registers
		//$4018 - $401F is APU and I/O functions that are normally disabled

		//} else if Address < 0x8000 {
		//	CartRAM[Address&0x1FFF] = Value
	} else {
		//Check what mapper chip we're using
		switch cartridge.MapperChipID {
		case 1: //MMC1
			mappers.MMC1_Write(Value, Address, CPU_TotalCycles)
		case 3: //CNROM
			prgValue := Read(Address)
			mappers.CNROM_Write(Value&prgValue, Address)
		//case 4: //MMC3
		default:
			OutsideCodeWrite = Address
		}
	}
	//MasterClockTick("WRITE")
}

func ReadPPU( /*Address uint16*/ ) byte {
	if ppuAddressBus < 0x2000 {
		//Read from pattern table.
		switch cartridge.MapperChipID {
		case 1: //MMC1
			return cartridge.CHRROM[mappers.MMC1_FetchPPUAddress(ppuAddressBus, cartridge.CHRROM_Size)]
		case 3: //CNROM
			return cartridge.CHRROM[mappers.CNROM_ReadAddress(ppuAddressBus)]
		//case 4: //MMC3
		//	return 0
		default:
			return cartridge.CHRROM[ppuAddressBus]
		}

		//else, nothing happens
	} else if ppuAddressBus < 0x3F00 {
		//Read from the Nametables
		switch cartridge.MapperChipID {
		case 1: //MMC1
			return cartridge.CartVRAM[mappers.MMC1_FetchNametable(ppuAddressBus)]
		//case 3: //CNROM
		//case 4: //MMC3
		default:
			if !cartridge.AltNametableLayout {
				if cartridge.IsNametableHorizontal {
					// Horizontal Mirroring
					return cartridge.VRAM[int(ppuAddressBus&0x3FF)|(int(ppuAddressBus&0x800)>>1)]
				} else {
					//Vertical Mirroring
					return cartridge.VRAM[int(ppuAddressBus&0x7FF)]
				}
			}
		}
		return 0

	} else {
		//Read from Palette RAM
		if (ppuAddressBus & 3) == 0 {
			return cartridge.PaletteRAM[ppuAddressBus&0x0F]
		} else {
			return cartridge.PaletteRAM[ppuAddressBus&0x1F]
		}
	}
}

func WritePPU(Value byte) {
	if VRAMAddress < 0x2000 {
		//Write to pattern table. (If the cartridge supports it)
		if cartridge.CHRROM_Size == 0 {
			switch cartridge.MapperChipID {
			case 1: //MMC1
				//if AltNametableLayout {
				//	mappers.MMC1_VRAM[mappers.MMC1_FetchPPUAddress(VRAMAddress, CHRROM_Size)] = Value
				//} else {
				cartridge.CHRROM[mappers.MMC1_FetchPPUAddress(VRAMAddress, cartridge.CHRROM_Size)] = Value
				//}
			case 3: //CNROM
				cartridge.CHRROM[mappers.CNROM_ReadAddress(VRAMAddress)] = Value
			//case 4: //MMC3
			default:
				cartridge.CHRROM[VRAMAddress] = Value
			}
		}
		//else, nothing happens because it's CHR-ROM
	} else if VRAMAddress < 0x3F00 {
		//Write to the Nametables

		switch cartridge.MapperChipID {
		case 1: //MMC1
			switch mappers.MMC1_GetNametableArrangement() {
			case 0: //Screen A only
				cartridge.CartVRAM[int(VRAMAddress&0x3FF)] = Value
			case 1: //Screen B Only
				cartridge.CartVRAM[int(VRAMAddress&0x3FF)+0x400] = Value
			case 2: //Horizontal
				cartridge.CartVRAM[int(VRAMAddress&0x7FF)] = Value
			case 3: //Vertical
				cartridge.CartVRAM[int(VRAMAddress&0x3FF)|int(VRAMAddress&0x800)>>1] = Value
			}
			//cartridge.CartVRAM[VRAMAddress&0xFFF] = Value
		//case 3: //CNROM
		//case 4: //MMC3
		default:
			if !cartridge.AltNametableLayout {
				if cartridge.IsNametableHorizontal {
					// Horizontal Mirroring
					cartridge.VRAM[int(VRAMAddress&0x3FF)|int(VRAMAddress&0x800)>>1] = Value
				} else {
					//Vertical Mirroring
					index := int(VRAMAddress & 0x7FF)
					cartridge.VRAM[index] = Value
				}
			}
		}

	} else {
		//Write to Palette RAM
		if (VRAMAddress & 3) == 0 {
			cartridge.PaletteRAM[VRAMAddress&0x0F] = Value
		} else {
			cartridge.PaletteRAM[VRAMAddress&0x1F] = Value
		}
	}
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

func PrebufferDMCSamples() {

	if apu.DMC.SampleLength > 0 && apu.DMC.SampleAddress > 0x7FFF {
		apu.DMC.SampleBuffer = nil

		//Get the entire sample upfront, because God is dead
		for i := uint16(0); i < apu.DMC.SampleLength; i++ {
			sampleAddr := apu.DMC.SampleAddress + i
			if sampleAddr < 0x8000 { //We hit 0xFFFF and wrapped around
				sampleAddr += 0x8000
			}

			apu.DMC.SampleBuffer = append(apu.DMC.SampleBuffer, Read(sampleAddr))
		}
		apu.DMC.SampleBufferPos = 0
	}

	//apu.DMC.SampleAddress = (0xC000 | (uint16(Value) << 6))

	//apu.DMC.SampleLength = ((uint16(Value) << 4) | 1)
}

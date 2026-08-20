package bus

import (
	"fmt"
	"mtt/timenes/common"
	"mtt/timenes/nes/apu"
	"mtt/timenes/nes/cartridge"
	"mtt/timenes/nes/cartridge/mappers"
	"mtt/timenes/nes/input"
	"mtt/timenes/nes/ppu"
)

type BUS struct {
}

var CPUBus byte
var OAMBusAddress byte
var OutsideCodeRead, OutsideCodeWrite uint16 = 0, 0

// Read from Address, and return that byte
func Read(Address uint16) byte {
	//var CPUBus byte
	if Address < 0x2000 {
		//Read from RAM (Accounting for RAM Mirroring)
		CPUBus = cartridge.RAM[Address&0x7FF]
	} else if Address < 0x4000 {
		//Reading a PPU Register
		Address &= 0x2007
		switch Address {
		case 0x2000: //PPUCTRL
		case 0x2001: //PPUMASK
		case 0x2002: //PPUSTATUS
			ppu.PPUBus &= 0x1F
			ppu.PPUBus |= byte(common.Ternary(ppu.PPUSTATUS_VBlank, 0x80, 0x00))
			ppu.PPUBus |= byte(common.Ternary(ppu.PPUSTATUS_SpriteZeroHit, 0x40, 0x00))
			ppu.PPUBus |= byte(common.Ternary(ppu.PPUSTATUS_Overflow, 0x20, 0x00))
			ppu.PPUSTATUS_VBlank = false
			ppu.WriteLatch = false
			ppu.UpdatePPUBus2002(ppu.PPUBus)
		case 0x2003: //OAM ADDR
		case 0x2004: //OAMDATA
			if ppu.PPUScanline < 240 && (ppu.PPUMASK_RenderBG || ppu.PPUMASK_RenderSprites) {
				//Return buffer?
				if ppu.PPUDot == 0 || ppu.PPUDot > 320 {
					ppu.PPUBus = ppu.SecondaryOAM[0]
				} else if ppu.PPUDot > 0 && ppu.PPUDot <= 64 {
					ppu.PPUBus = 0xFF
				} else {
					CPUBus = ppu.PPUBus
				}
			} else {
				ppu.PPUBus = ppu.OAM[OAMBusAddress]
			}
			ppu.UpdatePPUBus(ppu.PPUBus)
		case 0x2005: //PPUSCROLL
		case 0x2006: //PPUADDR
		case 0x2007: //PPUDATA
			ppu.PPUAddressBus = ppu.VRAMAddress
			if (ppu.VRAMAddress & 0x3FFF) > 0x3F00 {
				//Palette data
				data := ppu.ReadPPU()
				ppu.PPUBus = (ppu.PPUBus & 0xC0) | (data & byte(common.Ternary(ppu.PPUMASK_Greyscale, 0x30, 0x3F)))
				//ppu.PPUReadBuffer = ReadPPU((ppu.VRAMAddress & 0x2F00) | (ppu.VRAMAddress & 0xFF))
				//ppu.PPUReadBuffer = ppu.PPUBus
				ppu.UpdatePPUBus2007Palette(ppu.PPUBus)
			} else {
				ppu.PPUBus = ppu.PPUReadBuffer
				ppu.PPUReadBuffer = ppu.ReadPPU()
				ppu.UpdatePPUBus(ppu.PPUBus)
			}

			ppu.VRAMAddress += common.Ternary(ppu.PPUCTRL_VRAMInc32Mode, 0x20, 0x01)
			ppu.VRAMAddress &= 0x3FFF
		default:
			CPUBus = ppu.PPUBus
		}
		CPUBus = ppu.PPUBus
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
		CPUBus = apuStatus
	} else if Address == 0x4016 { //Controller 1
		cBit := byte((input.Controller1ShiftRegister & 0x80) >> 7)
		input.Controller1ShiftRegister <<= 1
		CPUBus = cBit
	} else if Address == 0x4017 { //Controller 2
		cBit := byte((input.Controller2ShiftRegister & 0x80) >> 7)
		input.Controller2ShiftRegister <<= 1
		CPUBus = cBit
	} else if Address >= 0x4018 && Address <= 0x401F {
		//APU and I/O functionality that is normally disabled.
	} else if Address < 0x7FFF {
		//Could be PRG-RAM on the cartridge
		switch cartridge.MapperChipID {
		case 1: //MMC1
			CPUBus = mappers.MMC1_PRGRAM[Address&0x1FFF]
		case 4: //MMC3
			CPUBus = mappers.MMC3_PRGRAM[Address&0x1FFF]
		default:
			//outsideCodeRead = Address

		}
	} else if Address >= 0x8000 {
		//Read from ROM
		switch cartridge.MapperChipID {
		case 1: //MMC1
			CPUBus = cartridge.PRGROM[mappers.MMC1_FetchCPUAddress(Address, cartridge.PRGROM_Size)]
		case 2: //UxROM
			CPUBus = cartridge.PRGROM[mappers.UxROM_FetchCPUAddress(Address, cartridge.PRGROM_Size)]
		//case 3: //CNROM
		case 4: //MMC3
			CPUBus = cartridge.PRGROM[mappers.MMC3_FetchCPUAddress(Address, cartridge.PRGROM_Size)]
		case 7: //AxROM
			CPUBus = cartridge.PRGROM[mappers.AxROM_FetchCPUAddress(Address, cartridge.PRGROM_Size)]
		default:
			CPUBus = cartridge.PRGROM[(Address-0x8000)&uint16(cartridge.PRGROM_Size-1)]
			//outsideCodeRead = Address
		}
	}
	//MasterClockTick("READ")
	return CPUBus
}

// Write the Value into the Address given (PPU may have extra steps)
func Write(Address uint16, Value byte) {
	if Address < 0x2000 {
		cartridge.RAM[Address&0x7FF] = Value
	} else if Address < 0x4000 {
		//Write to PPU Register
		ppu.UpdatePPUBus(Value)
		Address &= 0x2007
		switch Address {
		case 0x2000: //PPUCTRL
			ppu.PPUCTRL_NametableSelect = Value & 0x03
			ppu.PPUCTRL_VRAMInc32Mode = (Value & 0x04) != 0
			ppu.PPUCTRL_SpritePatternTable = (Value & 0x08) != 0
			ppu.PPUCTRL_BGPatternTable = (Value & 0x10) != 0
			ppu.PPUCTRL_Use8x16Sprites = (Value & 0x20) != 0
			ppu.PPUCTRL_EnableNMI = (Value & 0x80) != 0

			ppu.TransferAddress = (uint16(ppu.PPUCTRL_NametableSelect) << 10) | (uint16(ppu.TransferAddress) & 0x73FF)
		case 0x2001: //PPUMASK
			ppu.PPUMASK_Greyscale = (Value & 0x01) != 0
			ppu.PPUMASK_8pxMaskBG = (Value & 0x02) != 0
			ppu.PPUMASK_8pxMaskSprites = (Value & 0x04) != 0
			ppu.PPUMASK_RenderBG = (Value & 0x08) != 0
			ppu.PPUMASK_RenderSprites = (Value & 0x10) != 0
			//NTSC scanline stuff
			ppu.PPUMASK_EmphasisRed = (Value & 0x20) != 0
			ppu.PPUMASK_EmphasisGreen = (Value & 0x40) != 0
			ppu.PPUMASK_EmphasisBlue = (Value & 0x80) != 0
		case 0x2002: //PPUSTATUS
		case 0x2003: //OAMADDR
			OAMBusAddress = Value
		case 0x2004: //OAMDATA
			if ((ppu.PPUScanline >= 240 && ppu.PPUScanline < 261) && (ppu.PPUMASK_RenderBG || ppu.PPUMASK_RenderSprites)) || (!ppu.PPUMASK_RenderBG && !ppu.PPUMASK_RenderSprites) {
				if (OAMBusAddress & 3) == 2 {
					Value &= 0xE3
				}
				ppu.OAM[OAMBusAddress] = Value
				OAMBusAddress++
			} else {
				OAMBusAddress += 4
				OAMBusAddress &= 0xFC
			}
		case 0x2005: //PPUSCROLL
			if !ppu.WriteLatch {
				ppu.PPUScrollFineX = byte(Value & 7)
				ppu.TransferAddress = uint16((ppu.TransferAddress & 0b0111111111100000) | uint16(Value>>3))
			} else {
				ppu.TransferAddress = ((ppu.TransferAddress & 0b0000110000011111) | uint16(uint16(Value&0xF8)<<2) | uint16(uint16(Value&7)<<12) /*| (uint16(ppu.PPUCTRL_NametableSelect&1) << 10)*/)
			}
			ppu.WriteLatch = !ppu.WriteLatch
		case 0x2006: //PPUADDR
			if !ppu.WriteLatch {
				//First write sets the high byte
				//Tempppu.VRAMAddress = (uint16(Value&0x7F) << 8)
				ppu.TransferAddress = uint16((ppu.TransferAddress & 0xFF) | uint16(Value&0x3F)<<8)
				//The actual ppu.VRAMAddress isn't changed until the 2nd write
			} else {
				//Second write sets the low byte
				ppu.TransferAddress = ((ppu.TransferAddress & 0xFF00) | uint16(Value))
				ppu.VRAMAddress = ppu.TransferAddress /*& 0x3FFF*/
			}
			ppu.WriteLatch = !ppu.WriteLatch
		case 0x2007: //PPUDATA
			WritePPU(Value)

			ppu.VRAMAddress += common.Ternary(ppu.PPUCTRL_VRAMInc32Mode, 0x20, 0x01)
			ppu.VRAMAddress &= 0x3FFF
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
			apu.Pulse1.Sweep.Shift = (Value & 0x07)
			apu.Pulse1.Sweep.ReloadFlag = true
		case 0x4002:
			apu.Pulse1.TimerReloadValue = apu.SetTimerLow(Value, apu.Pulse1.TimerReloadValue)
			apu.Pulse1.Timer = apu.Pulse1.TimerReloadValue
		case 0x4003:
			apu.Pulse1.TimerReloadValue = apu.SetTimerHi(Value, apu.Pulse1.TimerReloadValue)
			apu.Pulse1.Timer = apu.Pulse1.TimerReloadValue
			if apu.Pulse1.Enabled {
				apu.Pulse1.LengthCounter.Counter = apu.LengthCounterLoad(Value >> 3)
			}
			apu.Pulse1.StartFlag = true
			apu.Pulse1.DutyPos = 0

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
			apu.Pulse2.Sweep.Shift = (Value & 0x07)
			apu.Pulse2.Sweep.ReloadFlag = true
		case 0x4006:
			apu.Pulse2.TimerReloadValue = apu.SetTimerLow(Value, apu.Pulse2.TimerReloadValue)
			apu.Pulse2.Timer = apu.Pulse2.TimerReloadValue
		case 0x4007:
			apu.Pulse2.TimerReloadValue = apu.SetTimerHi(Value, apu.Pulse2.TimerReloadValue)
			apu.Pulse2.Timer = apu.Pulse2.TimerReloadValue
			if apu.Pulse2.Enabled {
				apu.Pulse2.LengthCounter.Counter = apu.LengthCounterLoad(Value >> 3)
			}
			apu.Pulse2.StartFlag = true
			apu.Pulse1.DutyPos = 0

		// Triangle
		case 0x4008:
			apu.Triangle.LengthCounter.HaltFlag = ((Value & 0x80) >> 7) != 0
			if apu.Triangle.Enabled {
				apu.Triangle.LinearCounter.ReloadValue = (Value & 0x7F)
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
			apu.Noise.StartFlag = true

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
				ppu.OAM[OAMBusAddress] = Read((uint16(Value) << 8) + uint16(i))
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

			if apu.DMC.Enabled {
				if apu.DMC.BytesRemaining > 0 {
					apu.DMC.DMCRestartSample()
				}
			} else {
				apu.DMC.BytesRemaining = 0
			}

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
	} else if Address >= 0x6000 && Address < 0x8000 {
		//Check for PRG-RAM
		switch cartridge.MapperChipID {
		case 1: //MMC1
			mappers.MMC1_WriteToPRGRAM(Value, Address)
		case 3: //CNROM
			prgValue := Read(Address)
			mappers.CNROM_WriteToPRGRAM(Value&prgValue, Address) // Can have bus conflicts
		case 4: //MMC3
			mappers.MMC3_WriteToPRGRAM(Value, Address)
		default:
			//OutsideCodeWrite = Address
			fmt.Println("Write to unused memory addresss: $" + fmt.Sprintf("%04X", Address))
		}
	} else if Address >= 0x8000 { //Account for Mapper chips
		//Check what mapper chip we're using
		switch cartridge.MapperChipID {
		case 1: //MMC1
			mappers.MMC1_Write(Value, Address, common.CPU_TotalCycles)
		case 2: //UxROM
			mappers.UxROM_Write(Value, Address)
		case 3: //CNROM
			prgValue := Read(Address)
			mappers.CNROM_Write(Value&prgValue, Address) // Can have bus conflicts
		case 4: //MMC3
			mappers.MMC3_Write(Value, Address)
		case 7: //AxROM
			mappers.AxROM_Write(Value, Address)
		default:
			//OutsideCodeWrite = Address
			fmt.Println("Write to unused memory addresss: $" + fmt.Sprintf("%04X", Address))
		}
	}
	//MasterClockTick("WRITE")
}

func WritePPU(Value byte) {
	if ppu.VRAMAddress < 0x2000 {
		//Write to pattern table. (If the cartridge supports it)
		if cartridge.CHRROM_Size == 0 {
			switch cartridge.MapperChipID {
			case 1: //MMC1
				//if AltNametableLayout {
				//	mappers.MMC1_VRAM[mappers.MMC1_FetchPPUAddress(ppu.VRAMAddress, CHRROM_Size)] = Value
				//} else {
				cartridge.CHRROM[mappers.MMC1_FetchPPUAddress(ppu.VRAMAddress, cartridge.CHRROM_Size)] = Value
				//}
			case 3: //CNROM
				cartridge.CHRROM[mappers.CNROM_FetchPPUAddress(ppu.VRAMAddress)] = Value
			case 4: //MMC3
				cartridge.CHRROM[mappers.MMC3_FetchPPUAddress(ppu.VRAMAddress, cartridge.CHRROM_Size)] = Value
			default:
				cartridge.CHRROM[ppu.VRAMAddress] = Value
			}
			//For Pattern Table viewing
			common.PendingPatternUpdate = true
		}
		//else, nothing happens because it's CHR-ROM
	} else if ppu.VRAMAddress < 0x3F00 {
		//Write to the Nametables

		switch cartridge.MapperChipID {
		case 1: //MMC1
			switch mappers.MMC1_GetNametableArrangement() {
			case 0: //Screen A only
				cartridge.CartVRAM[int(ppu.VRAMAddress&0x3FF)] = Value
			case 1: //Screen B Only
				cartridge.CartVRAM[int(ppu.VRAMAddress&0x3FF)+0x400] = Value
			case 2: //Horizontal
				cartridge.CartVRAM[int(ppu.VRAMAddress&0x7FF)] = Value
			case 3: //Vertical
				cartridge.CartVRAM[int(ppu.VRAMAddress&0x3FF)|int(ppu.VRAMAddress&0x800)>>1] = Value
			}
			//cartridge.CartVRAM[ppu.VRAMAddress&0xFFF] = Value
		//case 3: //CNROM
		case 4: //MMC3
			//WriteToNametable(Value, mappers.MMC3_IsHorizontalNametable)
			cartridge.CartVRAM[int(ppu.VRAMAddress&0xFFF)] = Value
		default:
			WriteToNametable(Value, cartridge.IsNametableHorizontal)
		}

	} else {
		//Write to Palette RAM
		if (ppu.VRAMAddress & 3) == 0 {
			cartridge.PaletteRAM[ppu.VRAMAddress&0x0F] = Value
		} else {
			cartridge.PaletteRAM[ppu.VRAMAddress&0x1F] = Value
		}
	}
}

func WriteToNametable(Value byte, isHoriz bool) {
	if cartridge.AltNametableLayout {
	} else {
		if cartridge.IsNametableHorizontal {
			// Horizontal Mirroring
			cartridge.VRAM[int(ppu.VRAMAddress&0x3FF)|int(ppu.VRAMAddress&0x800)>>1] = Value
		} else {
			//Vertical Mirroring
			cartridge.VRAM[int(ppu.VRAMAddress&0x7FF)] = Value
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

package bus

import (
	"fmt"
	"mtt/timenes/common"
	"mtt/timenes/nes/cartridge"
	"mtt/timenes/nes/cartridge/mappers"
	"mtt/timenes/nes/input"
	"mtt/timenes/nes/ppu"
)

type BUS struct {
	cpu CPU
	apu APU
}

type CPU interface {
	DelayCPU(int)
}

type APU interface {
	ReadAPU(uint16) byte
	WriteAPU(uint16, byte)
}

var CPUBus byte
var OutsideCodeRead, OutsideCodeWrite uint16 = 0, 0

func NewBUS() *BUS {
	bus := BUS{}
	return &bus
}

func (b *BUS) SetCPU(c CPU) {
	b.cpu = c
}

func (b *BUS) SetAPU(a APU) {
	b.apu = a
}

// Read from Address, and return that byte
func (b *BUS) Read(Address uint16) byte {
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
				ppu.PPUBus = ppu.OAM[ppu.GetOAMBusAddress()]
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

		CPUBus = b.apu.ReadAPU(Address)
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
func (b *BUS) Write(Address uint16, Value byte) {
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
			ppu.SetOAMBusAddress(Value)
		case 0x2004: //OAMDATA
			oamAddr := ppu.GetOAMBusAddress()
			if ((ppu.PPUScanline >= 240 && ppu.PPUScanline < 261) && (ppu.PPUMASK_RenderBG || ppu.PPUMASK_RenderSprites)) || (!ppu.PPUMASK_RenderBG && !ppu.PPUMASK_RenderSprites) {
				if (oamAddr & 3) == 2 {
					Value &= 0xE3
				}
				ppu.OAM[oamAddr] = Value
				oamAddr++
			} else {
				oamAddr += 4
				oamAddr &= 0xFC
			}
			ppu.SetOAMBusAddress(oamAddr)
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
			ppu.WritePPU(Value)

			ppu.VRAMAddress += common.Ternary(ppu.PPUCTRL_VRAMInc32Mode, 0x20, 0x01)
			ppu.VRAMAddress &= 0x3FFF
		}

	} else if Address == 0x4014 { //OAMDMA
		oamAddr := ppu.GetOAMBusAddress()
		for i := 0; i < 256; i++ {
			ppu.OAM[oamAddr] = b.Read((uint16(Value) << 8) + uint16(i))
			oamAddr++
		}
		if common.CPU_TotalCycles%2 == 1 {
			b.cpu.DelayCPU(514)
		} else {
			b.cpu.DelayCPU(513)
		}
	} else if Address == 0x4016 { //Controller Input
		input.UpdateControllers()
	} else if Address < 0x4018 { //APU and I/O
		b.apu.WriteAPU(Address, Value)
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
			prgValue := b.Read(Address)
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
			prgValue := b.Read(Address)
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

//func (b *BUS) PrebufferDMCSamples() {
//
//	if apu.DMC.SampleLength > 0 && apu.DMC.SampleAddress > 0x7FFF {
//		apu.DMC.SampleBuffer = nil
//
//		//Get the entire sample upfront, because God is dead
//		for i := uint16(0); i < apu.DMC.SampleLength; i++ {
//			sampleAddr := apu.DMC.SampleAddress + i
//			if sampleAddr < 0x8000 { //We hit 0xFFFF and wrapped around
//				sampleAddr += 0x8000
//			}
//
//			apu.DMC.SampleBuffer = append(apu.DMC.SampleBuffer, b.Read(sampleAddr))
//		}
//		apu.DMC.SampleBufferPos = 0
//	}
//
//	//apu.DMC.SampleAddress = (0xC000 | (uint16(Value) << 6))
//
//	//apu.DMC.SampleLength = ((uint16(Value) << 4) | 1)
//}

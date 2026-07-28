package main

// Read from Address, and return that byte
func Read(Address uint16) byte {
	if Address < 0x2000 {
		//Read from RAM (Accounting for RAM Mirroring)
		return RAM[Address&0x7FF]
	} else if Address < 0x4000 {
		//Reading a PPU Register
		Address &= 0x2007
		switch Address {
		case 0x2000:
		case 0x2001:
		case 0x2002: //PPUSTATUS
			ppustatus := byte(0)
			ppustatus |= byte(ternary(ppuVBlank, 0x80, 0x00))
			ppustatus |= byte(ternary(ppuStatusSprZeroHit, 0x40, 0x00))
			ppustatus |= byte(ternary(ppuStatusOverflow, 0x20, 0x00))
			ppuVBlank = false
			WriteLatch = false
			return ppustatus
		case 0x2003:
		case 0x2004:
		case 0x2005:
		case 0x2006:
		case 0x2007: //PPUDATA
			temp := PPUReadBuffer
			if VRAMAddress > 0x3F00 {
				temp = ReadPPU(VRAMAddress)
			} else {
				PPUReadBuffer = ReadPPU(VRAMAddress)
			}

			VRAMAddress += ternary(ppuVRAMInc32Mode, 0x20, 0x01)
			VRAMAddress &= 0x3FFF
			return temp

		default:
			return 0
		}
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
	}
	return 0
}

// Write the Value into the Address given (PPU may have extra steps)
func Write(Address uint16, Value byte) {
	if Address < 0x2000 {
		RAM[Address&0x7FF] = Value
	} else if Address < 0x4000 {
		//Write to PPU Register
		Address &= 0x2007
		switch Address {
		case 0x2000: //PPUCTRL
			ppuNametableSelect = Value & 3
			ppuVRAMInc32Mode = (Value & 4) != 0
			ppuSpritePatternTable = (Value & 8) != 0
			ppuBGPatternTable = (Value & 0x10) != 0
			ppuUse8x16Sprites = (Value & 0x20) != 0
			ppuEnableNMI = (Value & 0x80) != 0

			TransferAddress |= (uint16(ppuNametableSelect) << 10)
		case 0x2001: //PPUMASK
			ppuMask_8pxMaskBG = (Value & 0x02) != 0
			ppuMask_8pxMaskSprites = (Value & 0x04) != 0
			ppuMask_RenderBG = (Value & 0x08) != 0
			ppuMask_RenderSprites = (Value & 0x10) != 0
			//ppuMask_EmphasisRed = (Value & 0x20) != 0
			//ppuMask_EmphasisGreen = (Value & 0x40) != 0
			//ppuMask_EmphasisBlue = (Value & 0x80) != 0
		case 0x2002: //PPUSTATUS
		case 0x2003: //OAMADDR
		case 0x2004: //OAMDATA
		case 0x2005: //PPUSCROLL
			if !WriteLatch {
				ppuScrollFineX = byte(Value & 7)
				TempVRAMAddress = uint16((TempVRAMAddress & 0b0111111111100000) | uint16(Value>>3))
			} else {
				TransferAddress = ((TempVRAMAddress & 0b0000110000011111) | uint16((uint16(Value&0xF8)<<2)|(uint16(Value&7)<<12)))
			}
			WriteLatch = !WriteLatch
		case 0x2006: //PPUADDR
			if !WriteLatch {
				//First write sets the high byte
				TempVRAMAddress = (uint16(Value&0x3F) << 8)
				//The actual VRAMAddress isn't changed until the 2nd write
			} else {
				//Second write sets the low byte
				VRAMAddress = (TempVRAMAddress | uint16(Value))
				TransferAddress = VRAMAddress
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

			VRAMAddress += ternary(ppuVRAMInc32Mode, 0x20, 0x01)
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
			apuPulse1.LengthCounterLoad = ((Value & 0xF8) >> 3)

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
			apuPulse2.LengthCounterLoad = ((Value & 0xF8) >> 3)

		// Triangle
		case 0x4008:
			apuTriangle.LengthCounterControl = ((Value & 0x80) >> 7) != 0
			apuTriangle.LinearCounterLoad = (Value & 0x7F)
		case 0x4009: //Unused
		case 0x400A:
			apuTriangle.Timer = uint16(Value)
		case 0x400B:
			apuTriangle.Timer |= (uint16(Value&0x7) << 8)
			apuTriangle.LengthCounterLoad = ((Value & 0xF8) >> 3)

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
			apuNoise.LengthCounterLoad = ((Value & 0xF8) >> 3)
		// DMC
		case 0x4010:
			apuDMC.IRQEnable = ((Value & 0x80) >> 7) != 0
		case 0x4011:
		case 0x4012:
		case 0x4013:

		case 0x4014: //OAM
			for i := 0; i < 256; i++ {
				OAM[i] = Read((uint16(Value) << 8) + uint16(i))
			}
		case 0x4015:
		case 0x4016: //Controller Input
			Controller1ShiftRegister = uint16(Controller1)
			Controller2ShiftRegister = uint16(Controller2)
		case 0x4017: //APU Frame Counter control
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

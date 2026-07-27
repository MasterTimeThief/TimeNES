package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var filepath string = "roms/smb.nes"

//var filepath string = "roms/nestest.nes"

//var filepath string = "roms/AccuracyCoin.nes"

//type App struct{ Clicks int }

const (
	screenWidth  = 256
	screenHeight = 240
	screenScale  = 3
)

type Game struct {
	gameScreen *image.RGBA
	keys       []ebiten.Key
}

var ProgramCounter uint16
var AddressBus, ppuAddressBus uint16
var StackPointer byte
var A, X, Y byte
var opcode byte
var operands []byte
var Cycles, TotalCycles int
var CPU_Halted = false
var flag_Carry, flag_Zero, flag_InterruptDisable, flag_Decimal, flag_Overflow, flag_Negative bool
var MagicConstant byte = 0xFF

var RAM [0x800]byte
var VRAM [0x800]byte
var PaletteRAM [0x20]byte
var ROM [0x8000]byte
var PRGROMSize uint16 = 0x4000

// var CHRData [0x2000]byte
var CHRROM [0x2000]byte
var Header [0x10]byte
var CartRAM [0x2000]byte

var LoggingCPU = false
var LoggingPPU = false
var LogCount = 1000
var InstructionCount int = 0
var frame int
var tracePC, traceVRAM uint16
var traceA, traceX, tracyY, traceSP byte
var traceFlagC, traceFlagZ, traceFlagI, traceFlagD, traceFlagV, traceFlagN bool
var traceCycles int

//var screenScale float32 = 3
//var nametable *image.RGBA = image.NewRGBA(image.Rect(0, 0, ScreenWidth, ScreenHeight))

var WriteLatch bool        //PPU's w register
var TransferAddress uint16 //PPU's t register
var VRAMAddress uint16     //PPU's v register
var TempVRAMAddress uint16 //PPU's v register (temporary)
var PPUReadBuffer byte
var NMILevelDetector, DoNMI bool

var ppuDot int      //The X position of the scanning beam
var ppuScanline int //The Y position of the scanning beam
var ppuVBlank, ppuHBlank bool
var ppuMask_8pxMaskBG, ppuMask_8pxMaskSprites, ppuMask_RenderBG, ppuMask_RenderSprites bool
var ppuNametableSelect byte
var ppuVRAMInc32Mode, ppuSpritePatternTable, ppuBGPatternTable, ppuUse8x16Sprites, ppuEnableNMI bool
var ppuShiftRegister_patternL, ppuShiftRegister_patternH, ppuShiftRegister_attributeL, ppuShiftRegister_attributeH uint16
var ppu8Step_patternLowBitPlane, ppu8Step_patternHighBitPlane, ppu8Step_attribute, ppu8Step_NextCharacter, ppu8Step_temp byte
var ppuScrollFineX byte
var DrawNewFrame bool = false

var screenCount int = 33

var Controller1, Controller2 byte
var Controller1ShiftRegister, Controller2ShiftRegister uint16

func (g *Game) Update() error {

	for {
		UpdateControllers(g)
		Emulate_CPU(g)
		if DrawNewFrame {
			DrawNewFrame = false
			return nil
		}
		if CPU_Halted {
			return ebiten.Termination
		}
	}

	//fmt.Println("CPU Paused, Drawing Window")

	//DrawWindow(ops, window)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// This graphics context is used for managing the rendering state.

	//if DrawNewFrame {
	screen.WritePixels(g.gameScreen.Pix)
	//DrawNewFrame = false
	//}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\tFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
	//ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth*screenScale, screenHeight*screenScale)
	ebiten.SetWindowTitle("TimeNES (MasterTimeThief made this :3)")
	ebiten.SetTPS(ebiten.SyncWithFPS)

	//Reset the console
	Reset()

	g := &Game{
		gameScreen: image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}

// Sets everything outside of operands for the tracelogger to run later
func prepTraceLogger() {
	tracePC = ProgramCounter
	traceVRAM = VRAMAddress

	traceA = A
	traceX = X
	tracyY = Y
	traceSP = StackPointer

	traceFlagC = flag_Carry
	traceFlagZ = flag_Zero
	traceFlagI = flag_InterruptDisable
	traceFlagD = flag_Decimal
	traceFlagV = flag_Overflow
	traceFlagN = flag_Negative

	traceCycles = TotalCycles
}

func TraceLogger() {
	OpCodeNames := [...]string{
		"BRK", "ORA", "HLT", "SLO", "NOP", "ORA", "ASL", "SLO", "PHP", "ORA",
		"ASL", "ANC", "NOP", "ORA", "ASL", "SLO", "BPL", "ORA", "HLT", "SLO",
		"NOP", "ORA", "ASL", "SLO", "CLC", "ORA", "NOP", "SLO", "NOP", "ORA",
		"ASL", "SLO", "JSR", "AND", "HLT", "RLA", "BIT", "AND", "ROL", "RLA",
		"PLP", "AND", "ROL", "ANC", "BIT", "AND", "ROL", "RLA", "BMI", "AND",
		"HLT", "RLA", "NOP", "AND", "ROL", "RLA", "SEC", "AND", "NOP", "RLA",
		"NOP", "AND", "ROL", "RLA", "RTI", "EOR", "HLT", "SRE", "NOP", "EOR",
		"LSR", "SRE", "PHA", "EOR", "LSR", "ALR", "JMP", "EOR", "LSR", "SRE",
		"BVC", "EOR", "HLT", "SRE", "NOP", "EOR", "LSR", "SRE", "CLI", "EOR",
		"NOP", "SRE", "NOP", "EOR", "LSR", "SRE", "RTS", "ADC", "HLT", "RRA",
		"NOP", "ADC", "ROR", "RRA", "PLA", "ADC", "ROR", "ARR", "JMP", "ADC",
		"ROR", "RRA", "BVS", "ADC", "HLT", "RRA", "NOP", "ADC", "ROR", "RRA",
		"SEI", "ADC", "NOP", "RRA", "NOP", "ADC", "ROR", "RRA", "NOP", "STA",
		"NOP", "SAX", "STY", "STA", "STX", "SAX", "DEY", "NOP", "TXA", "ANE",
		"STY", "STA", "STX", "SAX", "BCC", "STA", "HLT", "SHA", "STY", "STA",
		"STX", "SAX", "TYA", "STA", "TXS", "SHS", "SHY", "STA", "SHX", "SHA",
		"LDY", "LDA", "LDX", "LAX", "LDY", "LDA", "LDX", "LAX", "TAY", "LDA",
		"TAX", "LXA", "LDY", "LDA", "LDX", "LAX", "BCS", "LDA", "HLT", "LAX",
		"LDY", "LDA", "LDX", "LAX", "CLV", "LDA", "TSX", "LAE", "LDY", "LDA",
		"LDX", "LAX", "CPY", "CMP", "NOP", "DCP", "CPY", "CMP", "DEC", "DCP",
		"INY", "CMP", "DEX", "AXS", "CPY", "CMP", "DEC", "DCP", "BNE", "CMP",
		"HLT", "DCP", "NOP", "CMP", "DEC", "DPC", "CLD", "CMP", "NOP", "DCP",
		"NOP", "CMP", "DEC", "DCP", "CPX", "SBC", "NOP", "ISC", "CPX", "SBC",
		"INC", "ISC", "INX", "SBC", "NOP", "SBC", "CPX", "SBC", "INC", "ISC",
		"BEQ", "SBC", "HLT", "ISC", "NOP", "SBC", "INC", "ISC", "SED", "SBC",
		"NOP", "ISC", "NOP", "SBC", "INC", "ISC",
	}
	if LoggingCPU /*&& InstructionCount > 35000*/ {
		Traceline := "$" + fmt.Sprintf("%04X", tracePC)

		Traceline += "\t" + fmt.Sprintf("%02X ", opcode)
		opLength := len(operands)
		for i := range operands {
			Traceline += fmt.Sprintf("%02X ", operands[i])
		}

		if opLength < 2 {
			Traceline += "\t"
		}

		Traceline += "\t" + OpCodeNames[opcode]

		//Add operand / addresses after opcode

		if opcode == 0x99 {
			Traceline += " $" + fmt.Sprintf("%02X", operands[1]) + fmt.Sprintf("%02X", operands[0])
			Traceline += ", Y -> $"
			Traceline += fmt.Sprintf("%04X", BuildAddress(operands[0], operands[1])+uint16(Y))
		}

		Traceline += "\tA:" + fmt.Sprintf("%02X", traceA)
		Traceline += "\tX:" + fmt.Sprintf("%02X", traceX)
		Traceline += "\tY:" + fmt.Sprintf("%02X", tracyY)
		Traceline += "\tSP:" + fmt.Sprintf("%02X", traceSP)
		Traceline += "\t"

		//Processor Flags
		if traceFlagN {
			Traceline += "N"
		} else {
			Traceline += "n"
		}

		if traceFlagV {
			Traceline += "V"
		} else {
			Traceline += "v"
		}

		Traceline += "--"

		if traceFlagD {
			Traceline += "D"
		} else {
			Traceline += "d"
		}

		if traceFlagI {
			Traceline += "I"
		} else {
			Traceline += "i"
		}

		if traceFlagZ {
			Traceline += "Z"
		} else {
			Traceline += "z"
		}

		if traceFlagC {
			Traceline += "C"
		} else {
			Traceline += "c"
		}

		Traceline += "\tVRAM: " + fmt.Sprintf("%04X", traceVRAM)
		Traceline += "\tCycle: " + fmt.Sprintf("%d", traceCycles)
		Traceline += "\tInstructionCount: " + fmt.Sprintf("%d", InstructionCount)

		fmt.Println(Traceline)
		/*LogCount--
		if LogCount < 0 {
			CPU_Halted = true
		}*/
	}
}

func TraceLoggerPPU() {

	if LoggingPPU {

		Traceline := fmt.Sprintf("D: %-3d", ppuDot)
		Traceline += fmt.Sprintf("\tSL: %-3d", ppuScanline)

		//Add Shift registers

		Traceline += fmt.Sprintf("\tpL: %016b", ppuShiftRegister_patternL)
		Traceline += fmt.Sprintf("\tpH: %016b", ppuShiftRegister_patternH)
		Traceline += fmt.Sprintf("\taL: %016b", ppuShiftRegister_attributeL)
		Traceline += fmt.Sprintf("\taH: %016b", ppuShiftRegister_attributeH)

		//Print PalHi and PalLow
		col0 := byte((ppuShiftRegister_patternL >> (15 - ppuScrollFineX)) & 1)
		col1 := byte((ppuShiftRegister_patternH >> (15 - ppuScrollFineX)) & 1)

		pal0 := byte((ppuShiftRegister_attributeL >> (15 - ppuScrollFineX)) & 1)
		pal1 := byte((ppuShiftRegister_attributeH >> (15 - ppuScrollFineX)) & 1)

		Traceline += fmt.Sprintf("\tLo: %02b Hi: %02b", byte((col1<<1)|col0), byte((pal1<<1)|pal0))

		fmt.Println(Traceline)
		/*LogCount--
		if LogCount < 0 {
			CPU_Halted = true
		}*/
	}
}

var CartRamLastString string

func CartRAMLogger() {
	CartTestStatus := CartRAM[0]
	CartTestValid := CartRAM[1] == 0xDE && CartRAM[2] == 0xB0 && CartRAM[3] == 0x61

	newString := string(CartRAM[4:])

	if CartTestValid && CartTestStatus == 0x80 && (newString != CartRamLastString) {
		//Valid Test line
		CartRamLastString = newString
		fmt.Print(newString)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ternary(Condition bool, ValT, ValF uint16) uint16 {
	if Condition {
		return ValT
	}
	return ValF
}

//var Controller1State [8]bool
//var Controller2State [8]bool

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
			ppuMask_8pxMaskBG = (Value & 2) != 0
			ppuMask_8pxMaskSprites = (Value & 4) != 0
			ppuMask_RenderBG = (Value & 8) != 0
			ppuMask_RenderSprites = (Value & 0x10) != 0
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
	} else if Address == 0x4014 { //OAM
		for i := 0; i < 256; i++ {
			OAM[i] = Read((uint16(Value) << 8) + uint16(i))
		}
	} else if Address == 0x4016 { //Controller Input
		Controller1ShiftRegister = uint16(Controller1)
		Controller2ShiftRegister = uint16(Controller2)
	} else if Address < 0x4020 {
		//Audio Processing Unit stuff
		//$4000 - $4017 is APU and I/O registers
		//$4018 - $401F is APU and I/O functions that are normally disabled

	} else if Address < 0x8000 {
		CartRAM[Address&0x1FFF] = Value
	}
}

// Update states of controllers
func UpdateControllers(g *Game) {
	/*
		0 - B
		1 - A
		2 - Select
		3 - Start
		4 - Up
		5 - Down
		6 - Left
		7 - Right
	*/
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	Controller1 = 0
	Controller2 = 0
	for _, k := range g.keys {
		switch k {
		case ebiten.KeyX:
			Controller1 |= 0x80
		case ebiten.KeyZ:
			Controller1 |= 0x40
		case ebiten.KeyShiftRight:
			Controller1 |= 0x20
		case ebiten.KeyEnter:
			Controller1 |= 0x10
		case ebiten.KeyArrowUp:
			Controller1 |= 0x8
		case ebiten.KeyArrowDown:
			Controller1 |= 0x4
		case ebiten.KeyArrowLeft:
			Controller1 |= 0x2
		case ebiten.KeyArrowRight:
			Controller1 |= 0x1
		}
	}

	/*Controller2 |= byte(ternary(Controller2State[0], 0x80, 0x00))
	Controller2 |= byte(ternary(Controller2State[1], 0x40, 0x00))
	Controller2 |= byte(ternary(Controller2State[2], 0x20, 0x00))
	Controller2 |= byte(ternary(Controller2State[3], 0x10, 0x00))
	Controller2 |= byte(ternary(Controller2State[4], 0x08, 0x00))
	Controller2 |= byte(ternary(Controller2State[5], 0x04, 0x00))
	Controller2 |= byte(ternary(Controller2State[6], 0x02, 0x00))
	Controller2 |= byte(ternary(Controller2State[7], 0x01, 0x00))*/
}

func ReadPPU(Address uint16) byte {
	if Address < 0x2000 {
		//Read from pattern table.
		return CHRROM[Address]
		//else, nothing happens
	} else if Address < 0x3F00 {
		//Read from the Nametables
		if (Header[6] & 1) == 0 {
			// Horizontal Mirroring
			return VRAM[int(Address&0x3FF)|int(Address&0x800)>>1]
		} else {
			//Vertical Mirroring
			return VRAM[int(Address&0x7FF)]
		}
	} else {
		//Read from Palette RAM
		if (Address & 3) == 0 {
			return PaletteRAM[Address&0x0F]
		} else {
			return PaletteRAM[Address&0x1F]
		}
	}
}

func Reset() {
	//var HeaderedROM []byte := os.ReadFile()
	flag_InterruptDisable = true
	HeaderedROM, err := os.ReadFile(filepath)
	check(err)

	copy(Header[:], HeaderedROM[0x0:])
	size := uint16(Header[4])
	copy(ROM[:], HeaderedROM[0x10:])
	if Header[5] != 0 {
		copy(CHRROM[:], HeaderedROM[(0x0010+(0x4000*size)):])
	}
	//copy(CHRData[:], HeaderedROM[0x8010:])

	PCL := Read(0xFFFC)
	PCH := Read(0xFFFD)
	ProgramCounter = BuildAddress(PCL, PCH)
	//ProgramCounter = 0xC000
	StackPointer -= 3
	//fmt.Printf("%#x", ProgramCounter)
	//Run()
}

func BuildAddress(Value_Low, Value_High byte) uint16 {
	//b := []byte{Value_Low, Value_High}
	//return binary.LittleEndian.Uint16(b[0:])

	return (uint16(Value_High)<<8 | uint16(Value_Low))

}

func SetZNFlags(Value byte) {
	flag_Zero = (Value == 0x00)
	flag_Negative = (Value >= 0x80)
}

func ReadFromPC() byte {
	Value := Read(ProgramCounter)
	if LoggingCPU {
		operands = append(operands, Value)
	}
	ProgramCounter++
	return Value
}

//TODO: Check the ReadOperandAbsolute functions to add a CPU cycle if the high byte was altered

func ReadOperands_AbsoluteAddressed() uint16 {
	AddressBus := uint16(ReadFromPC())
	AddressBus = (uint16(ReadFromPC())<<8 | AddressBus)
	return AddressBus
}

func ReadOperands_IndirectAddressed() uint16 {
	/*Addr := ReadFromPC()
	TempAddress := Addr
	Addr = Read(uint16(TempAddress)) //Low byte of new address
	TempAddress++
	AddressBus = (uint16(Read(uint16(TempAddress)))<<8 | uint16(Addr)) //High byte
	return AddressBus*/
	AddressBus := uint16(ReadFromPC())
	AddressBus = (uint16(ReadFromPC())<<8 | AddressBus)
	//Now read from HERE
	indL := Read(AddressBus)
	indH := Read(AddressBus + 1)
	return BuildAddress(indL, indH)
}

func ReadOperands_AbsoluteAddressed_XIndexed() uint16 {
	AddressBus := uint16(ReadFromPC())
	AddressBus = (uint16(ReadFromPC())<<8 | AddressBus)
	AddressBus += uint16(X)
	return AddressBus
}

func ReadOperands_AbsoluteAddressed_YIndexed() uint16 {
	AddressBus := uint16(ReadFromPC())
	AddressBus = (uint16(ReadFromPC())<<8 | AddressBus)
	AddressBus += uint16(Y)
	return AddressBus
}

func ReadOperands_IndirectAddressed_XIndexed() uint16 {
	Addr := ReadFromPC() + X
	TempAddress := Addr
	Addr = Read(uint16(TempAddress)) //Low byte of new address
	TempAddress++
	AddressBus = (uint16(Read(uint16(TempAddress)))<<8 | uint16(Addr)) //High byte
	return AddressBus
}

func ReadOperands_IndirectAddressed_YIndexed() uint16 {
	Addr := ReadFromPC()
	TempAddress := Addr
	Addr = Read(uint16(TempAddress)) //Low byte of new address
	TempAddress++
	AddressBus = (uint16(Read(uint16(TempAddress)))<<8 | uint16(Addr)) //High byte
	AddressBus += uint16(Y)
	return AddressBus
}

func ReadOperands_ZeroPageAddressed() uint16 {
	AddressBus := ReadFromPC()
	return BuildAddress(AddressBus, 0x00)
}

func ReadOperands_ZeroPageAddressed_XIndexed() uint16 {
	AddressBus := ReadFromPC()
	return BuildAddress(AddressBus+X, 0x00)
}

func ReadOperands_ZeroPageAddressed_YIndexed() uint16 {
	AddressBus := ReadFromPC()
	return BuildAddress(AddressBus+Y, 0x00)
}

func Branch(Condition bool, Value byte) {
	if Condition {
		signedVal := int(Value)
		if signedVal > 127 {
			signedVal -= 256 //range from -128 to 127
		}
		ProgramCounter = uint16(ProgramCounter + uint16(signedVal))
		Cycles = 3
		//TODO: If taking branch changes the high byte, add a cycle
	} else {
		Cycles = 2
	}
}

func Push(Value byte) {
	//Store to the stack, and decrement the stack pointer
	Write(uint16(StackPointer)+0x100, Value)
	StackPointer--
}

func Pull() byte {
	//Increment the stack pointer, and read from the stack
	StackPointer++
	temp := Read(uint16(StackPointer) + 0x100)
	return temp
}

func BoolToInt(Flag bool) int {
	if Flag {
		return 1
	}
	return 0
}

func PushFlags() {

}

func PullFlags() {

}

// Performs Arithmetic Shift Left onto value at Address
func Op_ASL(Address uint16) {
	Value := Read(Address)
	flag_Carry = (Value >= 0x80)
	Value <<= 1
	Write(Address, Value)
	SetZNFlags(Value)

}

// Performs Arithmetic Shift Right onto value at Address
func Op_LSR(Address uint16) {
	Value := Read(Address)
	flag_Carry = (Value & 1) != 0
	Value >>= 1
	Write(Address, Value)
	SetZNFlags(Value)
}

// Perform Rotate Left onto value at Address
func Op_ROL(Address uint16) {
	Value := Read(Address)
	futureCarry := (Value >= 0x80)
	Value <<= 1
	if flag_Carry {
		Value |= 1
	}
	Write(Address, Value)
	flag_Carry = futureCarry
	SetZNFlags(Value)
}

// Perform Rotate Right onto value at Address
func Op_ROR(Address uint16) {
	Value := Read(Address)
	futureCarry := (Value & 1) != 0
	Value >>= 1
	if flag_Carry {
		Value |= 0x80
	}
	Write(Address, Value)
	flag_Carry = futureCarry
	SetZNFlags(Value)
}

// Increment Value, and save to Address
func Op_INC(Address uint16, Value byte) {
	Value++
	Write(Address, Value)
	SetZNFlags(Value)
}

// Decrement Value, and save to Address
func Op_DEC(Address uint16, Value byte) {
	Value--
	Write(Address, Value)
	SetZNFlags(Value)
}

// Bitwise OR with A
func Op_ORA(Value byte) {
	A |= Value
	SetZNFlags(A)
}

// Bitwise AND with A
func Op_AND(Value byte) {
	A &= Value
	SetZNFlags(A)
}

// Bitwise XOR with A
func Op_EOR(Value byte) {
	A ^= Value
	SetZNFlags(A)
}

// Add Value to A with Carry
func Op_ADC(Value byte) {
	IntSum := int(A) + int(Value) + BoolToInt(flag_Carry)
	flag_Overflow = (^int(A^Value) & (int(A) ^ IntSum) & 0x80) != 0
	flag_Carry = IntSum > 0xFF
	A = byte(IntSum)
	SetZNFlags(A)
}

// Subtract Value from A with Carry
func Op_SBC(Value byte) {
	IntSum := int(A) - int(Value) - BoolToInt(!flag_Carry)
	flag_Overflow = (int(A^Value) & (int(A) ^ IntSum) & 0x80) != 0
	flag_Carry = IntSum >= 0x00
	A = byte(IntSum)
	SetZNFlags(A)
}

// Compare Value with A
func Op_CMP(Value byte) {
	flag_Carry = Value <= A
	flag_Zero = (Value == A)
	flag_Negative = ((A - Value) >= 0x80)
}

// Compare Value with X
func Op_CPX(Value byte) {
	flag_Carry = Value <= X
	flag_Zero = (Value == X)
	flag_Negative = ((X - Value) >= 0x80)
}

// Compare Value with Y
func Op_CPY(Value byte) {
	flag_Carry = Value <= Y
	flag_Zero = (Value == Y)
	flag_Negative = ((Y - Value) >= 0x80)
}

// Uhh
func Op_BIT(Value byte) {
	//Bit Test
	flag_Zero = ((A & Value) == 0)
	flag_Negative = ((Value & 0x80) != 0)
	flag_Overflow = ((Value & 0x40) != 0)
}

func Emulate_CPU(g *Game) {
	//Non Maskable Interrupt check
	prevNMILevelDetector := NMILevelDetector
	NMILevelDetector = (ppuEnableNMI && ppuVBlank)
	if !prevNMILevelDetector && NMILevelDetector {
		DoNMI = true
	}

	//Get the opcode
	//var opcode byte
	if !DoNMI {
		//If we're not running an NMI
		opcode = Read(ProgramCounter)
		if LoggingCPU {
			prepTraceLogger()
		}
		ProgramCounter++
	} else {
		//If we're running an NMI, force opcode $00
		opcode = 0x00
	}

	Cycles = 0
	switch opcode {
	case 0x02: //HLT
		CPU_Halted = true
	case 0x12: //HLT
		CPU_Halted = true
	case 0x22: //HLT
		CPU_Halted = true
	case 0x32: //HLT
		CPU_Halted = true
	case 0x42: //HLT
		CPU_Halted = true
	case 0x52: //HLT
		CPU_Halted = true
	case 0x62: //HLT
		CPU_Halted = true
	case 0x72: //HLT
		CPU_Halted = true
	case 0x92: //HLT
		CPU_Halted = true
	case 0xB2: //HLT
		CPU_Halted = true
	case 0xD2: //HLT
		CPU_Halted = true
	case 0xF2: //HLT
		CPU_Halted = true
	case 0xA0: //LDY Immediate
		Y = ReadFromPC()
		SetZNFlags(Y)
		Cycles = 2
	case 0xA4: //LDY Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Y = Read(Address)
		SetZNFlags(Y)
		Cycles = 3
	case 0xAC: //LDY Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Y = Read(Address)
		SetZNFlags(Y)
		Cycles = 4
	case 0xB4: //LDY Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Y = Read(Address)
		SetZNFlags(Y)
		Cycles = 4
	case 0xBC: //LDY Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed()
		Y = Read(Address)
		SetZNFlags(Y)
		Cycles = 4

	case 0xA2: //LDX Immediate
		X = ReadFromPC()
		SetZNFlags(X)
		Cycles = 2
	case 0xA6: //LDX Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		X = Read(Address)
		SetZNFlags(X)
		Cycles = 3
	case 0xAE: //LDX Absolute
		Address := ReadOperands_AbsoluteAddressed()
		X = Read(Address)
		SetZNFlags(X)
		Cycles = 4
	case 0xB6: //LDX Zero Page, Y
		Address := ReadOperands_ZeroPageAddressed_YIndexed()
		X = Read(Address)
		SetZNFlags(X)
		Cycles = 4
	case 0xBE: //LDX Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed()
		X = Read(Address)
		SetZNFlags(X)
		Cycles = 4

	case 0x85: //STA Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Write(Address, A)
		Cycles = 3
	case 0x95: //STA Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Write(Address, A)
		Cycles = 4
	case 0x8D: //STA Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Write(Address, A)
		Cycles = 4
	case 0x9D: //STA Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed()
		Write(Address, A)
		Cycles = 5
	case 0x99: //STA Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed()
		Write(Address, A)
		Cycles = 5
	case 0x81: //STA Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Write(Address, A)
		Cycles = 6
	case 0x91: //STA Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed()
		Write(Address, A)
		Cycles = 6

	case 0xA9: //LDA Immediate
		A = ReadFromPC()
		SetZNFlags(A)
		Cycles = 2
	case 0xA5: //LDA Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		A = Read(Address)
		SetZNFlags(A)
		Cycles = 3
	case 0xB5: //LDA Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		A = Read(Address)
		SetZNFlags(A)
		Cycles = 4
	case 0xAD: //LDA Absolute
		Address := ReadOperands_AbsoluteAddressed()
		A = Read(Address)
		SetZNFlags(A)
		Cycles = 4
	case 0xBD: //LDA Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed()
		A = Read(Address)
		SetZNFlags(A)
		Cycles = 4
	case 0xB9: //LDA Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed()
		A = Read(Address)
		SetZNFlags(A)
		Cycles = 4
	case 0xA1: //LDA Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		A = Read(Address)
		SetZNFlags(A)
		Cycles = 6
	case 0xB1: //LDA Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed()
		A = Read(Address)
		SetZNFlags(A)
		Cycles = 5

	case 0x86: //STX Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Write(Address, X)
		Cycles = 3
	case 0x96: //STX Zero Page, Y
		Address := ReadOperands_ZeroPageAddressed_YIndexed()
		Write(Address, X)
		Cycles = 4
	case 0x8E: //STX Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Write(Address, X)
		Cycles = 4

	case 0x84: //STY Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Write(Address, Y)
		Cycles = 3
	case 0x94: //STY Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Write(Address, Y)
		Cycles = 4
	case 0x8C: //STY Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Write(Address, Y)
		Cycles = 4

	case 0x10: //BPL (Branch on Plus)
		temp := ReadFromPC()
		Branch(!flag_Negative, temp)
	case 0x30: //BMI (Branch on Minus)
		temp := ReadFromPC()
		Branch(flag_Negative, temp)
	case 0x50: //BVC (Branch on Overflow Clear)
		temp := ReadFromPC()
		Branch(!flag_Overflow, temp)
	case 0x70: //BVS (Branch on Overflow Set)
		temp := ReadFromPC()
		Branch(flag_Overflow, temp)
	case 0x90: //BCC (Branch on Carry Clear)
		temp := ReadFromPC()
		Branch(!flag_Carry, temp)
	case 0xB0: //BCS (Branch on Carry Set)
		temp := ReadFromPC()
		Branch(flag_Carry, temp)
	case 0xD0: //BNE (Branch on Not Equal)
		temp := ReadFromPC()
		Branch(!flag_Zero, temp)
	case 0xF0: //BEQ (Branch on Equal)
		temp := ReadFromPC()
		Branch(flag_Zero, temp)
	case 0x48: //PHA
		Push(A)
		Cycles = 3
	case 0x68: //PLA
		A = Pull()
		SetZNFlags(A)
		Cycles = 4
	case 0x20: //JSR
		temp_low := ReadFromPC()
		Push(byte(ProgramCounter / 0x100))
		Push(byte(ProgramCounter))
		temp_high := ReadFromPC()
		ProgramCounter = BuildAddress(temp_low, temp_high)
		Cycles = 6
	case 0x60: //RTS
		temp_low := Pull()
		temp_high := Pull()
		ProgramCounter = BuildAddress(temp_low, temp_high)
		ProgramCounter++
		Cycles = 6
	case 0x4C: //JMP
		Address := ReadOperands_AbsoluteAddressed()
		ProgramCounter = Address
		Cycles = 3
	case 0xE8: //INX
		X++
		SetZNFlags(X)
	case 0xCA: //DEX
		X--
		SetZNFlags(X)
	case 0xC8: //INY
		Y++
		SetZNFlags(Y)
	case 0x88: //DEY
		Y--
		SetZNFlags(Y)
	case 0xAA: //TAX
		X = A
		SetZNFlags(X)
	case 0x8A: //TXA
		A = X
		SetZNFlags(A)
	case 0xA8: //TAY
		Y = A
		SetZNFlags(Y)
	case 0x98: //TYA
		A = Y
		SetZNFlags(A)
	case 0x9A: //TXS
		StackPointer = X
		Cycles = 2
	case 0xBA: //TSX
		X = StackPointer
		SetZNFlags(X)
		Cycles = 2
	case 0x38: //SEC
		flag_Carry = true
	case 0x18: //CLC
		flag_Carry = false
	case 0xB8: //CLV
		flag_Overflow = false
	case 0x78: //SEI
		flag_InterruptDisable = true
	case 0x58: //CLI
		flag_InterruptDisable = false
	case 0xF8: //SED
		flag_Decimal = true
	case 0xD8: //CLD
		flag_Decimal = false
	case 0xEA: //NOP
		Cycles = 2
	//Non-official NOP codes
	case 0x1A: //NOP Implied
		Cycles = 2
	case 0x3A: //NOP Implied
		Cycles = 2
	case 0x5A: //NOP Implied
		Cycles = 2
	case 0x7A: //NOP Implied
		Cycles = 2
	case 0xDA: //NOP Implied
		Cycles = 2
	case 0xFA: //NOP Implied
		Cycles = 2
	case 0x80: //NOP Immediate
		ReadFromPC()
		Cycles = 2
	case 0x82: //NOP Immediate
		ReadFromPC()
		Cycles = 2
	case 0x89: //NOP Immediate
		ReadFromPC()
		Cycles = 2
	case 0xC2: //NOP Immediate
		ReadFromPC()
		Cycles = 2
	case 0xE2: //NOP Immediate
		ReadFromPC()
		Cycles = 2
	case 0x04: //NOP Zero Page
		Cycles = 3
	case 0x44: //NOP Zero Page
		Cycles = 3
	case 0x64: //NOP Zero Page
		Cycles = 3
	case 0x14: //NOP Zero Page, X
		Cycles = 4
	case 0x34: //NOP Zero Page, X
		Cycles = 4
	case 0x54: //NOP Zero Page, X
		Cycles = 4
	case 0x74: //NOP Zero Page, X
		Cycles = 4
	case 0xD4: //NOP Zero Page, X
		Cycles = 4
	case 0xF4: //NOP Zero Page, X
		Cycles = 4
	case 0x0C: //NOP Absolute
		Cycles = 4
	case 0x1C: //NOP Absolute, X
		Cycles = 4
	case 0x3C: //NOP Absolute, X
		Cycles = 4
	case 0x5C: //NOP Absolute, X
		Cycles = 4
	case 0x7C: //NOP Absolute, X
		Cycles = 4
	case 0xDC: //NOP Absolute, X
		Cycles = 4
	case 0xFC: //NOP Absolute, X
		Cycles = 4

	case 0x08: //PHP
		temp := byte(0)
		temp += byte(ternary(flag_Carry, 0x01, 0x00))
		temp += byte(ternary(flag_Zero, 0x02, 0x00))
		temp += byte(ternary(flag_InterruptDisable, 0x04, 0x00))
		temp += byte(ternary(flag_Decimal, 0x08, 0x00))
		temp += 0x10
		temp += 0x20
		temp += byte(ternary(flag_Overflow, 0x40, 0x00))
		temp += byte(ternary(flag_Negative, 0x80, 0x00))
		Push(temp)
		Cycles = 3
	case 0x28: //PLP
		temp := Pull()
		flag_Carry = (temp & 0x01) != 0
		flag_Zero = (temp & 0x02) != 0
		flag_InterruptDisable = (temp & 0x04) != 0
		flag_Decimal = (temp & 0x08) != 0
		flag_Overflow = (temp & 0x40) != 0
		flag_Negative = (temp & 0x80) != 0
		Cycles = 4

	case 0x0A: //ASL A
		flag_Carry = A > 127
		A <<= 1
		SetZNFlags(A)
		Cycles = 2
	case 0x06: //ASL Zero Page
		Op_ASL(ReadOperands_ZeroPageAddressed())
		Cycles = 5
	case 0x0E: //ASL Absolute
		Op_ASL(ReadOperands_AbsoluteAddressed())
		Cycles = 6
	case 0x16: //ASL Zero Page, X
		Op_ASL(ReadOperands_ZeroPageAddressed_XIndexed())
		Cycles = 6
	case 0x1E: //ASL Absolute, X
		Op_ASL(ReadOperands_AbsoluteAddressed_XIndexed())
		Cycles = 7

	case 0x2A: //ROL A
		futureCarry := (A >= 0x80)
		A <<= 1
		if flag_Carry {
			A |= 1
		}
		flag_Carry = futureCarry
		SetZNFlags(A)
		Cycles = 2
	case 0x26: //ROL Zero Page
		Op_ROL(ReadOperands_ZeroPageAddressed())
		Cycles = 5
	case 0x2E: //ROL Absolute
		Op_ROL(ReadOperands_AbsoluteAddressed())
		Cycles = 6
	case 0x36: //ROL Zero Page, X
		Op_ROL(ReadOperands_ZeroPageAddressed_XIndexed())
		Cycles = 6
	case 0x3E: //ROL Absolute, X
		Op_ROL(ReadOperands_AbsoluteAddressed_XIndexed())
		Cycles = 7

	case 0x4A: //LSR A
		flag_Carry = (A & 1) != 0
		A >>= 1
		SetZNFlags(A)
		Cycles = 2
	case 0x46: //LSR Zero Page
		Op_LSR(ReadOperands_ZeroPageAddressed())
		Cycles = 5
	case 0x4E: //LSR Absolute
		Op_LSR(ReadOperands_AbsoluteAddressed())
		Cycles = 6
	case 0x56: //LSR Zero Page, X
		Op_LSR(ReadOperands_ZeroPageAddressed_XIndexed())
		Cycles = 6
	case 0x5E: //LSR Absolute, X
		Op_LSR(ReadOperands_AbsoluteAddressed_XIndexed())
		Cycles = 7

	case 0x6A: //ROR A
		futureCarry := (A & 1) != 0
		A >>= 1
		if flag_Carry {
			A |= 0x80
		}
		flag_Carry = futureCarry
		SetZNFlags(A)
		Cycles = 2
	case 0x66: //ROR Zero Page
		Op_ROR(ReadOperands_ZeroPageAddressed())
		Cycles = 5
	case 0x6E: //ROR Absolute
		Op_ROR(ReadOperands_AbsoluteAddressed())
		Cycles = 6
	case 0x76: //ROR Zero Page, X
		Op_ROR(ReadOperands_ZeroPageAddressed_XIndexed())
		Cycles = 6
	case 0x7E: //ROR Absolute, X
		Op_ROR(ReadOperands_AbsoluteAddressed_XIndexed())
		Cycles = 7

	case 0xE6: //INC Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_INC(Address, Read(Address))
		Cycles = 5
	case 0xEE: //INC Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Op_INC(Address, Read(Address))
		Cycles = 6
	case 0xF6: //INC Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_INC(Address, Read(Address))
		Cycles = 6
	case 0xFE: //INC Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed()
		Op_INC(Address, Read(Address))
		Cycles = 7

	case 0xC6: //DEC Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_DEC(Address, Read(Address))
		Cycles = 5
	case 0xCE: //DEC Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Op_DEC(Address, Read(Address))
		Cycles = 6
	case 0xD6: //DEC Zeropage, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_DEC(Address, Read(Address))
		Cycles = 6
	case 0xDE: //DEC Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed()
		Op_DEC(Address, Read(Address))
		Cycles = 7

	case 0x09: //ORA Immediate
		Op_ORA(ReadFromPC())
		Cycles = 2
	case 0x05: //ORA Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_ORA(Read(Address))
		Cycles = 3
	case 0x0D: //ORA Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Op_ORA(Read(Address))
		Cycles = 4
	case 0x15: //ORA Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ORA(Read(Address))
		Cycles = 4
	case 0x1D: //ORA Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed()
		Op_ORA(Read(Address))
		Cycles = 4
	case 0x19: //ORA Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed()
		Op_ORA(Read(Address))
		Cycles = 4
	case 0x01: //ORA Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_ORA(Read(Address))
		Cycles = 6
	case 0x11: //ORA Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed()
		Op_ORA(Read(Address))
		Cycles = 5

	case 0x29: //AND Immediate
		Op_AND(ReadFromPC())
		Cycles = 2
	case 0x25: //AND Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_AND(Read(Address))
		Cycles = 3
	case 0x2D: //AND Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Op_AND(Read(Address))
		Cycles = 4
	case 0x35: //AND Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_AND(Read(Address))
		Cycles = 4
	case 0x3D: //AND Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed()
		Op_AND(Read(Address))
		Cycles = 4
	case 0x39: //AND Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed()
		Op_AND(Read(Address))
		Cycles = 4
	case 0x21: //AND Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_AND(Read(Address))
		Cycles = 6
	case 0x31: //AND Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed()
		Op_AND(Read(Address))
		Cycles = 5

	case 0x49: //EOR Immediate
		Op_EOR(ReadFromPC())
		Cycles = 2
	case 0x45: //EOR Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_EOR(Read(Address))
		Cycles = 3
	case 0x4D: //EOR Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Op_EOR(Read(Address))
		Cycles = 4
	case 0x55: //EOR Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_EOR(Read(Address))
		Cycles = 4
	case 0x5D: //EOR Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed()
		Op_EOR(Read(Address))
		Cycles = 4
	case 0x59: //EOR Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed()
		Op_EOR(Read(Address))
		Cycles = 4
	case 0x41: //EOR Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_EOR(Read(Address))
		Cycles = 6
	case 0x51: //EOR Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed()
		Op_EOR(Read(Address))
		Cycles = 5

	case 0x69: //ADC Immediate
		Op_ADC(ReadFromPC())
		Cycles = 2
	case 0x65: //ADC Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_ADC(Read(Address))
		Cycles = 3
	case 0x75: //ADC Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_ADC(Read(Address))
		Cycles = 4
	case 0x6D: //ADC Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Op_ADC(Read(Address))
		Cycles = 4
	case 0x7D: //ADC Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed()
		Op_ADC(Read(Address))
		Cycles = 4
	case 0x79: //ADC Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed()
		Op_ADC(Read(Address))
		Cycles = 4
	case 0x61: //ADC Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_ADC(Read(Address))
		Cycles = 6
	case 0x71: //ADC Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed()
		Op_ADC(Read(Address))
		Cycles = 5

	case 0xE9: //SBC Immediate
		Op_SBC(ReadFromPC())
		Cycles = 2
	case 0xEB: //SBC Immediate (Unofficial)
		Op_SBC(ReadFromPC())
		Cycles = 2
	case 0xE5: //SBC Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_SBC(Read(Address))
		Cycles = 3
	case 0xED: //SBC Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Op_SBC(Read(Address))
		Cycles = 4
	case 0xF5: //SBC Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_SBC(Read(Address))
		Cycles = 4
	case 0xFD: //SBC Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed()
		Op_SBC(Read(Address))
		Cycles = 4
	case 0xF9: //SBC Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed()
		Op_SBC(Read(Address))
		Cycles = 4
	case 0xE1: //SBC Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_SBC(Read(Address))
		Cycles = 6
	case 0xF1: //SBC Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed()
		Op_SBC(Read(Address))
		Cycles = 5

	case 0xC9: //CMP Immediate
		Op_CMP(ReadFromPC())
		Cycles = 2
	case 0xC5: //CMP Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_CMP(Read(Address))
		Cycles = 3
	case 0xCD: //CMP Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Op_CMP(Read(Address))
		Cycles = 4
	case 0xD5: //CMP Zero Page, X
		Address := ReadOperands_ZeroPageAddressed_XIndexed()
		Op_CMP(Read(Address))
		Cycles = 4
	case 0xDD: //CMP Absolute, X
		Address := ReadOperands_AbsoluteAddressed_XIndexed()
		Op_CMP(Read(Address))
		Cycles = 4
	case 0xD9: //CMP Absolute, Y
		Address := ReadOperands_AbsoluteAddressed_YIndexed()
		Op_CMP(Read(Address))
		Cycles = 4
	case 0xC1: //CMP Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		Op_CMP(Read(Address))
		Cycles = 6
	case 0xD1: //CMP Indirect, Y
		Address := ReadOperands_IndirectAddressed_YIndexed()
		Op_CMP(Read(Address))
		Cycles = 5

	case 0xE0: //CPX Immediate
		Op_CPX(ReadFromPC())
		Cycles = 2
	case 0xE4: //CPX Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_CPX(Read(Address))
		Cycles = 3
	case 0xEC: //CPX Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Op_CPX(Read(Address))
		Cycles = 4

	case 0xC0: //CPY Immediate
		Op_CPY(ReadFromPC())
		Cycles = 2
	case 0xC4: //CPY Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_CPY(Read(Address))
		Cycles = 3
	case 0xCC: //CPY Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Op_CPY(Read(Address))
		Cycles = 4

	case 0x24: //BIT Zero Page
		Address := ReadOperands_ZeroPageAddressed()
		Op_BIT(Read(Address))
		Cycles = 3
	case 0x2C: //BIT Absolute
		Address := ReadOperands_AbsoluteAddressed()
		Op_BIT(Read(Address))
		Cycles = 4

	case 0x00: //BRK
		if !DoNMI {
			ProgramCounter++
		}
		Push(byte(ProgramCounter >> 8))
		Push(byte(ProgramCounter))

		temp := byte(0)
		temp += byte(ternary(flag_Carry, 0x01, 0x00))
		temp += byte(ternary(flag_Zero, 0x02, 0x00))
		temp += byte(ternary(flag_InterruptDisable, 0x04, 0x00))
		temp += byte(ternary(flag_Decimal, 0x08, 0x00))
		temp += byte(ternary(DoNMI, 0x00, 0x10))
		temp += 0x20
		temp += byte(ternary(flag_Overflow, 0x40, 0x00))
		temp += byte(ternary(flag_Negative, 0x80, 0x00))
		Push(temp)
		//flag_InterruptDisable = true

		PCL := Read(ternary(DoNMI, 0xFFFA, 0xFFFE))
		PCH := Read(ternary(DoNMI, 0xFFFB, 0xFFFF))
		ProgramCounter = uint16((uint16(PCH) * 0x100) + uint16(PCL)) //BuildAddress(PCL, PCH)
		DoNMI = false
		Cycles = 7

	case 0x40: //RTI
		temp := Pull()
		flag_Carry = (temp & 1) != 0
		flag_Zero = (temp & 2) != 0
		flag_InterruptDisable = (temp & 4) != 0
		flag_Decimal = (temp & 8) != 0
		flag_Overflow = (temp & 64) != 0
		flag_Negative = (temp & 128) != 0
		temp_low := Pull()
		temp_high := Pull()
		ProgramCounter = BuildAddress(temp_low, temp_high)
		Cycles = 6

	case 0x6C: //JMP Indirect
		Address := ReadOperands_IndirectAddressed()
		ProgramCounter = Address
		Cycles = 5 //TODO: What the fuck
	case 0xA3: //LAX Indirect, X
		Address := ReadOperands_IndirectAddressed_XIndexed()
		A = Read(Address)
		X = A
		SetZNFlags(X)
		Cycles = 6
	case 0x07: // SLO Zero Page
		//ASL + ORA
		Addr := ReadOperands_ZeroPageAddressed()
		Op_ASL(Addr)
		Op_ORA(Read(Addr))
		Cycles = 5

	case 0x7F: //RRA Absolute, X
		//ROR + ADC
		Address := ReadOperands_AbsoluteAddressed_XIndexed()
		Op_ROR(Address)
		Op_ADC(Read(Address))
		Cycles = 7
	case 0xC7: //DCP Zero Page
		//DEC + CMP
		Address := ReadOperands_ZeroPageAddressed()
		Op_INC(Address, Read(Address))
		Op_CMP(Read(Address))
		Cycles = 5
	case 0x9E: //SHX Absolute, Y
		Addr := ReadOperands_AbsoluteAddressed_YIndexed()
		val := (X & byte(((Addr&0xFF00)>>8)+1))
		Write(Addr, val)
		Cycles = 5
	case 0x3B: //RLA Absolute, Y
		//ROL + AND
		Addr := ReadOperands_AbsoluteAddressed_YIndexed()
		Op_ROL(Addr)
		Op_AND(Read(Addr))
		Cycles = 7
	case 0x0B: //ANC Immediate
		//AND + Set Carry as ASL
		Op_AND(ReadFromPC())
		flag_Carry = flag_Negative
		Cycles = 2
	case 0x2B: //ANC2 Immediate
		//AND + Set Carry as ROL
		//Same as $0B
		Op_AND(ReadFromPC())
		flag_Carry = flag_Negative
		Cycles = 2
	case 0x4B: //ALR
		//AND + LSR
		Op_AND(Read(ProgramCounter))
		Op_LSR(ProgramCounter)
		ProgramCounter++
		Cycles = 2
	case 0x7B: //RRA Absolute, Y
		//ROR + ADC
		Address := ReadOperands_AbsoluteAddressed_YIndexed()
		Op_ROR(Address)
		Op_ADC(Read(Address))
		Cycles = 7
	case 0x6B: //ARR
		//AND + ROR
		Op_AND(Read(ProgramCounter))
		Op_ROR(ProgramCounter)
		ProgramCounter++
		Cycles = 2
	case 0x8B: //ANE
		//Highly unstable, Do Not Use
		A = (((A | 0xEE) & X) & MagicConstant)
		Cycles = 2
	case 0xAB: //LXA
		val := MagicConstant
		A = val
		X = A
		SetZNFlags(A)
		Cycles = 2
	case 0xFF: //ISC Absolute, X
		//INC + SBC
		Address := ReadOperands_AbsoluteAddressed_XIndexed()
		Op_INC(Address, Read(Address))
		Op_SBC(Read(Address))
		Cycles = 7

		/*
			case 0x00:
			case 0x00:
			case 0x00:
			case 0x00:
			case 0x00:
			case 0x00:
			case 0x00:
			case 0x00:
			case 0x00:
			case 0x00:
			case 0x00:
			case 0x00:
			case 0x00:
			case 0x00:
			case 0x00:
		*/

	/*
		case 0x07: //SLO Zero Page (Illegal?)
			//Do ASL and ORA???
			Address := ReadOperands_ZeroPageAddressed()
			Op_ASL(Address, Read(Address))
			Op_ORA(Read(Address))
			Cycles = 5
		case 0x1A: //NOP
			Cycles = 2
		case 0x1C: //NOP Absolute, X (?????????)
			ReadOperands_AbsoluteAddressed_XIndexed()
			Cycles = 4
		case 0x82: //NOP Immediate
			ReadFromPC()
			Cycles = 2
	*/
	default:
		fmt.Println("Unknown Opcode: " + fmt.Sprintf("%02X", opcode))
		LogCount--
		if LogCount <= 0 {
			CPU_Halted = true
		}
	}
	if LoggingCPU {
		TraceLogger()
	}
	//CartRAMLogger()
	operands = nil

	TotalCycles += Cycles
	for Cycles > 0 {
		Cycles--
		Emulate_PPU(g)
		Emulate_PPU(g)
		Emulate_PPU(g)
	}

	//Force Stop, just in case
	//if TotalCycles > 900000 {
	/*
		if InstructionCount > 500000 {
			fmt.Println("Too many cycles, end")
			CPU_Halted = true
		}
	*/
	InstructionCount++
}

func Emulate_PPU(g *Game) {

	if ppuDot == 1 && ppuScanline == 241 {
		ppuVBlank = true
		DrawNewFrame = true
	} else if ppuDot == 1 && ppuScanline == 261 {
		ppuVBlank = false
		ppuStatusOverflow = false
		ppuStatusSprZeroHit = false
	}

	SpriteEvaluation()
	if ppuScanline < 240 || ppuScanline == 261 {
		if (ppuDot > 0 && ppuDot <= 256) || (ppuDot > 320 && ppuDot <= 336) {
			//If this is a visible pixel, or preparing the start of the next scanline
			if ppuMask_RenderBG || ppuMask_RenderSprites {
				//If rendering is enabled
				if ppuMask_RenderBG { //If rendering the background, update the shift registers for the background
					ppuShiftRegister_patternL = ppuShiftRegister_patternL << 1     //Shift 1 bit to the left
					ppuShiftRegister_patternH = ppuShiftRegister_patternH << 1     //Shift 1 bit to the left
					ppuShiftRegister_attributeL = ppuShiftRegister_attributeL << 1 //Shift 1 bit to the left
					ppuShiftRegister_attributeH = ppuShiftRegister_attributeH << 1 //Shift 1 bit to the left
				}
				if ppuMask_RenderBG || ppuMask_RenderSprites { //If rendering at all, let's decrement the X position of the objects
					if ppuDot > 1 && ppuDot <= 256 { //Don't decrement until dot 1
						for i := 0; i < 8; i++ {
							if ppu_SpriteXposition[i] > 0 {
								ppu_SpriteXposition[i]-- //Decrement the position of all objects in secondary OAM. When this is zero, the PPU can draw it
							} else {
								ppu_SpriteShiftRegisterL[i] = byte(ppu_SpriteShiftRegisterL[i] << 1) //Shift 1 bit to the left
								ppu_SpriteShiftRegisterH[i] = byte(ppu_SpriteShiftRegisterH[i] << 1) //Shift 1 bit to the left
							}

						}
					}
				}
				PPU8Steps()
			}
		}
	}

	if (ppuScanline < 240 || ppuScanline == 261) && (ppuMask_RenderBG || ppuMask_RenderSprites) {
		//If this is a visible scanline and rendering sprites / background is enabled
		if ppuDot == 256 { //The Y Scroll is incremented on dot 256
			PPU_IncrementScrollY()
		} else if ppuDot == 257 { //The X scroll is reset on dot 257
			PPU_ResetXScroll()
		}
		if ppuDot >= 280 && ppuDot <= 304 && ppuScanline == 261 { //numbers from the nesdev wiki
			PPU_ResetYScroll() //The Y scroll is reset on every dot from 280 through 304 on the pre-render scanline
		}
	}

	//Drawing
	DrawScreen(g)

	ppuDot++
	if ppuDot > 341 {
		ppuDot = 0
		ppuScanline++
		if ppuScanline > 261 {
			ppuScanline = 0

		}
	}
}

func DrawScreen(g *Game) {
	Palette := [64]color.RGBA{
		{R: 0x65, G: 0x65, B: 0x65, A: 0xFF},
		{R: 0x00, G: 0x2A, B: 0x84, A: 0xFF},
		{R: 0x15, G: 0x13, B: 0xA2, A: 0xFF},
		{R: 0x3A, G: 0x01, B: 0x9E, A: 0xFF},
		{R: 0x59, G: 0x00, B: 0x7A, A: 0xFF},
		{R: 0x6A, G: 0x00, B: 0x3E, A: 0xFF},
		{R: 0x68, G: 0x08, B: 0x00, A: 0xFF},
		{R: 0x53, G: 0x1D, B: 0x00, A: 0xFF},
		{R: 0x32, G: 0x34, B: 0x00, A: 0xFF},
		{R: 0x0D, G: 0x46, B: 0x00, A: 0xFF},
		{R: 0x00, G: 0x4F, B: 0x00, A: 0xFF},
		{R: 0x00, G: 0x4C, B: 0x09, A: 0xFF},
		{R: 0x00, G: 0x3F, B: 0x4B, A: 0xFF},
		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},

		{R: 0xAE, G: 0xAE, B: 0xAE, A: 0xFF},
		{R: 0x17, G: 0x5F, B: 0xD6, A: 0xFF},
		{R: 0x43, G: 0x41, B: 0xFF, A: 0xFF},
		{R: 0x75, G: 0x29, B: 0xFA, A: 0xFF},
		{R: 0x9E, G: 0x1D, B: 0xCA, A: 0xFF},
		{R: 0xB4, G: 0x20, B: 0x7B, A: 0xFF},
		{R: 0xB1, G: 0x33, B: 0x22, A: 0xFF},
		{R: 0x96, G: 0x4E, B: 0x00, A: 0xFF},
		{R: 0x6A, G: 0x6C, B: 0x00, A: 0xFF},
		{R: 0x39, G: 0x84, B: 0x00, A: 0xFF},
		{R: 0x0F, G: 0x90, B: 0x00, A: 0xFF},
		{R: 0x00, G: 0x8D, B: 0x33, A: 0xFF},
		{R: 0x00, G: 0x7B, B: 0x8C, A: 0xFF},
		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},

		{R: 0xFE, G: 0xFE, B: 0xFE, A: 0xFF},
		{R: 0x66, G: 0xAF, B: 0xFF, A: 0xFF},
		{R: 0x93, G: 0x90, B: 0xFF, A: 0xFF},
		{R: 0xC5, G: 0x78, B: 0xFF, A: 0xFF},
		{R: 0xEE, G: 0x6C, B: 0xFF, A: 0xFF},
		{R: 0xFF, G: 0x6F, B: 0xCA, A: 0xFF},
		{R: 0xFF, G: 0x82, B: 0x71, A: 0xFF},
		{R: 0xE6, G: 0x9E, B: 0x25, A: 0xFF},
		{R: 0xBA, G: 0xBC, B: 0x00, A: 0xFF},
		{R: 0x88, G: 0xD5, B: 0x01, A: 0xFF},
		{R: 0x5E, G: 0xE1, B: 0x32, A: 0xFF},
		{R: 0x47, G: 0xDD, B: 0x82, A: 0xFF},
		{R: 0x4A, G: 0xCB, B: 0xDC, A: 0xFF},
		{R: 0x4E, G: 0x4E, B: 0x4E, A: 0xFF},
		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},

		{R: 0xFE, G: 0xFE, B: 0xFE, A: 0xFF},
		{R: 0xC0, G: 0xDE, B: 0xFF, A: 0xFF},
		{R: 0xD2, G: 0xD1, B: 0xFF, A: 0xFF},
		{R: 0xE7, G: 0xC7, B: 0xFF, A: 0xFF},
		{R: 0xF8, G: 0xC2, B: 0xFF, A: 0xFF},
		{R: 0xFF, G: 0xC3, B: 0xE9, A: 0xFF},
		{R: 0xFF, G: 0xCB, B: 0xC4, A: 0xFF},
		{R: 0xF5, G: 0xD7, B: 0xA5, A: 0xFF},
		{R: 0xE2, G: 0xE3, B: 0x94, A: 0xFF},
		{R: 0xCE, G: 0xED, B: 0x96, A: 0xFF},
		{R: 0xBC, G: 0xF2, B: 0xAA, A: 0xFF},
		{R: 0xB3, G: 0xF1, B: 0xCB, A: 0xFF},
		{R: 0xB4, G: 0xE9, B: 0xF0, A: 0xFF},
		{R: 0xB6, G: 0xB6, B: 0xB6, A: 0xFF},
		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
	}
	if ppuScanline < 240 && ppuDot > 0 && ppuDot <= 256 {
		var PalHi byte = 0  //Which color palette to use?
		var PalLow byte = 0 //Index into a color palette
		if ppuMask_RenderBG && (ppuDot > 8 || ppuMask_8pxMaskBG) {
			col0 := byte((ppuShiftRegister_patternL >> (15 - ppuScrollFineX)) & 1)
			col1 := byte((ppuShiftRegister_patternH >> (15 - ppuScrollFineX)) & 1)
			PalLow = byte(uint16(col1)<<1 | uint16(col0))

			pal0 := byte((ppuShiftRegister_attributeL >> (15 - ppuScrollFineX)) & 1)
			pal1 := byte((ppuShiftRegister_attributeH >> (15 - ppuScrollFineX)) & 1)
			PalHi = byte(uint16(pal1)<<1 | uint16(pal0))

			if PalLow == 0 && PalHi != 0 { //Color 0 of all palettes are mirrors of color 0 of palette 0
				PalHi = 0
			}
		}
		var SpritePalHi byte = 0        //Which color palette to use
		var SpritePalLow byte = 0       //Index into a color palette
		var SpritePriority bool = false //Is the sprite in front or behind the BG?
		if ppuMask_RenderSprites && (ppuDot > 8 || ppuMask_8pxMaskSprites) {
			for i := 0; i < 8; i++ {
				if ppu_SpriteXposition[i] == 0 && i < int(ppuSecondaryOAMSize/4) { //If the sprite X position == 0 (The x position is decremented every ppu cycle)
					SpixelL := ((ppu_SpriteShiftRegisterL[i]) & 0x80) != 0
					SpixelH := ((ppu_SpriteShiftRegisterH[i]) & 0x80) != 0
					SpritePalLow = 0
					if SpixelL {
						SpritePalLow = 1
					}
					if SpixelH {
						SpritePalLow |= 2
					}

					SpritePalHi = byte((ppu_SpriteAttribute[i] & 0x03) | 0x04)
					SpritePriority = ((ppu_SpriteAttribute[i] >> 5) & 1) == 0
				} else {
					continue
				}
				if SpritePalLow != 0 {
					if i == 0 && ppuScanlineContainsSpriteZero && SpritePalLow != 0 && PalLow != 0 && ppuMask_RenderBG && ppuDot < 256 {
						ppuStatusSprZeroHit = true
					}
					break
				}
			}
		}

		if (SpritePriority && SpritePalLow != 0) || PalLow == 0 {
			PalLow = SpritePalLow
			PalHi = SpritePalHi
			if PalLow == 0 {
				PalHi = 0
			}
		}

		color := Palette[PaletteRAM[(PalHi*4)+PalLow]]
		//RenderColor(g, color)
		pixIndex := uint64((((ppuScanline) * screenWidth) + (ppuDot - 1)) * 4)
		//pixIndex &= 0x3BFFF
		g.gameScreen.Pix[pixIndex] = color.R
		g.gameScreen.Pix[pixIndex+1] = color.G
		g.gameScreen.Pix[pixIndex+2] = color.B
		g.gameScreen.Pix[pixIndex+3] = color.A
		//g.gameScreen.Set(ppuDot-1, ppuScanline, color)

	}
}

/*func RenderColor(g *Game, color color.RGBA) {
	pixIndex := int((((ppuScanline-1)*screenWidth)+ppuDot)*4) - 1
	//pixIndex &= 0xEFFF
	g.gameScreen.Pix[pixIndex] = color.R
	g.gameScreen.Pix[pixIndex+1] = color.G
	g.gameScreen.Pix[pixIndex+2] = color.B
	g.gameScreen.Pix[pixIndex+3] = 0xFF
}*/

func PPU8Steps() {
	//What part of the 8-step process to run this cycle
	cycleTick := byte((ppuDot - 1) & 7)
	switch cycleTick {
	case 0:
		ppuShiftRegister_patternL = ((ppuShiftRegister_patternL & 0xFF00) | uint16(ppu8Step_patternLowBitPlane))
		ppuShiftRegister_patternH = ((ppuShiftRegister_patternH & 0xFF00) | uint16(ppu8Step_patternHighBitPlane))
		ppuShiftRegister_attributeL = ((ppuShiftRegister_attributeL & 0xFF00) | ternary((ppu8Step_attribute&1) == 1, 0xFF, 0x00))
		ppuShiftRegister_attributeH = ((ppuShiftRegister_attributeH & 0xFF00) | ternary((ppu8Step_attribute&2) == 2, 0xFF, 0x00))
		ppuAddressBus = (0x2000 + (VRAMAddress & 0x0FFF))
		ppu8Step_temp = ReadPPU(ppuAddressBus)
	case 1:
		ppu8Step_NextCharacter = ppu8Step_temp
	case 2:
		ppuAddressBus = (0x23C0 | (VRAMAddress & 0x0C00) | ((VRAMAddress >> 4) & 0x38) | ((VRAMAddress >> 2) & 0x07))
		ppu8Step_temp = ReadPPU(ppuAddressBus)
	case 3:
		ppu8Step_attribute = ppu8Step_temp
		//1 byte of attribute data covers 4 tiles. determine which tile this is for
		if (VRAMAddress & 3) >= 2 { //If this is in the right tile
			ppu8Step_attribute = byte(ppu8Step_attribute >> 2)
		}
		if (((VRAMAddress & 0b0000001111100000) >> 5) & 3) >= 2 { //If this is in the bottom tile
			ppu8Step_attribute = byte(ppu8Step_attribute >> 4)
		}
		ppu8Step_attribute = byte(ppu8Step_attribute & 3)
	case 4:
		ppuAddressBus = (((VRAMAddress & 0b0111000000000000) >> 12) | (uint16(ppu8Step_NextCharacter) * 16) | ternary(ppuBGPatternTable, 0x1000, 0))
		ppu8Step_temp = ReadPPU(ppuAddressBus)
	case 5:
		ppu8Step_patternLowBitPlane = ppu8Step_temp
		ppuAddressBus += 8
	case 6:
		ppu8Step_temp = ReadPPU(ppuAddressBus)
	case 7:
		ppu8Step_patternHighBitPlane = ppu8Step_temp
		//Increment VRAM with scrolling
		if (VRAMAddress & 0x001F) == 31 {
			VRAMAddress &= 0xFFE0 //Reset the scroll
			VRAMAddress ^= 0x0400 //Crossing into next nametable
		} else {
			VRAMAddress++
		}
	}
}

func PPU_IncrementScrollY() {
	if (VRAMAddress & 0x7000) != 0x7000 {
		VRAMAddress += 0x1000
	} else {
		VRAMAddress &= 0x0FFF
		y := uint16((VRAMAddress & 0x03E0) >> 5)
		if y == 29 {
			y = 0 // Reset the Y value and also flip some other bit in the 'v' register
			VRAMAddress ^= 0x0800
		} else {
			y++ //Increment the Y value
			y &= 0x1F
		}
		VRAMAddress = (uint16(VRAMAddress&0xFC1F) | uint16(y)<<5)
	}
}

func PPU_ResetXScroll() {
	VRAMAddress = ((VRAMAddress & 0b0111101111100000) | (TransferAddress & 0b0000010000011111))
}

func PPU_ResetYScroll() {
	VRAMAddress = ((VRAMAddress & 0b0000010000011111) | (TransferAddress & 0b0111101111100000))
}

var OAM [0x100]byte
var SecondaryOAM [0x20]byte
var ppuSpriteEvalTemp byte
var ppuOAMAddress, ppuSecondaryOAMAddress, ppuSecondaryOAMSize uint16
var ppuSecondaryOAMFull, ppuStatusOverflow, ppuStatusSprZeroHit, ppuScanlineContainsSpriteZero, ppuSpriteEvaluationOAMOverflowed bool
var ppuSpriteEvalTick int

var ppu_SpriteShiftRegisterL [8]byte
var ppu_SpriteShiftRegisterH [8]byte

var ppu_SpriteAttribute [8]byte
var ppu_SpritePattern [8]byte
var ppu_SpriteXposition [8]byte
var ppu_SpriteYposition [8]byte

func SpriteEvaluation() {
	if ppuDot == 0 { //Step 0: Reset Secondary OAM count
		ppuSecondaryOAMAddress = 0
		ppuSecondaryOAMFull = false
	} else if ppuDot > 0 && ppuDot <= 64 { //Step 1: Clear Secondary OAM
		if (ppuDot & 1) == 1 {
			//Odd PPU cycles load the value $FF
			ppuSpriteEvalTemp = 0xFF
		} else {
			//Even PPU cycles store the value in secondaryOAM
			SecondaryOAM[ppuSecondaryOAMAddress] = ppuSpriteEvalTemp
			ppuSecondaryOAMAddress++
			ppuSecondaryOAMAddress &= 0x1F //Keep this limited from $00 to $1F
		}
	} else if ppuDot > 64 && ppuDot <= 256 { //Step 2: Load OAM into Secondary OAM (If not full)
		if (ppuDot & 1) == 1 {
			//Odd PPU cycles load the value from OAM
			ppuSpriteEvalTemp = OAM[ppuOAMAddress&0xFF] //CHECK IF THIS IS LEGIT
		} else {
			if !ppuSpriteEvaluationOAMOverflowed {
				//Even PPU cycles store the value in secondaryOAM
				if !ppuSecondaryOAMFull { //If SecondaryOAM is not full yet
					// As long as secondaryOAM isn't full, this write *always* occurs, regardless of evaluation
					SecondaryOAM[ppuSecondaryOAMAddress] = ppuSpriteEvalTemp
				}
				if ppuSpriteEvalTick == 0 {
					//Reading index 0 of an object's set of 4 bytes
					if (ppuScanline-int(ppuSpriteEvalTemp) >= 0) && (ppuScanline-int(ppuSpriteEvalTemp) < int(ternary(ppuUse8x16Sprites, 16, 8))) {
						//This object *is* on this scanline!
						if !ppuSecondaryOAMFull {
							ppuSecondaryOAMAddress++ //Increment this for the next write to Secondary OAM
							ppuOAMAddress++          //Increment this for the next ream of Object Attribute Memory
							if ppuDot == 66 {
								// Rather than verifying that this is OAM index 0,
								// the PPU sets this flag if we found an object on this scanline
								// during ppuDot 66, which would be the PPU cycle evaluating index 0
								ppuScanlineContainsSpriteZero = true
							}

						} else {
							ppuStatusOverflow = true
						}
						ppuSpriteEvalTick++

					} else {
						ppuOAMAddress += 4
					}
				} else { //If ppuSpriteEvalTick != 0
					// Reading index 1, 2, or 3 of an object's OAM data.
					// We're not going to be making any checks for if things are on this scanline,
					// so we can just simply increment the OAM address
					ppuSecondaryOAMAddress++ //Increment this for the next write to Secondary OAM
					ppuOAMAddress++          //Increment this for the next ream of Object Attribute Memory
					if ppuSecondaryOAMAddress == 0x20 {
						ppuSecondaryOAMFull = true
					}
					ppuSpriteEvalTick++
					ppuSpriteEvalTick &= 3 // Wrap around to tick 0 after tick 3
				}
				if ppuOAMAddress == 0 {
					// If we overflow the OAM address, we want to stop running the sprite evaluation checks until dot 257
					ppuSpriteEvaluationOAMOverflowed = true
				}
			}
		}
	} else if ppuDot > 256 && ppuDot <= 320 { //Step 3:
		ppuOAMAddress = 0 //This is set to $00 during every one of these cycles
		if ppuDot == 257 {
			ppuSecondaryOAMSize = ppuSecondaryOAMAddress
			ppuSecondaryOAMAddress = 0
			ppuSpriteEvalTick = 0
		}

		switch ppuSpriteEvalTick {
		case 0:
			//Set this object's Y position in the array
			ppu_SpriteYposition[ppuSecondaryOAMAddress/4] = SecondaryOAM[ppuSecondaryOAMAddress]
			ppuSecondaryOAMAddress++
		case 1:
			//Set this object's Y position in the array
			ppu_SpritePattern[ppuSecondaryOAMAddress/4] = SecondaryOAM[ppuSecondaryOAMAddress]
			ppuSecondaryOAMAddress++
		case 2:
			//Set this object's Y position in the array
			ppu_SpriteAttribute[ppuSecondaryOAMAddress/4] = SecondaryOAM[ppuSecondaryOAMAddress]
			ppuSecondaryOAMAddress++
		case 3:
			//Set this object's Y position in the array
			ppu_SpriteXposition[ppuSecondaryOAMAddress/4] = SecondaryOAM[ppuSecondaryOAMAddress]
		case 4:
			ppuAddressBus = ppuFindSpritePatternData(ppuSecondaryOAMAddress / 4)
		case 5:
			ppuSpriteEvalTemp = ReadPPU(ppuAddressBus)
			if ppuScanline == 261 {
				ppuSpriteEvalTemp = 0 //Clear this if this is the pre-render line
			}
			if ((ppu_SpriteAttribute[ppuSecondaryOAMAddress/4] >> 6) & 1) == 1 { //Attributes are set up to flip X
				// Real nice way to change the order of bits from 76543210 to 01234567
				ppuSpriteEvalTemp = byte(((ppuSpriteEvalTemp & 0xF0) >> 4) | ((ppuSpriteEvalTemp & 0xF) << 4))
				ppuSpriteEvalTemp = byte(((ppuSpriteEvalTemp & 0xCC) >> 2) | ((ppuSpriteEvalTemp & 0x33) << 2))
				ppuSpriteEvalTemp = byte(((ppuSpriteEvalTemp & 0xAA) >> 1) | ((ppuSpriteEvalTemp & 0x55) << 1))
			}
			ppu_SpriteShiftRegisterL[ppuSecondaryOAMAddress/4] = ppuSpriteEvalTemp
		case 6:
			ppuAddressBus += 8
		case 7:
			ppuSpriteEvalTemp = ReadPPU(ppuAddressBus)
			if ppuScanline == 261 {
				ppuSpriteEvalTemp = 0 //Clear this if this is the pre-render line
			}
			if ((ppu_SpriteAttribute[ppuSecondaryOAMAddress/4] >> 6) & 1) == 1 { //Attributes are set up to flip X
				// Real nice way to change the order of bits from 76543210 to 01234567
				ppuSpriteEvalTemp = byte(((ppuSpriteEvalTemp & 0xF0) >> 4) | ((ppuSpriteEvalTemp & 0xF) << 4))
				ppuSpriteEvalTemp = byte(((ppuSpriteEvalTemp & 0xCC) >> 2) | ((ppuSpriteEvalTemp & 0x33) << 2))
				ppuSpriteEvalTemp = byte(((ppuSpriteEvalTemp & 0xAA) >> 1) | ((ppuSpriteEvalTemp & 0x55) << 1))
			}
			ppu_SpriteShiftRegisterH[ppuSecondaryOAMAddress/4] = ppuSpriteEvalTemp
			ppuSecondaryOAMAddress++
		}
		ppuSpriteEvalTick++
		ppuSpriteEvalTick &= 7 // And reset at 8

	}
}

func ppuFindSpritePatternData(SecondaryOAMSlot uint16) uint16 {
	if !ppuUse8x16Sprites { //8x8 sprites
		// The address is $0000 or $1000 depending on the pattern table
		// plus the pattern value from OAM * 16
		// plus the number of scanlines from the top of the object
		// if the attributes are set to flip Y, it's 7 - the number of scanlines from the top of the object
		if ((ppu_SpriteAttribute[SecondaryOAMSlot] >> 7) & 1) == 0 { //Attributes are not set up to flip Y
			return uint16(ternary(ppuSpritePatternTable, 0x1000, 0) + (uint16(ppu_SpritePattern[SecondaryOAMSlot]) << 4) + uint16(ppuScanline-int(ppu_SpriteYposition[SecondaryOAMSlot])))
		} else { //Attributes are set up to flip Y
			return uint16(ternary(ppuSpritePatternTable, 0x1000, 0) + (uint16(ppu_SpritePattern[SecondaryOAMSlot]) << 4) + uint16((7-(ppuScanline-int(ppu_SpriteYposition[SecondaryOAMSlot])))&7))

		}

	} else { //8x16 sprites
		// in 8x16 mode, instead of using ppu_SpritePattern to deternime which pattern table to fetch data from...
		// these sprites instead use bit 0 of the object's pattern information from OAM

		// The address is $0000 or $1000 depending on the pattern table
		// plus (the pattern value from OAM, clearing bit 0) * 16
		// plus the number of scanlines from the top of the object
		// if the attributes are set to flip Y, it's 7 - the number of scanlines from the top of the object

		//If we're drawing the bottom half of the sprite, add 16
		if ((ppu_SpriteAttribute[SecondaryOAMSlot] >> 7) & 1) == 0 { //Attributes are not set up to flip Y
			if ppuScanline-int(ppu_SpriteYposition[SecondaryOAMSlot]) < 8 {
				return uint16(ternary((ppu_SpritePattern[SecondaryOAMSlot]&1) == 1, 0x1000, 0) | (uint16(ppu_SpritePattern[SecondaryOAMSlot]&0xFE) << 4) + uint16(ppuScanline-int(ppu_SpriteYposition[SecondaryOAMSlot])))
			} else {
				return uint16(ternary((ppu_SpritePattern[SecondaryOAMSlot]&1) == 1, 0x1000, 0) + ((uint16(ppu_SpritePattern[SecondaryOAMSlot]&0xFE) << 4) + 16) + uint16((ppuScanline-int(ppu_SpriteYposition[SecondaryOAMSlot]))&7))
			}
		} else { //Attributes are set up to flip Y
			if ppuScanline-int(ppu_SpriteYposition[SecondaryOAMSlot]) < 8 {
				return uint16(ternary((ppu_SpritePattern[SecondaryOAMSlot]&1) == 1, 0x1000, 0) + ((uint16(ppu_SpritePattern[SecondaryOAMSlot]&0xFE) << 4) + 16) + uint16(((ppuScanline-int(ppu_SpriteYposition[SecondaryOAMSlot]))&7)+7))
			} else {
				return uint16(ternary((ppu_SpritePattern[SecondaryOAMSlot]&1) == 1, 0x1000, 0) + ((uint16(ppu_SpritePattern[SecondaryOAMSlot]&0xFE) << 4) + 7) + uint16((ppuScanline-int(ppu_SpriteYposition[SecondaryOAMSlot]))&7))
			}
		}
	}
}

package debug

import (
	"fmt"
	"log"
	"mtt/timenes/common"
	"mtt/timenes/nes/bus"
	"mtt/timenes/nes/ppu"
	"os"
	"path/filepath"
	"time"
)

type opcodeDetail struct {
	name     string
	operands int         // No. of operands
	adMode   AddressMode // Addressing mode
}

type AddressMode int

const (
	Address_None AddressMode = iota
	Address_Immediate
	Address_Absolute
	Address_AbsoluteX
	Address_AbsoluteY
	Address_Indirect
	Address_IndirectX
	Address_IndirectY
	Address_ZeroPage
	Address_ZeroPageX
	Address_ZeroPageY
)

var opcodeInfo = map[byte]opcodeDetail{
	0x00: {"BRK", 0, Address_None},
	0x01: {"ORA", 0, Address_None},
	0x02: {"HLT", 0, Address_None},
	0x03: {"SLO", 0, Address_None},
	0x04: {"NOP", 0, Address_None},
	0x05: {"ORA", 0, Address_None},
	0x06: {"ASL", 0, Address_None},
	0x07: {"SLO", 0, Address_None},
	0x08: {"PHP", 0, Address_None},
	0x09: {"ORA", 0, Address_None},
	0x0A: {"ASL", 0, Address_None},
	0x0B: {"ANC", 0, Address_None},
	0x0C: {"NOP", 0, Address_None},
	0x0D: {"ORA", 0, Address_None},
	0x0E: {"ASL", 0, Address_None},
	0x0F: {"SLO", 0, Address_None},

	0x10: {"BPL", 0, Address_None},
	0x11: {"ORA", 0, Address_None},
	0x12: {"HLT", 0, Address_None},
	0x13: {"SLO", 0, Address_None},
	0x14: {"NOP", 0, Address_None},
	0x15: {"ORA", 0, Address_None},
	0x16: {"ASL", 0, Address_None},
	0x17: {"SLO", 0, Address_None},
	0x18: {"CLC", 0, Address_None},
	0x19: {"ORA", 0, Address_None},
	0x1A: {"NOP", 0, Address_None},
	0x1B: {"SLO", 0, Address_None},
	0x1C: {"NOP", 0, Address_None},
	0x1D: {"ORA", 0, Address_None},
	0x1E: {"ASL", 0, Address_None},
	0x1F: {"SLO", 0, Address_None},

	0x20: {"JSR", 0, Address_None},
	0x21: {"AND", 0, Address_None},
	0x22: {"HLT", 0, Address_None},
	0x23: {"RLA", 0, Address_None},
	0x24: {"BIT", 0, Address_None},
	0x25: {"AND", 0, Address_None},
	0x26: {"ROL", 0, Address_None},
	0x27: {"RLA", 0, Address_None},
	0x28: {"PLP", 0, Address_None},
	0x29: {"AND", 0, Address_None},
	0x2A: {"ROL", 0, Address_None},
	0x2B: {"ANC", 0, Address_None},
	0x2C: {"BIT", 0, Address_None},
	0x2D: {"AND", 0, Address_None},
	0x2E: {"ROL", 0, Address_None},
	0x2F: {"RLA", 0, Address_None},

	0x30: {"BMI", 0, Address_None},
	0x31: {"AND", 0, Address_None},
	0x32: {"HLT", 0, Address_None},
	0x33: {"RLA", 0, Address_None},
	0x34: {"NOP", 0, Address_None},
	0x35: {"AND", 0, Address_None},
	0x36: {"ROL", 0, Address_None},
	0x37: {"RLA", 0, Address_None},
	0x38: {"SEC", 0, Address_None},
	0x39: {"AND", 0, Address_None},
	0x3A: {"NOP", 0, Address_None},
	0x3B: {"RLA", 0, Address_None},
	0x3C: {"NOP", 0, Address_None},
	0x3D: {"AND", 0, Address_None},
	0x3E: {"ROL", 0, Address_None},
	0x3F: {"RLA", 0, Address_None},

	0x40: {"RTI", 0, Address_None},
	0x41: {"EOR", 0, Address_None},
	0x42: {"HLT", 0, Address_None},
	0x43: {"SRE", 0, Address_None},
	0x44: {"NOP", 0, Address_None},
	0x45: {"EOR", 0, Address_None},
	0x46: {"LSR", 0, Address_None},
	0x47: {"SRE", 0, Address_None},
	0x48: {"PHA", 0, Address_None},
	0x49: {"EOR", 0, Address_None},
	0x4A: {"LSR", 0, Address_None},
	0x4B: {"ALR", 0, Address_None},
	0x4C: {"JMP", 0, Address_None},
	0x4D: {"EOR", 0, Address_None},
	0x4E: {"LSR", 0, Address_None},
	0x4F: {"SRE", 0, Address_None},

	0x50: {"BVC", 0, Address_None},
	0x51: {"EOR", 0, Address_None},
	0x52: {"HLT", 0, Address_None},
	0x53: {"SRE", 0, Address_None},
	0x54: {"NOP", 0, Address_None},
	0x55: {"EOR", 0, Address_None},
	0x56: {"LSR", 0, Address_None},
	0x57: {"SRE", 0, Address_None},
	0x58: {"CLI", 0, Address_None},
	0x59: {"EOR", 0, Address_None},
	0x5A: {"NOP", 0, Address_None},
	0x5B: {"SRE", 0, Address_None},
	0x5C: {"NOP", 0, Address_None},
	0x5D: {"EOR", 0, Address_None},
	0x5E: {"LSR", 0, Address_None},
	0x5F: {"SRE", 0, Address_None},

	0x60: {"RTS", 0, Address_None},
	0x61: {"ADC", 0, Address_None},
	0x62: {"HLT", 0, Address_None},
	0x63: {"RRA", 0, Address_None},
	0x64: {"NOP", 0, Address_None},
	0x65: {"ADC", 0, Address_None},
	0x66: {"ROR", 0, Address_None},
	0x67: {"RRA", 0, Address_None},
	0x68: {"PLA", 0, Address_None},
	0x69: {"ADC", 0, Address_None},
	0x6A: {"ROR", 0, Address_None},
	0x6B: {"ARR", 0, Address_None},
	0x6C: {"JMP", 0, Address_None},
	0x6D: {"ADC", 0, Address_None},
	0x6E: {"ROR", 0, Address_None},
	0x6F: {"RRA", 0, Address_None},

	0x70: {"BVS", 0, Address_None},
	0x71: {"ADC", 0, Address_None},
	0x72: {"HLT", 0, Address_None},
	0x73: {"RRA", 0, Address_None},
	0x74: {"NOP", 0, Address_None},
	0x75: {"ADC", 0, Address_None},
	0x76: {"ROR", 0, Address_None},
	0x77: {"RRA", 0, Address_None},
	0x78: {"SEI", 0, Address_None},
	0x79: {"ADC", 0, Address_None},
	0x7A: {"NOP", 0, Address_None},
	0x7B: {"RRA", 0, Address_None},
	0x7C: {"NOP", 0, Address_None},
	0x7D: {"ADC", 0, Address_None},
	0x7E: {"ROR", 0, Address_None},
	0x7F: {"RRA", 0, Address_None},

	0x80: {"NOP", 0, Address_None},
	0x81: {"STA", 0, Address_None},
	0x82: {"NOP", 0, Address_None},
	0x83: {"SAX", 0, Address_None},
	0x84: {"STY", 0, Address_None},
	0x85: {"STA", 0, Address_None},
	0x86: {"STX", 0, Address_None},
	0x87: {"SAX", 0, Address_None},
	0x88: {"DEY", 0, Address_None},
	0x89: {"NOP", 0, Address_None},
	0x8A: {"TXA", 0, Address_None},
	0x8B: {"ANE", 0, Address_None},
	0x8C: {"STY", 0, Address_None},
	0x8D: {"STA", 0, Address_None},
	0x8E: {"STX", 0, Address_None},
	0x8F: {"SAX", 0, Address_None},

	0x90: {"BCC", 0, Address_None},
	0x91: {"STA", 0, Address_None},
	0x92: {"HLT", 0, Address_None},
	0x93: {"SHA", 0, Address_None},
	0x94: {"STY", 0, Address_None},
	0x95: {"STA", 0, Address_None},
	0x96: {"STX", 0, Address_None},
	0x97: {"SAX", 0, Address_None},
	0x98: {"TYA", 0, Address_None},
	0x99: {"STA", 0, Address_None},
	0x9A: {"TXS", 0, Address_None},
	0x9B: {"SHS", 0, Address_None},
	0x9C: {"SHY", 0, Address_None},
	0x9D: {"STA", 0, Address_None},
	0x9E: {"SHX", 0, Address_None},
	0x9F: {"SHA", 0, Address_None},

	0xA0: {"LDY", 0, Address_None},
	0xA1: {"LDA", 0, Address_None},
	0xA2: {"LDX", 0, Address_None},
	0xA3: {"LAX", 0, Address_None},
	0xA4: {"LDY", 0, Address_None},
	0xA5: {"LDA", 0, Address_None},
	0xA6: {"LDX", 0, Address_None},
	0xA7: {"LAX", 0, Address_None},
	0xA8: {"TAY", 0, Address_None},
	0xA9: {"LDA", 0, Address_None},
	0xAA: {"TAX", 0, Address_None},
	0xAB: {"LXA", 0, Address_None},
	0xAC: {"LDY", 0, Address_None},
	0xAD: {"LDA", 0, Address_None},
	0xAE: {"LDX", 0, Address_None},
	0xAF: {"LAX", 0, Address_None},

	0xB0: {"BCS", 0, Address_None},
	0xB1: {"LDA", 0, Address_None},
	0xB2: {"HLT", 0, Address_None},
	0xB3: {"LAX", 0, Address_None},
	0xB4: {"LDY", 0, Address_None},
	0xB5: {"LDA", 0, Address_None},
	0xB6: {"LDX", 0, Address_None},
	0xB7: {"LAX", 0, Address_None},
	0xB8: {"CLV", 0, Address_None},
	0xB9: {"LDA", 0, Address_None},
	0xBA: {"TSX", 0, Address_None},
	0xBB: {"LAE", 0, Address_None},
	0xBC: {"LDY", 0, Address_None},
	0xBD: {"LDA", 0, Address_None},
	0xBE: {"LDX", 0, Address_None},
	0xBF: {"LAX", 0, Address_None},

	0xC0: {"CPY", 0, Address_None},
	0xC1: {"CMP", 0, Address_None},
	0xC2: {"NOP", 0, Address_None},
	0xC3: {"DCP", 0, Address_None},
	0xC4: {"CPY", 0, Address_None},
	0xC5: {"CMP", 0, Address_None},
	0xC6: {"DEC", 0, Address_None},
	0xC7: {"DCP", 0, Address_None},
	0xC8: {"INY", 0, Address_None},
	0xC9: {"CMP", 0, Address_None},
	0xCA: {"DEX", 0, Address_None},
	0xCB: {"AXS", 0, Address_None},
	0xCC: {"CPY", 0, Address_None},
	0xCD: {"CMP", 0, Address_None},
	0xCE: {"DEC", 0, Address_None},
	0xCF: {"DCP", 0, Address_None},

	0xD0: {"BNE", 0, Address_None},
	0xD1: {"CMP", 0, Address_None},
	0xD2: {"HLT", 0, Address_None},
	0xD3: {"DCP", 0, Address_None},
	0xD4: {"NOP", 0, Address_None},
	0xD5: {"CMP", 0, Address_None},
	0xD6: {"DEC", 0, Address_None},
	0xD7: {"DPC", 0, Address_None},
	0xD8: {"CLD", 0, Address_None},
	0xD9: {"CMP", 0, Address_None},
	0xDA: {"NOP", 0, Address_None},
	0xDB: {"DCP", 0, Address_None},
	0xDC: {"NOP", 0, Address_None},
	0xDD: {"CMP", 0, Address_None},
	0xDE: {"DEC", 0, Address_None},
	0xDF: {"DCP", 0, Address_None},

	0xE0: {"CPX", 0, Address_None},
	0xE1: {"SBC", 0, Address_None},
	0xE2: {"NOP", 0, Address_None},
	0xE3: {"ISC", 0, Address_None},
	0xE4: {"CPX", 0, Address_None},
	0xE5: {"SBC", 0, Address_None},
	0xE6: {"INC", 0, Address_None},
	0xE7: {"ISC", 0, Address_None},
	0xE8: {"INX", 0, Address_None},
	0xE9: {"SBC", 0, Address_None},
	0xEA: {"NOP", 0, Address_None},
	0xEB: {"SBC", 0, Address_None},
	0xEC: {"CPX", 0, Address_None},
	0xED: {"SBC", 0, Address_None},
	0xEE: {"INC", 0, Address_None},
	0xEF: {"ISC", 0, Address_None},

	0xF0: {"BEQ", 0, Address_None},
	0xF1: {"SBC", 0, Address_None},
	0xF2: {"HLT", 0, Address_None},
	0xF3: {"ISC", 0, Address_None},
	0xF4: {"NOP", 0, Address_None},
	0xF5: {"SBC", 0, Address_None},
	0xF6: {"INC", 0, Address_None},
	0xF7: {"ISC", 0, Address_None},
	0xF8: {"SED", 0, Address_None},
	0xF9: {"SBC", 0, Address_None},
	0xFA: {"NOP", 0, Address_None},
	0xFB: {"ISC", 0, Address_None},
	0xFC: {"NOP", 0, Address_None},
	0xFD: {"SBC", 0, Address_None},
	0xFE: {"INC", 0, Address_None},
	0xFF: {"ISC", 0, Address_None},
}

func TraceLogger(opcode, A, X, Y, SP, status byte, pc uint16) {
	//if LoggingCPU /*&& InstructionCount > 35000*/ {
	OP := opcodeInfo[opcode]
	// Program Counter address
	Traceline := "$" + fmt.Sprintf("%04X", pc)

	// Current opcode byte
	Traceline += "\t" + fmt.Sprintf("%02X ", opcode)

	// operands, if any
	bus := bus.NewBUS()
	for i := 0; i < OP.operands; i++ {
		Traceline += fmt.Sprintf("%02X ", bus.Read(pc+uint16(i)))
	}

	if OP.operands < 2 {
		Traceline += "\t"
	}

	// Opcode name
	Traceline += "\t" + OP.name

	// Write out instruction based on addressing and type
	switch opcode {

	}

	/*if opcode == 0x99 {
		Traceline += " $" + fmt.Sprintf("%02X", operands[1]) + fmt.Sprintf("%02X", operands[0])
		Traceline += ", Y -> $"
		Traceline += fmt.Sprintf("%04X", BuildAddress(operands[0], operands[1])+uint16(Y))
	}*/

	Traceline += "\tA:" + fmt.Sprintf("%02X", A)
	Traceline += "\tX:" + fmt.Sprintf("%02X", X)
	Traceline += "\tY:" + fmt.Sprintf("%02X", Y)
	Traceline += "\tSP:" + fmt.Sprintf("%02X", SP)
	Traceline += "\t"

	//Processor Flags
	if (status & 0x80) != 0 { // Negative flag
		Traceline += "N"
	} else {
		Traceline += "n"
	}

	if (status & 0x40) != 0 { // Overflow flag
		Traceline += "V"
	} else {
		Traceline += "v"
	}

	Traceline += "-"

	if (status & 0x10) != 0 { // "B" flag
		Traceline += "B"
	} else {
		Traceline += "b"
	}

	if (status & 0x08) != 0 { // Decimal flag
		Traceline += "D"
	} else {
		Traceline += "d"
	}

	if (status & 0x04) != 0 { // Interrupt Disable flag
		Traceline += "I"
	} else {
		Traceline += "i"
	}

	if (status & 0x02) != 0 { // Zero flag
		Traceline += "Z"
	} else {
		Traceline += "z"
	}

	if (status & 0x01) != 0 { // Carry flag
		Traceline += "C"
	} else {
		Traceline += "c"
	}

	Traceline += "\tVRAM: " + fmt.Sprintf("%04X", ppu.VRAMAddress)
	Traceline += "\tCycle: " + fmt.Sprintf("%d", common.CPU_TotalCycles)
	//Traceline += "\tInstructionCount: " + fmt.Sprintf("%d", InstructionCount)

	//fmt.Println(Traceline)
	log.Println(Traceline)

	/*LogCount--
	if LogCount < 0 {
		CPU_Halted = true
	}*/
	//}
}

func setupLogging(serviceName string) (*os.File, error) {
	logDir := "logs"

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02")
	logPath := filepath.Join(logDir, fmt.Sprintf("%s-%s.log", serviceName, timestamp))

	return os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
}

func SetupTraceLogger() {
	if logFile, err := setupLogging("tracelog"); err == nil {
		log.SetOutput(logFile)
	} else {
		log.SetOutput(os.Stderr) // Always have a fallback
		log.Printf("Failed to open log file, using stderr: %v", err)
	}
}

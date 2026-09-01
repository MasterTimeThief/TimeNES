package cpu

import (
	"mtt/timenes/common"
)

func (cpu *CPU) BuildAddress(low, high byte) uint16 {
	//cpu.AddressBus = (uint16(Value_High)<<8 | uint16(Value_Low))
	cpu.AddressBus = common.Combine2Bytes(low, high)
	return cpu.AddressBus
}

// Read from address
func (cpu *CPU) Read(Address uint16) byte {
	return cpu.bus.Read(Address)
}

// Write to address
func (cpu *CPU) Write(Address uint16, Value byte) {
	cpu.bus.Write(Address, Value)
}

// Read from the Program Counter, and return the result
func (cpu *CPU) ReadFromPC() byte {
	Value := cpu.Read(cpu.PC)
	cpu.PC++
	return Value
}

// Read from the Address Bus, and return the result
func (cpu *CPU) ReadFromAB() byte {
	return cpu.Read(cpu.AddressBus)
}

func (cpu *CPU) WriteToPC(Value byte) {
	cpu.Write(cpu.PC, Value)
}

func (cpu *CPU) WriteToAB(Value byte) {
	cpu.Write(cpu.AddressBus, Value)
}

func (cpu *CPU) SetZNFlags(Value byte) {
	cpu.flag_Zero = (Value == 0x00)
	cpu.flag_Negative = (Value >= 0x80)
}

// Branch to new address if condition is true
//
// 2-4 Steps
func (cpu *CPU) Branch(condition bool) {
	switch cpu.subCycle {
	case 1:
		cpu.PollInterrupts()
		cpu.DL = cpu.ReadFromPC()
		cpu.AddressBus = cpu.PC
		if !condition {
			cpu.CompleteInstruction()
		}
	case 2:
		cpu.ReadFromAB() // Dummy read
		signedVal := int(cpu.DL)
		if signedVal >= 128 {
			signedVal -= 256 //range from -128 to 127
		}
		cpu.TempAddress = uint16(cpu.PC + uint16(signedVal))
		cpu.PC = (cpu.PC & 0xFF00) | ((cpu.PC + uint16(cpu.DL)) & 0xFF)
		cpu.AddressBus = cpu.PC
		if (cpu.TempAddress & 0xFF00) == (cpu.PC & 0xFF00) {
			cpu.CompleteInstruction()
		}
	case 3:
		cpu.PollInterrupts_CantDisableIRQ() // If the first poll detected an IRQ, this second poll should not be allowed to un-set the IRQ.
		cpu.ReadFromAB()                    // Dummy read
		cpu.PC = (cpu.TempAddress & 0xFF00) | (cpu.PC & 0xFF)
		cpu.CompleteInstruction()
	}
}

// The RESET instruction has unique behavior where it reads from the stack, and decrements the stack pointer.
func (cpu *CPU) ResetReadPush() {
	cpu.Read(uint16(cpu.SP) + 0x100)
	cpu.SP--
}

// Store to the stack, and decrement the stack pointer
func (cpu *CPU) Push(Value byte) {
	cpu.Write(uint16(cpu.SP)+0x100, Value)
	cpu.SP--
}

// Increment the stack pointer, and read from the stack
func (cpu *CPU) Pull() byte {
	cpu.SP++
	return cpu.Read(uint16(cpu.SP) + 0x100)
}

// Push the Status Register to the Stack
func (cpu *CPU) PushFlags() {
	temp := byte(0)
	temp += byte(common.Ternary(cpu.flag_Carry, 0x01, 0x00))
	temp += byte(common.Ternary(cpu.flag_Zero, 0x02, 0x00))
	temp += byte(common.Ternary(cpu.flag_InterruptDisable, 0x04, 0x00))
	temp += byte(common.Ternary(cpu.flag_Decimal, 0x08, 0x00))
	temp += byte(common.Ternary(cpu.flag_B, 0x10, 0x00)) //B Flag
	temp += 0x20
	temp += byte(common.Ternary(cpu.flag_Overflow, 0x40, 0x00))
	temp += byte(common.Ternary(cpu.flag_Negative, 0x80, 0x00))
	cpu.Push(temp)
}

// Pull the Status Register from the Stack
func (cpu *CPU) PullFlags() {
	temp := cpu.Pull()
	cpu.flag_Carry = (temp & 0x01) != 0
	cpu.flag_Zero = (temp & 0x02) != 0
	cpu.flag_InterruptDisable = (temp & 0x04) != 0
	cpu.flag_Decimal = (temp & 0x08) != 0
	cpu.flag_B = (temp & 0x10) != 0
	cpu.flag_Overflow = (temp & 0x40) != 0
	cpu.flag_Negative = (temp & 0x80) != 0
}

// Performs Arithmetic Shift Left onto value at Address
func (cpu *CPU) Op_ASL() {
	Value := cpu.Read(cpu.AddressBus)
	cpu.flag_Carry = (Value >= 0x80)
	Value <<= 1
	cpu.Write(cpu.AddressBus, Value)
	cpu.SetZNFlags(Value)
}

// Performs Arithmetic Shift Right onto value at Address
func (cpu *CPU) Op_LSR() {
	Value := cpu.Read(cpu.AddressBus)
	cpu.flag_Carry = (Value & 1) != 0
	Value >>= 1
	cpu.Write(cpu.AddressBus, Value)
	cpu.SetZNFlags(Value)
}

// Perform Rotate Left onto value at Address
func (cpu *CPU) Op_ROL() {
	Value := cpu.Read(cpu.AddressBus)
	futureCarry := (Value >= 0x80)
	Value <<= 1
	if cpu.flag_Carry {
		Value |= 1
	}
	cpu.Write(cpu.AddressBus, Value)
	cpu.flag_Carry = futureCarry
	cpu.SetZNFlags(Value)
}

// Perform Rotate Right onto value at Address
func (cpu *CPU) Op_ROR() {
	Value := cpu.Read(cpu.AddressBus)
	futureCarry := (Value & 1) != 0
	Value >>= 1
	if cpu.flag_Carry {
		Value |= 0x80
	}
	cpu.Write(cpu.AddressBus, Value)
	cpu.flag_Carry = futureCarry
	cpu.SetZNFlags(Value)
}

// Increment Value, and save to Address
func (cpu *CPU) Op_INC(Value byte) {
	Value++
	cpu.Write(cpu.AddressBus, Value)
	cpu.SetZNFlags(Value)
}

// Decrement Value, and save to Address
func (cpu *CPU) Op_DEC(Value byte) {
	Value--
	cpu.Write(cpu.AddressBus, Value)
	cpu.SetZNFlags(Value)
}

// Bitwise OR with A
func (cpu *CPU) Op_ORA(Value byte) {
	cpu.A |= Value
	cpu.SetZNFlags(cpu.A)
}

// Bitwise AND with A
func (cpu *CPU) Op_AND(Value byte) {
	cpu.A &= Value
	cpu.SetZNFlags(cpu.A)
}

// Bitwise XOR with A
func (cpu *CPU) Op_EOR(Value byte) {
	cpu.A ^= Value
	cpu.SetZNFlags(cpu.A)
}

// Add Value to A with Carry
func (cpu *CPU) Op_ADC(Value byte) {
	IntSum := int(cpu.A) + int(Value) + common.BoolToInt(cpu.flag_Carry)
	cpu.flag_Overflow = (^int(cpu.A^Value) & (int(cpu.A) ^ IntSum) & 0x80) != 0
	cpu.flag_Carry = IntSum > 0xFF
	cpu.A = byte(IntSum)
	cpu.SetZNFlags(cpu.A)
}

// Subtract Value from A with Carry
func (cpu *CPU) Op_SBC(Value byte) {
	IntSum := int(cpu.A) - int(Value) - common.BoolToInt(!cpu.flag_Carry)
	cpu.flag_Overflow = (int(cpu.A^Value) & (int(cpu.A) ^ IntSum) & 0x80) != 0
	cpu.flag_Carry = IntSum >= 0x00
	cpu.A = byte(IntSum)
	cpu.SetZNFlags(cpu.A)
}

// Compare Value with A
func (cpu *CPU) Op_CMP(Value byte) {
	cpu.flag_Carry = Value <= cpu.A
	cpu.flag_Zero = (Value == cpu.A)
	cpu.flag_Negative = ((cpu.A - Value) >= 0x80)
}

// Compare Value with X
func (cpu *CPU) Op_CPX(Value byte) {
	cpu.flag_Carry = Value <= cpu.X
	cpu.flag_Zero = (Value == cpu.X)
	cpu.flag_Negative = ((cpu.X - Value) >= 0x80)
}

// Compare Value with Y
func (cpu *CPU) Op_CPY(Value byte) {
	cpu.flag_Carry = Value <= cpu.Y
	cpu.flag_Zero = (Value == cpu.Y)
	cpu.flag_Negative = ((cpu.Y - Value) >= 0x80)
}

// Bit Test
func (cpu *CPU) Op_BIT(Value byte) {
	cpu.flag_Zero = ((cpu.A & Value) == 0)
	cpu.flag_Negative = ((Value & 0x80) != 0)
	cpu.flag_Overflow = ((Value & 0x40) != 0)
}

//
// Illegal operations
//

// SLO: ASL + ORA
func (cpu *CPU) Op_SLO() {
	cpu.Op_ASL()
	cpu.Op_ORA(cpu.DL)
}

// CPU Kill function
func (cpu *CPU) Kill() {
	CPU_Halted = true
	common.Filepath = ""
	common.ROMExists = false
	common.ROMLoaded = false
	common.SetUIMessage("Game Crashed!")
}

func (cpu *CPU) DelayCPU(delay int) {
	cpu.DelayCounter = delay
}

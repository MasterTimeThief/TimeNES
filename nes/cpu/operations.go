package cpu

import (
	"mtt/timenes/common"
)

func (c *CPU) BuildAddress(low, high byte) uint16 {
	//c.AddressBus = (uint16(Value_High)<<8 | uint16(Value_Low))
	c.AddressBus = common.Combine2Bytes(low, high)
	return c.AddressBus
}

// Read from address
func (c *CPU) Read(Address uint16) byte {
	return c.bus.Read(Address)
}

// Write to address
func (c *CPU) Write(Address uint16, Value byte) {
	c.bus.Write(Address, Value)
}

// Read from the Program Counter, and return the result
func (c *CPU) ReadFromPC() byte {
	Value := c.Read(c.PC)
	c.PC++
	return Value
}

// Read from the Address Bus, and return the result
func (c *CPU) ReadFromAB() byte {
	return c.Read(c.AddressBus)
}

func (c *CPU) WriteToPC(Value byte) {
	c.Write(c.PC, Value)
}

func (c *CPU) WriteToAB(Value byte) {
	c.Write(c.AddressBus, Value)
}

func (c *CPU) SetZNFlags(Value byte) {
	c.flag_Zero = (Value == 0x00)
	c.flag_Negative = (Value >= 0x80)
}

// Branch to new address if condition is true
//
// 2-4 Steps
func (c *CPU) Branch(condition bool) {
	switch c.InstructionCycle {
	case 1:
		c.PollInterrupts()
		c.DL = c.ReadFromPC()
		c.AddressBus = c.PC
		if !condition {
			c.CompleteInstruction()
		}
	case 2:
		c.AddressBus = c.PC
		c.ReadFromAB() // Dummy read
		signedVal := int(c.DL)
		if signedVal >= 128 {
			signedVal -= 256 //range from -128 to 127
		}
		c.TempAddress = uint16(c.PC + uint16(signedVal))
		c.PC = (c.PC & 0xFF00) | ((c.PC + uint16(c.DL)) & 0xFF)
		if (c.TempAddress & 0xFF00) == (c.PC & 0xFF00) {
			c.CompleteInstruction()
		}
	case 3:
		c.AddressBus = c.PC
		c.PollInterrupts_CantDisableIRQ() // If the first poll detected an IRQ, this second poll should not be allowed to un-set the IRQ.
		c.ReadFromAB()                    // Dummy read
		c.PC = (c.TempAddress & 0xFF00) | (c.PC & 0xFF)
		c.CompleteInstruction()
	}
}

// The RESET instruction has unique behavior where it reads from the stack, and decrements the stack pointer.
func (c *CPU) ResetReadPush() {
	c.Read(uint16(c.SP) + 0x100)
	c.SP--
}

// Store to the stack, and decrement the stack pointer
func (c *CPU) Push(Value byte) {
	c.Write(uint16(c.SP)+0x100, Value)
	c.SP--
}

// Increment the stack pointer, and read from the stack
func (c *CPU) Pull() byte {
	c.SP++
	return c.Read(uint16(c.SP) + 0x100)
}

// Push the Status Register to the Stack
func (c *CPU) PushFlags() {
	temp := byte(0)
	temp += byte(common.Ternary(c.flag_Carry, 0x01, 0x00))
	temp += byte(common.Ternary(c.flag_Zero, 0x02, 0x00))
	temp += byte(common.Ternary(c.flag_InterruptDisable, 0x04, 0x00))
	temp += byte(common.Ternary(c.flag_Decimal, 0x08, 0x00))
	temp += byte(common.Ternary(c.flag_B, 0x10, 0x00)) //B Flag
	temp += 0x20
	temp += byte(common.Ternary(c.flag_Overflow, 0x40, 0x00))
	temp += byte(common.Ternary(c.flag_Negative, 0x80, 0x00))
	c.Push(temp)
}

// Pull the Status Register from the Stack
func (c *CPU) PullFlags() {
	temp := c.Pull()
	c.flag_Carry = (temp & 0x01) != 0
	c.flag_Zero = (temp & 0x02) != 0
	c.flag_InterruptDisable = (temp & 0x04) != 0
	c.flag_Decimal = (temp & 0x08) != 0
	c.flag_B = (temp & 0x10) != 0
	c.flag_Overflow = (temp & 0x40) != 0
	c.flag_Negative = (temp & 0x80) != 0
}

// Performs Arithmetic Shift Left onto value at Address
func (c *CPU) Op_ASL() {
	//Value := c.Read(c.AddressBus)
	c.flag_Carry = (c.DL >= 0x80)
	c.DL <<= 1
	c.WriteToAB(c.DL)
	c.SetZNFlags(c.DL)
}

// Performs Arithmetic Shift Right onto value at Address
func (c *CPU) Op_LSR() {
	//Value := c.Read(c.AddressBus)
	c.flag_Carry = (c.DL & 1) != 0
	c.DL >>= 1
	c.WriteToAB(c.DL)
	c.SetZNFlags(c.DL)
}

// Perform Rotate Left onto value at Address
func (c *CPU) Op_ROL() {
	//Value := c.Read(c.AddressBus)
	futureCarry := (c.DL >= 0x80)
	c.DL <<= 1
	if c.flag_Carry {
		c.DL |= 1
	}
	c.WriteToAB(c.DL)
	c.flag_Carry = futureCarry
	c.SetZNFlags(c.DL)
}

// Perform Rotate Right onto value at Address
func (c *CPU) Op_ROR() {
	//Value := c.Read(c.AddressBus)
	futureCarry := (c.DL & 1) != 0
	c.DL >>= 1
	if c.flag_Carry {
		c.DL |= 0x80
	}
	c.WriteToAB(c.DL)
	c.flag_Carry = futureCarry
	c.SetZNFlags(c.DL)
}

// Increment Value, and save to Address
func (c *CPU) Op_INC(Value byte) {
	Value++
	c.WriteToAB(Value)
	c.SetZNFlags(Value)
}

// Decrement Value, and save to Address
func (c *CPU) Op_DEC(Value byte) {
	Value--
	c.WriteToAB(Value)
	c.SetZNFlags(Value)
}

// Bitwise OR with A
func (c *CPU) Op_ORA(Value byte) {
	c.A |= Value
	c.SetZNFlags(c.A)
}

// Bitwise AND with A
func (c *CPU) Op_AND(Value byte) {
	c.A &= Value
	c.SetZNFlags(c.A)
}

// Bitwise XOR with A
func (c *CPU) Op_EOR(Value byte) {
	c.A ^= Value
	c.SetZNFlags(c.A)
}

// Add Value to A with Carry
func (c *CPU) Op_ADC(Value byte) {
	IntSum := int(c.A) + int(Value) + common.BoolToInt(c.flag_Carry)
	c.flag_Overflow = (^int(c.A^Value) & (int(c.A) ^ IntSum) & 0x80) != 0
	c.flag_Carry = IntSum > 0xFF
	c.A = byte(IntSum)
	c.SetZNFlags(c.A)
}

// Subtract Value from A with Carry
func (c *CPU) Op_SBC(Value byte) {
	IntSum := int(c.A) - int(Value) - common.BoolToInt(!c.flag_Carry)
	c.flag_Overflow = (int(c.A^Value) & (int(c.A) ^ IntSum) & 0x80) != 0
	c.flag_Carry = IntSum >= 0x00
	c.A = byte(IntSum)
	c.SetZNFlags(c.A)
}

// Compare Value with A
func (c *CPU) Op_CMP(Value byte) {
	c.flag_Carry = Value <= c.A
	c.flag_Zero = (Value == c.A)
	c.flag_Negative = ((c.A - Value) >= 0x80)
}

// Compare Value with X
func (c *CPU) Op_CPX(Value byte) {
	c.flag_Carry = Value <= c.X
	c.flag_Zero = (Value == c.X)
	c.flag_Negative = ((c.X - Value) >= 0x80)
}

// Compare Value with Y
func (c *CPU) Op_CPY(Value byte) {
	c.flag_Carry = Value <= c.Y
	c.flag_Zero = (Value == c.Y)
	c.flag_Negative = ((c.Y - Value) >= 0x80)
}

// Bit Test
func (c *CPU) Op_BIT(Value byte) {
	c.flag_Zero = ((c.A & Value) == 0)
	c.flag_Negative = ((Value & 0x80) != 0)
	c.flag_Overflow = ((Value & 0x40) != 0)
}

//
// Illegal operations
//

// SLO: ASL + ORA
func (c *CPU) Op_SLO() {
	c.Op_ASL()
	c.Op_ORA(c.DL)
}

// DCP: DEC + CMP
func (c *CPU) Op_DCP() {
	c.Op_DEC(c.ReadFromAB())
	c.Op_CMP(c.ReadFromAB())
}

// RLA: ROL + AND
func (c *CPU) Op_RLA() {
	c.Op_ROL()
	c.Op_AND(c.ReadFromAB())
}

// SRE: LSR + EOR
func (c *CPU) Op_SRE() {
	c.Op_LSR()
	c.Op_EOR(c.ReadFromAB())
}

// RRA: ROR + ADC
func (c *CPU) Op_RRA() {
	c.Op_ROR()
	c.Op_ADC(c.ReadFromAB())
}

// ISC: INC + SBC
func (c *CPU) Op_ISC() {
	c.Op_INC(c.ReadFromAB())
	c.Op_SBC(c.ReadFromAB())
}

// SH* Unstability
func (c *CPU) UnstableAddressBus(Value byte) {
	c.AddressBus = c.AddressBus | ((c.AddressBus>>8)&uint16(Value))<<8
}

// CPU Kill function
func (c *CPU) Kill() {
	CPU_Halted = true
	common.Filepath = ""
	common.ROMExists = false
	common.ROMLoaded = false
	common.SetUIMessage("Game Crashed!")
}

func (c *CPU) DelayCPU(delay int) {
	c.DelayCounter = delay
}

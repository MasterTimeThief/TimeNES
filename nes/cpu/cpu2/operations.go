package cpu2

import (
	"mtt/timenes/common"
	"mtt/timenes/nes/bus"
)

func (cpu *CPU) BuildAddress(low, high byte) uint16 {
	//cpu.AddressBus = (uint16(Value_High)<<8 | uint16(Value_Low))
	cpu.AddressBus = common.Combine2Bytes(low, high)
	return cpu.AddressBus
}

// Read from the Program Counter, and return the result
func (cpu *CPU) ReadFromPC() byte {
	Value := bus.Read(cpu.PC)
	//if LoggingCPU {
	//	operands = append(operands, Value)
	//}
	cpu.PC++
	return Value
}

// Read from the Address Bus, and return the result
func (cpu *CPU) ReadFromAB() byte {
	return bus.Read(cpu.AddressBus)
}

func (cpu *CPU) WriteToPC(Value byte) {
	bus.Write(cpu.PC, Value)
}

func (cpu *CPU) WriteToAB(Value byte) {
	bus.Write(cpu.AddressBus, Value)
}

func (cpu *CPU) SetZNFlags(Value byte) {
	cpu.flag_Zero = (Value == 0x00)
	cpu.flag_Negative = (Value >= 0x80)
}

func (cpu *CPU) Branch(Condition bool, Value byte) {
	if Condition {
		//fmt.Println("branch taken")
		signedVal := int(Value)
		if signedVal > 127 {
			signedVal -= 256 //range from -128 to 127
		}
		//CPU_Cycles = 3
		if byte((cpu.PC&0xFF00)>>4) != byte(((cpu.PC+uint16(signedVal))&0xFF00)>>4) {
			//CPU_Cycles++ //Extra cycle for crossing page boundary
			//MasterClockTick("branch page cross")
		}
		cpu.PC = uint16(cpu.PC + uint16(signedVal))
		//MasterClockTick("branch taken")
	} else {
		//fmt.Println("branch not taken")
		//CPU_Cycles = 2
	}
}

func (cpu *CPU) Push(Value byte) {
	//Store to the stack, and decrement the stack pointer
	bus.Write(uint16(cpu.SP)+0x100, Value)
	////MasterClockTick("push")
	cpu.SP--
}

func (cpu *CPU) Pull() byte {
	//Increment the stack pointer, and read from the stack
	cpu.SP++
	//MasterClockTick("pull SP++")
	temp := bus.Read(uint16(cpu.SP) + 0x100)
	//MasterClockTick("pull")
	return temp
}

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

func (cpu *CPU) PullFlags() {
	temp := cpu.Pull()
	cpu.flag_Carry = (temp & 0x01) != 0
	cpu.flag_Zero = (temp & 0x02) != 0
	cpu.flag_InterruptDisable = (temp & 0x04) != 0
	cpu.flag_Decimal = (temp & 0x08) != 0
	cpu.flag_B = (temp & 0x10) != 0
	cpu.flag_Overflow = (temp & 0x40) != 0
	cpu.flag_Negative = (temp & 0x80) != 0
	//MasterClockTick("pull flags")
}

// Performs Arithmetic Shift Left onto value at Address
func (cpu *CPU) Op_ASL(Address uint16) {
	Value := bus.Read(Address)
	cpu.flag_Carry = (Value >= 0x80)
	Value <<= 1
	bus.Write(Address, Value)
	cpu.SetZNFlags(Value)

}

// Performs Arithmetic Shift Right onto value at Address
func (cpu *CPU) Op_LSR(Address uint16) {
	Value := bus.Read(Address)
	cpu.flag_Carry = (Value & 1) != 0
	Value >>= 1
	bus.Write(Address, Value)
	cpu.SetZNFlags(Value)
}

// Perform Rotate Left onto value at Address
func (cpu *CPU) Op_ROL(Address uint16) {
	Value := bus.Read(Address)
	futureCarry := (Value >= 0x80)
	Value <<= 1
	if cpu.flag_Carry {
		Value |= 1
	}
	//MasterClockTick("rol")
	bus.Write(Address, Value)
	cpu.flag_Carry = futureCarry
	cpu.SetZNFlags(Value)
}

// Perform Rotate Right onto value at Address
func (cpu *CPU) Op_ROR(Address uint16) {
	Value := bus.Read(Address)
	futureCarry := (Value & 1) != 0
	Value >>= 1
	if cpu.flag_Carry {
		Value |= 0x80
	}
	//MasterClockTick("ror")
	bus.Write(Address, Value)
	cpu.flag_Carry = futureCarry
	cpu.SetZNFlags(Value)
}

// Increment Value, and save to Address
func (cpu *CPU) Op_INC(Address uint16, Value byte) {
	Value++
	bus.Write(Address, Value)
	cpu.SetZNFlags(Value)
}

// Decrement Value, and save to Address
func (cpu *CPU) Op_DEC(Address uint16, Value byte) {
	Value--
	bus.Write(Address, Value)
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

// Uhh
func (cpu *CPU) Op_BIT(Value byte) {
	//Bit Test
	cpu.flag_Zero = ((cpu.A & Value) == 0)
	cpu.flag_Negative = ((Value & 0x80) != 0)
	cpu.flag_Overflow = ((Value & 0x40) != 0)
}

// CPU Kill function
func (cpu *CPU) Kill() {
	cpu.CPU_Halted = true
	common.Filepath = ""
	common.ROMExists = false
	common.ROMLoaded = false
	common.SetUIMessage("Game Crashed!")
}

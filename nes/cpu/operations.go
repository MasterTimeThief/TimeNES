package cpu

import (
	"mtt/timenes/common"
	"mtt/timenes/nes/bus"
)

// Performs Arithmetic Shift Left onto value at Address
func Op_ASL(Address uint16) {
	Value := bus.Read(Address)
	flag_Carry = (Value >= 0x80)
	Value <<= 1
	bus.Write(Address, Value)
	SetZNFlags(Value)

}

// Performs Arithmetic Shift Right onto value at Address
func Op_LSR(Address uint16) {
	Value := bus.Read(Address)
	flag_Carry = (Value & 1) != 0
	Value >>= 1
	bus.Write(Address, Value)
	SetZNFlags(Value)
}

// Perform Rotate Left onto value at Address
func Op_ROL(Address uint16) {
	Value := bus.Read(Address)
	futureCarry := (Value >= 0x80)
	Value <<= 1
	if flag_Carry {
		Value |= 1
	}
	//MasterClockTick("rol")
	bus.Write(Address, Value)
	flag_Carry = futureCarry
	SetZNFlags(Value)
}

// Perform Rotate Right onto value at Address
func Op_ROR(Address uint16) {
	Value := bus.Read(Address)
	futureCarry := (Value & 1) != 0
	Value >>= 1
	if flag_Carry {
		Value |= 0x80
	}
	//MasterClockTick("ror")
	bus.Write(Address, Value)
	flag_Carry = futureCarry
	SetZNFlags(Value)
}

// Increment Value, and save to Address
func Op_INC(Address uint16, Value byte) {
	Value++
	bus.Write(Address, Value)
	SetZNFlags(Value)
}

// Decrement Value, and save to Address
func Op_DEC(Address uint16, Value byte) {
	Value--
	bus.Write(Address, Value)
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
	IntSum := int(A) + int(Value) + common.BoolToInt(flag_Carry)
	flag_Overflow = (^int(A^Value) & (int(A) ^ IntSum) & 0x80) != 0
	flag_Carry = IntSum > 0xFF
	A = byte(IntSum)
	SetZNFlags(A)
}

// Subtract Value from A with Carry
func Op_SBC(Value byte) {
	IntSum := int(A) - int(Value) - common.BoolToInt(!flag_Carry)
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

package cpu

import (
	"mtt/timenes/common"
	"mtt/timenes/nes/bus"
)

func BuildAddress(low, high byte) uint16 {
	//AddressBus = (uint16(Value_High)<<8 | uint16(Value_Low))
	AddressBus = common.Combine2Bytes(low, high)
	return AddressBus
}

// Read from the Program Counter, and return the result
func ReadFromPC() byte {
	Value := bus.Read(PC)
	//if LoggingCPU {
	//	operands = append(operands, Value)
	//}
	PC++
	return Value
}

// Read from the Address Bus, and return the result
func ReadFromAB() byte {
	return bus.Read(AddressBus)
}

func WriteToPC(Value byte) {
	bus.Write(PC, Value)
}

func WriteToAB(Value byte) {
	bus.Write(AddressBus, Value)
}

func SetZNFlags(Value byte) {
	flag_Zero = (Value == 0x00)
	flag_Negative = (Value >= 0x80)
}

func Branch(Condition bool, Value byte) {
	if Condition {
		//fmt.Println("branch taken")
		signedVal := int(Value)
		if signedVal > 127 {
			signedVal -= 256 //range from -128 to 127
		}
		CPU_Cycles = 3
		if byte((PC&0xFF00)>>4) != byte(((PC+uint16(signedVal))&0xFF00)>>4) {
			CPU_Cycles++ //Extra cycle for crossing page boundary
			//MasterClockTick("branch page cross")
		}
		PC = uint16(PC + uint16(signedVal))
		//MasterClockTick("branch taken")
	} else {
		//fmt.Println("branch not taken")
		CPU_Cycles = 2
	}
}

func Push(Value byte) {
	//Store to the stack, and decrement the stack pointer
	bus.Write(uint16(SP)+0x100, Value)
	////MasterClockTick("push")
	SP--
}

func Pull() byte {
	//Increment the stack pointer, and read from the stack
	SP++
	//MasterClockTick("pull SP++")
	temp := bus.Read(uint16(SP) + 0x100)
	//MasterClockTick("pull")
	return temp
}

func PushFlags() {
	temp := byte(0)
	temp += byte(common.Ternary(flag_Carry, 0x01, 0x00))
	temp += byte(common.Ternary(flag_Zero, 0x02, 0x00))
	temp += byte(common.Ternary(flag_InterruptDisable, 0x04, 0x00))
	temp += byte(common.Ternary(flag_Decimal, 0x08, 0x00))
	temp += byte(common.Ternary(flag_B, 0x10, 0x00)) //B Flag
	temp += 0x20
	temp += byte(common.Ternary(flag_Overflow, 0x40, 0x00))
	temp += byte(common.Ternary(flag_Negative, 0x80, 0x00))
	Push(temp)

}

func PullFlags() {
	temp := Pull()
	flag_Carry = (temp & 0x01) != 0
	flag_Zero = (temp & 0x02) != 0
	flag_InterruptDisable = (temp & 0x04) != 0
	flag_Decimal = (temp & 0x08) != 0
	flag_B = (temp & 0x10) != 0
	flag_Overflow = (temp & 0x40) != 0
	flag_Negative = (temp & 0x80) != 0
	//MasterClockTick("pull flags")
}

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

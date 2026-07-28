package main

type Pulse1Channel struct {
	Duty              byte
	LengthCounterHalt bool
	ConstantVolume    bool
	Envelope          byte
	SweepUnitEnabled  bool
	Period            byte
	Negate            bool
	ShiftRegister     byte
	Timer             uint16
	LengthCounterLoad byte
}

type Pulse2Channel struct {
	Duty              byte
	LengthCounterHalt bool
	ConstantVolume    bool
	Envelope          byte
	SweepUnitEnabled  bool
	Period            byte
	Negate            bool
	ShiftRegister     byte
	Timer             uint16
	LengthCounterLoad byte
}

type TriangleChannel struct {
	LengthCounterControl bool
	LinearCounterLoad    byte
	Timer                uint16
	LengthCounterLoad    byte
}

type NoiseChannel struct {
	LengthCounterHalt bool
	ConstantVolume    bool
	Envelope          byte
	Mode              bool
	Period            byte
	LengthCounterLoad byte
}
type DeltaModChannel struct {
	IRQEnable bool
	Loop      bool
	Frequency byte
	Address   uint16
	Length    uint16
}

var apuPulse1 Pulse1Channel
var apuPulse2 Pulse2Channel
var apuTriangle TriangleChannel
var apuNoise NoiseChannel
var apuDMC DeltaModChannel

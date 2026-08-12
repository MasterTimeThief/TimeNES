package apu

type NoiseChannel struct {
	Enabled          bool
	Timer            uint16
	TimerReloadValue uint16
	Output           byte
	Mode             bool
	ShiftRegister    uint16
	LengthCounter
	Envelope
}

var apuNoiseTimerLUT = [16]uint16{4, 8, 16, 32, 64, 96, 128, 160, 202, 254, 380, 508, 762, 1016, 2034, 4068}

func (n *NoiseChannel) SetNoiseTimer(Value byte) {
	n.TimerReloadValue = apuNoiseTimerLUT[Value]
	n.Timer = n.TimerReloadValue
}

func (n *NoiseChannel) ResetNoise() {
	n.Enabled = false
	n.ResetLengthCounter()
	n.Timer = 0

	n.ResetEnvelope()

	n.Mode = false
	n.Output = 0
	n.ShiftRegister = 1
}

func (n *NoiseChannel) ClockNoiseTimer() {
	if n.Timer == 0 {
		n.ClockShiftRegister()
		n.Timer = n.TimerReloadValue
	} else {
		n.Timer--
	}
}

func (n *NoiseChannel) ClockShiftRegister() {
	//Feedback is calculated as the exclusive-OR of bit 0 and one other bit: bit 6 if Mode flag is set, otherwise bit 1.
	firstBit := byte(n.ShiftRegister & 1)
	var secondBit byte
	if n.Mode {
		secondBit = byte(n.ShiftRegister&0x40) >> 6
	} else {
		secondBit = byte(n.ShiftRegister&0x02) >> 1
	}
	feedback := firstBit ^ secondBit

	//The shift register is shifted right by one bit.
	n.ShiftRegister >>= 1
	n.ShiftRegister &= 0x7FFF

	//Bit 14, the leftmost bit, is set to the feedback calculated earlier.
	n.ShiftRegister = (n.ShiftRegister & 0b0011111111111111) | (uint16(feedback&1) << 14)
}

func (n *NoiseChannel) UpdateNoiseOutput() byte {
	if !apuEnabled || !n.Enabled || (n.ShiftRegister&1) == 1 || n.LengthCounter.Counter == 0 {
		return 0
	} else if n.ConstantVolume {
		return n.Volume
	} else {
		return n.Decay
	}
}

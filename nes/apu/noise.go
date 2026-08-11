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

func (noise *NoiseChannel) SetNoiseTimer(Value byte) {
	noise.TimerReloadValue = apuNoiseTimerLUT[Value]
	noise.Timer = noise.TimerReloadValue
}

func (noise *NoiseChannel) ClockShiftRegister() {
	//Feedback is calculated as the exclusive-OR of bit 0 and one other bit: bit 6 if Mode flag is set, otherwise bit 1.
	firstBit := byte(noise.ShiftRegister & 1)
	var secondBit byte
	if noise.Mode {
		secondBit = byte(noise.ShiftRegister&0x40) >> 6
	} else {
		secondBit = byte(noise.ShiftRegister&0x02) >> 1
	}
	feedback := firstBit ^ secondBit

	//The shift register is shifted right by one bit.
	noise.ShiftRegister >>= 1
	noise.ShiftRegister &= 0x7FFF

	//Bit 14, the leftmost bit, is set to the feedback calculated earlier.
	noise.ShiftRegister = (noise.ShiftRegister & 0b0011111111111111) | (uint16(feedback&1) << 14)
}

func (noise *NoiseChannel) UpdateNoiseOutput() {
	if !apuEnabled || !noise.Enabled || (noise.ShiftRegister&1) == 1 || noise.LengthCounter.Counter == 0 {
		noise.Output = 0
	} else {
		if noise.ConstantVolume {
			noise.Output = noise.Volume
		} else {
			noise.Output = noise.Decay
		}
	}
}

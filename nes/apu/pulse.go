package apu

type PulseChannel struct {
	Enabled          bool
	Timer            uint16
	TimerReloadValue uint16
	Output           byte
	Duty             byte
	DutyPos          byte
	LengthCounter
	Envelope
	Sweep
}

var apuDutySequences = [4][8]byte{
	{0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 1, 1},
	{0, 0, 0, 0, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 0, 0},
}

func (pulse *PulseChannel) IsPulseMuted() bool {
	return pulse.Timer < 8 || (!pulse.Sweep.Negate && pulse.Sweep.Period > 0x7FF)
}

func (pulse *PulseChannel) UpdatePulseOutput() {
	duty := apuDutySequences[pulse.Duty][pulse.DutyPos]
	var volume byte
	if pulse.ConstantVolume {
		volume = pulse.Volume
	} else {
		volume = pulse.Decay
	}

	if !apuEnabled || duty == 0 || pulse.LengthCounter.Counter == 0 || pulse.Timer < 8 {
		pulse.Output = 0
	} else {
		pulse.Output = (duty * volume)
	}
}

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

func (p *PulseChannel) IsPulseMuted() bool {
	return p.Timer < 8 || (!p.Sweep.Negate && p.Sweep.Period > 0x7FF)
}

func (p *PulseChannel) ResetPulse() {
	p.Enabled = false
	p.ResetLengthCounter()
	p.Timer = 0

	p.ResetEnvelope()

	p.Duty = 0
	p.DutyPos = 0
	p.ResetSweep()
}

func (p *PulseChannel) ClockPulseTimer() {
	if p.Timer == 0 {
		p.DutyPos = (p.DutyPos + 1) & 7
		p.Timer = p.TimerReloadValue
	} else {
		p.Timer--
	}
}

func (p *PulseChannel) UpdatePulseOutput() {
	duty := apuDutySequences[p.Duty][p.DutyPos]
	var volume byte
	if p.ConstantVolume {
		volume = p.Volume
	} else {
		volume = p.Decay
	}

	if !apuEnabled || duty == 0 || p.LengthCounter.Counter == 0 || p.Timer < 8 {
		p.Output = 0
	} else {
		p.Output = (duty * volume)
	}
}

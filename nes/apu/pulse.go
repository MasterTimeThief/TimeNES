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
	ForceMute bool
}

var apuDutySequences = [4][8]byte{
	{0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 1, 1},
	{0, 0, 0, 0, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 0, 0},
}

func (Sw *Sweep) ResetSweep() {
	Sw.Enabled = false
	Sw.ResetDivider()
	Sw.ReloadFlag = false
	Sw.Negate = false
	Sw.Shift = 0
}

func (p *PulseChannel) ClockSweep(isPulse1 bool) {
	if p.Sweep.ReloadFlag {
		if p.Sweep.Enabled && p.Sweep.Divider.Counter == 0 {
			p.UpdateTargetPeriod(isPulse1)
		}
		p.Sweep.Divider.Counter = p.Sweep.Divider.Period
		p.Sweep.ReloadFlag = false
	} else if p.Sweep.Divider.Counter > 0 {
		p.Sweep.Divider.Counter--
	} else {
		if p.Sweep.Enabled {
			p.UpdateTargetPeriod(isPulse1)
		}
		p.Sweep.Divider.Counter = p.Sweep.Divider.Period
	}
	/*
		if Sw.Counter == 0 {
			//Sweep is enabled, the shift count is nonzero
			if Sw.Enabled {
				//Pulse's period is set to target period
				Sw.Period = Sw.TargetPeriod
			} else {
				//Pulse's period remains unchanged, but the sweep unit's divider continues to count down  and reload the divider's period as normal
			}
		}

		if Sw.Counter == 0 || Sw.ReloadFlag {
			//Divider counter is set to P and the reload flag is cleared
			Sw.Counter = Sw.Period
			Sw.ReloadFlag = false
		} else {
			//Otherwise, the divider's counter is decremented
			Sw.Counter--
		}
	*/
}

func (p *PulseChannel) UpdateTargetPeriod(isPulse1 bool) {

	shift := p.TimerReloadValue >> p.Sweep.Shift

	if p.Sweep.Negate {
		shift = -shift
		if isPulse1 {
			shift--
		}
	}

	p.TimerReloadValue += shift
	/*
		if newTarget < 0 {
			newTarget = 0
		}

		Sw.TargetPeriod = uint16(newTarget)*/

}

func (p *PulseChannel) SetEnabled(enable bool) {
	p.Enabled = enable
	if !enable {
		p.LengthCounter.Counter = 0
	}
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
	p.ForceMute = false
}

func (p *PulseChannel) ClockPulseTimer() {
	if p.Timer == 0 {
		p.DutyPos = (p.DutyPos + 1) & 7
		p.Timer = p.TimerReloadValue
	} else {
		p.Timer--
	}
}

func (p *PulseChannel) UpdatePulseOutput() byte {
	duty := apuDutySequences[p.Duty][p.DutyPos]
	if MuteEmulator || p.ForceMute || !p.Enabled || duty == 0 || p.LengthCounter.Counter == 0 || p.TimerReloadValue < 8 || p.TimerReloadValue > 0x7FF {
		return 0
	} else if p.Envelope.ConstantVolume {
		return p.Envelope.Volume
	} else {
		return p.Envelope.Decay
	}
}

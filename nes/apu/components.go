package apu

type Envelope struct {
	ConstantVolume bool
	Volume         byte
	StartFlag      bool
	LoopFlag       bool
	Decay          byte
	Divider
}

type Divider struct {
	Period  uint16
	Counter uint16
}

type LengthCounter struct { //TODO: Maybe part of Envelope?
	Counter     byte
	HaltFlag    bool
	ReloadFlag  bool
	ReloadValue byte
}

type Sweep struct {
	Enabled bool
	Divider
	ReloadFlag   bool
	Negate       bool
	Shift        byte
	TargetPeriod uint16
}

func (LC *LengthCounter) ResetLengthCounter() {
	LC.Counter = 0
	LC.HaltFlag = false
	LC.ReloadFlag = false
	LC.ReloadValue = 0
}

func (LC *LengthCounter) ReloadLengthCounter() {
	if LC.ReloadFlag && LC.Counter == 0 {
		LC.Counter = LC.ReloadValue
	} else {
		LC.ReloadFlag = false
	}
}

func (LC *LengthCounter) ClockLengthCounter() {
	if LC.Counter != 0 && !LC.HaltFlag && !LC.ReloadFlag {
		LC.Counter--
	}
}

func (LC *LengthCounter) ClockLinearCounter(control bool) {
	if LC.ReloadFlag {
		LC.Counter = LC.ReloadValue
	} else if LC.Counter != 0 {
		LC.Counter--
	}

	if !control {
		LC.ReloadFlag = false
	}
}

func (Env *Envelope) ResetEnvelope() {
	Env.ConstantVolume = false
	Env.Volume = 0
	Env.StartFlag = false
	Env.ResetDivider()
	Env.Decay = 0
}

func (Env *Envelope) ClockEnvelope() {
	if !Env.StartFlag {
		if Env.ClockDivider() {
			Env.Divider.Counter = uint16(Env.Volume)
			if Env.Decay != 0 {
				Env.Decay--
			} else if Env.LoopFlag {
				Env.Decay = 15
			}
		}

	} else {
		Env.StartFlag = false
		Env.Decay = 15
		Env.Divider.Counter = Env.Divider.Period
	}
}

func (Sw *Sweep) ResetSweep() {
	Sw.Enabled = false
	Sw.ResetDivider()
	Sw.ReloadFlag = false
	Sw.Negate = false
	Sw.Shift = 0
}

func (Sw *Sweep) ClockSweep() {
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
}

func (Sw *Sweep) UpdateTargetPeriod(isPulse1 bool) {

	shift := int(Sw.Period >> Sw.Shift)

	if Sw.Negate {
		shift *= -1
		if isPulse1 {
			shift--
		}
	}

	newTarget := int(Sw.Period) + shift

	if newTarget < 0 {
		newTarget = 0
	}

	Sw.TargetPeriod = uint16(newTarget)

}

func (Dv *Divider) ResetDivider() {
	Dv.Counter = 0
	Dv.Period = 0
}

func (Dv *Divider) ClockDivider() bool {
	if Dv.Counter == 0 {
		Dv.Counter = Dv.Period
		return true
	} else {
		Dv.Counter--
		return false
	}

}

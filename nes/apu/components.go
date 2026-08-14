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

type LinearCounter struct {
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

func (LC *LinearCounter) ResetLinearCounter() {
	LC.Counter = 0
	LC.HaltFlag = false
	LC.ReloadFlag = false
	LC.ReloadValue = 0
}

/*
	func (LC *LengthCounter) ReloadLengthCounter() {
		if LC.ReloadFlag && LC.Counter == 0 {
			LC.Counter = LC.ReloadValue
		} else {
			LC.ReloadFlag = false
		}
	}
*/
func (LC *LengthCounter) ClockLengthCounter() {
	if LC.Counter > 0 && !LC.HaltFlag {
		LC.Counter--
	}
}

func (LC *LinearCounter) ClockLinearCounter(control bool) {
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

package apu

import "mtt/timenes/common"

type TriangleChannel struct {
	Enabled          bool
	Timer            uint16
	TimerReloadValue uint16
	Output           byte
	LinearCounter    LengthCounter
	SeqPos           byte
	LengthCounter
}

var apuTriangleSequences = [32]byte{
	15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
}

func (t *TriangleChannel) ResetTriangle() {
	t.Enabled = false
	t.LengthCounter.ResetLengthCounter()
	t.LinearCounter.ResetLengthCounter()
	t.Timer = 0
	t.TimerReloadValue = 0
	t.SeqPos = 0
}

func (t *TriangleChannel) ClockTriangleTimer() {
	if t.Timer == 0 {
		t.Timer = t.TimerReloadValue
		if t.LengthCounter.Counter != 0 && t.LinearCounter.Counter != 0 {
			t.SeqPos = (t.SeqPos + 1) & 31
		}
	} else {
		t.Timer--
	}
}

func (t *TriangleChannel) UpdateTriangleOutput() byte {
	if common.MuteEmulator || !t.Enabled || t.LengthCounter.Counter == 0 || t.LinearCounter.Counter == 0 {
		return 0
	}
	return apuTriangleSequences[t.SeqPos]
}

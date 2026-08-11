package apu

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

func (triangle *TriangleChannel) UpdateTriangleOutput() {
	triangle.Output = apuTriangleSequences[triangle.SeqPos]

	if !apuEnabled || triangle.LengthCounter.Counter == 0 || triangle.LinearCounter.Counter == 0 {
		triangle.Output = 0
	}
}

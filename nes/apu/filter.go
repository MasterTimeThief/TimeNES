package apu

import (
	"math"
)

// https://stackoverflow.com/questions/28291582/implementing-a-high-pass-filter-to-an-audio-signal
type Filter struct {
	resonance  float32
	frequency  float32
	SampleRate int
	IsHighPass bool

	Value                 float32
	c, a1, a2, a3, b1, b2 float32
	InputHistory          [2]float32
	OutputHistory         [3]float32
}

var HighFilter90, HighFilter440, LowFilter14k Filter

func InitFilters() {
	HighFilter90.InitFilter(90, apuSampleRate, true, 1)
	HighFilter440.InitFilter(440, apuSampleRate, true, 1)
	LowFilter14k.InitFilter(14000, apuSampleRate, false, 1)
}

func (f *Filter) InitFilter(freq float32, sampleRate int, isHighPass bool, resonance float32) {
	//Add High / Low Passes
	f.resonance = resonance
	f.frequency = freq
	f.SampleRate = sampleRate
	f.IsHighPass = isHighPass

	if f.IsHighPass {
		f.c = float32(math.Tan(float64(math.Pi * f.frequency / float32(f.SampleRate))))
		f.a1 = float32(1.0) / (float32(1.0) + (resonance * f.c) + (f.c * f.c))
		f.a2 = float32(-2.0) * f.a1
		f.a3 = f.a1
		f.b1 = float32(2.0) * ((f.c * f.c) - float32(1.0)) * f.a1
		f.b2 = (float32(1.0) - (f.resonance * f.c) + (f.c * f.c)) * f.a1
	} else {
		f.c = float32(1.0) / float32(math.Tan(float64(math.Pi*f.frequency/float32(f.SampleRate))))
		f.a1 = float32(1.0) / (float32(1.0) + (f.resonance * f.c) + (f.c * f.c))
		f.a2 = float32(2.0) * f.a1
		f.a3 = f.a1
		f.b1 = float32(2.0) * (float32(1.0) - (f.c * f.c)) * f.a1
		f.b2 = (float32(1.0) - (resonance * f.c) + (f.c * f.c)) * f.a1
	}

}

func (f *Filter) FilterUpdate(NewValue float32) float32 {
	newOutput := (f.a1 * float32(NewValue)) + (f.a2 * f.InputHistory[0]) + (f.a3 * f.InputHistory[1]) - (f.b1 * f.OutputHistory[0]) - (f.b2 * f.OutputHistory[1])

	f.InputHistory[1] = f.InputHistory[0]
	f.InputHistory[0] = NewValue

	f.OutputHistory[2] = f.OutputHistory[1]
	f.OutputHistory[1] = f.OutputHistory[0]
	f.OutputHistory[0] = newOutput

	return newOutput
}

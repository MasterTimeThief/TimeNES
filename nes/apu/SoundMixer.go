package apu

import (
	"math"
	"mtt/timenes/common"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// var apuMixerInputBuffer [apuMaxSamplesPerFrame]uint16
// var apuMixerOutputBuffer [apuMaxSampleRate]uint16
// var apuClockRate, apuSampleRate uint32
//var apuSample float32
//var apuBuffer bytes.Buffer

//var audioBuffer chan byte

type AudioBufferStruct struct {
	SampleRate float64
	ringBuffer *ringBuffer
	sample     float32
}

var audBuf *AudioBufferStruct

func (abs *AudioBufferStruct) Read(p []byte) (int, error) {

	n := abs.ringBuffer.Read(p)
	if n == 0 {
		clear(p)
		return len(p), nil
	}
	return n, nil
}

func InitAudioOutput() {
	audBuf = &AudioBufferStruct{
		ringBuffer: newRingBuffer(apuSampleRate),
		SampleRate: float64(apuMaxSamplesPerFrame),
	}
}

func NewAudioContext() *audio.Context {
	return audio.NewContext(apuSampleRate)
}

func NewAudioPlayer(context *audio.Context) *audio.Player {
	player, err := context.NewPlayerF32(audBuf)
	common.Check(err)

	player.SetBufferSize(time.Second / 60)
	player.SetVolume(0.15)
	go func() {
		player.Play()
	}()
	return player
}

const (
	apuCycleLength        int = 10000
	apuBitsPerCycle       int = 16
	apuSampleRate         int = 60000
	apuMaxSamplesPerFrame int = apuSampleRate / 60 * 4 * 2

	FrameCounterRate  = float64(common.CPU_Frequency) / 240.0
	DefaultSampleRate = float64(common.CPU_Frequency) / float64(apuSampleRate)
)

func (a *AudioBufferStruct) sendSample() {
	result := a.sample /*/ float32(a.SampleRate)*/
	a.sample = 0
	b := math.Float32bits(result)
	newSample := []byte{
		byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24),
		byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24),
	}
	a.ringBuffer.Write(newSample)
}

func AudioOutput() {
	//Pulse channels
	Pulse1.UpdatePulseOutput()
	Pulse2.UpdatePulseOutput()
	pulse_out := squareTable[Pulse1.Output+Pulse2.Output]

	Triangle.UpdateTriangleOutput()
	Noise.UpdateNoiseOutput()

	tnd_out := tndTable[(3*Triangle.Output)+(2*Noise.Output)+DMC.Output]

	audBuf.sample += pulse_out + tnd_out

	if apuDMAGetCycle {
		audBuf.sendSample()
	}
}

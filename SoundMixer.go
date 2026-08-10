package main

import (
	"math"
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

	if g.audioContext == nil {
		g.audioContext = audio.NewContext(apuSampleRate)
	}
	if g.player == nil {
		// Pass the (infinite) stream to NewPlayer.
		// After calling Play, the stream never ends as long as the player object lives.
		var err error
		g.player, err = g.audioContext.NewPlayerF32(audBuf)
		check(err)

		g.player.SetBufferSize(time.Second / 60)
		g.player.SetVolume(0.15)
		go func() {
			g.player.Play()
		}()
	}
}

const (
	apuCycleLength        int = 10000
	apuBitsPerCycle       int = 16
	apuSampleRate         int = 60000
	apuMaxSamplesPerFrame int = apuSampleRate / 60 * 4 * 2

	FrameCounterRate  = float64(CPU_Frequency) / 240.0
	DefaultSampleRate = float64(CPU_Frequency) / float64(apuSampleRate)
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
	apuPulse1.UpdatePulseOutput()
	apuPulse2.UpdatePulseOutput()
	pulse_out := squareTable[apuPulse1.Output+apuPulse2.Output]

	apuTriangle.UpdateTriangleOutput()
	apuNoise.UpdateNoiseOutput()

	tnd_out := tndTable[(3*apuTriangle.Output)+(2*apuNoise.Output)+apuDMC.Output]

	audBuf.sample += pulse_out + tnd_out

	if apuDMAGetCycle {
		audBuf.sendSample()
	}
}

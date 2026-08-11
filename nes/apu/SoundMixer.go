package apu

import (
	"math"
	"mtt/timenes/common"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	//apuCycleLength        int = 10000
	//apuBitsPerCycle       int = 16
	apuInputSamples       int = 14934
	apuRingBufferSize     int = apuInputSamples * 8
	apuSampleRate         int = 48000
	apuMaxSamplesPerFrame int = apuSampleRate / 60

	//FrameCounterRate  = float64(common.CPU_Frequency) / 240.0
	//DefaultSampleRate = (common.CPU_Frequency) / (apuSampleRate)
)

type AudioBufferStruct struct {
	//SampleRate  float64
	ringBuffer  *ringBuffer
	frameBuffer [][]byte
	sample      float32
}

var audBuf *AudioBufferStruct
var squareTable [31]float32
var tndTable [203]float32

func (abs *AudioBufferStruct) Read(p []byte) (int, error) {

	n := abs.ringBuffer.Read(p)
	if n == 0 {
		clear(p)
		return len(p), nil
	}
	return n, nil
}

func InitAudioOutput() {
	for i := range squareTable {
		squareTable[i] = float32(95.52 / (8128.0/float64(i) + 100))
	}
	for i := range tndTable {
		tndTable[i] = float32(163.67 / (24329.0/float64(i) + 100))
	}
	audBuf = &AudioBufferStruct{
		ringBuffer: newRingBuffer(apuRingBufferSize),
		//SampleRate: float64(apuMaxSamplesPerFrame),
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

func (a *AudioBufferStruct) sendSample() {
	result := a.sample /*/ float32(a.SampleRate)*/
	a.sample = 0
	b := math.Float32bits(result)
	newSample := []byte{
		byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24),
		byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24),
	}
	a.frameBuffer = append(a.frameBuffer, newSample)
	//a.ringBuffer.Write(newSample)
}

func AudioOutput() {
	//Pulse channels
	Pulse1.UpdatePulseOutput()
	Pulse2.UpdatePulseOutput()
	pulse_out := squareTable[Pulse1.Output+Pulse2.Output]

	Triangle.UpdateTriangleOutput()
	Noise.UpdateNoiseOutput()

	tnd_out := tndTable[(3*Triangle.Output)+(2*Noise.Output)+DMC.Output]

	if apuDMAGetCycle {
		audBuf.sample += pulse_out + tnd_out
		audBuf.sendSample()
	}
}

func TransferBuffer() {

	inputBufferSize := len(audBuf.frameBuffer)
	NearestNeighborRatio := float32(inputBufferSize) / float32(apuMaxSamplesPerFrame)
	NearestNeighborIndex := float32(0)

	//7466-7468
	//fmt.Printf("Buffer size: %d\n", inputBufferSize)

	for i := 0; i < apuMaxSamplesPerFrame; i++ {
		NearestNeighborIndex += NearestNeighborRatio

		newIndex := int(NearestNeighborIndex)

		if int(NearestNeighborIndex) < inputBufferSize {
			audBuf.ringBuffer.Write(audBuf.frameBuffer[newIndex])
		}
	}

	audBuf.frameBuffer = nil
}

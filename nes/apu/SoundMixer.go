package apu

import (
	"encoding/binary"
	"math"
	"mtt/timenes/common"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	//apuCycleLength        int = 10000
	//apuBitsPerCycle       int = 16
	apuInputSamples       int = 14935
	apuRingBufferSize     int = apuInputSamples * 8
	apuSampleRate         int = 48000
	apuMaxSamplesPerFrame int = apuSampleRate / 60

	//FrameCounterRate  = float64(common.CPU_Frequency) / 240.0
	//DefaultSampleRate = (common.CPU_Frequency) / (apuSampleRate)
)

type AudioBufferStruct struct {
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
		ringBuffer:  newRingBuffer(apuRingBufferSize),
		frameBuffer: make([][]byte, apuInputSamples),
	}
	InitFilters()
}

func NewAudioContext() *audio.Context {
	return audio.NewContext(apuSampleRate)
}

func NewAudioPlayer(context *audio.Context) *audio.Player {
	player, err := context.NewPlayerF32(audBuf)
	common.Check(err)

	player.SetBufferSize(time.Second / 60)
	player.SetVolume(0.5)
	go func() {
		player.Play()
	}()
	return player
}

func (a *AudioBufferStruct) sendSample() {
	result := a.sample /*/ float32(a.SampleRate)*/
	a.sample = 0
	newSample := BytesFromFloat32(result)
	a.frameBuffer = append(a.frameBuffer, newSample)
	//a.ringBuffer.Write(newSample)
}

func AudioOutput() {
	//Pulse channels
	pulse1 := Pulse1.UpdatePulseOutput()
	pulse2 := Pulse2.UpdatePulseOutput()
	pulse_out := squareTable[pulse1+pulse2]

	triangle := Triangle.UpdateTriangleOutput()
	noise := Noise.UpdateNoiseOutput()
	//dmc := byte(common.Ternary(DMC.Enabled, uint16(DMC.Output), 0))
	dmc := DMC.Output

	tnd_out := tndTable[(3*triangle)+(2*noise)+dmc]

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

			//Filter
			FilteredSample := audBuf.frameBuffer[newIndex]

			if len(audBuf.frameBuffer[newIndex]) >= 4 {
				InputSampleFloat := Float32FromBytes(audBuf.frameBuffer[newIndex])
				FilteredSample = BytesFromFloat32(LowFilter14k.FilterUpdate(HighFilter440.FilterUpdate(HighFilter90.FilterUpdate(InputSampleFloat))))
				//FilteredSample = BytesFromFloat32(HighFilter440.FilterUpdate(InputSampleFloat))
			}

			audBuf.ringBuffer.Write(FilteredSample)
		}
	}

	audBuf.frameBuffer = nil
}

func Float32FromBytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func BytesFromFloat32(float float32) []byte {
	b := math.Float32bits(float)
	return []byte{
		byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24),
		byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24),
	}
}

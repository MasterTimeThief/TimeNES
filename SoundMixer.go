package main

import (
	"bytes"
	"math"
	"math/rand/v2"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
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
	buffer     bytes.Buffer
	sample     float32
	streamer   beep.Streamer
}

var audBuf *AudioBufferStruct
var beepStreamer beep.Streamer
var audioQueue Queue

func (abs *AudioBufferStruct) Read(p []byte) (int, error) {
	/*const bytesPerSample = 8

	n := len(buf) / bytesPerSample * bytesPerSample

	//v := APU_Mixer()

	for i := 0; i < n; i++ {
		buf[i], _ = apuBuffer.ReadByte()
	}

	apuBuffer.Reset()

	//abs.pos %= length * bytesPerSample

	return n, nil*/

	/*cap := apuBuffer.Cap()
	length := apuBuffer.Len()
	if cap != 0 && length != 0 {
		fmt.Printf("Cap: %d\tLen: %d\tBuf: %d\t\n", cap, length, len(buf))
	}*/

	n := abs.ringBuffer.Read(p)
	if n == 0 {
		clear(p)
		return len(p), nil
	}
	return n, nil
}

type Queue struct {
	streamers []beep.Streamer
}

func (q *Queue) Add(streamers ...beep.Streamer) {
	q.streamers = append(q.streamers, streamers...)
}

func (q *Queue) Stream(samples [][2]float64) (n int, ok bool) {
	// We use the filled variable to track how many samples we've
	// successfully filled already. We loop until all samples are filled.
	filled := 0
	for filled < len(samples) {
		// There are no streamers in the queue, so we stream silence.
		if len(q.streamers) == 0 {
			for i := range samples[filled:] {
				samples[i][0] = 0
				samples[i][1] = 0
			}
			break
		}

		// We stream from the first streamer in the queue.
		n, ok := q.streamers[0].Stream(samples[filled:])
		// If it's drained, we pop it from the queue, thus continuing with
		// the next streamer.
		if !ok {
			q.streamers = q.streamers[1:]
		}
		// We update the number of filled samples.
		filled += n
	}
	return len(samples), true
}

func (q *Queue) Err() error {
	return nil
}

func InitAudioOutput() {
	speaker.Init(beep.SampleRate(int(apuSampleRate)), int(apuMaxSamplesPerFrame))

	speaker.Play(&audioQueue)

	audBuf = &AudioBufferStruct{
		ringBuffer: newRingBuffer(40960),
		SampleRate: 44100,
	}
	//audioBuffer = make(chan byte, apuSampleRate)
	/*if g.audioContext == nil {
		g.audioContext = audio.NewContext(apuSampleRate)
	}
	if g.player == nil {
		// Pass the (infinite) stream to NewPlayer.
		// After calling Play, the stream never ends as long as the player object lives.
		var err error
		g.player, err = g.audioContext.NewPlayerF32(audBuf)
		check(err)

		g.player.SetBufferSize(time.Second / 20)
		//g.player.SetVolume(0.5)
		go func() {
			g.player.Play()
		}()
	}*/
}

const (
	apuCycleLength        int = 10000
	apuBitsPerCycle       int = 16
	apuSampleRate         int = 44100
	apuMaxSamplesPerFrame int = apuSampleRate / 60 * 4 * 2
)

func InitMixer() {

}

func (a *AudioBufferStruct) sendSample() {
	result := a.sample /*/ float32(a.SampleRate)*/
	a.sample = 0
	b := math.Float32bits(result)
	newSample := []byte{
		byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24),
		byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24),
	}
	a.ringBuffer.Write(newSample)

	newStream := beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {

		random := rand.Float64()*2 - 1
		samples[0][0] = random
		samples[0][1] = random

		return len(samples), true
	})

	speaker.Lock()
	audioQueue.Add(newStream)
	speaker.Unlock()
}

func AudioOutput() {
	//Pulse channels
	apuPulse1.UpdatePulseOutput()
	apuPulse2.UpdatePulseOutput()
	pulse_out := squareTable[apuPulse1.Output+apuPulse2.Output]

	if apuPulse1.Output != 0 {
		print("")
	}

	tnd_out := tndTable[0] //For now...

	audBuf.sample += pulse_out + tnd_out

	if apuDMAGetCycle {
		audBuf.sendSample()
	}
}

func APU_Mixer() uint32 {

	/*b := math.Float32bits(result)

	audBuf.buf.Write([]byte{
		byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24),
		byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24),
	})*/

	//io.Copy(bytes.NewBuffer(newBuffer), &apuBuffer)
	//for i := range newBuffer {
	//apuBuffer.Write(newBuffer)
	//audioBuffer = append(audioBuffer, newBuffer[i])
	//}
	return 0
	/*sr := beep.SampleRate(apuMaxSampleRate)
	var f float64
	if pulse.Timer < 8 {
		f = 0
	} else {
		f = (1789773 / (16 * (float64(pulse.Timer+1) + 1)))
	}
	square, err := generators.SquareTone(sr, f)
	if err != nil {
		panic(err)
	}

	volume := &effects.Volume{
		Streamer: square,
		Base:     2,
		Volume:   0,
		Silent:   pulse.IsPulseMuted(),
	}
	speaker.Play(volume)*/
}

func OutputSound() {

}

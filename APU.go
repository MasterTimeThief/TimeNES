package main

type Pulse1Channel struct {
	Duty              byte
	LengthCounterHalt bool
	ConstantVolume    bool
	Envelope          byte
	SweepUnitEnabled  bool
	Period            byte
	Negate            bool
	ShiftRegister     byte
	Timer             uint16
	LengthCounter     byte
	Enabled           bool
}

type Pulse2Channel struct {
	Duty              byte
	LengthCounterHalt bool
	ConstantVolume    bool
	Envelope          byte
	SweepUnitEnabled  bool
	Period            byte
	Negate            bool
	ShiftRegister     byte
	Timer             uint16
	LengthCounter     byte
	Enabled           bool
}

type TriangleChannel struct {
	Control             bool
	LinearCounterReload bool
	LinearCounter       byte
	Timer               uint16
	LengthCounter       byte
	Enabled             bool
}

type NoiseChannel struct {
	LengthCounterHalt bool
	ConstantVolume    bool
	Envelope          byte
	Mode              bool
	Period            byte
	LengthCounter     byte
	Enabled           bool
}
type DeltaModChannel struct {
	IRQEnable   bool
	Loop        bool
	Frequency   byte
	LoadCounter byte
	Address     byte
	Length      int
}

var apuPulse1 Pulse1Channel
var apuPulse2 Pulse2Channel
var apuTriangle TriangleChannel
var apuNoise NoiseChannel
var apuDMC DeltaModChannel

var apuFrameCounterCycles = 0
var apuEnablePulse1, apuEnablePulse2, apuEnableTriangle, apuEnableNoise, apuEnableDMC bool
var apuDMCInterrupt, apuDMCDelayed, apuFrameInterrupt, apuDMAGetCycle, apuInhibitIRQ, apuFrameCounterMode bool
var apuFrameCounter int

var apu4017ResetTimer int = 0

func Emulate_APU(g *Game) {
	//ClockFrameCounter()

}

func DMA_Get() {

}

func DMA_Put() {

}

func LengthCounterLoad(Value byte) byte {

	//Short n' Easy
	//LengthLookupTable := [32]int{10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14, 12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30}
	//return byte(LengthLookupTable[Value])

	//Flip it around?
	/*newValue := byte(0)

	newValue |= ((Value & 1) << 4)
	newValue |= (Value >> 1)*/

	length := byte(0)
	if (Value & 1) == 1 { //Odd values are linear lengths
		if Value == 1 {
			length = 0xFE
		} else {
			length = Value - 1
		}
	} else { //Even values are note lengths
		switch Value {
		//Notes with base length 12 (4/4 at 75 bpm):
		case 0x1E:
			length = 32 // (96 times 1/3, quarter note triplet)
		case 0x1C:
			length = 16 // (48 times 1/3, eighth note triplet)
		case 0x1A:
			length = 72 // (48 times 1 1/2, dotted quarter)
		case 0x18:
			length = 192 // (Whole Note)
		case 0x16:
			length = 96 // (Half Note)
		case 0x14:
			length = 48 // (Quarter Note)
		case 0x12:
			length = 24 // (Eight Note)
		case 0x10:
			length = 12 // (Sixteenth)

		//Notes with base length 10 (4/4 at 90 bpm, with relative durations being the same as above):
		case 0x0E:
			length = 26 // (Approx. 80 times 1/3, quarter note triplet)
		case 0x0C:
			length = 14 // (Approx. 40 times 1/3, eighth note triplet)
		case 0x0A:
			length = 60 // (40 times 1 1/2, dotted quarter)
		case 0x08:
			length = 160 // (Whole Note)
		case 0x06:
			length = 80 // (Half Note)
		case 0x04:
			length = 40 // (Quarter Note)
		case 0x02:
			length = 20 // (Eight Note)
		case 0x00:
			length = 10 // (Sixteenth)
		}
	}
	return length
}

func ClockFrameCounter() { //Also called Frame Sequencer
	//Do frame counter shit
	apuFrameCounterCycles++

	if !apuFrameCounterMode { //4-Cycle mode
		switch apuFrameCounterCycles {
		case 3728:
			ClockFrameCounterQuarterFrame()
		case 7456:
			ClockFrameCounterQuarterFrame()
			ClockFrameCounterHalfFrame()
		case 11185:
			ClockFrameCounterQuarterFrame()
		case 14914:
			if !apuDMAGetCycle {
				ClockFrameCounterQuarterFrame()
				ClockFrameCounterHalfFrame()
			}
			apuFrameInterrupt = !apuInhibitIRQ
		case 14915:
			apuFrameCounterCycles = 0
			apuFrameInterrupt = !apuInhibitIRQ
		}
	} else { //5-Cycle mode
		switch apuFrameCounterCycles {
		case 3728:
			ClockFrameCounterQuarterFrame()
		case 7456:
			ClockFrameCounterQuarterFrame()
			ClockFrameCounterHalfFrame()
		case 11185:
			ClockFrameCounterQuarterFrame()
		case 14914:
			//Nothing?
		case 18640:
			ClockFrameCounterQuarterFrame()
			ClockFrameCounterHalfFrame()
		case 18641:
			apuFrameCounterCycles = 0
		}
	}

}

func ClockFrameCounterQuarterFrame() {
	apuPulse1.Envelope--
	apuPulse2.Envelope--
	apuNoise.Envelope--
	apuTriangle.LinearCounter--

}

func ClockFrameCounterHalfFrame() {
	if !apuPulse1.LengthCounterHalt && apuPulse1.LengthCounter != 0 {
		apuPulse1.LengthCounter--
	}
	if !apuPulse2.LengthCounterHalt && apuPulse2.LengthCounter != 0 {
		apuPulse2.LengthCounter--
	}
	if !apuNoise.LengthCounterHalt && apuNoise.LengthCounter != 0 {
		apuNoise.LengthCounter--
	}
	if !apuTriangle.Control && apuTriangle.LengthCounter != 0 {
		apuTriangle.LengthCounter--
	}

	//Sweep???

}

/*func soundLoop() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}

	f, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		usage()
		os.Exit(0)
	}

	sr := beep.SampleRate(48000)
	speaker.Init(sr, 4800)

	sine, err := generators.SineTone(sr, f)
	if err != nil {
		panic(err)
	}

	triangle, err := generators.TriangleTone(sr, f)
	if err != nil {
		panic(err)
	}

	square, err := generators.SquareTone(sr, f)
	if err != nil {
		panic(err)
	}

	sawtooth, err := generators.SawtoothTone(sr, f)
	if err != nil {
		panic(err)
	}

	sawtoothReversed, err := generators.SawtoothToneReversed(sr, f)
	if err != nil {
		panic(err)
	}

	// Play 2 seconds of each tone
	two := sr.N(2 * time.Second)

	ch := make(chan struct{})
	sounds := []beep.Streamer{
		beep.Callback(print("sine")),
		beep.Take(two, sine),
		beep.Callback(print("triangle")),
		beep.Take(two, triangle),
		beep.Callback(print("square")),
		beep.Take(two, square),
		beep.Callback(print("sawtooth")),
		beep.Take(two, sawtooth),
		beep.Callback(print("sawtooth reversed")),
		beep.Take(two, sawtoothReversed),
		beep.Callback(func() {
			ch <- struct{}{}
		}),
	}
	speaker.Play(beep.Seq(sounds...))

	<-ch
}


func usage() {
	fmt.Printf("usage: %s freq\n", os.Args[0])
	fmt.Println("where freq must be a float between 1 and 24000")
	fmt.Println("24000 because samplerate of 48000 is hardcoded")
}

// a simple clousure to wrap fmt.Println
func print(s string) func() {
	return func() {
		fmt.Println(s)
	}
}

*/

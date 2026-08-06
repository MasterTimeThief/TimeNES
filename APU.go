package main

type Envelope struct {
	ConstantVolume bool
	Volume         byte
	StartFlag      bool
	Divider        Divider
	Counter        byte
}

type Sweep struct {
	Enabled    bool
	Divider    Divider
	ReloadFlag bool
}

type Divider struct {
	Period  byte //P
	Counter byte
}

type LengthCounter struct {
	Counter     byte
	HaltFlag    bool
	ReloadFlag  bool
	ReloadValue byte
}

type PulseChannel struct {
	Enabled                  bool
	LengthCounter            byte
	LengthCounterHalt        bool
	LengthCounterReload      bool
	LengthCounterReloadValue byte
	Timer                    uint16
	Envelope                 Envelope
	//Channel Specific Variables
	Duty byte
	//ConstantVolume   bool
	SweepUnitEnabled bool
	Period           byte
	Negate           bool
	ShiftRegister    byte
}

type TriangleChannel struct {
	Enabled                  bool
	LengthCounter            byte
	LengthCounterHalt        bool
	LengthCounterReload      bool
	LengthCounterReloadValue byte
	Timer                    uint16
	//Channel Specific Variables
	LinearCounterReload bool
	LinearCounter       byte
}

type NoiseChannel struct {
	Enabled                  bool
	LengthCounter            byte
	LengthCounterHalt        bool
	LengthCounterReload      bool
	LengthCounterReloadValue byte
	Timer                    uint16
	Envelope                 Envelope
	//Channel Specific Variables
	//ConstantVolume bool
	Mode   bool
	Period byte
}
type DeltaModChannel struct {
	Enabled              bool
	Timer                uint16
	IRQEnable            bool
	Loop                 bool
	SampleRate           uint16 // AKA Frequency
	Output               byte   //AKA LoadCounter
	SampleAddress        uint16
	SampleLength         uint16
	BytesRemaining       uint16
	Buffer               byte
	SampleAddressCounter uint16
	Shifter              byte
	ShifterBitsRemaining byte
	DPCM_Up              bool
	//Channel Specific Variables
	//Frequency   byte
	//LoadCounter byte
}

var apuDutySequences = [4][8]int{
	{0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 1, 1},
	{0, 0, 0, 0, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 0, 0},
}

var apuPulse1 PulseChannel
var apuPulse2 PulseChannel
var apuTriangle TriangleChannel
var apuNoise NoiseChannel
var apuDMC DeltaModChannel

// var apuEnablePulse1, apuEnablePulse2, apuEnableTriangle, apuEnableNoise, apuEnableDMC bool
var apuDMCInterrupt, apuDMCDelayed, apuFrameInterrupt, apuDMAGetCycle, apuInhibitIRQ, apuFrameCounterMode, apuSilent, apuIsHalfFrame bool
var apuFrameCounter int = 0
var IRQLevelDetector, DoIRQ bool

var apuEnvelopeStartFlag bool
var apuEnvelopeDivider bool
var apuEnvelopeDecay byte

var apuDoDMCDMA, apuDMCDMAHalt bool
var apuDMCDMADelay, apuCannotDMCDMARightNow byte

var apu4017ResetTimer int = 0
var apuLengthCounterLUT = [32]byte{10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14, 12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30}
var apuDMCSampleRateLUT = [16]uint16{428, 380, 340, 320, 286, 254, 226, 214, 190, 160, 142, 128, 106, 84, 72, 54}

func Emulate_APU(g *Game) {
	//ClockFrameCounter()
	if apuDMAGetCycle { //DMA Get Cycle
		DMA_Get()
	} else { //DMA Put Cycle
		DMA_Put()
	}
	/*
		if (APU_DelayedDMC4015 > 0)
		{
			APU_DelayedDMC4015--;
			if (APU_DelayedDMC4015 == 0)
			{
				APU_Status_DMC = APU_Status_DelayedDMC;
				if (!APU_Status_DMC)
				{
					APU_DMC_BytesRemaining = 0;
				}
			}
		}
	*/
	apuTriangle.Timer-- //Every CPU Cycle

	//Clock sequencer
	if (apu4017ResetTimer & 0x80) == 0 {
		apu4017ResetTimer--
		if (apu4017ResetTimer & 0x80) != 0 {
			apuFrameCounter = 0
		}
	}

	ClockFrameCounter()

	//If this isn't a Frame Counter half-frame
	if !apuIsHalfFrame {
		if apuPulse1.LengthCounterReload {
			apuPulse1.LengthCounter = apuPulse1.LengthCounterReloadValue
		}
		if apuPulse2.LengthCounterReload {
			apuPulse2.LengthCounter = apuPulse2.LengthCounterReloadValue
		}
		if apuTriangle.LengthCounterReload {
			apuTriangle.LengthCounter = apuTriangle.LengthCounterReloadValue
		}
		if apuNoise.LengthCounterReload {
			apuNoise.LengthCounter = apuNoise.LengthCounterReloadValue
		}
		apuPulse1.LengthCounterReload = false
		apuPulse2.LengthCounterReload = false
		apuTriangle.LengthCounterReload = false
		apuNoise.LengthCounterReload = false
	}
}

func DMA_Get() {
	//Clock timers
	apuPulse1.Timer--
	apuPulse2.Timer--
	apuNoise.Timer--

	apuDMC.Timer--
	apuDMC.Timer-- // the table is in CPU cycles, but the count is in APU cycles

	if apuDMC.Timer == 0 {
		apuDMC.Timer = apuDMC.SampleRate
		apuDMC.DPCM_Up = (apuDMC.Shifter & 1) == 1
		if apuDMC.DPCM_Up {
			if apuDMC.Output <= 125 { // this is 7 bit, and cannot go above 127

				apuDMC.Output += 2
			}
		} else {
			if apuDMC.Output >= 2 { // this is 7 bit, and cannot go below 0

				apuDMC.Output -= 2
			}
		}
		apuDMC.Shifter >>= 1                  // shift the bits in the shift register
		apuDMC.ShifterBitsRemaining--         // and decrement the "bits remaining" counter.
		if apuDMC.ShifterBitsRemaining == 0 { // If there are no bits left,

			apuDMC.ShifterBitsRemaining = 8 // it's time for a DMC DMA!

			if apuDMC.BytesRemaining > 0 /*|| APU_SetImplicitAbortDMC4015*/ {
				if !apuDoDMCDMA && apuCannotDMCDMARightNow != 2 {
					// if playing a sample:
					apuDoDMCDMA = true
					apuDMCDMAHalt = true
				}
				//if APU_SetImplicitAbortDMC4015 {
				//	APU_ImplicitAbortDMC4015 = true // check for weird DMA abort behavior
				//	APU_SetImplicitAbortDMC4015 = false
				//}
				apuDMC.Shifter = apuDMC.Buffer // and set up the shifter with the new values.
				apuSilent = false              // The APU is not silent.

			} else {
				apuSilent = true
			}
		}
	}
	if apuCannotDMCDMARightNow > 0 {
		apuCannotDMCDMARightNow -= 2
	}
}

func DMA_Put() {
	/*if Clearing_APU_FrameInterrupt {
		Clearing_APU_FrameInterrupt = false
		APU_Status_FrameInterrupt = false
		IRQ_LevelDetector = false
	}*/
	// DMC load from 4015
	if apuDMCDMADelay > 0 {
		apuDMCDMADelay--                         // there's a small delay beetween the write occurring and the DMA beginning
		if apuDMCDMADelay == 0 && !apuDoDMCDMA { // if the DMA is already happening because of the timer

			apuDoDMCDMA = true
			apuDMCDMAHalt = true
			apuDMC.Shifter = apuDMC.Buffer
			apuSilent = false
		}
	}
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
	apuFrameCounter++
	apuIsHalfFrame = false

	if !apuFrameCounterMode { //4-Cycle mode
		switch apuFrameCounter {
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
			apuFrameInterrupt = true
			if !IRQLevelDetector {
				IRQLevelDetector = !apuInhibitIRQ
			}
		case 14915:
			apuFrameInterrupt = !apuInhibitIRQ
			if !IRQLevelDetector {
				IRQLevelDetector = !apuInhibitIRQ
			}
			apuFrameCounter = 0
		}
	} else { //5-Cycle mode
		switch apuFrameCounter {
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
			apuFrameCounter = 0
		}
	}

}

func ClockFrameCounterQuarterFrame() {

	if apuEnvelopeStartFlag {
		apuEnvelopeStartFlag = false
		apuEnvelopeDecay = 15
	} else {
		apuEnvelopeDivider = true
	}

	apuPulse1.Envelope.Volume--
	apuPulse2.Envelope.Volume--
	apuNoise.Envelope.Volume--
	apuTriangle.LinearCounter--

}

func ClockFrameCounterHalfFrame() {

	if apuPulse1.LengthCounterReload && apuPulse1.LengthCounter == 0 {
		apuPulse1.LengthCounter = apuPulse1.LengthCounterReloadValue
	} else {
		apuPulse1.LengthCounterReload = false
	}

	if apuPulse2.LengthCounterReload && apuPulse2.LengthCounter == 0 {
		apuPulse2.LengthCounter = apuPulse2.LengthCounterReloadValue
	} else {
		apuPulse2.LengthCounterReload = false
	}

	if apuTriangle.LengthCounterReload && apuTriangle.LengthCounter == 0 {
		apuTriangle.LengthCounter = apuTriangle.LengthCounterReloadValue
	} else {
		apuTriangle.LengthCounterReload = false
	}

	if apuNoise.LengthCounterReload && apuNoise.LengthCounter == 0 {
		apuNoise.LengthCounter = apuNoise.LengthCounterReloadValue
	} else {
		apuNoise.LengthCounterReload = false
	}

	// length counters and sweep
	if !apuPulse1.Enabled {
		apuPulse1.LengthCounter = 0
	}
	if !apuPulse2.Enabled {
		apuPulse2.LengthCounter = 0
	}
	if !apuTriangle.Enabled {
		apuTriangle.LengthCounter = 0
	}
	if !apuNoise.Enabled {
		apuNoise.LengthCounter = 0
	}

	if apuPulse1.LengthCounter != 0 && !apuPulse1.LengthCounterHalt && !apuPulse1.LengthCounterReload {
		apuPulse1.LengthCounter--
	}
	if apuPulse2.LengthCounter != 0 && !apuPulse2.LengthCounterHalt && !apuPulse2.LengthCounterReload {
		apuPulse2.LengthCounter--
	}
	if apuTriangle.LengthCounter != 0 && !apuTriangle.LengthCounterHalt && !apuTriangle.LengthCounterReload {
		apuTriangle.LengthCounter--
	}
	if apuNoise.LengthCounter != 0 && !apuNoise.LengthCounterHalt && !apuNoise.LengthCounterReload {
		apuNoise.LengthCounter--
	}

	apuIsHalfFrame = true

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

package main

type Envelope struct {
	ConstantVolume bool
	Volume         byte
	StartFlag      bool
	LoopFlag       bool
	Decay          byte
	Divider
}

type Sweep struct {
	Enabled bool
	Divider
	ReloadFlag   bool
	Negate       bool
	Shift        byte
	TargetPeriod uint16
}

type Divider struct {
	Period  uint16
	Counter uint16
}

type LengthCounter struct { //TODO: Maybe part of Envelope?
	Counter     byte
	HaltFlag    bool
	ReloadFlag  bool
	ReloadValue byte
}

type LinearCounter struct { //TODO: Maybe part of Envelope?
	Counter     byte
	HaltFlag    bool
	ReloadFlag  bool
	ReloadValue byte
}

type PulseChannel struct {
	Enabled bool
	LengthCounter
	Timer            uint16
	TimerReloadValue uint16
	Envelope
	//Channel Specific Variables
	Duty    byte
	DutyPos byte
	//SweepEnabled  bool
	//SweepPeriod   byte
	//SweepNegate   bool
	Sweep
	//ShiftRegister byte
	Output byte
}

type TriangleChannel struct {
	Enabled bool
	LengthCounter
	Timer            uint16
	TimerReloadValue uint16
	//Channel Specific Variables
	LinearCounter
	SeqPos byte
	Output byte
	//LinearCounterReload bool
	//LinearCounter       byte
}

type NoiseChannel struct {
	Enabled bool
	LengthCounter
	Timer            uint16
	TimerReloadValue uint16
	Envelope
	//Channel Specific Variables
	Mode          bool
	ShiftRegister uint16
	Output        byte
}
type DeltaModChannel struct {
	Enabled              bool
	Timer                uint16
	IRQEnable            bool
	Loop                 bool
	SampleRate           uint16 // AKA Frequency
	Output               byte   // AKA LoadCounter
	SampleAddress        uint16
	SampleLength         uint16
	BytesRemaining       uint16
	Buffer               byte
	SampleAddressCounter uint16
	Shifter              byte
	ShifterBitsRemaining byte
	DPCM_Up              bool
}

var apuDutySequences = [4][8]byte{
	{0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 1, 1},
	{0, 0, 0, 0, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 0, 0},
}

var apuTriangleSequences = [32]byte{
	15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
}

var apuPulse1 PulseChannel
var apuPulse2 PulseChannel
var apuTriangle TriangleChannel
var apuNoise NoiseChannel
var apuDMC DeltaModChannel

// var apuEnablePulse1, apuEnablePulse2, apuEnableTriangle, apuEnableNoise, apuEnableDMC bool
var apuDMCInterrupt, apuDMCDelayed, apuFrameInterrupt, apuInhibitIRQ, apuFrameCounterMode, apuIsHalfFrame bool
var apuFrameCounter int = 0
var IRQLevelDetector, DoIRQ bool

var apuEnvelopeStartFlag bool
var apuEnvelopeDivider bool
var apuEnvelopeDecay byte

var apuDMAGetCycle, apuDoDMCDMA, apuDMCDMAHalt bool
var apuDMCDMADelay, apuCannotDMCDMARightNow byte

var apu4017ResetTimer int = 0
var apuLengthCounterLUT = [32]byte{10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14, 12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30}
var apuDMCSampleRateLUT = [16]uint16{428, 380, 340, 320, 286, 254, 226, 214, 190, 160, 142, 128, 106, 84, 72, 54}
var apuNoiseTimerLUT = [16]uint16{4, 8, 16, 32, 64, 96, 128, 160, 202, 254, 380, 508, 762, 1016, 2034, 4068}

var squareTable [31]float32
var tndTable [203]float32

var apuSilent bool = true

func ResetAPU() {

	//Reset Pulse 1
	apuPulse1.Enabled = false
	apuPulse1.ResetLengthCounter()
	apuPulse1.Timer = 0

	apuPulse1.ResetEnvelope()

	apuPulse1.Duty = 0
	apuPulse1.DutyPos = 0
	apuPulse1.ResetSweep()
	//apuPulse1.ShiftRegister = 0

	//Reset Pulse 2
	apuPulse2.Enabled = false
	apuPulse2.ResetLengthCounter()
	apuPulse2.Timer = 0

	apuPulse2.ResetEnvelope()

	apuPulse2.Duty = 0
	apuPulse2.DutyPos = 0
	apuPulse2.ResetSweep()
	//apuPulse2.ShiftRegister = 0

	//Reset Triangle
	apuTriangle.Enabled = false
	apuTriangle.ResetLengthCounter()
	apuTriangle.Timer = 0
	apuTriangle.TimerReloadValue = 0
	apuTriangle.SeqPos = 0

	apuTriangle.LinearCounter.ReloadFlag = false
	apuTriangle.LinearCounter.Counter = 0
	apuTriangle.LinearCounter.ReloadValue = 0
	apuTriangle.LinearCounter.HaltFlag = false

	//Reset Noise
	apuNoise.Enabled = false
	apuNoise.ResetLengthCounter()
	apuNoise.Timer = 0

	apuNoise.ResetEnvelope()

	apuNoise.Mode = false
	apuNoise.Output = 0
	apuNoise.ShiftRegister = 1

	//Reset DMC
	apuDMC.Enabled = false
	apuDMC.Timer = 0
	apuDMC.IRQEnable = false
	apuDMC.Loop = false
	apuDMC.SampleRate = 0
	apuDMC.Output = 0
	apuDMC.SampleAddress = 0
	apuDMC.SampleLength = 0
	apuDMC.BytesRemaining = 0
	apuDMC.Buffer = 0
	apuDMC.SampleAddressCounter = 0
	apuDMC.Shifter = 0
	apuDMC.ShifterBitsRemaining = 0
	apuDMC.DPCM_Up = false

	//APU Variables
	apuDMCInterrupt = false
	apuDMCDelayed = false
	apuFrameInterrupt = false
	apuInhibitIRQ = false
	apuFrameCounterMode = false
	//apuSilent = false
	apuIsHalfFrame = false
	apuFrameCounter = 0
	IRQLevelDetector = false
	DoIRQ = false

	apuEnvelopeStartFlag = false
	apuEnvelopeDivider = false
	apuEnvelopeDecay = 0

	apuDMAGetCycle = true
	apuDoDMCDMA, apuDMCDMAHalt = false, false
	apuDMCDMADelay, apuCannotDMCDMARightNow = 0, 0

	apu4017ResetTimer = 0
	//apuMixerInputBuffer = [apuMaxSamplesPerFrame]uint16{}

}

func InitAPU() {
	for i := range squareTable {
		squareTable[i] = float32(95.52 / (8128.0/float64(i) + 100))
	}
	for i := range tndTable {
		tndTable[i] = float32(163.67 / (24329.0/float64(i) + 100))
	}

}

func Emulate_APU(g *Game) {
	AudioOutput()

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

	//Clock triangle timer every cycle
	if apuTriangle.Timer == 0 {
		apuTriangle.SeqPos = (apuTriangle.SeqPos + 1) & 31
		apuTriangle.Timer = apuTriangle.TimerReloadValue
	} else {
		apuTriangle.Timer--
	}

	//Clock sequencer
	if (apu4017ResetTimer & 0x80) == 0 {
		apu4017ResetTimer--
		if (apu4017ResetTimer & 0x80) != 0 {
			apuFrameCounter = 0
		}
	}

	ClockFrameCounter()

	//If this isn't a Frame Counter half-frame
	/*if !apuIsHalfFrame {
		if apuPulse1.LengthCounter.ReloadFlag {
			apuPulse1.LengthCounter.Counter = apuPulse1.LengthCounter.ReloadValue
		}
		if apuPulse2.LengthCounter.ReloadFlag {
			apuPulse2.LengthCounter.Counter = apuPulse2.LengthCounter.ReloadValue
		}
		if apuTriangle.LengthCounter.ReloadFlag {
			apuTriangle.LengthCounter.Counter = apuTriangle.LengthCounter.ReloadValue
		}
		if apuNoise.LengthCounter.ReloadFlag {
			apuNoise.LengthCounter.Counter = apuNoise.LengthCounter.ReloadValue
		}
		apuPulse1.LengthCounter.ReloadFlag = false
		apuPulse2.LengthCounter.ReloadFlag = false
		apuTriangle.LengthCounter.ReloadFlag = false
		apuNoise.LengthCounter.ReloadFlag = false
	}*/

	apuDMAGetCycle = !apuDMAGetCycle
}

func DMA_Get() {
	//Clock timers
	if apuPulse1.Timer == 0 {
		apuPulse1.DutyPos = (apuPulse1.DutyPos + 1) & 7
		apuPulse1.Timer = apuPulse1.TimerReloadValue
	} else {
		apuPulse1.Timer--
	}

	if apuPulse2.Timer == 0 {
		apuPulse2.DutyPos = (apuPulse2.DutyPos + 1) & 7
		apuPulse2.Timer = apuPulse2.TimerReloadValue
	} else {
		apuPulse2.Timer--
	}

	if apuNoise.Timer == 0 {
		apuNoise.ClockShiftRegister()
		apuNoise.Timer = apuNoise.TimerReloadValue
	} else {
		apuNoise.Timer--
	}

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
				//apuSilent = false              // The APU is not silent.

			} else {
				//apuSilent = true
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

	/*if apuPulse1.Envelope.Volume != 0 {
		apuPulse1.Envelope.Volume--
	}
	if apuPulse2.Envelope.Volume != 0 {
		apuPulse2.Envelope.Volume--
	}*/
	apuPulse1.ClockEnvelope()
	apuPulse2.ClockEnvelope()
	apuNoise.ClockEnvelope()

	apuTriangle.ClockLinearCounter(apuTriangle.LengthCounter.HaltFlag)

}

func ClockFrameCounterHalfFrame() {

	//TODO: Split off Clocks into seperate functions for readability
	apuPulse1.ReloadLengthCounter()
	apuPulse2.ReloadLengthCounter()
	apuTriangle.ReloadLengthCounter()
	apuNoise.ReloadLengthCounter()

	// length counters and sweep
	if !apuPulse1.Enabled {
		apuPulse1.LengthCounter.Counter = 0
	}
	if !apuPulse2.Enabled {
		apuPulse2.LengthCounter.Counter = 0
	}
	if !apuTriangle.Enabled {
		apuTriangle.LengthCounter.Counter = 0
	}
	if !apuNoise.Enabled {
		apuNoise.LengthCounter.Counter = 0
	}

	apuPulse1.ClockSweep()
	apuPulse2.ClockSweep()

	apuPulse1.ClockLengthCounter()
	apuPulse2.ClockLengthCounter()
	apuTriangle.ClockLengthCounter()
	apuNoise.ClockLengthCounter()

	apuIsHalfFrame = true

}

func (LC *LengthCounter) ResetLengthCounter() {
	LC.Counter = 0
	LC.HaltFlag = false
	LC.ReloadFlag = false
	LC.ReloadValue = 0
}

func (LC *LengthCounter) ReloadLengthCounter() {
	if LC.ReloadFlag && LC.Counter == 0 {
		LC.Counter = LC.ReloadValue
	} else {
		LC.ReloadFlag = false
	}
}

func (LC *LengthCounter) ClockLengthCounter() {
	if LC.Counter != 0 && !LC.HaltFlag && !LC.ReloadFlag {
		LC.Counter--
	}
}

func (LC *LinearCounter) ClockLinearCounter(control bool) {
	if LC.ReloadFlag {
		LC.Counter = LC.ReloadValue
	} else if LC.Counter != 0 {
		LC.Counter--
	}

	if !control {
		LC.ReloadFlag = false
	}
}

func (Env *Envelope) ResetEnvelope() {
	Env.ConstantVolume = false
	Env.Volume = 0
	Env.StartFlag = false
	Env.ResetDivider()
	Env.Decay = 0
}

func (Env *Envelope) ClockEnvelope() {
	if !Env.StartFlag {
		if Env.ClockDivider() {
			Env.Divider.Counter = uint16(Env.Volume)
			if Env.Decay != 0 {
				Env.Decay--
			} else if Env.LoopFlag {
				Env.Decay = 15
			}
		}

	} else {
		Env.StartFlag = false
		Env.Decay = 15
		Env.Divider.Counter = Env.Divider.Period
	}
}

func (Sw *Sweep) ResetSweep() {
	Sw.Enabled = false
	Sw.ResetDivider()
	Sw.ReloadFlag = false
	Sw.Negate = false
	Sw.Shift = 0
}

func (Sw *Sweep) ClockSweep() {
	if Sw.Counter == 0 {
		//Sweep is enabled, the shift count is nonzero
		if Sw.Enabled {
			//Pulse's period is set to target period
			Sw.Period = Sw.TargetPeriod
		} else {
			//Pulse's period remains unchanged, but the sweep unit's divider continues to count down  and reload the divider's period as normal
		}
	}

	if Sw.Counter == 0 || Sw.ReloadFlag {
		//Divider counter is set to P and the reload flag is cleared
		Sw.Counter = Sw.Period
		Sw.ReloadFlag = false
	} else {
		//Otherwise, the divider's counter is decremented
		Sw.Counter--
	}
}

func (Sw *Sweep) UpdateTargetPeriod(isPulse1 bool) {

	shift := int(Sw.Period >> Sw.Shift)

	if Sw.Negate {
		shift *= -1
		if isPulse1 {
			shift--
		}
	}

	newTarget := int(Sw.Period) + shift

	if newTarget < 0 {
		newTarget = 0
	}

	Sw.TargetPeriod = uint16(newTarget)

}

func (Dv *Divider) ResetDivider() {
	Dv.Counter = 0
	Dv.Period = 0
}

func (Dv *Divider) ClockDivider() bool {
	if Dv.Counter == 0 {
		Dv.Counter = Dv.Period
		return true
	} else {
		Dv.Counter--
		return false
	}

}

func (pulse *PulseChannel) IsPulseMuted() bool {
	return pulse.Timer < 8 || (!pulse.Sweep.Negate && pulse.Sweep.Period > 0x7FF)
}

func SetTimerLow(Value byte, currTimer uint16) uint16 {
	return (currTimer & 0x700) | uint16(Value)
}

func SetTimerHi(Value byte, currTimer uint16) uint16 {
	return (currTimer & 0xFF) | (uint16(Value&0x7) << 8)
}

func (pulse *PulseChannel) UpdatePulseOutput() {
	duty := apuDutySequences[pulse.Duty][pulse.DutyPos]
	var volume byte
	if pulse.ConstantVolume {
		volume = pulse.Volume
	} else {
		volume = pulse.Decay
	}

	if apuSilent || duty == 0 || pulse.LengthCounter.Counter == 0 || pulse.Timer < 8 {
		pulse.Output = 0
	} else {
		pulse.Output = (duty * volume)
	}
}

func (triangle *TriangleChannel) UpdateTriangleOutput() {
	triangle.Output = apuTriangleSequences[triangle.SeqPos]

	if apuSilent || triangle.LengthCounter.Counter == 0 || triangle.LinearCounter.Counter == 0 {
		triangle.Output = 0
	}
}

func (noise *NoiseChannel) SetNoiseTimer(Value byte) {
	noise.TimerReloadValue = apuNoiseTimerLUT[Value]
	noise.Timer = noise.TimerReloadValue
}

func (noise *NoiseChannel) ClockShiftRegister() {
	//Feedback is calculated as the exclusive-OR of bit 0 and one other bit: bit 6 if Mode flag is set, otherwise bit 1.
	firstBit := byte(noise.ShiftRegister & 1)
	secondBit := byte(noise.ShiftRegister&ternary(noise.Mode, 0x40, 0x02)) >> ternary(noise.Mode, 6, 1)
	feedback := firstBit ^ secondBit

	//The shift register is shifted right by one bit.
	noise.ShiftRegister >>= 1
	noise.ShiftRegister &= 0x7FFF

	//Bit 14, the leftmost bit, is set to the feedback calculated earlier.
	noise.ShiftRegister = (noise.ShiftRegister & 0b0011111111111111) | (uint16(feedback&1) << 14)
}

func (noise *NoiseChannel) UpdateNoiseOutput() {
	if apuSilent || !noise.Enabled || (noise.ShiftRegister&1) == 1 || noise.LengthCounter.Counter == 0 {
		noise.Output = 0
	} else {
		if noise.ConstantVolume {
			noise.Output = noise.Volume
		} else {
			noise.Output = noise.Decay
		}
	}
}

/*
func soundLoop() {
	f := float64(20000)

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
*/

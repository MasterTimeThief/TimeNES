package apu

import "mtt/timenes/common"

type APU struct {
}

//TODO: Fix variable names for channels
var Pulse1 PulseChannel
var Pulse2 PulseChannel
var Triangle TriangleChannel
var Noise NoiseChannel
var DMC DeltaModChannel

// var apuEnablePulse1, apuEnablePulse2, apuEnableTriangle, apuEnableNoise, apuEnableDMC bool
var APUDMCInterrupt, apuDMCDelayed, APUFrameInterrupt, APUInhibitIRQ, APUFrameCounterMode, apuIsHalfFrame bool
var apuFrameCounter int = 0
var IRQLevelDetector, DoIRQ bool

var apuDMAGetCycle, apuDoDMCDMA, apuDMCDMAHalt bool
var apuDMCDMADelay, apuCannotDMCDMARightNow byte

var apu4017ResetTimer int = 0
var MuteEmulator bool = false

//var apuLengthCounterLUT = [32]byte{10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14, 12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30}

func NewAPU() *APU {
	InitAudioOutput()
	apu := APU{}
	return &apu
}

func ResetAPU() {

	//Reset Channels
	Pulse1.ResetPulse()
	Pulse2.ResetPulse()
	Triangle.ResetTriangle()
	Noise.ResetNoise()
	DMC.ResetDMC()

	//APU Variables
	APUDMCInterrupt = false
	apuDMCDelayed = false
	APUFrameInterrupt = false
	APUInhibitIRQ = false
	APUFrameCounterMode = false
	//apuSilent = false
	apuIsHalfFrame = false
	apuFrameCounter = 0
	IRQLevelDetector = false
	DoIRQ = false

	apuDMAGetCycle = true
	apuDoDMCDMA, apuDMCDMAHalt = false, false
	apuDMCDMADelay, apuCannotDMCDMARightNow = 0, 0

	apu4017ResetTimer = 0
	//apuMixerInputBuffer = [apuMaxSamplesPerFrame]uint16{}
}

func APU_Cycle() {

	if apuDMAGetCycle { //DMA Get Cycle
		AudioOutput()
		DMA_Get()

	} else { //DMA Put Cycle
		ClockFrameCounter()
		DMA_Put()
	}
	//Clock triangle timer every cycle
	Triangle.ClockTriangleTimer()

	//Clock sequencer
	if (apu4017ResetTimer & 0x80) == 0 {
		apu4017ResetTimer--
		if (apu4017ResetTimer & 0x80) != 0 {
			apuFrameCounter = 0
		}
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
	if common.PendingMute {
		MuteEmulator = !MuteEmulator
		common.PendingMute = false
	}
}

func DMA_Get() {
	//Clock timers
	Pulse1.ClockPulseTimer()
	Pulse2.ClockPulseTimer()
	Noise.ClockNoiseTimer()
	DMC.ClockDMCTimer()

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
			DMC.Shifter = DMC.Buffer
			DMC.Enabled = true
		}
	}
}

func Set4017ResetTimer() {
	if apuDMAGetCycle {
		apu4017ResetTimer = 4
	} else {
		apu4017ResetTimer = 3
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

	if !APUFrameCounterMode { //4-Cycle mode
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
			APUFrameInterrupt = true
			if !IRQLevelDetector {
				IRQLevelDetector = !APUInhibitIRQ
			}
		case 14915:
			APUFrameInterrupt = !APUInhibitIRQ
			if !IRQLevelDetector {
				IRQLevelDetector = !APUInhibitIRQ
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
	Pulse1.ClockEnvelope()
	Pulse2.ClockEnvelope()
	Noise.ClockEnvelope()
	Triangle.ClockLinearCounter(Triangle.LengthCounter.HaltFlag)
}

func ClockFrameCounterHalfFrame() {

	//TODO: Split off Clocks into seperate functions for readability
	//Pulse1.ReloadLengthCounter()
	//Pulse2.ReloadLengthCounter()
	//Triangle.ReloadLengthCounter()
	//Noise.ReloadLengthCounter()

	// length counters and sweep
	/*if !Pulse1.Enabled {
		Pulse1.LengthCounter.Counter = 0
	}
	if !Pulse2.Enabled {
		Pulse2.LengthCounter.Counter = 0
	}
	if !Triangle.Enabled {
		Triangle.LengthCounter.Counter = 0
	}
	if !Noise.Enabled {
		Noise.LengthCounter.Counter = 0
	}*/

	Pulse1.ClockSweep(true)
	Pulse2.ClockSweep(false)

	Pulse1.ClockLengthCounter()
	Pulse2.ClockLengthCounter()
	Triangle.ClockLengthCounter()
	Noise.ClockLengthCounter()

	apuIsHalfFrame = true

}

func SetTimerLow(Value byte, currTimer uint16) uint16 {
	return (currTimer & 0x700) | uint16(Value)
}

func SetTimerHi(Value byte, currTimer uint16) uint16 {
	return (currTimer & 0xFF) | (uint16(Value&0x7) << 8)
}

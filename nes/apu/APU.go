package apu

import (
	"mtt/timenes/common"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type APU struct {
	cpu CPU

	//TODO: Fix variable names for channels
	Pulse1   PulseChannel
	Pulse2   PulseChannel
	Triangle TriangleChannel
	Noise    NoiseChannel
	DMC      DeltaModChannel

	audioContext *audio.Context
	player       *audio.Player
}

type CPU interface {
	Read(uint16) byte
	DelayCPU(int)
}

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
	apu := APU{}
	apu.InitAudioOutput()
	return &apu
}

func (a *APU) SetCPU(c CPU) {
	a.cpu = c
	a.DMC.SetCPU(c)
}

func (a *APU) ResetAPU() {

	//Reset Channels
	a.Pulse1.ResetPulse()
	a.Pulse2.ResetPulse()
	a.Triangle.ResetTriangle()
	a.Noise.ResetNoise()
	a.DMC.ResetDMC()

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

func (a *APU) APU_Cycle() {

	if apuDMAGetCycle { //DMA Get Cycle
		a.AudioOutput()
		a.DMA_Get()

	} else { //DMA Put Cycle
		a.ClockFrameCounter()
		a.DMA_Put()
	}
	//Clock triangle timer every cycle
	a.Triangle.ClockTriangleTimer()

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
		if apua.Pulse1.LengthCounter.ReloadFlag {
			apua.Pulse1.LengthCounter.Counter = apua.Pulse1.LengthCounter.ReloadValue
		}
		if apua.Pulse2.LengthCounter.ReloadFlag {
			apua.Pulse2.LengthCounter.Counter = apua.Pulse2.LengthCounter.ReloadValue
		}
		if apua.Triangle.LengthCounter.ReloadFlag {
			apua.Triangle.LengthCounter.Counter = apua.Triangle.LengthCounter.ReloadValue
		}
		if apua.Noise.LengthCounter.ReloadFlag {
			apua.Noise.LengthCounter.Counter = apua.Noise.LengthCounter.ReloadValue
		}
		apua.Pulse1.LengthCounter.ReloadFlag = false
		apua.Pulse2.LengthCounter.ReloadFlag = false
		apua.Triangle.LengthCounter.ReloadFlag = false
		apua.Noise.LengthCounter.ReloadFlag = false
	}*/

	apuDMAGetCycle = !apuDMAGetCycle
	if common.PendingMute {
		MuteEmulator = !MuteEmulator
		common.PendingMute = false
	}
}

func (a *APU) DMA_Get() {
	//Clock timers
	a.Pulse1.ClockPulseTimer()
	a.Pulse2.ClockPulseTimer()
	a.Noise.ClockNoiseTimer()
	a.DMC.ClockDMCTimer()

	if apuCannotDMCDMARightNow > 0 {
		apuCannotDMCDMARightNow -= 2
	}
}

func (a *APU) DMA_Put() {
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
			a.DMC.Shifter = a.DMC.Buffer
			a.DMC.Enabled = true
		}
	}
}

func (a *APU) ReadAPU(Address uint16) byte {
	status := byte(0)
	status |= byte(common.Ternary(APUDMCInterrupt, 0x80, 0x00))                                                                      //DMC Interrupt
	status |= byte(common.Ternary(APUFrameInterrupt, 0x40, 0x00))                                                                    //Frame Interrupt
	status |= byte(common.Ternary(a.DMC.BytesRemaining > 0, 0x10, 0x00))                                                             //DMC Active
	status |= byte(common.Ternary(!(a.Noise.LengthCounter.Counter == 0 /*|| apua.Noise.LengthCounter.HaltFlag*/), 0x08, 0x00))       //Noise Active
	status |= byte(common.Ternary(!(a.Triangle.LengthCounter.Counter == 0 /*|| apua.Triangle.LengthCounter.HaltFlag*/), 0x04, 0x00)) //Triangle Active
	status |= byte(common.Ternary(!(a.Pulse2.LengthCounter.Counter == 0 /*|| apua.Pulse2.LengthCounter.HaltFlag*/), 0x02, 0x00))     //Pulse 2 Active
	status |= byte(common.Ternary(!(a.Pulse1.LengthCounter.Counter == 0 /*|| apua.Pulse1.LengthCounter.HaltFlag*/), 0x01, 0x00))     //Pulse 1 Active
	APUFrameInterrupt = false
	return status
}

func (a *APU) WriteAPU(Address uint16, Value byte) {
	switch Address {
	// Pulse 1
	case 0x4000:
		a.Pulse1.Duty = (Value & 0xC0) >> 6
		a.Pulse1.LengthCounter.HaltFlag = ((Value & 0x20) >> 5) != 0
		a.Pulse1.Envelope.ConstantVolume = ((Value & 0x10) >> 4) != 0
		a.Pulse1.Envelope.Volume = (Value & 0xF)
		a.Pulse1.LoopFlag = a.Pulse1.HaltFlag
	case 0x4001:
		a.Pulse1.Sweep.Enabled = ((Value & 0x80) >> 7) != 0
		a.Pulse1.Sweep.Period = uint16(Value&0x70) >> 4
		a.Pulse1.Sweep.Negate = ((Value & 0x8) >> 3) != 0
		a.Pulse1.Sweep.Shift = (Value & 0x07)
		a.Pulse1.Sweep.ReloadFlag = true
	case 0x4002:
		a.Pulse1.TimerReloadValue = SetTimerLow(Value, a.Pulse1.TimerReloadValue)
		a.Pulse1.Timer = a.Pulse1.TimerReloadValue
	case 0x4003:
		a.Pulse1.TimerReloadValue = SetTimerHi(Value, a.Pulse1.TimerReloadValue)
		a.Pulse1.Timer = a.Pulse1.TimerReloadValue
		if a.Pulse1.Enabled {
			a.Pulse1.LengthCounter.Counter = LengthCounterLoad(Value >> 3)
		}
		a.Pulse1.StartFlag = true
		a.Pulse1.DutyPos = 0

	// Pulse 2
	case 0x4004:
		a.Pulse2.Duty = (Value & 0xC0) >> 6
		a.Pulse2.LengthCounter.HaltFlag = ((Value & 0x20) >> 5) != 0
		a.Pulse2.Envelope.ConstantVolume = ((Value & 0x10) >> 4) != 0
		a.Pulse2.Envelope.Volume = (Value & 0xF)
		a.Pulse2.LoopFlag = a.Pulse2.HaltFlag
	case 0x4005:
		a.Pulse2.Sweep.Enabled = ((Value & 0x80) >> 7) != 0
		a.Pulse2.Sweep.Period = uint16(Value&0x70) >> 4
		a.Pulse2.Sweep.Negate = ((Value & 0x8) >> 3) != 0
		a.Pulse2.Sweep.Shift = (Value & 0x07)
		a.Pulse2.Sweep.ReloadFlag = true
	case 0x4006:
		a.Pulse2.TimerReloadValue = SetTimerLow(Value, a.Pulse2.TimerReloadValue)
		a.Pulse2.Timer = a.Pulse2.TimerReloadValue
	case 0x4007:
		a.Pulse2.TimerReloadValue = SetTimerHi(Value, a.Pulse2.TimerReloadValue)
		a.Pulse2.Timer = a.Pulse2.TimerReloadValue
		if a.Pulse2.Enabled {
			a.Pulse2.LengthCounter.Counter = LengthCounterLoad(Value >> 3)
		}
		a.Pulse2.StartFlag = true
		a.Pulse1.DutyPos = 0

	// Triangle
	case 0x4008:
		a.Triangle.LengthCounter.HaltFlag = ((Value & 0x80) >> 7) != 0
		if a.Triangle.Enabled {
			a.Triangle.LinearCounter.ReloadValue = (Value & 0x7F)
		}
	case 0x4009: //Unused
	case 0x400A:
		//apua.Triangle.Timer = uint16(Value)
		a.Triangle.TimerReloadValue = SetTimerLow(Value, a.Triangle.TimerReloadValue)
		a.Triangle.Timer = a.Triangle.TimerReloadValue
	case 0x400B:
		//apua.Triangle.Timer |= (uint16(Value&0x7) << 8)
		a.Triangle.TimerReloadValue = SetTimerHi(Value, a.Triangle.TimerReloadValue)
		a.Triangle.Timer = a.Triangle.TimerReloadValue
		if a.Triangle.Enabled {
			a.Triangle.LengthCounter.Counter = LengthCounterLoad(Value >> 3)
		}
		a.Triangle.LinearCounter.ReloadFlag = true

	// Noise
	case 0x400C:
		a.Noise.LengthCounter.HaltFlag = ((Value & 0x20) >> 5) != 0
		a.Noise.Envelope.ConstantVolume = ((Value & 0x10) >> 4) != 0
		a.Noise.Envelope.Volume = (Value & 0xF)
	case 0x400D: //Unused
	case 0x400E:
		a.Noise.Mode = ((Value & 0x80) >> 7) != 0
		a.Noise.SetNoiseTimer(Value & 0xF)
	case 0x400F:
		if a.Noise.Enabled {
			a.Noise.LengthCounter.Counter = LengthCounterLoad(Value >> 3)
		}
		a.Noise.StartFlag = true

	// DMC
	case 0x4010:
		a.DMC.IRQEnable = ((Value & 0x80) >> 7) != 0
		a.DMC.Loop = ((Value & 0x40) >> 6) != 0
		a.DMC.SampleRate = APUDMCSampleRateLUT[Value&0xF]
		if !a.DMC.IRQEnable {
			APUDMCInterrupt = false
			IRQLevelDetector = false
		}
	case 0x4011:
		a.DMC.Output = (Value & 0x7F)
	case 0x4012:
		a.DMC.SampleAddress = (0xC000 | (uint16(Value) << 6))
		a.DMC.CurrentAddress = a.DMC.SampleAddress
	case 0x4013:
		a.DMC.SampleLength = ((uint16(Value) << 4) | 1)
		a.DMC.BytesRemaining = a.DMC.SampleLength

	case 0x4015: //APU Status
		//apua.DMC.BytesRemaining = int((Value & 0x10) >> 4)
		a.Noise.LengthCounter.Counter &= (0xFF * ((Value & 0x08) >> 3))
		a.Triangle.LengthCounter.Counter &= (0xFF * ((Value & 0x04) >> 2))
		a.Triangle.LinearCounter.Counter &= (0xFF * ((Value & 0x04) >> 2))
		a.Pulse2.LengthCounter.Counter &= (0xFF * ((Value & 0x02) >> 1))
		a.Pulse1.LengthCounter.Counter &= (0xFF * (Value & 0x01))

		a.DMC.Enabled = (Value & 0x10) != 0
		a.Noise.Enabled = (Value & 0x08) != 0
		a.Triangle.Enabled = (Value & 0x04) != 0
		a.Pulse2.Enabled = (Value & 0x02) != 0
		a.Pulse1.Enabled = (Value & 0x01) != 0

		if a.DMC.Enabled {
			if a.DMC.BytesRemaining > 0 {
				a.DMC.DMCRestartSample()
			}
		} else {
			a.DMC.BytesRemaining = 0
		}

		APUDMCInterrupt = false
		IRQLevelDetector = false

	case 0x4017: //APU Frame Counter control
		//modeFlagPrev := apuFrameCounterMode
		APUFrameCounterMode = ((Value & 0x80) >> 7) != 0
		APUInhibitIRQ = ((Value & 0x40) >> 6) != 0
		if APUInhibitIRQ {
			APUFrameInterrupt = false
			IRQLevelDetector = false
		}
		if /*!modeFlagPrev &&*/ APUFrameCounterMode {
			a.ClockFrameCounterQuarterFrame()
			a.ClockFrameCounterHalfFrame()
		}
		Set4017ResetTimer()

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

func (a *APU) ClockFrameCounter() { //Also called Frame Sequencer
	//Do frame counter shit
	apuFrameCounter++
	apuIsHalfFrame = false

	if !APUFrameCounterMode { //4-Cycle mode
		switch apuFrameCounter {
		case 3728:
			a.ClockFrameCounterQuarterFrame()
		case 7456:
			a.ClockFrameCounterQuarterFrame()
			a.ClockFrameCounterHalfFrame()
		case 11185:
			a.ClockFrameCounterQuarterFrame()
		case 14914:
			if !apuDMAGetCycle {
				a.ClockFrameCounterQuarterFrame()
				a.ClockFrameCounterHalfFrame()
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
			a.ClockFrameCounterQuarterFrame()
		case 7456:
			a.ClockFrameCounterQuarterFrame()
			a.ClockFrameCounterHalfFrame()
		case 11185:
			a.ClockFrameCounterQuarterFrame()
		case 14914:
			//Nothing?
		case 18640:
			a.ClockFrameCounterQuarterFrame()
			a.ClockFrameCounterHalfFrame()
		case 18641:
			apuFrameCounter = 0
		}
	}

}

func (a *APU) ClockFrameCounterQuarterFrame() {
	a.Pulse1.ClockEnvelope()
	a.Pulse2.ClockEnvelope()
	a.Noise.ClockEnvelope()
	a.Triangle.ClockLinearCounter(a.Triangle.LengthCounter.HaltFlag)
}

func (a *APU) ClockFrameCounterHalfFrame() {

	//TODO: Split off Clocks into seperate functions for readability
	//a.Pulse1.ReloadLengthCounter()
	//a.Pulse2.ReloadLengthCounter()
	//a.Triangle.ReloadLengthCounter()
	//a.Noise.ReloadLengthCounter()

	// length counters and sweep
	/*if !a.Pulse1.Enabled {
		a.Pulse1.LengthCounter.Counter = 0
	}
	if !a.Pulse2.Enabled {
		a.Pulse2.LengthCounter.Counter = 0
	}
	if !a.Triangle.Enabled {
		a.Triangle.LengthCounter.Counter = 0
	}
	if !a.Noise.Enabled {
		a.Noise.LengthCounter.Counter = 0
	}*/

	a.Pulse1.ClockSweep(true)
	a.Pulse2.ClockSweep(false)

	a.Pulse1.ClockLengthCounter()
	a.Pulse2.ClockLengthCounter()
	a.Triangle.ClockLengthCounter()
	a.Noise.ClockLengthCounter()

	apuIsHalfFrame = true

}

func SetTimerLow(Value byte, currTimer uint16) uint16 {
	return (currTimer & 0x700) | uint16(Value)
}

func SetTimerHi(Value byte, currTimer uint16) uint16 {
	return (currTimer & 0xFF) | (uint16(Value&0x7) << 8)
}

func (a *APU) GetPulse1Mute() *bool {
	return &a.Pulse1.ForceMute
}
func (a *APU) GetPulse2Mute() *bool {
	return &a.Pulse2.ForceMute
}
func (a *APU) GetTriangleMute() *bool {
	return &a.Triangle.ForceMute
}
func (a *APU) GetNoiseMute() *bool {
	return &a.Noise.ForceMute
}
func (a *APU) GetDMCMute() *bool {
	return &a.DMC.ForceMute
}

func (a *APU) SetEmulatorvolume(vol float64) {
	a.player.SetVolume(vol)
}

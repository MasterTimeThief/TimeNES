package apu

import "mtt/timenes/nes/cartridge"

type DeltaModChannel struct {
	Enabled              bool
	Timer                uint16
	IRQEnable            bool
	Loop                 bool
	SampleRate           uint16 // AKA Frequency
	Output               byte   // AKA LoadCounter
	SampleAddress        uint16
	CurrentAddress       uint16
	SampleLength         uint16
	BytesRemaining       uint16
	Buffer               byte
	SampleAddressCounter uint16
	Shifter              byte
	ShifterBitsRemaining byte
	//DPCM_Up              bool

	ForceMute bool

	// This will get the whole sample at once,
	// so I don't have to go back and read it
	// each time I need a new sample.
	// Also because Go is being a baby about it.
	//SampleBuffer    []byte
	//SampleBufferPos byte

	cpu CPU
}

func (d *DeltaModChannel) SetCPU(c CPU) {
	d.cpu = c
}

var APUDMCSampleRateLUT = [16]uint16{214, 190, 170, 160, 143, 127, 113, 107, 95, 80, 71, 64, 53, 42, 36, 27}

func (d *DeltaModChannel) ResetDMC() {
	d.Enabled = false
	d.Timer = 0
	d.IRQEnable = false
	d.Loop = false
	d.SampleRate = 0
	d.Output = 0
	d.SampleAddress = 0
	d.SampleLength = 0
	d.BytesRemaining = 0
	d.Buffer = 0
	d.SampleAddressCounter = 0
	d.Shifter = 0
	d.ShifterBitsRemaining = 0
	//d.DPCM_Up = false

	d.ForceMute = false
}

func (d *DeltaModChannel) ClockDMCTimer() {
	if d.Timer == 0 {
		d.ClockDMCOutputUnit()
		d.Timer = d.SampleRate
	} else {
		d.Timer--
	}
}

func (d *DeltaModChannel) ClockDMCOutputUnit() {
	if d.Enabled {
		if (d.Shifter & 1) == 1 {
			if d.Output <= 125 { // this is 7 bit, and cannot go above 127
				d.Output += 2
			}
		} else {
			if d.Output >= 2 { // this is 7 bit, and cannot go below 0
				d.Output -= 2
			}
		}
	}

	d.Shifter >>= 1                  // shift the bits in the shift register
	d.ShifterBitsRemaining--         // and decrement the "bits remaining" counter.
	if d.ShifterBitsRemaining == 0 { // If there are no bits left, a new output cycle is started
		d.DMCOutputCycle()
	}
}

func (d *DeltaModChannel) DMCOutputCycle() {
	d.ShifterBitsRemaining = 8 // it's time for a DMC DMA!

	if d.BytesRemaining > 0 /*|| APU_SetImplicitAbortDMC4015*/ {
		//if !apuDoDMCDMA && apuCannotDMCDMARightNow != 2 {
		//	// if playing a sample:
		//	apuDoDMCDMA = true
		//	apuDMCDMAHalt = true
		//}
		d.Shifter = d.Buffer // and set up the shifter with the new values.
		d.Buffer = 0
		d.DMCMemoryReader()
		d.Enabled = true // The DMC is not silent.
	} else {
		d.Enabled = false
	}
}

func (d *DeltaModChannel) DMCMemoryReader() {
	//Check for Mappers
	switch cartridge.MapperChipID {
	default:
		if d.Buffer == 0 && d.BytesRemaining > 0 {
			d.cpu.DelayCPU(4)
			d.Buffer = d.cpu.Read(d.CurrentAddress)

			//Advance the address
			d.CurrentAddress++
			if d.CurrentAddress == 0 { //We hit 0xFFFF and wrapped around
				d.CurrentAddress = 0x8000
			}
			if d.BytesRemaining > 0 {
				d.BytesRemaining--
			}

			if d.BytesRemaining == 0 {
				if d.Loop {
					d.DMCRestartSample()
				} else {
					d.Enabled = false
					if d.IRQEnable {
						APUDMCInterrupt = true
					}
				}
			}
		} else {
			d.Buffer = 0
		}
	}
}

func (d *DeltaModChannel) DMCRestartSample() {
	//d.SampleBufferPos = 0
	d.CurrentAddress = d.SampleAddress
	d.BytesRemaining = d.SampleLength
}

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
	DPCM_Up              bool

	// This will get the whole sample at once,
	// so I don't have to go back and read it
	// each time I need a new sample.
	// Also because Go is being a baby about it.
	SampleBuffer    []byte
	SampleBufferPos byte
}

var APUDMCSampleRateLUT = [16]uint16{428, 380, 340, 320, 286, 254, 226, 214, 190, 160, 142, 128, 106, 84, 72, 54}

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
	d.DPCM_Up = false
}

func (d *DeltaModChannel) ClockDMCTimer() {

	//Refill the sample buffer if empty here?
	if d.Buffer == 0 {
		d.DMCMemoryReader()
	}

	if d.Timer == 0 {
		if d.Enabled {
			d.Timer = d.SampleRate
			d.DPCM_Up = (d.Shifter & 1) == 1
			if d.DPCM_Up {
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
		if d.ShifterBitsRemaining == 0 { // If there are no bits left,

			d.ShifterBitsRemaining = 8 // it's time for a DMC DMA!

			if d.BytesRemaining > 0 /*|| APU_SetImplicitAbortDMC4015*/ {
				if !apuDoDMCDMA && apuCannotDMCDMARightNow != 2 {
					// if playing a sample:
					apuDoDMCDMA = true
					apuDMCDMAHalt = true
				}
				//if APU_SetImplicitAbortDMC4015 {
				//	APU_ImplicitAbortDMC4015 = true // check for weird DMA abort behavior
				//	APU_SetImplicitAbortDMC4015 = false
				//}
				d.Shifter = d.Buffer // and set up the shifter with the new values.
				d.Buffer = 0
				d.Enabled = true // The DMC is not silent.

			} else {
				d.Enabled = false
			}
		}
	} else {
		d.Timer--
	}
}

func (d *DeltaModChannel) DMCMemoryReader() {
	//Check for Mappers
	switch cartridge.MapperChipID {
	default:
		//TODO: Stall for 1-4 CPU cycles (?)
		//d.Buffer = nes.Read(d.SampleAddress)
		if len(d.SampleBuffer) > 0 {
			d.Buffer = d.SampleBuffer[d.SampleBufferPos]
			d.SampleBufferPos++

			//Advance the address, even if technically we don't need to
			d.CurrentAddress++
			if d.CurrentAddress < 0x8000 { //We hit 0xFFFF and wrapped around
				d.CurrentAddress += 0x8000
			}
			d.BytesRemaining--
			if d.BytesRemaining == 0 && d.Loop {
				d.SampleBufferPos = 0
				d.CurrentAddress = d.SampleAddress
				d.BytesRemaining = d.SampleLength
			} else if d.BytesRemaining == 0 && d.IRQEnable {
				APUDMCInterrupt = true
			}
		} else {
			d.Buffer = 0
		}
	}
}

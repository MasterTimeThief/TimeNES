package apu

import (
	"mtt/timenes/common"
)

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

	//d.Timer--
	d.Timer-- // the table is in CPU cycles, but the count is in APU cycles

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
				d.Enabled = true     // The APU is not silent.

			} else {
				d.Enabled = false
			}
		}
	}
}

func (d *DeltaModChannel) FillDMCBuffer() {
	//Check for Mappers
	switch common.MapperChipID {
	default:
		//d.Buffer = nes.Read(d.SampleAddress)
	}
}

package apu

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

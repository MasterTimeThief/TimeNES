package common

// Header Variables

var PRGROM_Size uint32 // Size of PRG ROM
var CHRROM_Size uint32 // Size of CHR ROM (value 0 means the board uses CHR RAM)
var IsNametableHorizontal bool
var HasBatteryRAM bool
var AltNametableLayout bool
var NES2_Header bool //Is the header in NES 2.0 format, rather than iNES

var MapperChipID byte

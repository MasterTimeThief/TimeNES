package common

var Filepath string = ""
var ROMExists bool = false
var ROMLoaded bool = false

func SelectROM(fp string) {
	Filepath = fp
	ROMExists = true
	ROMLoaded = false
}

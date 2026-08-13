package common

import "github.com/sqweek/dialog"

var Filepath string = ""
var ROMExists bool = false
var ROMLoaded bool = false

func SelectROM(fp string) {
	Filepath = fp
	ROMExists = true
	ROMLoaded = false
}

func FileSelectDialog() {
	file, err := dialog.File().Title("Select NES ROM File").Filter("NES ROM files (*.nes)", "nes").Load()
	if err == nil {
		SelectROM(file)
	}
}

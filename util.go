package main

import "os"

func BoolToInt(Flag bool) int {
	if Flag {
		return 1
	}
	return 0
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ternary(Condition bool, ValT, ValF uint16) uint16 {
	if Condition {
		return ValT
	}
	return ValF
}

func GetCurrentDirectory(dir string) {
	dirEntries, err := os.ReadDir(dir)
	check(err)

	g.directory = dirEntries
}

func firstN(str string, n int) string {
	v := []rune(str)
	if n >= len(v) {
		return str
	}
	return string(v[:n])
}

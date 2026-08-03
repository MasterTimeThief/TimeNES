package main

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

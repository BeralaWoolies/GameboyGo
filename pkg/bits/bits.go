package bits

func Set(val uint8, pos uint8) uint8 {
	return val | (1 << pos)
}

func Reset(val uint8, pos uint8) uint8 {
	return val & ^(1 << pos)
}

func HiByte(val uint16) uint8 {
	return uint8(val >> 8)
}

func LoByte(val uint16) uint8 {
	return uint8(val & 0xFF)
}

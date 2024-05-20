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

func IsSetInBoth(val1 uint8, val2 uint8, pos uint8) bool {
	return IsSet(val1, pos) && IsSet(val2, pos)
}

func IsSet(val uint8, pos uint8) bool {
	return val&(1<<pos) != 0
}

func GetBit(val uint8, pos uint8) uint8 {
	return (val >> pos) & 1
}

package bits

import "fmt"

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

func SetCond(val uint8, pos uint8, cond uint8) uint8 {
	if cond != 0 {
		return Set(val, pos)
	}

	return Reset(val, pos)
}

func BoolToUint8(val bool) uint8 {
	if val {
		return 1
	}

	return 0
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

func HexString(val int) string {
	return fmt.Sprintf("0x%02x", val)
}

func NBitMask(n uint32) uint32 {
	res := uint32(0)
	for bit := 0; bit < int(n); bit++ {
		res |= (1 << bit)
	}

	return res
}

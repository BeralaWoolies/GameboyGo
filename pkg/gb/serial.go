package gb

import "log"

type SerialPort struct {
	ic *IntruptController
	sb uint8
	sc uint8
}

const (
	SB_ADDR = 0xFF01
	SC_ADDR = 0xFF02
)

func (s *SerialPort) init(ic *IntruptController) {
	s.ic = ic
	s.sb = 0xFF
}

func (s *SerialPort) contains(addr uint16) bool {
	return inRange(addr, SB_ADDR, SC_ADDR)
}

func (s *SerialPort) read(addr uint16) uint8 {
	switch addr {
	case SB_ADDR:
		return s.sb
	case SC_ADDR:
		return s.sc
	default:
		log.Fatalf("MMU mapped an illegal read address: 0x%02x to Serial Port", addr)
		return 0xFF
	}
}

func (s *SerialPort) write(addr uint16, data uint8) {
	switch addr {
	case SB_ADDR:
		return
	case SC_ADDR:
		s.sc = data
		if data == 0x81 {
			s.ic.requestIntrupt(SERIAL_INTRUPT_BIT)
		}
	default:
		log.Fatalf("MMU mapped an illegal write address: 0x%02x to Serial Port", addr)
	}
}

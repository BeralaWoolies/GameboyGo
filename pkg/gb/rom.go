package gb

import (
	"log"
	"os"
)

type ROM struct {
	rom [ROM_SIZE]byte
}

const (
	ROM_SIZE = 0x8000
	ROM_BASE = 0x0
	ROM_TOP  = 0x7FFF
)

func newROM(filename string) *ROM {
	r := &ROM{}

	rom, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	copy(r.rom[:], rom)

	return r
}

func (r *ROM) contains(addr uint16) bool {
	return inRange(addr, ROM_BASE, ROM_TOP)
}

func (r *ROM) read(addr uint16) uint8 {
	return r.rom[addr]
}

func (r *ROM) write(addr uint16, data uint8) {
	return
}

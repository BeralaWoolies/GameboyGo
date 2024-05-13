package gb

import (
	"log"
	"os"
)

type ROM struct {
	rom [0x8000]byte
}

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
	return addr <= 0x7FFF
}

func (r *ROM) read(addr uint16) uint8 {
	return r.rom[addr]
}

func (r *ROM) write(addr uint16, data uint8) {
	r.rom[addr] = data
}

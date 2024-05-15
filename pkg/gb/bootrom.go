package gb

import (
	"fmt"
	"log"
	"os"
)

type BootRom struct {
	rom       []byte
	enableReg uint8
	mmu       *MMU
}

// boot rom will take precedence before I/O address space, for R/W to 0xFF50
const (
	BOOT_ROM_BASE        = 0x0
	BOOR_ROM_TOP         = 0xFF
	BOOT_ROM_ENABLE_ADDR = 0xFF50
)

func newBootROM(filename string, mmu *MMU) *BootRom {
	b := &BootRom{enableReg: 0x0, mmu: mmu}

	var err error
	b.rom, err = os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

func (b *BootRom) contains(addr uint16) bool {
	return inRange(addr, BOOT_ROM_BASE, BOOR_ROM_TOP) || addr == BOOT_ROM_ENABLE_ADDR
}

func (b *BootRom) read(addr uint16) uint8 {
	switch addr {
	case BOOT_ROM_ENABLE_ADDR:
		return b.enableReg
	default:
		return b.rom[addr]
	}
}

func (b *BootRom) write(addr uint16, data uint8) {
	switch addr {
	case BOOT_ROM_ENABLE_ADDR:
		b.enableReg = data
		if b.enableReg != 0 {
			fmt.Println("Disabling BOOT ROM...")

			// unmap boot rom and let data "fall" to the next addrSpace
			b.mmu.unmapAddrSpace(b)
			b.mmu.write(BOOT_ROM_ENABLE_ADDR, data)
		}
	default:
		b.rom[addr] = data
	}
}

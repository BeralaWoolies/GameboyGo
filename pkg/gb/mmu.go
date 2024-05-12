package gb

import "fmt"

type MMU struct {
	memory [0x10000]uint8
}

func (mmu *MMU) init() {
	mmu.memory = [0x10000]uint8{}
	mmu.memory[0xFF05] = 0x00
	mmu.memory[0xFF06] = 0x00
	mmu.memory[0xFF07] = 0x00
	mmu.memory[0xFF10] = 0x80
	mmu.memory[0xFF11] = 0xBF
	mmu.memory[0xFF12] = 0xF3
	mmu.memory[0xFF14] = 0xBF
	mmu.memory[0xFF16] = 0x3F
	mmu.memory[0xFF17] = 0x00
	mmu.memory[0xFF19] = 0xBF
	mmu.memory[0xFF1A] = 0x7F
	mmu.memory[0xFF1B] = 0xFF
	mmu.memory[0xFF1C] = 0x9F
	mmu.memory[0xFF1E] = 0xBF
	mmu.memory[0xFF20] = 0xFF
	mmu.memory[0xFF21] = 0x00
	mmu.memory[0xFF22] = 0x00
	mmu.memory[0xFF23] = 0xBF
	mmu.memory[0xFF24] = 0x77
	mmu.memory[0xFF25] = 0xF3
	mmu.memory[0xFF26] = 0xF1
	mmu.memory[0xFF40] = 0x91
	mmu.memory[0xFF42] = 0x00
	mmu.memory[0xFF43] = 0x00
	mmu.memory[0xFF45] = 0x00
	mmu.memory[0xFF47] = 0xFC
	mmu.memory[0xFF48] = 0xFF
	mmu.memory[0xFF49] = 0xFF
	mmu.memory[0xFF4A] = 0x00
	mmu.memory[0xFF4B] = 0x00
	mmu.memory[0xFFFF] = 0x00
}

func (mmu *MMU) mapRom(bootRom []byte) {
	copy(mmu.memory[:], bootRom)
	for _, b := range mmu.memory {
		fmt.Printf("%02x ", b)
	}
}

func (mmu *MMU) read(address uint16) uint8 {
	if address <= 0xFFFF {
		return mmu.memory[address]
	}
	return 0xFF
}

func (mmu *MMU) write(address uint16, data uint8) {
	if address <= 0xFFFF {
		mmu.memory[address] = data
		fmt.Printf("Write at: [0x%04x] = 0x%02x\n", address, data)
	}
}

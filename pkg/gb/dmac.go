package gb

import (
	"fmt"
	"log"
)

type DMAC struct {
	src      uint8
	active   bool
	currByte uint16
	delayed  bool
	mmu      *MMU
}

const OAM_DMA_TRANSFER_ADDR = 0xFF46

func (dmac *DMAC) init(mmu *MMU) {
	dmac.src = 0
	dmac.active = false
	dmac.currByte = 0
	dmac.delayed = false
	dmac.mmu = mmu
}

func (dmac *DMAC) contains(addr uint16) bool {
	return addr == OAM_DMA_TRANSFER_ADDR
}

func (dmac *DMAC) write(addr uint16, data uint8) {
	if addr == OAM_DMA_TRANSFER_ADDR {
		dmac.initOAMTransfer(data)
		return
	}

	log.Fatalf("MMU mapped an illegal write address: 0x%02x to DMA controller", addr)
}

func (dmac *DMAC) read(addr uint16) uint8 {
	if addr == OAM_DMA_TRANSFER_ADDR {
		fmt.Println("Reading from DMA")
		return dmac.src
	}

	log.Fatalf("MMU mapped an illegal read address: 0x%02x to DMA controller", addr)
	return 0xFF
}

func (dmac *DMAC) step(cTicks int) {
	if !dmac.active {
		return
	}

	if dmac.delayed {
		cTicks -= 4
		dmac.delayed = false

		if cTicks <= 0 {
			return
		}
	}

	for i := 0; i < cTicks; i += 4 {
		dmac.transferOAM()
	}
}

func (dmac *DMAC) transferOAM() {
	if !dmac.active {
		return
	}

	dest := OAM_BASE + dmac.currByte
	src := dmac.mmu.read((uint16(dmac.src) * 0x100) + dmac.currByte)
	dmac.mmu.write(dest, src)

	dmac.currByte++
	dmac.active = bool(dmac.currByte < OAM_SIZE)
	if !dmac.active {
		fmt.Printf("DMA OAM transfer completed, %d bytes transferred\n", dmac.currByte)
	}
}

func (dmac *DMAC) initOAMTransfer(data uint8) {
	fmt.Printf("DMA OAM transfer starting at 0x%02x\n", data)
	dmac.src = data
	dmac.active = true
	dmac.currByte = 0
	dmac.delayed = true
}

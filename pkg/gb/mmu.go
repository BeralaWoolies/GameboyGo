package gb

import (
	"fmt"
	"log"
	"slices"
)

type Addressable interface {
	contains(address uint16) bool
	read(address uint16) uint8
	write(address uint16, data uint8)
}

type MMU struct {
	addrSpaces []Addressable
}

func (mmu *MMU) mapAddrSpace(addrSpace Addressable) {
	mmu.addrSpaces = append(mmu.addrSpaces, addrSpace)
}

func (mmu *MMU) unmapAddrSpace(addrSpace Addressable) {
	if idx := slices.Index(mmu.addrSpaces, addrSpace); idx != -1 {
		fmt.Printf("Unmapping address space at index: %d\n", idx)
		mmu.addrSpaces = append(mmu.addrSpaces[:idx], mmu.addrSpaces[idx+1:]...)
	}
}

func (mmu *MMU) addrSpace(addr uint16) Addressable {
	for _, space := range mmu.addrSpaces {
		if space.contains(addr) {
			return space
		}
	}

	return nil
}

func (mmu *MMU) read(addr uint16) uint8 {
	if space := mmu.addrSpace(addr); space != nil {
		return space.read(addr)
	}

	log.Fatalf("MMU has no mapping read address: 0x%02x", addr)
	return 0xFF
}

func (mmu *MMU) write(addr uint16, data uint8) {
	if space := mmu.addrSpace(addr); space != nil {
		space.write(addr, data)
		return
	}

	log.Fatalf("MMU has no mapping write address: 0x%02x", addr)
}

func inRange(addr uint16, base uint16, top uint16) bool {
	return addr >= base && addr <= top
}

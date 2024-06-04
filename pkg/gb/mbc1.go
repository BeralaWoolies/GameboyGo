package gb

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
)

type MBC1 struct {
	cart        *Cart
	ramEnabled  bool
	numBanks    uint32
	romBankMask uint32
	ramBankNum  uint32
	romLo       uint32
	mode        BankMode
}

type BankMode uint8

const (
	MODE0 = 0
	MODE1 = 1
)

func (mbc *MBC1) init(cart *Cart) {
	mbc.cart = cart
	mbc.romLo = 1
	mbc.mode = MODE0
	mbc.numBanks = uint32(cart.romSize / 0x4000)
	nBits := uint32(math.Log2(float64(mbc.numBanks)))
	mbc.romBankMask = bits.NBitMask(nBits)

	fmt.Println("MBC1 INFO:")
	fmt.Println("numBanks: ", mbc.numBanks)
	fmt.Println("Bits to address banks: ", nBits)
	fmt.Println("Rom bank mask: ", "0b"+strconv.FormatInt(int64(mbc.romBankMask), 2))
	fmt.Println("====================================")
}

func (mbc *MBC1) contains(address uint16) bool {
	return inRange(address, ROM_BASE, ROM_TOP) || inRange(address, EXT_RAM_BASE, EXT_RAM_TOP)
}

func (mbc *MBC1) read(addr uint16) uint8 {
	if inRange(addr, ROM_BASE, ROM_TOP) {
		if inRange(addr, ROM_BASE, 0x3FFF) {
			// bank 00 of rom
			if mbc.bigROM() && mbc.mode == MODE1 {
				return mbc.cart.rom[(((mbc.ramBankNum<<5)&mbc.romBankMask)*0x4000)+uint32(addr)]
			}

			return mbc.cart.rom[addr]
		}

		if inRange(addr, 0x4000, ROM_TOP) {
			// switchable bank of rom
			return mbc.cart.rom[((((mbc.ramBankNum<<5)|mbc.romLo)&mbc.romBankMask)*0x4000)+uint32(addr-0x4000)]
		}
	} else if inRange(addr, EXT_RAM_BASE, EXT_RAM_TOP) {
		if !mbc.ramEnabled {
			return 0xFF
		}

		bank := uint32(0)
		if mbc.bigRam() && mbc.mode == MODE1 {
			bank = mbc.ramBankNum
		}

		return mbc.cart.ram[(uint32(bank)*0x2000+uint32(addr-EXT_RAM_BASE))%mbc.cart.ramSize]
	}

	log.Fatalf("MMU mapped an illegal read address: 0x%02x to MBC1", addr)
	return 0xFF
}

func (mbc *MBC1) write(addr uint16, data uint8) {
	if inRange(addr, ROM_BASE, ROM_TOP) {
		if inRange(addr, ROM_BASE, 0x1FFF) {
			// trap to modify ram enable register
			mbc.ramEnabled = (data&0xF == 0xA)
			return
		}

		if inRange(addr, 0x2000, 0x3FFF) {
			// set lower 5 bits of rom bank register
			val := data & 0x1F
			if val == 0 {
				val = 1
			}

			mbc.romLo = uint32(val) & mbc.romBankMask
			return
		}

		if inRange(addr, 0x4000, 0x5FFF) {
			// set either the upper 2 bits of rom bank register or select ram banks 0-3
			mbc.ramBankNum = uint32(data & 0x3)
			return
		}

		if inRange(addr, 0x6000, 0x7FFF) {
			// set the banking mode
			mbc.mode = BankMode(data & 0x1)
			return
		}
	} else if inRange(addr, EXT_RAM_BASE, EXT_RAM_TOP) {
		if !mbc.ramEnabled {
			return
		}

		bank := uint32(0)
		if mbc.bigRam() && mbc.mode == MODE1 {
			bank = mbc.ramBankNum
		}

		mbc.cart.ram[(uint32(bank)*0x2000+uint32(addr-EXT_RAM_BASE))%mbc.cart.ramSize] = data
		return
	}

	log.Fatalf("MMU mapped an illegal write address: 0x%02x to MBC1", addr)
}

func (mbc *MBC1) bigROM() bool {
	// ROM can make use of the 2-bit register
	return mbc.cart.romSize >= 0x100000
}

func (mbc *MBC1) bigRam() bool {
	// RAM can make use of the 2-bit register
	return mbc.cart.ramSize > 0x2000
}

package gb

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
)

type MBC1 struct {
	cart         *Cart
	ram          []byte
	ramEnabled   bool
	romBankNum   uint8
	romBankMask  uint8
	romBankNBits uint8
	ramBankNum   uint8
	numBanks     uint8
	mode         BankMode
}

type BankMode uint8

const (
	MODE0 = 0
	MODE1 = 1
)

func (mbc *MBC1) init(cart *Cart) {
	mbc.cart = cart
	mbc.romBankNum = 1
	mbc.ramEnabled = false
	mbc.mode = MODE0
	mbc.numBanks = uint8(cart.romSize / 0x4000)
	mbc.romBankNBits = uint8(math.Log2(float64(mbc.numBanks)))
	mbc.romBankMask = bits.NBitMask(mbc.romBankNBits)

	fmt.Println("MBC1 INFO:")
	fmt.Println("numBanks: ", mbc.numBanks)
	fmt.Println("Bits to address banks: ", mbc.romBankNBits)
	fmt.Println("Rom bank mask: ", "0b"+strconv.FormatInt(int64(mbc.romBankMask), 2))
	fmt.Println("====================================")

	if cart.ramSize != 0 {
		mbc.ram = make([]byte, cart.ramSize, cart.ramSize)
	}
}

func (mbc *MBC1) contains(address uint16) bool {
	return inRange(address, ROM_BASE, ROM_TOP) || inRange(address, EXT_RAM_BASE, EXT_RAM_TOP)
}

func (mbc *MBC1) read(addr uint16) uint8 {
	if inRange(addr, ROM_BASE, ROM_TOP) {
		if inRange(addr, ROM_BASE, 0x3FFF) {
			// bank 00 of rom
			return mbc.cart.rom[addr]
		}

		if inRange(addr, 0x4000, ROM_TOP) {
			// switchable bank of rom
			bank := mbc.romBankNum
			// if mbc.mode == MODE1 {
			// 	bank &= 0x1F
			// }

			return mbc.cart.rom[(uint16(bank)*0x4000)+addr-0x4000]
		}
	} else if inRange(addr, EXT_RAM_BASE, EXT_RAM_TOP) {
		if !mbc.ramEnabled {
			return 0xFF
		}

		bank := mbc.ramBankNum
		// if mbc.mode == MODE0 {
		// 	bank = 0
		// }

		return mbc.ram[(uint16(bank)*0x2000)+addr-EXT_RAM_BASE]
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
			if data == 0 {
				// bank 0x00 -> 0x01 translation considers the 5 bits
				data = 1
			}

			val := data & 0x1F

			if mbc.romBankNBits > 5 {
				mbc.romBankNum = (mbc.romBankNum & 0x60) | (val & mbc.romBankMask)
			} else {
				mbc.romBankNum = (val & mbc.romBankMask)
			}
			return
		}

		if inRange(addr, 0x4000, 0x5FFF) {
			// set either the upper 2 bits of rom bank register or select ram banks 0-3
			val := data & 0x3

			if mbc.mode == MODE0 {
				if mbc.cart.romSize >= 0x100000 {
					// only for 1 MiB or larger ROM carts
					mbc.romBankNum = (mbc.romBankNum & 0x1F) | (val << 5)
				}
			} else {
				if mbc.cart.ramSize == 0x8000 {
					// only for 32 Kib RAM carts
					mbc.ramBankNum = val
				}
			}

			return
		}

		if inRange(addr, 0x6000, 0x7FFF) {
			// set the banking mode
			if mbc.cart.ramSize <= 0x2000 && mbc.cart.romSize <= 0x80000 {
				return
			}

			val := data & 0x1
			mbc.mode = BankMode(val)

			if mbc.mode == MODE0 {
				// only ram bank 0 can be accessed
				mbc.ramBankNum = 0
			} else {
				// clear upper 2 bits
				mbc.romBankNum &= 0x1F
			}

			return
		}
	} else if inRange(addr, EXT_RAM_BASE, EXT_RAM_TOP) {
		if !mbc.ramEnabled {
			return
		}

		mbc.ram[(uint16(mbc.ramBankNum)*0x2000)+addr-EXT_RAM_BASE] = data
		return
	}

	log.Fatalf("MMU mapped an illegal write address: 0x%02x to MBC1", addr)
}

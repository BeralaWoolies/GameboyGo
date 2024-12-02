package gb

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
)

type MBC3 struct {
	cart         *Cart
	ramEnabled   bool
	timerEnabled bool
	currRTC      uint8
	rtcReg       [5]uint8
	numBanks     uint32
	romBankMask  uint32
	ramBankNum   uint32
	romLo        uint32
	mode         SelectMode
}

type SelectMode uint8

const (
	RAM_SELECT = 0
	RTC_SELECT = 1

	RTC_S  = 0x08
	RTC_M  = 0x09
	RTC_H  = 0x0A
	RTC_DL = 0x0B
	RTC_DH = 0x0C
)

func (mbc *MBC3) init(cart *Cart) {
	mbc.cart = cart
	mbc.romLo = 1
	mbc.mode = RAM_SELECT
	mbc.numBanks = uint32(cart.romSize / 0x4000)
	nBits := uint32(math.Log2(float64(mbc.numBanks)))
	mbc.romBankMask = bits.NBitMask(nBits)
	mbc.rtcReg = [5]uint8{}

	fmt.Println("MBC3 INFO:")
	fmt.Println("numBanks: ", mbc.numBanks)
	fmt.Println("Bits to address banks: ", nBits)
	fmt.Println("Rom bank mask: ", "0b"+strconv.FormatInt(int64(mbc.romBankMask), 2))
	fmt.Println("====================================")
}

func (mbc *MBC3) contains(address uint16) bool {
	return inRange(address, ROM_BASE, ROM_TOP) || inRange(address, EXT_RAM_BASE, EXT_RAM_TOP)
}

func (mbc *MBC3) read(addr uint16) uint8 {
	if inRange(addr, ROM_BASE, ROM_TOP) {
		if inRange(addr, ROM_BASE, 0x3FFF) {
			// bank 00 of rom
			return mbc.cart.rom[addr]
		}

		if inRange(addr, 0x4000, ROM_TOP) {
			// switchable bank of rom
			return mbc.cart.rom[((mbc.romLo&mbc.romBankMask)*0x4000)+uint32(addr-0x4000)]
		}
	} else if inRange(addr, EXT_RAM_BASE, EXT_RAM_TOP) {
		if mbc.mode == RAM_SELECT {
			if !mbc.ramEnabled {
				return 0xFF
			}

			bank := uint32(0)
			if mbc.bigRam() {
				bank = mbc.ramBankNum
			}

			return mbc.cart.ram[(uint32(bank)*0x2000+uint32(addr-EXT_RAM_BASE))%mbc.cart.ramSize]
		} else {
			if !mbc.timerEnabled {
				return 0xFF
			}

			return mbc.rtcReg[mbc.currRTC]
		}
	}

	log.Fatalf("MMU mapped an illegal read address: 0x%02x to MBC3", addr)
	return 0xFF
}

func (mbc *MBC3) write(addr uint16, data uint8) {
	if inRange(addr, ROM_BASE, ROM_TOP) {
		if inRange(addr, ROM_BASE, 0x1FFF) {
			// trap to modify ram enable register
			mbc.ramEnabled = (data&0xF == 0xA)
			mbc.timerEnabled = (data&0xF == 0xA)
			return
		}

		if inRange(addr, 0x2000, 0x3FFF) {
			// set lower 5 bits of rom bank register
			val := data & 0x7F
			if val == 0 {
				val = 1
			}

			mbc.romLo = uint32(val) & mbc.romBankMask
			return
		}

		if inRange(addr, 0x4000, 0x5FFF) {
			// set either the upper 2 bits of rom bank register or select ram banks 0-3
			if data <= 0x03 {
				mbc.ramBankNum = uint32(data & 0x3)
				mbc.mode = SelectMode(RAM_SELECT)
				return
			}

			if data >= RTC_S && data <= RTC_DH {
				mbc.currRTC = data - RTC_S
				mbc.mode = SelectMode(RTC_SELECT)
				return
			}

			return
		}

		if inRange(addr, 0x6000, 0x7FFF) {
			// TODO: latch clock data
			return
		}
	} else if inRange(addr, EXT_RAM_BASE, EXT_RAM_TOP) {
		if mbc.mode == RAM_SELECT {
			if !mbc.ramEnabled {
				return
			}

			bank := uint32(0)
			if mbc.bigRam() {
				bank = mbc.ramBankNum
			}

			mbc.cart.ram[(uint32(bank)*0x2000+uint32(addr-EXT_RAM_BASE))%mbc.cart.ramSize] = data
			return
		} else {
			if !mbc.timerEnabled {
				return
			}

			mbc.rtcReg[mbc.currRTC] = data
			return
		}
	}

	log.Fatalf("MMU mapped an illegal write address: 0x%02x to MBC3", addr)
}

func (mbc *MBC3) bigRam() bool {
	// RAM can make use of the 2-bit register
	return mbc.cart.ramSize > 0x2000
}

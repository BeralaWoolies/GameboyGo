package gb

import (
	"fmt"
	"log"
	"os"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
)

type Cart struct {
	rom     []byte
	romSize uint32
	ramSize uint32

	mbc      *MBC1
	cartType uint8

	ram     bool
	battery bool
}

const (
	ROM_BASE = 0x0
	ROM_TOP  = 0x7FFF

	EXT_RAM_BASE = 0xA000
	EXT_RAM_TOP  = 0xBFFF
)

var cartTypes = map[int]string{
	0x00: "ROM ONLY",
	0x01: "MBC1",
	0x02: "MBC1+RAM",
	0x03: "MBC1+RAM+BATTERY",
	0x05: "MBC2",
	0x06: "MBC2+BATTERY",
	0x08: "ROM+RAM",
	0x09: "ROM+RAM+BATTERY",
	0x0B: "MMM01",
	0x0C: "MMM01+RAM",
	0x0D: "MMM01+RAM+BATTERY",
	0x0F: "MBC3+TIMER+BATTERY",
	0x10: "MBC3+TIMER+RAM+BATTERY",
	0x11: "MBC3",
	0x12: "MBC3+RAM",
	0x13: "MBC3+RAM+BATTERY",
	0x19: "MBC5",
	0x1A: "MBC5+RAM",
	0x1B: "MBC5+RAM+BATTERY",
	0x1C: "MBC5+RUMBLE",
	0x1D: "MBC5+RUMBLE+RAM",
	0x1E: "MBC5+RUMBLE+RAM+BATTERY",
	0x20: "MBC6",
	0x22: "MBC7+SENSOR+RUMBLE+RAM+BATTERY",
	0xFC: "POCKET CAMERA",
	0xFD: "BANDAI TAMA5",
	0xFE: "HuC3",
	0xFF: "HuC1+RAM+BATTERY",
}

var destCodes = map[int]string{
	0x00: "Japan (and possibly overseas)",
	0x01: "Overseas only",
}

var ramSizes = map[uint8]uint32{
	0x00: 0,       // 0       bytes
	0x01: 0x800,   // 2048    bytes (unofficial)
	0x02: 0x2000,  // 8192    bytes
	0x03: 0x8000,  // 32768   bytes
	0x04: 0x20000, // 131072  bytes
	0x05: 0x10000, // 65536   bytes
}

func (c *Cart) load(filename string) {
	rom, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	c.rom = rom
	c.mbc = &MBC1{}

	c.parseHeader()
	c.mbc.init(c)
}

func (c *Cart) parseHeader() {
	if len(c.rom) < 0x150 {
		log.Fatal("Invalid gameboy cartridge")
	}

	c.cartType = c.rom[0x0147]
	c.romSize = 32 * (1 << c.rom[0x0148]) * 1024
	c.ramSize = ramSizes[c.rom[0x0149]]
	c.ram = c.ramSize != 0

	fmt.Println("========= Cartridge Header =========")
	fmt.Println("Title: ", string(c.rom[0x0134:0x0144]))
	fmt.Println("Manufacturer code: ", string(c.rom[0x013F:0x0143]))
	fmt.Println("CGB Flag: ", bits.HexString(int(c.rom[0x0143])))
	fmt.Println("Licensee code: ", string(c.rom[0x0144:0x0146]))
	fmt.Println("SGB Flag: ", bits.HexString(int(c.rom[0x0146])))
	fmt.Println("Cartridge type: ", cartTypes[int(c.cartType)])
	fmt.Println("ROM size: ", c.romSize/1024, "KiB")
	fmt.Println("SRAM size: ", c.ramSize/1024, "KiB")
	fmt.Println("Destination code: ", destCodes[int(c.rom[0x014A])])
	fmt.Println("Old licensee code: ", bits.HexString(int(c.rom[0x014B])))
	fmt.Println("Mask ROM version number: ", bits.HexString(int(c.rom[0x014C])))
	fmt.Println("Header checksum: ", bits.HexString(int(c.rom[0x014D])))
	fmt.Println("Global checksum: ", bits.HexString(int((uint16(c.rom[0x014E])<<8)|uint16(c.rom[0x014F]))))
	fmt.Println("====================================")
}

func (c *Cart) contains(addr uint16) bool {
	if c.romOnly() {
		return inRange(addr, ROM_BASE, ROM_TOP)
	}

	return c.mbc.contains(addr)
}

func (c *Cart) read(addr uint16) uint8 {
	if c.romOnly() {
		return c.rom[addr]
	}

	return c.mbc.read(addr)
}

func (c *Cart) write(addr uint16, data uint8) {
	if c.romOnly() {
		return
	}

	c.mbc.write(addr, data)
}

func (c *Cart) romOnly() bool {
	return c.cartType == 0x00
}

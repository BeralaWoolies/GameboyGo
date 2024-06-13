package gb

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
	"github.com/edsrzf/mmap-go"
)

type Cart struct {
	rom         []byte
	ram         mmap.MMap
	romSize     uint32
	ramSize     uint32
	title       string
	savFilePath string

	mbc      MemoryBankController
	cartType uint8

	hasRam  bool
	battery bool
}

type MemoryBankController interface {
	Addressable
	init(cart *Cart)
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
	if len(c.rom) < 0x150 {
		log.Fatal("Invalid gameboy cartridge")
	}

	c.title = strings.ReplaceAll(string(c.rom[0x0134:0x0144]), "\x00", "")
	if c.title == "" {
		c.title = "Unknown Title"
	}
	c.cartType = c.rom[0x0147]
	c.battery = strings.Contains(strings.ToLower(cartTypes[int(c.cartType)]), "battery")
	c.romSize = 32 * (1 << c.rom[0x0148]) * 1024
	c.ramSize = ramSizes[c.rom[0x0149]]
	c.hasRam = c.ramSize != 0

	c.ram = make([]byte, c.ramSize, c.ramSize)
	if c.cartType >= 0x01 && c.cartType <= 0x03 {
		c.mbc = &MBC1{}
	} else if c.cartType >= 0x0F && c.cartType <= 0x13 {
		c.mbc = &MBC3{}
	}

	if c.battery {
		c.ram = c.loadSave(filename)
	}

	c.printHeader()
	if !c.romOnly() {
		c.mbc.init(c)
	}
}

func (c *Cart) loadSave(filename string) mmap.MMap {
	savDir := "saves"
	if err := os.MkdirAll(savDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	savFilename := fmt.Sprintf("%s.sav", strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename)))
	c.savFilePath = filepath.Join(savDir, savFilename)

	sav, err := os.OpenFile(c.savFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer sav.Close()

	if err := sav.Truncate(int64(c.ramSize)); err != nil {
		log.Fatal(err)
	}

	sram, err := mmap.Map(sav, mmap.RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}

	return sram
}

func (c *Cart) syncSave() {
	if !c.battery {
		return
	}

	c.ram.Flush()
	c.ram.Unmap()
	fmt.Printf("Flushed save to %s\n", c.savFilePath)
}

func (c *Cart) printHeader() {
	fmt.Println("========= Cartridge Header =========")
	fmt.Println("Title: ", c.title)
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

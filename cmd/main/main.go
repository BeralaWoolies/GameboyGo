package main

import (
	"fmt"
	"os"

	"github.com/BeralaWoolies/GameboyGo/pkg/gb"
)

func hexDump(bytes []byte) {
	for _, b := range bytes {
		fmt.Printf("%02x ", b)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Printf("Usage: go run %s rom-file [boot-rom]\n", os.Args[0])
		os.Exit(1)
	}

	gameboy := gb.NewGameboy()
	gameboy.Start()

	// boot_rom, err := os.ReadFile("boot_rom.bin")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// hexDump(boot_rom)
}

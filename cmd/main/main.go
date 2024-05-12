package main

import (
	"log"
	"os"

	"github.com/BeralaWoolies/GameboyGo/pkg/gb"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalf("Usage: go run %s myrom.gb", os.Args[0])
	}

	rom := args[0]
	gameboy := gb.NewGameboy(rom)
	gameboy.Start()
}

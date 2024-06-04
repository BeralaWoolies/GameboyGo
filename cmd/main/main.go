package main

import (
	"flag"

	"github.com/BeralaWoolies/GameboyGo/pkg/gb"
)

func main() {
	rom := flag.String("rom", "", "specify a .gb rom")
	debugMode := flag.Bool("d", false, "specify -d to enable debug mode")
	flag.Parse()

	gameboy := gb.NewGameboy(*rom, *debugMode)
	gameboy.Start()
}

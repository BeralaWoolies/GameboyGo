package main

import (
	"flag"
	"log"
	"os"

	"github.com/BeralaWoolies/GameboyGo/pkg/gb"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalf("Usage: go run %s myrom.gb", os.Args[0])
	}

	rom := flag.String("rom", "", "specify a .gb rom")
	debugMode := flag.Bool("d", false, "specify -d to enable debug mode")
	fpsCounter := flag.Bool("fps", false, "specify -fps to show FPS counter")
	flag.Parse()

	gameboy := gb.NewGameboy(*rom, *debugMode, *fpsCounter)
	gameboy.Start()
}

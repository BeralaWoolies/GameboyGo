package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/BeralaWoolies/GameboyGo/pkg/gb"
)

var cpuprofile *string = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile *string = flag.String("memprofile", "", "write memory profile to `file`")
var rom *string = flag.String("rom", "", "must specify a .gb or .gbc rom")
var bootrom *string = flag.String("bootrom", "", "optionally play the boot rom")
var debugMode *bool = flag.Bool("d", false, "optionally enable debug mode")

func main() {
	parseArgs()

	gameboy := gb.NewGameboy(gb.GameboyOptions{
		Filename:        *rom,
		DebugMode:       *debugMode,
		BootRomFilename: *bootrom,
	})
	gameboy.Start()
}

func parseArgs() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	if *rom == "" {
		fmt.Println("Must specify a .gb or .gbc rom")
		flag.Usage()
		os.Exit(2)
	}
}

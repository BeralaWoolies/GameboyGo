package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/BeralaWoolies/GameboyGo/pkg/gb"
)

var cpuprofile *string = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile *string = flag.String("memprofile", "", "write memory profile to `file`")
var rom *string = flag.String("rom", "", "specify a .gb rom")
var debugMode *bool = flag.Bool("d", false, "specify -d to enable debug mode")

func main() {
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

	gameboy := gb.NewGameboy(*rom, *debugMode)
	gameboy.Start()
}

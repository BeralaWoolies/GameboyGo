package gb

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Gameboy struct {
	mmu   *MMU
	cpu   *CPU
	ppu   *PPU
	timer *Timer
	dmac  *DMAController
	ic    *IntruptController
}

const (
	CLOCK_SPEED       = 4194304
	FPS               = 60
	C_TICKS_PER_FRAME = CLOCK_SPEED / FPS

	SCREEN_WIDTH  = GB_SCREEN_WIDTH + DEBUG_SCREEN_WIDTH
	SCREEN_HEIGHT = max(GB_SCREEN_HEIGHT, DEBUG_SCREEN_HEIGHT)

	GB_SCREEN_WIDTH  = 160
	GB_SCREEN_HEIGHT = 144

	DEBUG_SCREEN_WIDTH  = 128
	DEBUG_SCREEN_HEIGHT = 192

	WINDOW_WIDTH  = SCREEN_WIDTH * 3
	WINDOW_HEIGHT = SCREEN_HEIGHT * 3
)

func NewGameboy(filename string) *Gameboy {
	gb := &Gameboy{}
	gb.init(filename)
	return gb
}

func (gb *Gameboy) init(filename string) {
	gb.initHardware()
	gb.initMemoryMap(filename)
}

func (gb *Gameboy) initHardware() {
	gb.mmu = &MMU{}
	gb.cpu = &CPU{}
	gb.ppu = &PPU{}
	gb.timer = &Timer{}
	gb.dmac = &DMAController{}
	gb.ic = &IntruptController{}

	gb.cpu.init(gb.mmu)
	gb.ppu.init(gb.mmu, gb.dmac, gb.ic)
	gb.dmac.init(gb.mmu)
	gb.timer.init(gb.mmu, gb.ic)
	gb.ic.init(gb.mmu, gb.cpu)
}

func (gb *Gameboy) initMemoryMap(filename string) {
	gb.mmu.mapAddrSpace(newBootROM("boot_rom.bin", gb.mmu))
	gb.mmu.mapAddrSpace(newROM(filename))
	gb.mmu.mapAddrSpace(gb.ppu)
	gb.mmu.mapAddrSpace(gb.timer)
	gb.mmu.mapAddrSpace(gb.ic)

	// for now have our generic RAM be last in precedence to "catch" unimplemented addresses
	gb.mmu.mapAddrSpace(newGenericRAM())

	// manually disable boot rom for now
	gb.mmu.write(BOOT_ROM_ENABLE_ADDR, 1)
}

func (gb *Gameboy) printRegisters() {
	fmt.Printf("A: 0x%02x | %d\n", gb.cpu.reg.A, gb.cpu.reg.A)
	fmt.Printf("B: 0x%02x | %d\n", gb.cpu.reg.B, gb.cpu.reg.B)
	fmt.Printf("C: 0x%02x | %d\n", gb.cpu.reg.C, gb.cpu.reg.C)
	fmt.Printf("D: 0x%02x | %d\n", gb.cpu.reg.D, gb.cpu.reg.D)
	fmt.Printf("E: 0x%02x | %d\n", gb.cpu.reg.E, gb.cpu.reg.E)
	fmt.Printf("F: 0x%02x | %d\n", gb.cpu.reg.F, gb.cpu.reg.F)
	fmt.Printf("H: 0x%02x | %d\n", gb.cpu.reg.H, gb.cpu.reg.H)
	fmt.Printf("L: 0x%02x | %d\n", gb.cpu.reg.L, gb.cpu.reg.L)
	fmt.Printf("SP: 0x%04x | %d\n", gb.cpu.reg.SP, gb.cpu.reg.SP)
	fmt.Printf("PC: 0x%04x | %d\n", gb.cpu.reg.PC, gb.cpu.reg.PC)
	fmt.Println()
}

func (gb *Gameboy) Start() {
	ebiten.SetWindowSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	ebiten.SetWindowTitle("GameboyGo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(gb); err != nil {
		log.Fatal(err)
	}
}

func (gb *Gameboy) Update() error {
	for cTicksThisUpdate := 0; cTicksThisUpdate < C_TICKS_PER_FRAME; {
		cTicks := 4
		if !gb.cpu.halted {
			cTicks = gb.cpu.step()
		}

		cTicksThisUpdate += cTicks
		gb.ppu.step(cTicks)
		gb.timer.step(cTicks)
		gb.dmac.step(cTicks)
		cTicksThisUpdate += gb.ic.handleIntrupts()
	}

	return nil
}

func (gb *Gameboy) Draw(screen *ebiten.Image) {
	gb.ppu.updateDebugScreen(screen, GB_SCREEN_WIDTH, 0)
}

func (gb *Gameboy) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

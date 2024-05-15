package gb

import (
	"fmt"
	"log"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
	"github.com/hajimehoshi/ebiten/v2"
)

type Gameboy struct {
	mmu   *MMU
	cpu   *CPU
	ppu   *PPU
	dmac  *DMAC
	timer *Timer
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

	ISR_CLOCK_TICKS = 20
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
	gb.dmac = &DMAC{}
	gb.timer = &Timer{}

	gb.cpu.init(gb.mmu)
	gb.ppu.init(gb.mmu)
	gb.dmac.init(gb.mmu)
	gb.timer.init(gb.mmu)
}

func (gb *Gameboy) initMemoryMap(filename string) {
	gb.mmu.mapAddrSpace(newBootROM("boot_rom.bin", gb.mmu))
	gb.mmu.mapAddrSpace(newROM(filename))
	gb.mmu.mapAddrSpace(gb.ppu)
	gb.mmu.mapAddrSpace(gb.dmac)
	gb.mmu.mapAddrSpace(gb.timer)

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
		gb.timer.step(cTicks)
		gb.dmac.step(cTicks)
		cTicksThisUpdate += gb.handleIntrupts()
	}

	return nil
}

func (gb *Gameboy) Draw(screen *ebiten.Image) {
	gb.ppu.updateDebugScreen(screen, GB_SCREEN_WIDTH, 0)
}

func (gb *Gameboy) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func (gb *Gameboy) handleIntrupts() int {
	if gb.cpu.IMEDelay {
		gb.cpu.clearIMEDelay()
		gb.cpu.setIME()
		return 0
	}

	if !gb.cpu.IME && !gb.cpu.halted {
		return 0
	}

	IE := gb.readIE()
	IF := gb.readIF()

	if bits.IsSetInBoth(IE, IF, VBLANK_INTRUPT_BIT) {
		gb.serviceIntrupt(VBLANK_INTRUPT_BIT)
	} else if bits.IsSetInBoth(IE, IF, LCD_INTRUPT_BIT) {
		gb.serviceIntrupt(LCD_INTRUPT_BIT)
	} else if bits.IsSetInBoth(IE, IF, TIMER_INTRUPT_BIT) {
		gb.serviceIntrupt(TIMER_INTRUPT_BIT)
	} else if bits.IsSetInBoth(IE, IF, SERIAL_INTRUPT_BIT) {
		gb.serviceIntrupt(SERIAL_INTRUPT_BIT)
	} else if bits.IsSetInBoth(IE, IF, JOYPAD_INTRUPT_BIT) {
		gb.serviceIntrupt(JOYPAD_INTRUPT_BIT)
	} else {
		return 0
	}

	return ISR_CLOCK_TICKS
}

func (gb *Gameboy) serviceIntrupt(intruptBit uint8) {
	if !gb.cpu.IME && gb.cpu.halted {
		gb.cpu.exitHaltedState()
		return
	}

	gb.cpu.exitHaltedState()
	gb.cpu.clearIME()
	gb.clearIFBit(intruptBit)
	gb.cpu.pushStack(gb.cpu.reg.PC)

	switch intruptBit {
	case VBLANK_INTRUPT_BIT:
		gb.cpu.setPC(VBLANK_INTRUPT_VEC)
	case LCD_INTRUPT_BIT:
		gb.cpu.setPC(STAT_INTRUPT_VEC)
	case TIMER_INTRUPT_BIT:
		gb.cpu.setPC(TIMER_INTRUPT_VEC)
	case SERIAL_INTRUPT_BIT:
		gb.cpu.setPC(SERIAL_INTRUPT_VEC)
	case JOYPAD_INTRUPT_BIT:
		gb.cpu.setPC(JOYPAD_INTRUPT_VEC)
	}
}

func (gb *Gameboy) clearIFBit(pos uint8) {
	gb.mmu.write(IF_ADDR, bits.Reset(gb.readIF(), pos))
}

func (gb *Gameboy) readIE() uint8 {
	return gb.mmu.read(IE_ADDR)
}

func (gb *Gameboy) readIF() uint8 {
	return gb.mmu.read(IF_ADDR)
}

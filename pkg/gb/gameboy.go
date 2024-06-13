package gb

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type Gameboy struct {
	mmu          *MMU
	cpu          *CPU
	ppu          *PPU
	joyp         *Joypad
	serial       *SerialPort
	timer        *Timer
	cart         *Cart
	dmac         *DMAController
	ic           *IntruptController
	btnMappings  map[ebiten.Key]func(pressed bool)
	opts         GameboyOptions
	screenWidth  int
	screenHeight int
	windowWidth  int
	windowHeight int
	speedUp      bool
}

type GameboyOptions struct {
	Filename        string
	DebugMode       bool
	BootRomFilename string
	Stats           bool
}

const (
	FPS             = 60
	TICKS_PER_FRAME = TICKS_PER_SCANLINE * SCANLINES_PER_FRAME

	GB_SCREEN_WIDTH  = 160
	GB_SCREEN_HEIGHT = 144

	TILE_DATA_SCREEN_WIDTH  = 128
	TILE_DATA_SCREEN_HEIGHT = 192

	TILE_MAP_SCREEN_WIDTH  = 256
	TILE_MAP_SCREEN_HEIGHT = 256
)

func NewGameboy(opts GameboyOptions) *Gameboy {
	gb := &Gameboy{opts: opts}
	gb.init(opts.Filename)

	if gb.opts.DebugMode {
		gb.screenWidth = GB_SCREEN_WIDTH + TILE_DATA_SCREEN_WIDTH + (2 * TILE_MAP_SCREEN_WIDTH)
		gb.screenHeight = max(GB_SCREEN_HEIGHT, TILE_DATA_SCREEN_HEIGHT, TILE_MAP_SCREEN_HEIGHT)
		gb.windowWidth = gb.screenWidth * 3
		gb.windowHeight = gb.screenHeight * 3
	} else {
		gb.screenWidth = GB_SCREEN_WIDTH
		gb.screenHeight = GB_SCREEN_HEIGHT
		gb.windowWidth = gb.screenWidth * 4
		gb.windowHeight = gb.screenHeight * 4
	}

	return gb
}

func (gb *Gameboy) init(filename string) {
	gb.initHardware(filename)
	gb.initMemoryMap()
	gb.bindUIEvents()
}

func (gb *Gameboy) initHardware(filename string) {
	gb.mmu = &MMU{}
	gb.cpu = &CPU{}
	gb.ppu = &PPU{}
	gb.joyp = &Joypad{}
	gb.serial = &SerialPort{}
	gb.timer = &Timer{}
	gb.cart = &Cart{}
	gb.dmac = &DMAController{}
	gb.ic = &IntruptController{}

	gb.cpu.init(gb.mmu)
	gb.ppu.init(gb.mmu, gb.dmac, gb.ic)
	gb.joyp.init(gb.ic)
	gb.serial.init(gb.ic)
	gb.timer.init(gb.mmu, gb.ic)
	gb.cart.load(filename)
	gb.dmac.init(gb.mmu)
	gb.ic.init(gb.mmu, gb.cpu)
}

func (gb *Gameboy) initMemoryMap() {
	if gb.hasBootRom() {
		gb.mmu.mapAddrSpace(newBootROM(gb.opts.BootRomFilename, gb.mmu))
	}
	gb.mmu.mapAddrSpace(gb.cart)
	gb.mmu.mapAddrSpace(gb.ppu)
	gb.mmu.mapAddrSpace(gb.joyp)
	gb.mmu.mapAddrSpace(gb.serial)
	gb.mmu.mapAddrSpace(gb.timer)
	gb.mmu.mapAddrSpace(gb.ic)

	// for now have our generic RAM be last in precedence to "catch" unimplemented addresses
	gb.mmu.mapAddrSpace(newGenericRAM())
}

func (gb *Gameboy) hasBootRom() bool {
	return gb.opts.BootRomFilename != ""
}

func (gb *Gameboy) bindUIEvents() {
	gb.btnMappings = map[ebiten.Key]func(pressed bool){
		ebiten.KeyArrowUp:    gb.joyp.up.press,
		ebiten.KeyArrowDown:  gb.joyp.down.press,
		ebiten.KeyArrowRight: gb.joyp.right.press,
		ebiten.KeyArrowLeft:  gb.joyp.left.press,
		ebiten.KeyA:          gb.joyp.a.press,
		ebiten.KeyS:          gb.joyp.b.press,
		ebiten.KeySpace:      gb.joyp.sel.press,
		ebiten.KeyEnter:      gb.joyp.start.press,
		ebiten.KeyD:          gb.setSpeed,
	}
}

func (gb *Gameboy) setSpeed(speedUp bool) {
	gb.speedUp = speedUp

	if gb.speedUp {
		ebiten.SetTPS(FPS * 2)
	} else {
		ebiten.SetTPS(FPS)
	}
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
	fmt.Println("Starting...")

	ebiten.SetWindowSize(gb.windowWidth, gb.windowHeight)
	ebiten.SetWindowTitle(fmt.Sprintf("GameboyGo - %s", gb.cart.title))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(FPS)
	if !gb.opts.DebugMode {
		ebiten.SetVsyncEnabled(true)
	} else {
		ebiten.SetVsyncEnabled(false)
	}

	defer gb.cart.syncSave()

	if !gb.hasBootRom() {
		gb.powerUpSequence()
	}

	if err := ebiten.RunGame(gb); err != nil {
		log.Fatal(err)
	}
}

func (gb *Gameboy) Update() error {
	gb.handleUIEvents()

	for gb.cpu.ticks < TICKS_PER_FRAME {
		ticksThisUpdate := 4
		if !gb.cpu.halted {
			ticksThisUpdate = gb.cpu.step()
		}

		gb.ppu.step(ticksThisUpdate)
		gb.timer.step(ticksThisUpdate)
		gb.dmac.step(ticksThisUpdate)

		gb.cpu.ticks += (ticksThisUpdate + gb.ic.handleIntrupts())
	}

	gb.cpu.ticks -= TICKS_PER_FRAME

	return nil
}

func (gb *Gameboy) handleUIEvents() {
	for kbKey, press := range gb.btnMappings {
		press(ebiten.IsKeyPressed(kbKey))
	}
}

func (gb *Gameboy) powerUpSequence() {
	// cpu registers
	gb.cpu.setA(0x01)
	gb.cpu.setFlag(ZERO_FLAG_BIT, true)
	gb.cpu.setFlag(SUB_FLAG_BIT, false)
	gb.cpu.setFlag(HALF_CARRY_FLAG_BIT, true)
	gb.cpu.setFlag(CARRY_FLAG_BIT, true)
	gb.cpu.setB(0x00)
	gb.cpu.setC(0x13)
	gb.cpu.setD(0x00)
	gb.cpu.setE(0xD8)
	gb.cpu.setH(0x01)
	gb.cpu.setL(0x4D)
	gb.cpu.setPC(0x0100)
	gb.cpu.setSP(0xFFFE)

	// IO registers
	gb.mmu.write(JOYP_ADDR, 0xCF)
	gb.mmu.write(SB_ADDR, 0x00)
	gb.mmu.write(SC_ADDR, 0x7E)
	gb.timer.setDivDirect(0xAB) // set directly, since writes resets DIV
	gb.mmu.write(TIMA_ADDR, 0x00)
	gb.mmu.write(TMA_ADDR, 0x00)
	gb.mmu.write(TAC_ADDR, 0xF8)
	gb.mmu.write(IF_ADDR, 0xE1)

	// audio registers
	gb.mmu.write(0xFF10, 0x80)
	gb.mmu.write(0xFF11, 0xBF)
	gb.mmu.write(0xFF12, 0xF3)
	gb.mmu.write(0xFF13, 0xFF)
	gb.mmu.write(0xFF14, 0xBF)
	gb.mmu.write(0xFF16, 0x3F)
	gb.mmu.write(0xFF17, 0x00)
	gb.mmu.write(0xFF18, 0xFF)
	gb.mmu.write(0xFF19, 0xBF)
	gb.mmu.write(0xFF1A, 0x7F)
	gb.mmu.write(0xFF1B, 0xFF)
	gb.mmu.write(0xFF1C, 0x9F)
	gb.mmu.write(0xFF1D, 0xFF)
	gb.mmu.write(0xFF1E, 0xBF)
	gb.mmu.write(0xFF20, 0xFF)
	gb.mmu.write(0xFF21, 0x00)
	gb.mmu.write(0xFF22, 0x00)
	gb.mmu.write(0xFF23, 0xBF)
	gb.mmu.write(0xFF24, 0x77)
	gb.mmu.write(0xFF25, 0xF3)
	gb.mmu.write(0xFF26, 0xF1)

	// ppu registers
	gb.mmu.write(LCDC_ADDR, 0x91)
	gb.ppu.setStatDirect(0x85)
	gb.mmu.write(SCY_ADDR, 0x00)
	gb.mmu.write(SCX_ADDR, 0x00)
	gb.mmu.write(LY_ADDR, 0x00)
	gb.mmu.write(LYC_ADDR, 0x00)
	gb.ppu.setDMADirect(0xFF)
	gb.mmu.write(BG_PALETTE_ADDR, 0xFC)
	gb.mmu.write(WY_ADDR, 0x00)
	gb.mmu.write(WX_ADDR, 0x00)

	fmt.Println("Finished power up sequence...")
}

func (gb *Gameboy) Draw(screen *ebiten.Image) {
	gb.updateWindow()

	if !gb.opts.DebugMode {
		gb.ppu.updateGBScreen(screen, &ebiten.DrawImageOptions{})
	} else {
		opt := ebiten.DrawImageOptions{}
		dbgOpt := ebiten.DrawImageOptions{}

		opt.GeoM.Translate(0, (TILE_DATA_SCREEN_HEIGHT-GB_SCREEN_HEIGHT)/2)
		gb.ppu.updateGBScreen(screen, &opt)

		dbgOpt.GeoM.Translate(GB_SCREEN_WIDTH, 0)
		gb.ppu.updateTileDataScreen(screen, &dbgOpt)

		dbgOpt.GeoM.Translate(TILE_DATA_SCREEN_WIDTH, 0)
		gb.ppu.updateTileMaps(screen, &dbgOpt)
	}
}

func (gb *Gameboy) updateWindow() {
	emu := fmt.Sprintf("GameboyGo - %s", gb.cart.title)

	stats := ""
	if gb.opts.Stats {
		stats = fmt.Sprintf("(FPS: %s, SPEED: 1x)", strconv.Itoa(int(ebiten.ActualFPS())))
		if gb.speedUp {
			stats = strings.Replace(stats, "1x", "2x", 1)
		}
	}

	ebiten.SetWindowTitle(strings.Join([]string{emu, stats}, " "))
}

func (gb *Gameboy) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return gb.screenWidth, gb.screenHeight
}

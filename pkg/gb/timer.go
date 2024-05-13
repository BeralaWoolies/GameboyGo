package gb

import (
	"log"

	"github.com/BeralaWoolies/GameboyGo/pkg/bits"
)

type Timer struct {
	div         uint8
	tima        uint8
	tma         uint8
	tac         uint8
	divCounter  int
	timaCounter int
	mmu         *MMU
}

const (
	DIV_ADDR  = 0xFF04
	TIMA_ADDR = 0xFF05
	TMA_ADDR  = 0xFF06
	TAC_ADDR  = 0xFF07

	TAC_TIMER_ENABLE_BIT = 2
	TAC_FREQ_DIV_MSK     = 0x3
	TAC_MSK              = 0x7

	HZ_4096   = 0
	HZ_262144 = 1
	HZ_65536  = 2
	HZ_16386  = 3
)

func (t *Timer) init(mmu *MMU) {
	t.div = 0xAB
	t.tima = 0x0
	t.tma = 0x0
	t.tac = 0x0
	t.divCounter = 0xCC
	t.timaCounter = 0x0
	t.mmu = mmu
}

func (t *Timer) contains(addr uint16) bool {
	switch addr {
	case DIV_ADDR:
		return true
	case TIMA_ADDR:
		return true
	case TMA_ADDR:
		return true
	case TAC_ADDR:
		return true
	default:
		return false
	}
}

func (t *Timer) write(addr uint16, data uint8) {
	switch addr {
	case DIV_ADDR:
		// not allowed to write to Divider Register
		t.div = 0
	case TIMA_ADDR:
		t.tima = data
	case TMA_ADDR:
		t.tma = data
	case TAC_ADDR:
		t.tac = data & TAC_MSK
	default:
		// mmu should never map an illegal address here
		log.Fatalf("MMU mapped an illegal write address: 0x%02x to Timer", addr)
	}
}

func (t *Timer) read(addr uint16) uint8 {
	switch addr {
	case DIV_ADDR:
		return t.div
	case TIMA_ADDR:
		return t.tima
	case TMA_ADDR:
		return t.tma
	case TAC_ADDR:
		return t.tac & TAC_MSK
	default:
		// mmu should never map an illegal address here
		log.Fatalf("MMU mapped an illegal read address: 0x%02x to Timer", addr)
	}

	return 0xFF
}

func (t *Timer) step(cTicks int) {
	t.stepDIV(cTicks)
	if t.timerEnabled() {
		t.stepTIMA(cTicks)
	}
}

func (t *Timer) stepDIV(cTicks int) {
	t.divCounter += cTicks
	if t.divCounter >= 256 {
		t.divCounter -= 256
		t.div++
	}
}

func (t *Timer) stepTIMA(cTicks int) {
	t.timaCounter += cTicks

	ovflowTicks := 0
	switch freqDiv := t.tac & TAC_FREQ_DIV_MSK; freqDiv {
	case HZ_4096:
		ovflowTicks = 1024
	case HZ_262144:
		ovflowTicks = 16
	case HZ_65536:
		ovflowTicks = 64
	case HZ_16386:
		ovflowTicks = 256
	}

	for t.timaCounter >= ovflowTicks {
		t.timaCounter -= ovflowTicks
		t.tima++

		if t.tima == 0x00 {
			t.mmu.write(IF_ADDR, bits.Set(t.mmu.read(IF_ADDR), TIMER_INTRUPT_BIT))
			t.tima = t.tma
		}
	}
}

func (t *Timer) timerEnabled() bool {
	return bits.IsSet(t.tac, TAC_TIMER_ENABLE_BIT)
}

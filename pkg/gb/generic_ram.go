package gb

type GenericRAM struct {
	memory [RAM_SIZE]uint8 // generic ram that covers 0xA000 - 0xFFFF
}

const (
	RAM_SIZE = 0x6000
	RAM_BASE = 0xA000
	RAM_TOP  = 0xFFFF
)

func newGenericRAM() *GenericRAM {
	r := &GenericRAM{}
	r.init()

	return r
}

func (r *GenericRAM) init() {
	r.memory = [RAM_SIZE]uint8{}
	r.memory[0xFF10-RAM_BASE] = 0x80
	r.memory[0xFF11-RAM_BASE] = 0xBF
	r.memory[0xFF12-RAM_BASE] = 0xF3
	r.memory[0xFF14-RAM_BASE] = 0xBF
	r.memory[0xFF16-RAM_BASE] = 0x3F
	r.memory[0xFF17-RAM_BASE] = 0x00
	r.memory[0xFF19-RAM_BASE] = 0xBF
	r.memory[0xFF1A-RAM_BASE] = 0x7F
	r.memory[0xFF1B-RAM_BASE] = 0xFF
	r.memory[0xFF1C-RAM_BASE] = 0x9F
	r.memory[0xFF1E-RAM_BASE] = 0xBF
	r.memory[0xFF20-RAM_BASE] = 0xFF
	r.memory[0xFF21-RAM_BASE] = 0x00
	r.memory[0xFF22-RAM_BASE] = 0x00
	r.memory[0xFF23-RAM_BASE] = 0xBF
	r.memory[0xFF24-RAM_BASE] = 0x77
	r.memory[0xFF25-RAM_BASE] = 0xF3
	r.memory[0xFF26-RAM_BASE] = 0xF1
	r.memory[0xFF40-RAM_BASE] = 0x91
	r.memory[0xFF42-RAM_BASE] = 0x00
	r.memory[0xFF43-RAM_BASE] = 0x00
	r.memory[0xFF45-RAM_BASE] = 0x00
	r.memory[0xFF47-RAM_BASE] = 0xFC
	r.memory[0xFF48-RAM_BASE] = 0xFF
	r.memory[0xFF49-RAM_BASE] = 0xFF
	r.memory[0xFF4A-RAM_BASE] = 0x00
	r.memory[0xFF4B-RAM_BASE] = 0x00
}

func (r *GenericRAM) contains(addr uint16) bool {
	return inRange(addr, RAM_BASE, RAM_TOP)
}

func (r *GenericRAM) read(addr uint16) uint8 {
	return r.memory[addr-RAM_BASE]
}

func (r *GenericRAM) write(addr uint16, data uint8) {
	r.memory[addr-RAM_BASE] = data
}

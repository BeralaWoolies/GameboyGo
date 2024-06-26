package gb

type DMAController struct {
	mmu *MMU

	src      uint8
	active   bool
	currByte uint16
	delayed  bool
}

func (dmac *DMAController) init(mmu *MMU) {
	dmac.src = 0
	dmac.active = false
	dmac.currByte = 0
	dmac.delayed = false
	dmac.mmu = mmu
}

func (dmac *DMAController) step(cTicks int) {
	if !dmac.active {
		return
	}

	if dmac.delayed {
		cTicks -= 4
		dmac.delayed = false

		if cTicks <= 0 {
			return
		}
	}

	for i := 0; i < cTicks; i += 4 {
		dmac.transferOAM()
	}
}

func (dmac *DMAController) transferOAM() {
	if !dmac.active {
		return
	}

	dest := OAM_BASE + dmac.currByte
	src := dmac.mmu.read((uint16(dmac.src) * 0x100) + dmac.currByte)
	dmac.mmu.write(dest, src)

	dmac.currByte++
	dmac.active = dmac.currByte < OAM_SIZE
}

func (dmac *DMAController) initOAMTransfer(data uint8) {
	dmac.src = data
	dmac.active = true
	dmac.currByte = 0
	dmac.delayed = true
}

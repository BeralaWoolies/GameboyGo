package gb

import (
	"fmt"

	"github.com/BeralaWoolies/GameboyGo/pkg/queue"
)

type PixelFIFO struct {
	ppu *PPU

	ticks int

	currState       PixelFIFOState
	bgFIFO          *queue.Queue[PixelFIFOItem]
	spriteFIFO      *queue.Queue[PixelFIFOItem]
	spriteFetchFIFO *queue.Queue[Sprite]

	tileLoByte      uint8
	tileHiByte      uint8
	cacheTileLoByte uint8
	cacheTileHiByte uint8

	tileId        uint8
	tileMapY      uint16
	tileLine      uint16
	cacheTileLine uint16
	bgFetchX      uint8
	winFetchX     uint8

	BGWinComplete bool
	spriteFetch   bool
	sprite        Sprite

	windowFetch bool
	scxDropped  uint8
}
type PixelFIFOState uint8

type PixelFIFOItem struct {
	color      uint8
	bgPriority bool // only applies for sprites
	palette    Palette
}
type Palette uint8
type PixelType uint8

const (
	ReadTileID     PixelFIFOState = 0
	ReadTileDataLo PixelFIFOState = 1
	ReadTileDataHi PixelFIFOState = 2
	IDLE           PixelFIFOState = 3
	PushFIFO       PixelFIFOState = 4

	BGP  Palette = 0
	OBP0 Palette = 1
	OBP1 Palette = 2

	SPRITE PixelType = 0
	BG     PixelType = 1
	WINDOW PixelType = 2
)

func (pxF *PixelFIFO) init(ppu *PPU) {
	pxF.ppu = ppu
	pxF.bgFIFO = queue.New[PixelFIFOItem]()
	pxF.spriteFIFO = queue.New[PixelFIFOItem]()
	pxF.spriteFetchFIFO = queue.New[Sprite]()
}

func (pxF *PixelFIFO) setState(state PixelFIFOState) {
	pxF.currState = state
}

func (pxF *PixelFIFO) tick() {
	pxF.ticks++

	if pxF.ticks&1 == 1 {
		return
	}

	switch pxF.currState {
	case ReadTileID:
		if pxF.spriteFetch {
			pxF.sprite = pxF.spriteFetchFIFO.Remove()
			pxF.tileId = pxF.sprite.tileId
			pxF.tileLine = uint16(pxF.ppu.ly-pxF.sprite.y-16) % TILE_WIDTH

			if pxF.ppu.getSpriteHeight() == 16 {
				if pxF.sprite.flippedY() {
					if pxF.ppu.ly <= pxF.sprite.y-8 {
						pxF.tileId |= 0x01
					} else {
						pxF.tileId &= 0xFE
					}
				} else {
					if pxF.ppu.ly <= pxF.sprite.y-8 {
						pxF.tileId &= 0xFE
					} else {
						pxF.tileId |= 0x01
					}
				}
			}
		} else {
			var tileMapYOff uint16
			var tileMapXOff uint16
			var tileMapAddr uint16

			if pxF.windowFetch {
				tileMapYOff = (uint16(pxF.ppu.wly) / TILE_WIDTH) * TILE_MAP_WIDTH
				tileMapXOff = uint16(pxF.winFetchX) % TILE_MAP_WIDTH
				tileMapAddr = pxF.ppu.getWinTileMap()
			} else {
				tileMapYOff = (pxF.tileMapY / TILE_WIDTH) * TILE_MAP_WIDTH
				tileMapXOff = (uint16(pxF.bgFetchX) + (uint16(pxF.ppu.scx) / TILE_WIDTH)) % TILE_MAP_WIDTH
				tileMapAddr = pxF.ppu.getBGTileMap()
			}

			tileMapAddr += (tileMapYOff + tileMapXOff)

			pxF.tileId = pxF.ppu.vram[tileMapAddr-VRAM_BASE]
		}

		pxF.setState(ReadTileDataLo)
	case ReadTileDataLo:
		tileLine := pxF.tileLine

		if pxF.spriteFetch && pxF.sprite.flippedY() {
			tileLine = 7 - tileLine
		}

		pxF.tileLoByte, _ = pxF.getTileLine(pxF.tileId, tileLine)

		pxF.setState(ReadTileDataHi)
	case ReadTileDataHi:
		tileLine := pxF.tileLine

		if pxF.spriteFetch && pxF.sprite.flippedY() {
			tileLine = 7 - tileLine
		}

		_, pxF.tileHiByte = pxF.getTileLine(pxF.tileId, tileLine)

		pxF.setState(IDLE)
	case IDLE:
		pxF.setState(PushFIFO)
	case PushFIFO:
		// Check if sprite fetch is pending, attempt to finish up BG/Win fetch, and cache if we can't then switch into
		// sprite fetching mode
		if pxF.spriteFetch {
			for pxF.spriteFIFO.Length() < 8 {
				pxF.spriteFIFO.Add(PixelFIFOItem{
					color:      0,
					bgPriority: true,
					palette:    OBP0,
				})
			}

			for i := 0; i < 8; i++ {
				if pxF.spriteFIFO.Get(i).color == 0 {
					bit := 7 - i

					if pxF.sprite.flippedX() {
						bit = i
					}

					pxF.spriteFIFO.Replace(i, PixelFIFOItem{
						color:      getColor(pxF.tileLoByte, pxF.tileHiByte, uint8(bit)),
						bgPriority: pxF.sprite.getBGPriority(),
						palette:    pxF.sprite.getPalette(),
					})
				}
			}

			if pxF.sprite.x < 8 {
				for i := 0; i < int(8-pxF.sprite.x); i++ {
					pxF.spriteFIFO.Remove()
				}
			}

			// Check if we need to fetch another sprite
			pxF.setState(ReadTileID)
			if pxF.spriteFetchFIFO.IsEmpty() {
				// else we restart BG/Win fetch at either the next fetchX (if it successfully pushed)
				// or restart back at FIFO for it to continue to push to FIFO
				pxF.spriteFetch = false

				pxF.tileLoByte = pxF.cacheTileLoByte
				pxF.tileHiByte = pxF.cacheTileHiByte
				pxF.tileLine = pxF.cacheTileLine

				if !pxF.BGWinComplete {
					pxF.setState(PushFIFO)
				}
			}
		} else {
			pushed := false
			if pxF.bgFIFO.IsEmpty() {
				for bit := 7; bit >= 0; bit-- {
					if pxF.dropSCX() {
						continue
					}

					pxF.bgFIFO.Add(PixelFIFOItem{
						color:      getColor(pxF.tileLoByte, pxF.tileHiByte, uint8(bit)),
						bgPriority: false,
						palette:    BGP,
					})
				}

				pushed = true
				if pxF.windowFetch {
					pxF.winFetchX++
				} else {
					pxF.bgFetchX++
				}
			}

			if !pxF.spriteFetchFIFO.IsEmpty() {
				pxF.spriteFetch = true

				pxF.BGWinComplete = pushed
				pxF.cacheTileLoByte = pxF.tileLoByte
				pxF.cacheTileHiByte = pxF.tileHiByte
				pxF.cacheTileLine = pxF.tileLine
			}

			pxF.setState(ReadTileID)
		}
	}
}

func (pxF *PixelFIFO) start() {
	pxF.currState = ReadTileID
	pxF.bgFIFO.Clear()
	pxF.spriteFIFO.Clear()
	pxF.spriteFetchFIFO.Clear()
	pxF.spriteFetch = false
	pxF.tileMapY = uint16(pxF.ppu.ly + pxF.ppu.scy)
	pxF.tileLine = pxF.tileMapY % TILE_WIDTH
	pxF.bgFetchX = 0
	pxF.ticks = 0
	pxF.scxDropped = 0
}

func (pxF *PixelFIFO) dropSCX() bool {
	if pxF.scxDropped < pxF.ppu.scx%8 {
		pxF.scxDropped++
		return true
	}

	return false
}

func (pxF *PixelFIFO) startWinFetch() {
	pxF.currState = ReadTileID
	pxF.bgFIFO.Clear()
	pxF.winFetchX = 0
	pxF.windowFetch = true
	pxF.tileMapY = uint16(pxF.ppu.wly)
	pxF.tileLine = pxF.tileMapY % TILE_WIDTH
}

func (pxF *PixelFIFO) startBGFetch() {
	pxF.currState = ReadTileID
	pxF.bgFIFO.Clear()
	pxF.windowFetch = false
	pxF.tileMapY = uint16(pxF.ppu.ly + pxF.ppu.scy)
	pxF.tileLine = pxF.tileMapY % TILE_WIDTH
}

func (pxF *PixelFIFO) getTileLine(tileId uint8, tileLine uint16) (loByte uint8, hiByte uint8) {
	addr, unsig := pxF.ppu.getTileDataArea()
	addr -= VRAM_BASE

	if unsig {
		addr += (uint16(tileId) * TILE_SIZE) + (tileLine * 2)
	} else {
		addr += (uint16(int(int8(tileId)) * TILE_SIZE)) + (tileLine * 2)
	}

	return pxF.ppu.vram[addr], pxF.ppu.vram[addr+1]
}

func (pxF *PixelFIFO) pop() (PixelFIFOItem, error) {
	if pxF.spriteFetch || !pxF.spriteFetchFIFO.IsEmpty() {
		return PixelFIFOItem{}, fmt.Errorf("Pixel fetching is suspended")
	}

	if !pxF.bgFIFO.IsEmpty() {
		// mix here with spriteFIFO
		bgPixel := pxF.bgFIFO.Remove()

		if !pxF.ppu.spritesEnabled() || pxF.spriteFIFO.IsEmpty() {
			if !pxF.ppu.bgWinEnabled() {
				bgPixel.color = 0
			}

			return bgPixel, nil
		}

		return pxF.mix(bgPixel, pxF.spriteFIFO.Remove()), nil
	}

	return PixelFIFOItem{}, fmt.Errorf("FIFO is empty")
}

func (pxF *PixelFIFO) mix(bgPixel PixelFIFOItem, spritePixel PixelFIFOItem) PixelFIFOItem {
	if !pxF.ppu.bgWinEnabled() {
		return spritePixel
	}

	if (spritePixel.bgPriority && bgPixel.color != 0) || spritePixel.color == 0 {
		return bgPixel
	}

	return spritePixel
}

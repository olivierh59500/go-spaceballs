package demo

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

// BlocksEffect shares the contour/ring-mask renderer with the neighboring
// scenes, while retaining the authored corner palettes and two-plane shadow.
type BlocksEffect struct {
	model     source.BlockModel
	clock     *source.BlocksClock
	animation source.Animation
	masks     *maskRing
	grid      *composite.PaletteGrid
	colors    [2][17 * 15]color.NRGBA
	material  composite.PaletteGridState
	closed    bool
	dirty     bool
}

func NewBlocksEffect() (*BlocksEffect, error) {
	bytes, err := assets.Files.ReadFile("raw/second-effects.bin")
	if err != nil {
		return nil, err
	}
	model, err := source.ReadBlockModel(bytes)
	if err != nil {
		return nil, err
	}
	animation, err := source.LoadAnimation(assets.Files, "blocks")
	if err != nil {
		return nil, err
	}
	e, err := newBlocksRenderer()
	if err != nil {
		return nil, err
	}
	e.model, e.clock, e.animation = model, source.NewBlocksClock(model), animation
	return e, nil
}

func newBlocksRenderer() (*BlocksEffect, error) {
	black := color.NRGBA{A: 255}
	grid, err := composite.NewPaletteGrid(composite.PaletteGridConfig{
		Width: Width, Height: Height, Columns: 15, Rows: 17,
		X:                composite.PaletteGridAxis{CellSize: 24, Offset: 8},
		Y:                composite.PaletteGridAxis{CellSize: 16, FirstSpan: 36, Indices: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 14, 15, 16}},
		ControlThreshold: .5, BodyColor: &black,
	})
	if err != nil {
		return nil, err
	}
	e := &BlocksEffect{grid: grid, dirty: true}
	e.masks = newMaskRing()
	return e, nil
}

func (e *BlocksEffect) Update(f kit.Frame) error {
	e.SetTick(int(math.Round(f.Time * FPS)))
	return nil
}

func (e *BlocksEffect) SetTick(tick int) {
	tick = max(0, min(tick, source.BlocksTicks-1))
	if e.clock.Tick > tick+1 {
		e.clock = source.NewBlocksClock(e.model)
		e.masks.clear()
	}
	for e.clock.Tick <= tick {
		if e.clock.Step() {
			e.dirty = true
			e.masks.prepare(e.clock.Current, e.animation.Frames[e.clock.Frame], false)
		}
	}
}

func (e *BlocksEffect) Draw(dst *ebiten.Image) {
	e.drawGrid(dst, e.clock.Background, e.clock.Shadow, e.clock.Display)
}

func (e *BlocksEffect) drawGrid(dst *ebiten.Image, background, shadow [17][15]uint16, display [2]int) {
	if e.dirty {
		for row := 0; row < 17; row++ {
			for col := 0; col < 15; col++ {
				for bank, word := range [2]uint16{background[row][col], shadow[row][col]} {
					c := source.RGB12(word)
					e.colors[bank][row*15+col] = color.NRGBA{R: c.R, G: c.G, B: c.B, A: 255}
				}
			}
		}
		if err := e.grid.SetColors(e.colors[0][:], e.colors[1][:]); err != nil {
			panic(err)
		}
		e.dirty = false
	}
	if err := e.grid.Draw(dst, e.masks.planes[display[1]], e.masks.planes[display[0]], e.material); err != nil {
		panic(err)
	}
}

func (e *BlocksEffect) Close() error {
	if !e.closed {
		e.masks.close()
		e.grid.Close()
		e.closed = true
	}
	return nil
}

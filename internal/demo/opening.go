package demo

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

// OpeningEffect keeps the original held-plane and triple-buffer composition.
// It samples an absolute time; seeking backward resets its retained masks.
type OpeningEffect struct {
	clock   *source.OpeningClock
	banks   map[string]source.Animation
	planes  [3][2]*ebiten.Image
	bank    *composite.ContourBank
	lookup  *composite.BitplanePalette
	palette [4]color.NRGBA
	closed  bool
}

func NewOpeningEffect() (*OpeningEffect, error) {
	e := &OpeningEffect{clock: source.NewOpeningClock(), banks: make(map[string]source.Animation)}
	for _, name := range []string{"hands", "opening", "first"} {
		bank, err := source.LoadAnimation(assets.Files, name)
		if err != nil {
			return nil, err
		}
		e.banks[name] = bank
	}
	e.bank = newContourBank(Width, Height, 3, 2, ebiten.BlendXor)
	var err error
	e.lookup, err = composite.NewBitplanePalette(composite.BitplanePaletteConfig{Width: Width, Height: Height, Planes: 2, Palette: e.palette[:]})
	if err != nil {
		e.bank.Close()
		return nil, err
	}
	for i := range e.planes {
		for j := range e.planes[i] {
			e.planes[i][j] = e.bank.Image(i, j)
		}
	}
	return e, nil
}

func (e *OpeningEffect) Update(f kit.Frame) error {
	e.SetTick(int(math.Round(f.Time * FPS)))
	return nil
}

func (e *OpeningEffect) SetTick(tick int) {
	tick = max(0, min(tick, source.OpeningTicks-1))
	if e.clock.Tick > tick+1 {
		e.clock = source.NewOpeningClock()
		e.bank.Clear()
	}
	for e.clock.Tick <= tick {
		if e.clock.Step() {
			e.prepare()
		}
	}
}

func (e *OpeningEffect) prepare() {
	c := e.clock
	if c.SecondOnly {
		paintContours(e.bank, c.Current, 1, true, nil)
	} else {
		paintContours(e.bank, c.Current, 0, true, nil)
		if c.FillSecond {
			paintContours(e.bank, c.Current, 1, true, nil)
		}
	}
	if !c.Draw {
		return
	}
	frame := e.banks[c.Bank].Frames[c.Frame]
	point := func(p source.Point) ebiten.Vertex {
		x, y := int(p.X)*351/256, int(p.Y)*289/204
		if c.SecondOnly {
			x -= 15
			y = min(289, int(p.Y)*350/204+8)
		}
		return render.Vertex(float64(x), float64(y), 0, 0, color.White)
	}
	for plane := 0; plane < 2; plane++ {
		bit := plane
		if c.SecondOnly {
			bit--
		}
		if bit < 0 {
			continue
		}
		paintContours(e.bank, c.Current, plane, false, func(batch *render.Batch) {
			for _, polygon := range frame {
				if polygon.Mask>>uint(bit)&1 == 0 || len(polygon.Points) < 3 {
					continue
				}
				at := func(i int) ebiten.Vertex { return point(polygon.Points[i]) }
				if plane == 1 && !c.FillSecond {
					batch.StrokeContour(len(polygon.Points), at, render.ContourStroke{})
				} else {
					batch.Fan(len(polygon.Points), at)
				}
			}
		})
	}
}

func (e *OpeningEffect) Draw(dst *ebiten.Image) {
	for i, word := range e.clock.Palette {
		e.palette[i] = source.RGB12(word)
	}
	planes := [2]*ebiten.Image{e.planes[e.clock.Display[0]][0], e.planes[e.clock.Display[1]][1]}
	drawBitplanes(e.lookup, dst, planes[:], e.palette[:])
}

func (e *OpeningEffect) Close() error {
	if !e.closed {
		e.bank.Close()
		e.lookup.Close()
		e.closed = true
	}
	return nil
}

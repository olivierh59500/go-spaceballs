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

type RibbonEffect struct {
	model   source.RibbonModel
	clock   *source.RibbonClock
	poses   source.Animation
	planes  [4][2]*ebiten.Image
	white   *ebiten.Image
	batch   *render.Batch
	lookup  *composite.BitplanePalette
	palette [16]color.NRGBA
	closed  bool
}

func NewRibbonEffect() (*RibbonEffect, error) {
	data, err := assets.Files.ReadFile("raw/seventh-effects.bin")
	if err != nil {
		return nil, err
	}
	model, err := source.ReadRibbonModel(data)
	if err != nil {
		return nil, err
	}
	poses, err := source.LoadAnimation(assets.Files, "close-up")
	if err != nil {
		return nil, err
	}
	e := &RibbonEffect{model: model, clock: source.NewRibbonClock(model), poses: poses,
		white: ebiten.NewImage(1, 1), batch: render.NewBatch(4096)}
	e.lookup, err = composite.NewBitplanePalette(composite.BitplanePaletteConfig{
		Width: Width, Height: Height, Planes: 4, Palette: e.palette[:],
	})
	if err != nil {
		e.white.Deallocate()
		return nil, err
	}
	e.white.Fill(color.White)
	e.batch.Options.FillRule = ebiten.FillRuleEvenOdd
	for i := range e.planes {
		for j := range e.planes[i] {
			e.planes[i][j] = render.NewSurface(Width, Height)
		}
	}
	return e, nil
}

func (e *RibbonEffect) Update(f kit.Frame) error {
	tick := max(0, min(source.RibbonTicks-1, int(math.Round(f.Time*FPS))))
	if e.clock.Tick > tick+1 {
		e.clock = source.NewRibbonClock(e.model)
		for _, pair := range e.planes {
			for _, plane := range pair {
				plane.Clear()
			}
		}
	}
	for e.clock.Tick <= tick {
		old := e.clock.Current
		if e.clock.Step() {
			e.drawPose(old, e.clock.Frame, e.clock.Mirror)
		}
		if e.clock.Tick&1 == 0 {
			for _, plane := range e.planes[e.clock.Current] {
				plane.Clear()
			}
		}
	}
	return nil
}

func (e *RibbonEffect) drawPose(group, frame, mirror int) {
	point := func(p source.Point) ebiten.Vertex {
		x, y := int(p.X), int(p.Y)
		if mirror == 0 {
			y = 204 - y
		} else {
			x = 256 - x
		}
		return render.Vertex(float64(x*351/256), float64(y*289/204), 0, 0, color.White)
	}
	for bit := 0; bit < 2; bit++ {
		e.batch.Begin(e.planes[group][bit], e.white)
		for _, polygon := range e.poses.Frames[frame] {
			if polygon.Mask>>uint(bit)&1 == 0 || len(polygon.Points) < 3 {
				continue
			}
			first := point(polygon.Points[0])
			for i := 1; i+1 < len(polygon.Points); i++ {
				e.batch.Triangle(first, point(polygon.Points[i]), point(polygon.Points[i+1]))
			}
		}
		e.batch.Flush()
	}
}

func (e *RibbonEffect) Draw(dst *ebiten.Image) {
	for i, word := range e.clock.Palette {
		e.palette[i] = source.RGB12(word)
	}
	a, b := e.clock.Display[0], e.clock.Display[1]
	planes := [4]*ebiten.Image{e.planes[a][0], e.planes[a][1], e.planes[b][0], e.planes[b][1]}
	drawBitplanes(e.lookup, dst, planes[:], e.palette[:])
}

func (e *RibbonEffect) Close() error {
	if !e.closed {
		for _, pair := range e.planes {
			for _, plane := range pair {
				plane.Deallocate()
			}
		}
		e.white.Deallocate()
		e.lookup.Close()
		e.closed = true
	}
	return nil
}

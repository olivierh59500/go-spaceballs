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

// TrailsEffect draws the original retained-mask program, including its bank
// cues, clipped close-up lookup and periodically inverted 32-color palette.
type TrailsEffect struct {
	model   source.TrailModel
	clock   *source.TrailsClock
	banks   map[string]source.Animation
	masks   *maskRing
	lookup  *composite.BitplanePalette
	palette [32]color.NRGBA
	closed  bool
}

func NewTrailsEffect() (*TrailsEffect, error) {
	data, err := assets.Files.ReadFile("raw/first-effects.bin")
	if err != nil {
		return nil, err
	}
	model, err := source.ReadTrailModel(data)
	if err != nil {
		return nil, err
	}
	e := &TrailsEffect{model: model, clock: source.NewTrailsClock(model),
		banks: make(map[string]source.Animation)}
	for _, name := range []string{"trails-a", "trails-b", "trails-c", "trails-d"} {
		e.banks[name], err = source.LoadAnimation(assets.Files, name)
		if err != nil {
			return nil, err
		}
	}
	e.lookup, err = composite.NewBitplanePalette(composite.BitplanePaletteConfig{
		Width: Width, Height: Height, Planes: 5, Palette: e.palette[:],
	})
	if err != nil {
		return nil, err
	}
	e.masks = newMaskRing()
	return e, nil
}

func (e *TrailsEffect) Update(f kit.Frame) error {
	e.SetTick(int(math.Round(f.Time * FPS)))
	return nil
}

func (e *TrailsEffect) SetTick(tick int) {
	tick = max(0, min(tick, source.TrailsTicks-1))
	if e.clock.Tick > tick+1 {
		e.clock = source.NewTrailsClock(e.model)
		e.masks.clear()
	}
	for e.clock.Tick <= tick {
		if e.clock.Step() {
			var frame []source.Polygon
			if e.clock.Draw {
				frame = e.banks[e.clock.Bank].Frames[e.clock.Frame]
			}
			e.masks.prepare(e.clock.Current, frame, e.clock.Zoom)
		}
	}
}

func (e *TrailsEffect) Draw(dst *ebiten.Image) {
	var planes [5]*ebiten.Image
	for i, index := range e.clock.Display {
		planes[i] = e.masks.planes[index]
	}
	for i, word := range e.clock.Palette {
		e.palette[i] = source.RGB12(word)
	}
	drawBitplanes(e.lookup, dst, planes[:], e.palette[:])
}

func (e *TrailsEffect) Close() error {
	if !e.closed {
		e.masks.close()
		e.lookup.Close()
		e.closed = true
	}
	return nil
}

package demo

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

type SlicedEffect struct {
	model  source.SlicedModel
	clock  *source.SlicedClock
	banks  map[string]source.Animation
	raster *NoiseEffect
}

func NewSlicedEffect() (*SlicedEffect, error) {
	read := func(name string) ([]byte, error) { return assets.Files.ReadFile("raw/" + name + ".bin") }
	angular, err := read("fifth-effects")
	if err != nil {
		return nil, err
	}
	material, err := read("fifth-animation")
	if err != nil {
		return nil, err
	}
	previous, err := read("third-animation")
	if err != nil {
		return nil, err
	}
	retained, err := source.ReadAngularModel(angular, material, previous)
	if err != nil {
		return nil, err
	}
	controller, err := read("sixth-effects")
	if err != nil {
		return nil, err
	}
	model, err := source.ReadSlicedModel(controller, retained.Material)
	if err != nil {
		return nil, err
	}
	raster, err := newNoiseRenderer(model.Material)
	if err != nil {
		return nil, err
	}
	e := &SlicedEffect{model: model, clock: source.NewSlicedClock(model), raster: raster, banks: make(map[string]source.Animation)}
	for _, name := range []string{"outline-a", "outline-b"} {
		e.banks[name], err = source.LoadAnimation(assets.Files, name)
		if err != nil {
			raster.Close()
			return nil, err
		}
	}
	return e, nil
}

func (e *SlicedEffect) Update(f kit.Frame) error {
	tick := max(0, min(source.SlicedTicks-1, int(math.Round(f.Time*FPS))))
	if e.clock.Tick > tick+1 {
		e.clock = source.NewSlicedClock(e.model)
		e.raster.masks.clear()
	}
	for e.clock.Tick <= tick {
		if e.clock.Step() {
			var frame []source.Polygon
			if e.clock.Draw {
				frame = e.banks[e.clock.Bank].Frames[e.clock.Frame]
			}
			e.raster.masks.prepareSliced(e.clock.Current, frame, e.clock.Zoom, e.clock.Bank == "outline-b")
		}
	}
	return nil
}

func (e *SlicedEffect) Draw(dst *ebiten.Image) {
	for i, word := range e.clock.Palette {
		c := source.RGB12(word)
		e.raster.palette[3*i], e.raster.palette[3*i+1], e.raster.palette[3*i+2] = float32(c.R)/255, float32(c.G)/255, float32(c.B)/255
	}
	e.raster.drawPlanes(dst, e.clock.Display, e.clock.Texture)
}
func (e *SlicedEffect) Close() error { return e.raster.Close() }

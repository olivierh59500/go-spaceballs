package demo

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

// AngularEffect keeps the authored shape/cue program and reuses the textured
// five-mask composition of the earlier retained-dancer passage.
type AngularEffect struct {
	model  source.AngularModel
	clock  *source.AngularClock
	banks  map[string]source.Animation
	raster *NoiseEffect
}

func NewAngularEffect() (*AngularEffect, error) {
	controller, err := assets.Files.ReadFile("raw/fifth-effects.bin")
	if err != nil {
		return nil, err
	}
	material, err := assets.Files.ReadFile("raw/fifth-animation.bin")
	if err != nil {
		return nil, err
	}
	preceding, err := assets.Files.ReadFile("raw/third-animation.bin")
	if err != nil {
		return nil, err
	}
	model, err := source.ReadAngularModel(controller, material, preceding)
	if err != nil {
		return nil, err
	}
	raster, err := newNoiseRenderer(model.Material)
	if err != nil {
		return nil, err
	}
	e := &AngularEffect{model: model, clock: source.NewAngularClock(model), raster: raster, banks: make(map[string]source.Animation)}
	for _, name := range []string{"zoom-a", "zoom-b", "zoom-c", "zoom-d", "zoom-e"} {
		bank, err := source.LoadAnimation(assets.Files, name)
		if err != nil {
			raster.Close()
			return nil, err
		}
		e.banks[name] = bank
	}
	return e, nil
}

func (e *AngularEffect) Update(f kit.Frame) error {
	tick := max(0, min(source.AngularTicks-1, int(math.Round(f.Time*FPS))))
	if e.clock.Tick > tick+1 {
		e.clock = source.NewAngularClock(e.model)
		e.raster.masks.clear()
	}
	for e.clock.Tick <= tick {
		if e.clock.Step() {
			var frame []source.Polygon
			if e.clock.Draw {
				frame = e.banks[e.clock.Bank].Frames[e.clock.Frame]
			}
			e.raster.masks.prepare(e.clock.Current, frame, false)
		}
	}
	return nil
}

func (e *AngularEffect) Draw(dst *ebiten.Image) {
	e.raster.drawPlanes(dst, e.clock.Display, e.clock.Texture)
}
func (e *AngularEffect) Close() error { return e.raster.Close() }

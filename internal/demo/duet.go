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

type DuetEffect struct {
	model  source.DuetModel
	clock  *source.DuetClock
	banks  map[string]source.Animation
	raster *NoiseEffect
	ending *composite.BitplanePalette
	colors [2]color.NRGBA
}

func retainedMaterial() (source.NoiseModel, error) {
	read := func(name string) ([]byte, error) { return assets.Files.ReadFile("raw/" + name + ".bin") }
	controller, err := read("fifth-effects")
	if err != nil {
		return source.NoiseModel{}, err
	}
	data, err := read("fifth-animation")
	if err != nil {
		return source.NoiseModel{}, err
	}
	previous, err := read("third-animation")
	if err != nil {
		return source.NoiseModel{}, err
	}
	model, err := source.ReadAngularModel(controller, data, previous)
	return model.Material, err
}

func NewDuetEffect() (*DuetEffect, error) {
	controller, err := assets.Files.ReadFile("raw/eighth-effects.bin")
	if err != nil {
		return nil, err
	}
	material, err := retainedMaterial()
	if err != nil {
		return nil, err
	}
	model, err := source.ReadDuetModel(controller, material)
	if err != nil {
		return nil, err
	}
	raster, err := newNoiseRenderer(model.Material)
	if err != nil {
		return nil, err
	}
	ending, err := composite.NewBitplanePalette(composite.BitplanePaletteConfig{Width: Width, Height: Height, Planes: 1, Palette: make([]color.NRGBA, 2), Channels: []composite.BitplaneChannel{composite.BitplaneRed}})
	if err != nil {
		raster.Close()
		return nil, err
	}
	e := &DuetEffect{model: model, clock: source.NewDuetClock(model), raster: raster, ending: ending,
		banks: make(map[string]source.Animation)}
	for _, name := range []string{"two-body-a", "two-body-b"} {
		e.banks[name], err = source.LoadAnimation(assets.Files, name)
		if err != nil {
			e.Close()
			return nil, err
		}
	}
	return e, nil
}

func (e *DuetEffect) Update(f kit.Frame) error {
	tick := max(0, min(source.DuetTicks-1, int(math.Round(f.Time*FPS))))
	if e.clock.Tick > tick+1 {
		e.clock = source.NewDuetClock(e.model)
		e.raster.masks.clear()
	}
	for e.clock.Tick <= tick {
		if e.clock.Step() {
			e.raster.masks.prepare(e.clock.Current, e.banks[e.clock.Bank].Frames[e.clock.Frame], false)
		}
	}
	return nil
}

func (e *DuetEffect) Draw(dst *ebiten.Image) {
	if e.clock.Ending {
		for i, word := range e.clock.EndingColors {
			e.colors[i] = source.RGB12(word)
		}
		if e.clock.EndingTexture < 0 {
			dst.Fill(e.colors[0])
		} else {
			planes := [1]*ebiten.Image{e.raster.texture[e.clock.EndingTexture]}
			drawBitplanes(e.ending, dst, planes[:], e.colors[:])
		}
		return
	}
	for i, word := range e.clock.Palette {
		e.raster.palette[i] = source.RGB12(word)
	}
	e.raster.drawPlanes(dst, e.clock.Display, e.clock.Texture)
}
func (e *DuetEffect) Close() error { e.ending.Close(); return e.raster.Close() }

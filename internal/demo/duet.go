package demo

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

type DuetEffect struct {
	model   source.DuetModel
	clock   *source.DuetClock
	banks   map[string]source.Animation
	raster  *NoiseEffect
	ending  *ebiten.Shader
	colors  []float32
	options ebiten.DrawRectShaderOptions
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
	ending, err := ebiten.NewShader([]byte(duetEndingShader))
	if err != nil {
		raster.Close()
		return nil, err
	}
	e := &DuetEffect{model: model, clock: source.NewDuetClock(model), raster: raster, ending: ending,
		banks: make(map[string]source.Animation), colors: make([]float32, 6)}
	for _, name := range []string{"two-body-a", "two-body-b"} {
		e.banks[name], err = source.LoadAnimation(assets.Files, name)
		if err != nil {
			e.Close()
			return nil, err
		}
	}
	e.options.Uniforms = map[string]any{"Colors": e.colors}
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
			c := source.RGB12(word)
			e.colors[3*i], e.colors[3*i+1], e.colors[3*i+2] = float32(c.R)/255, float32(c.G)/255, float32(c.B)/255
		}
		e.options.Images[0] = nil
		if e.clock.EndingTexture >= 0 {
			e.options.Images[0] = e.raster.texture[e.clock.EndingTexture]
		}
		dst.DrawRectShader(Width, Height, e.ending, &e.options)
		return
	}
	for i, word := range e.clock.Palette {
		c := source.RGB12(word)
		e.raster.palette[3*i], e.raster.palette[3*i+1], e.raster.palette[3*i+2] = float32(c.R)/255, float32(c.G)/255, float32(c.B)/255
	}
	e.raster.drawPlanes(dst, e.clock.Display, e.clock.Texture)
}
func (e *DuetEffect) Close() error { e.ending.Deallocate(); return e.raster.Close() }

const duetEndingShader = `//kage:unit pixels
package main
var Colors [2]vec3
func Fragment(position vec4,texCoord vec2,color vec4) vec4 {
	bit:=int(step(0.5,imageSrc0At(texCoord).r))
	return vec4(Colors[bit],1)
}
`

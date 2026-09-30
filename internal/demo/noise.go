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

// NoiseEffect shares the retained-mask drawing and four-bit packing pass with
// TrailsEffect. Its sixth plane selects the original half-bright palette.
type NoiseEffect struct {
	model   source.NoiseModel
	clock   *source.NoiseClock
	poses   source.Animation
	masks   *maskRing
	texture [4]*ebiten.Image
	lookup  *composite.BitplanePalette
	palette [64]color.NRGBA
	closed  bool
}

func NewNoiseEffect() (*NoiseEffect, error) {
	controller, err := assets.Files.ReadFile("raw/third-effects.bin")
	if err != nil {
		return nil, err
	}
	material, err := assets.Files.ReadFile("raw/third-animation.bin")
	if err != nil {
		return nil, err
	}
	model, err := source.ReadNoiseModel(controller, material)
	if err != nil {
		return nil, err
	}
	poses, err := source.LoadAnimation(assets.Files, "middle-turn")
	if err != nil {
		return nil, err
	}
	e, err := newNoiseRenderer(model)
	if err != nil {
		return nil, err
	}
	e.poses = poses
	e.initialize()
	return e, nil
}

func newNoiseRenderer(model source.NoiseModel) (*NoiseEffect, error) {
	e := &NoiseEffect{model: model, masks: newMaskRing()}
	for i, bits := range model.Bits {
		pixels, err := source.MonochromeImage(bits)
		if err != nil {
			e.Close()
			return nil, err
		}
		e.texture[i] = ebiten.NewImageFromImage(pixels)
	}
	for i, word := range model.Palette {
		e.palette[i] = source.RGB12(word)
	}
	var err error
	e.lookup, err = composite.NewBitplanePalette(composite.BitplanePaletteConfig{
		Width: Width, Height: Height, Planes: 6, Palette: e.palette[:],
		Channels: []composite.BitplaneChannel{composite.BitplaneAlpha, composite.BitplaneAlpha,
			composite.BitplaneAlpha, composite.BitplaneAlpha, composite.BitplaneAlpha, composite.BitplaneRed},
	})
	if err != nil {
		e.Close()
		return nil, err
	}
	return e, nil
}

func (e *NoiseEffect) initialize() {
	e.clock = source.NewNoiseClock()
	e.masks.clear()
	for e.clock.Prime() {
		e.masks.prepare(e.clock.Current, e.poses.Frames[e.clock.Frame], false)
	}
}

func (e *NoiseEffect) Update(f kit.Frame) error {
	e.SetTick(int(math.Round(f.Time * FPS)))
	return nil
}

func (e *NoiseEffect) SetTick(tick int) {
	tick = max(0, min(tick, source.NoiseTicks-1))
	if e.clock.Tick > tick+1 {
		e.initialize()
	}
	for e.clock.Tick <= tick {
		if e.clock.Step() {
			e.masks.prepare(e.clock.Current, e.poses.Frames[e.clock.Frame], false)
		}
	}
}

func (e *NoiseEffect) Draw(dst *ebiten.Image) {
	e.drawPlanes(dst, e.clock.Display, e.clock.Texture)
}

func (e *NoiseEffect) drawPlanes(dst *ebiten.Image, display [5]int, texture int) {
	var planes [6]*ebiten.Image
	for i, index := range display {
		planes[i] = e.masks.planes[index]
	}
	planes[5] = e.texture[texture]
	drawBitplanes(e.lookup, dst, planes[:], e.palette[:])
}

func (e *NoiseEffect) Close() error {
	if !e.closed {
		e.masks.close()
		for _, texture := range e.texture {
			if texture != nil {
				texture.Deallocate()
			}
		}
		e.lookup.Close()
		e.closed = true
	}
	return nil
}

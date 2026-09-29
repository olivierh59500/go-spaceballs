package demo

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/render"
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
	packed  *ebiten.Image
	pack    *ebiten.Shader
	shader  *ebiten.Shader
	palette []float32
	packOp  ebiten.DrawRectShaderOptions
	finalOp ebiten.DrawRectShaderOptions
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
	pack, err := ebiten.NewShader([]byte(trailsPackShader))
	if err != nil {
		return nil, err
	}
	shader, err := ebiten.NewShader([]byte(noiseShader))
	if err != nil {
		pack.Deallocate()
		return nil, err
	}
	e := &NoiseEffect{model: model, poses: poses, clock: source.NewNoiseClock(),
		pack: pack, shader: shader, masks: newMaskRing(), packed: render.NewSurface(Width, Height), palette: make([]float32, 192)}
	for i, bits := range model.Bits {
		pixels, err := source.MonochromeImage(bits)
		if err != nil {
			e.Close()
			return nil, err
		}
		e.texture[i] = ebiten.NewImageFromImage(pixels)
	}
	for i, word := range model.Palette {
		c := source.RGB12(word)
		e.palette[3*i], e.palette[3*i+1], e.palette[3*i+2] = float32(c.R)/255, float32(c.G)/255, float32(c.B)/255
	}
	e.packOp.Blend = ebiten.BlendCopy
	e.finalOp.Uniforms = map[string]any{"Palette": e.palette}
	e.initialize()
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
	for i := range e.packOp.Images {
		e.packOp.Images[i] = e.masks.planes[e.clock.Display[i]]
	}
	e.packed.DrawRectShader(Width, Height, e.pack, &e.packOp)
	e.finalOp.Images[0], e.finalOp.Images[1], e.finalOp.Images[2] = e.packed,
		e.masks.planes[e.clock.Display[4]], e.texture[e.clock.Texture]
	dst.DrawRectShader(Width, Height, e.shader, &e.finalOp)
}

func (e *NoiseEffect) Close() error {
	if !e.closed {
		e.masks.close()
		for _, texture := range e.texture {
			if texture != nil {
				texture.Deallocate()
			}
		}
		e.packed.Deallocate()
		e.pack.Deallocate()
		e.shader.Deallocate()
		e.closed = true
	}
	return nil
}

const noiseShader = `//kage:unit pixels
package main

var Palette [64]vec3

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	bits := imageSrc0At(texCoord)
	fifth := step(0.5, imageSrc1At(texCoord).a)
	dither := step(0.5, imageSrc2At(texCoord).r)
	index := int(bits.r + bits.g*2 + bits.b*4 + bits.a*8 + fifth*16 + dither*32 + 0.5)
	return vec4(Palette[index], 1)
}
`

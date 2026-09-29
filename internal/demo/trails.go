package demo

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

// TrailsEffect draws the original retained-mask program, including its bank
// cues, clipped close-up lookup and periodically inverted 32-color palette.
type TrailsEffect struct {
	model   source.TrailModel
	clock   *source.TrailsClock
	banks   map[string]source.Animation
	planes  [6]*ebiten.Image
	packed  *ebiten.Image
	white   *ebiten.Image
	pack    *ebiten.Shader
	shader  *ebiten.Shader
	batch   *render.Batch
	palette []float32
	packOp  ebiten.DrawRectShaderOptions
	finalOp ebiten.DrawRectShaderOptions
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
		banks: make(map[string]source.Animation), batch: render.NewBatch(4096), palette: make([]float32, 96)}
	for _, name := range []string{"trails-a", "trails-b", "trails-c", "trails-d"} {
		e.banks[name], err = source.LoadAnimation(assets.Files, name)
		if err != nil {
			return nil, err
		}
	}
	e.pack, err = ebiten.NewShader([]byte(trailsPackShader))
	if err != nil {
		return nil, err
	}
	e.shader, err = ebiten.NewShader([]byte(trailsShader))
	if err != nil {
		e.pack.Deallocate()
		return nil, err
	}
	for i := range e.planes {
		e.planes[i] = render.NewSurface(Width, Height)
	}
	e.packed = render.NewSurface(Width, Height)
	e.white = ebiten.NewImage(1, 1)
	e.white.Fill(color.White)
	e.batch.Options.FillRule = ebiten.FillRuleEvenOdd
	e.packOp.Blend = ebiten.BlendCopy
	e.finalOp.Uniforms = map[string]any{"Palette": e.palette}
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
		for _, plane := range e.planes {
			plane.Clear()
		}
	}
	for e.clock.Tick <= tick {
		if e.clock.Step() {
			plane := e.planes[e.clock.Current]
			plane.Clear()
			if e.clock.Draw {
				e.batch.Begin(plane, e.white)
				for _, polygon := range e.banks[e.clock.Bank].Frames[e.clock.Frame] {
					filledContour(e.batch, polygon.Points, e.clock.Zoom)
				}
				e.batch.Flush()
			}
		}
	}
}

func (e *TrailsEffect) Draw(dst *ebiten.Image) {
	// Pack four independent mask bits into RGBA. This keeps the final material
	// lookup within Ebitengine's four-source limit, without CPU pixel readback.
	for i := range e.packOp.Images {
		e.packOp.Images[i] = e.planes[e.clock.Display[i]]
	}
	e.packed.DrawRectShader(Width, Height, e.pack, &e.packOp)
	for i, word := range e.clock.Palette {
		c := source.RGB12(word)
		e.palette[3*i], e.palette[3*i+1], e.palette[3*i+2] = float32(c.R)/255, float32(c.G)/255, float32(c.B)/255
	}
	e.finalOp.Images[0], e.finalOp.Images[1] = e.packed, e.planes[e.clock.Display[4]]
	dst.DrawRectShader(Width, Height, e.shader, &e.finalOp)
}

func (e *TrailsEffect) Close() error {
	if !e.closed {
		for _, plane := range e.planes {
			plane.Deallocate()
		}
		e.packed.Deallocate()
		e.white.Deallocate()
		e.pack.Deallocate()
		e.shader.Deallocate()
		e.closed = true
	}
	return nil
}

const trailsPackShader = `//kage:unit pixels
package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	p := texCoord - imageSrc0Origin()
	return vec4(
		step(0.5, imageSrc0At(imageSrc0Origin() + p).a),
		step(0.5, imageSrc1At(imageSrc1Origin() + p).a),
		step(0.5, imageSrc2At(imageSrc2Origin() + p).a),
		step(0.5, imageSrc3At(imageSrc3Origin() + p).a))
}
`

const trailsShader = `//kage:unit pixels
package main

var Palette [32]vec3

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	p := texCoord - imageSrc0Origin()
	bits := imageSrc0At(texCoord)
	fifth := step(0.5, imageSrc1At(imageSrc1Origin() + p).a)
	index := int(bits.r + bits.g*2 + bits.b*4 + bits.a*8 + fifth*16 + 0.5)
	return vec4(Palette[index], 1)
}
`

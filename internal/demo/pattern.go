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

// PatternEffect reconstructs the first musical scene from the authored material,
// contour banks, integer motion tables and RGB12 palette program.
type PatternEffect struct {
	material source.PatternMaterial
	frames   [][]source.Polygon
	state    source.PatternState
	mask     *ebiten.Image
	canvas   *ebiten.Image
	texture  *ebiten.Image
	white    *ebiten.Image
	shader   *ebiten.Shader
	batch    *render.Batch
	palette  []float32
	first    []float32
	second   []float32
	options  ebiten.DrawRectShaderOptions
}

func NewPatternEffect() (*PatternEffect, error) {
	bytes, err := assets.Files.ReadFile("raw/later-animation.bin")
	if err != nil {
		return nil, err
	}
	material, err := source.ReadPatternMaterial(bytes)
	if err != nil {
		return nil, err
	}
	var frames [][]source.Polygon
	for _, name := range []string{"turn", "long-turn"} {
		bank, err := source.LoadAnimation(assets.Files, name)
		if err != nil {
			return nil, err
		}
		frames = append(frames, bank.Frames...)
	}
	shader, err := ebiten.NewShader([]byte(patternShader))
	if err != nil {
		return nil, err
	}
	e := &PatternEffect{material: material, frames: frames, shader: shader,
		mask: render.NewSurface(640, 512), canvas: render.NewSurface(342, Height), texture: ebiten.NewImageFromImage(material.Image()),
		white: ebiten.NewImage(1, 1), batch: render.NewBatch(4096),
		palette: make([]float32, 8*3), first: make([]float32, 2), second: make([]float32, 2)}
	e.white.Fill(color.White)
	e.batch.Options.FillRule = ebiten.FillRuleEvenOdd
	e.options.Images[0], e.options.Images[1] = e.mask, e.texture
	e.options.Uniforms = map[string]any{"Palette": e.palette, "First": e.first, "Second": e.second}
	e.SetTick(0)
	return e, nil
}

func (e *PatternEffect) Update(f kit.Frame) error {
	e.SetTick(int(math.Round(f.Time * FPS)))
	return nil
}

func (e *PatternEffect) SetTick(tick int) { e.state = e.material.State(tick) }

func (e *PatternEffect) Draw(dst *ebiten.Image) {
	e.mask.Clear()
	if e.state.Frame >= 0 && e.state.Frame < len(e.frames) {
		e.batch.Begin(e.mask, e.white)
		point := func(p source.Point) ebiten.Vertex {
			return render.Vertex(float64(int(p.X)*351/256), float64(int(p.Y)*289/204), 0, 0, color.White)
		}
		// This scene merges all contour commands into one XOR-filled body plane.
		for _, polygon := range e.frames[e.state.Frame] {
			if len(polygon.Points) < 3 {
				continue
			}
			first := point(polygon.Points[0])
			for i := 1; i+1 < len(polygon.Points); i++ {
				e.batch.Triangle(first, point(polygon.Points[i]), point(polygon.Points[i+1]))
			}
		}
		e.batch.Flush()
	}
	for i, word := range e.state.Palette {
		c := source.RGB12(word)
		e.palette[3*i], e.palette[3*i+1], e.palette[3*i+2] = float32(c.R)/255, float32(c.G)/255, float32(c.B)/255
	}
	// Display fetch begins 32 low-resolution pixels before the visible window.
	// The first material plane also keeps the source's fine-scroll delay.
	e.first[0], e.first[1] = float32(e.state.X+17), float32(e.state.Y)
	e.second[0], e.second[1] = 32, float32(e.state.SecondY)
	e.canvas.DrawRectShader(640, 512, e.shader, &e.options)
	var op ebiten.DrawImageOptions
	// DIWSTRT $1c71 and DIWSTOP $3ec7 expose 342 low-resolution pixels.
	op.GeoM.Scale(float64(Width)/342, 1)
	dst.DrawImage(e.canvas, &op)
}

func (e *PatternEffect) Close() error {
	e.mask.Deallocate()
	e.canvas.Deallocate()
	e.texture.Deallocate()
	e.white.Deallocate()
	e.shader.Deallocate()
	return nil
}

const patternShader = `//kage:unit pixels
package main

var First vec2
var Second vec2
var Palette [8]vec3

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	p := texCoord - imageSrc0Origin()
	body := int(step(0.5, imageSrc0At(imageSrc0Origin() + p + vec2(32, 0)).a))
	one := int(step(0.5, imageSrc1At(imageSrc1Origin() + p + First).r))
	two := int(step(0.5, imageSrc1At(imageSrc1Origin() + p + Second).r))
	return vec4(Palette[body + one*2 + two*4], 1)
}
`

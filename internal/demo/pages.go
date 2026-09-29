package demo

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
	"image"
	"image/color"
	"math"
)

// PageEffect renders recovered palette indices with source RGB12 transitions.
type PageEffect struct {
	model                source.IndexedPage
	image, canvas        *ebiten.Image
	shader               *ebiten.Shader
	palette              []float32
	options              ebiten.DrawRectShaderOptions
	from, to             []uint16
	fadeStart, fadeTicks int
	tick                 int
	closed               bool
	artwork              image.Rectangle
}

func NewPageEffect(name string) (*PageEffect, error) {
	read := func(name string) ([]byte, error) { return assets.Files.ReadFile("raw/" + name + ".bin") }
	var model source.IndexedPage
	var err error
	switch name {
	case "state", "of", "the", "art":
		packed, e := read("pattern")
		if e != nil {
			return nil, e
		}
		director, e := read("director")
		if e != nil {
			return nil, e
		}
		offset := map[string]int{"state": 0, "of": 0x14c8, "the": 0x1fc8, "art": 0x220c}[name]
		model, err = source.TitlePage(packed, director, offset)
	case "credits":
		director, e := read("director")
		if e != nil {
			return nil, e
		}
		model, err = source.CreditPage(director)
	case "dragon":
		packed, e := read("credits-packed")
		if e != nil {
			return nil, e
		}
		director, e := read("second-director")
		if e != nil {
			return nil, e
		}
		model, err = source.DragonPage(packed, director)
	case "closing":
		controller, e := read("final-effects")
		if e != nil {
			return nil, e
		}
		model, err = source.ClosingPage(controller)
	default:
		return nil, fmt.Errorf("demo: unknown recovered page %q", name)
	}
	if err != nil {
		return nil, err
	}
	shader, err := ebiten.NewShader([]byte(pageShader))
	if err != nil {
		return nil, err
	}
	e := &PageEffect{model: model, image: ebiten.NewImageFromImage(model.Pixels), canvas: render.NewSurface(model.Pixels.Bounds().Dx(), model.Pixels.Bounds().Dy()), shader: shader, palette: make([]float32, 48), fadeStart: -1}
	e.options.Images[0] = e.image
	e.options.Uniforms = map[string]any{"Palette": e.palette}
	e.from, e.to = append([]uint16(nil), model.Palette...), append([]uint16(nil), model.Palette...)
	if model.Hires {
		r := model.Pixels.Bounds()
		minX, maxX := r.Max.X, 0
		for y := r.Min.Y; y < r.Max.Y; y++ {
			for x := r.Min.X; x < r.Max.X; x++ {
				if model.Pixels.NRGBAAt(x, y).R != 0 {
					minX = min(minX, x)
					maxX = max(maxX, x+1)
				}
			}
		}
		e.artwork = image.Rect(minX, 0, maxX, r.Dy())
	}
	return e, nil
}

func (e *PageEffect) Fade(from, to []uint16, start, ticks int) {
	e.from, e.to = append([]uint16(nil), from...), append([]uint16(nil), to...)
	e.fadeStart, e.fadeTicks = start, ticks
}
func (e *PageEffect) Update(f kit.Frame) error {
	e.tick = max(0, int(math.Round(f.Time*FPS)))
	return nil
}
func (e *PageEffect) Draw(dst *ebiten.Image) {
	for i := range e.model.Palette {
		word := e.from[i]
		if e.fadeStart >= 0 && e.tick >= e.fadeStart {
			step := min(e.fadeTicks-1, e.tick-e.fadeStart)
			word = source.BlendRGB12(e.from[i], e.to[i], step, max(1, e.fadeTicks-1))
		}
		c := source.RGB12(word)
		e.palette[3*i], e.palette[3*i+1], e.palette[3*i+2] = float32(c.R)/255, float32(c.G)/255, float32(c.B)/255
	}
	dst.Fill(color.Black)
	b := e.canvas.Bounds()
	e.canvas.DrawRectShader(b.Dx(), b.Dy(), e.shader, &e.options)
	var op ebiten.DrawImageOptions
	canvas := e.canvas
	if e.model.Hires {
		// Preserve the recorded hires artwork aspect within the common viewport.
		canvas = e.canvas.SubImage(e.artwork).(*ebiten.Image)
		op.GeoM.Scale(226.0/float64(e.artwork.Dx()), float64(Height)/float64(e.artwork.Dy()))
		op.GeoM.Translate(65, 0)
	} else if e.canvas.Bounds().Dy() != Height {
		op.GeoM.Scale(1, float64(Height)/float64(e.canvas.Bounds().Dy()))
	}
	dst.DrawImage(canvas, &op)
}
func (e *PageEffect) Close() error {
	if !e.closed {
		e.image.Deallocate()
		e.canvas.Deallocate()
		e.shader.Deallocate()
		e.closed = true
	}
	return nil
}

const pageShader = `//kage:unit pixels
package main
var Palette [16]vec3
func Fragment(position vec4,texCoord vec2,color vec4)vec4{index:=int(imageSrc0At(texCoord).r*15+.5);return vec4(Palette[index],1)}
`

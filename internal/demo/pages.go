package demo

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/effects"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

// PageEffect renders recovered palette indices with source RGB12 transitions.
type PageEffect struct {
	model    source.IndexedPage
	image    *ebiten.Image
	renderer *effects.IndexedImage
	closed   bool
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
	e := &PageEffect{model: model, image: ebiten.NewImageFromImage(model.Pixels)}
	config := effects.IndexedImageConfig{Image: e.image, Palette: packedWords(model.Palette), FPS: FPS,
		Channel: composite.BitplaneRed, Scale: 15}
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
		config.Crop = image.Rect(minX, 0, maxX, r.Dy())
		// Keep the recovered hires artwork aspect in the original viewport.
		config.Options.GeoM.Scale(226.0/float64(config.Crop.Dx()), float64(Height)/float64(config.Crop.Dy()))
		config.Options.GeoM.Translate(65, 0)
	} else if model.Pixels.Bounds().Dy() != Height {
		config.Options.GeoM.Scale(1, float64(Height)/float64(model.Pixels.Bounds().Dy()))
	}
	e.renderer, err = effects.NewIndexedImage(config)
	if err != nil {
		e.image.Deallocate()
		return nil, err
	}
	return e, nil
}

func (e *PageEffect) Fade(from, to []uint16, start, ticks int) {
	if err := e.renderer.Fade(packedWords(from), packedWords(to), start, ticks); err != nil {
		panic(err)
	}
}
func (e *PageEffect) Update(f kit.Frame) error {
	return e.renderer.Update(f)
}
func (e *PageEffect) Draw(dst *ebiten.Image) {
	dst.Fill(color.Black)
	e.renderer.Draw(dst)
}
func (e *PageEffect) Close() error {
	if !e.closed {
		e.image.Deallocate()
		e.renderer.Close()
		e.closed = true
	}
	return nil
}

func packedWords(words []uint16) []uint32 {
	out := make([]uint32, len(words))
	for i, word := range words {
		out[i] = uint32(word)
	}
	return out
}

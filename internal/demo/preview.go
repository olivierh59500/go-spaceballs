// Package demo assembles native effects from the recovered production data.
package demo

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

const Width, Height, FPS = 352, 290, 50

// Preview is an isolated source-asset host, not the complete scene director.
type Preview struct {
	animation source.Animation
	frame     int
	white     *ebiten.Image
	picture   *ebiten.Image
	batch     *render.Batch
	hands     bool
	hires     bool
	pattern   *PatternEffect
	opening   *OpeningEffect
}

func NewPreview(bank string, frame, pictureOffset int) (*Preview, error) {
	p := &Preview{frame: frame, batch: render.NewBatch(4096)}
	if bank == "opening-effect" {
		var err error
		p.opening, err = NewOpeningEffect()
		return p, err
	}
	if bank == "pattern" {
		var err error
		p.pattern, err = NewPatternEffect()
		return p, err
	}
	if bank == "dragon" {
		data, _ := assets.Files.ReadFile("raw/credits-packed.bin")
		director, _ := assets.Files.ReadFile("raw/second-director.bin")
		pixels, err := source.Dragon(data, director)
		if err != nil {
			return nil, err
		}
		p.picture, p.hires = ebiten.NewImageFromImage(pixels), true
		return p, nil
	}
	if bank == "picture" {
		data, _ := assets.Files.ReadFile("raw/pattern.bin")
		director, _ := assets.Files.ReadFile("raw/director.bin")
		pixels, err := source.TitlePicture(data, director, pictureOffset)
		if err != nil {
			return nil, err
		}
		p.picture = ebiten.NewImageFromImage(pixels)
		return p, nil
	}
	name := bank
	if bank == "intro" {
		name = "opening"
	}
	if bank == "first-animation" {
		name = "first"
	}
	var err error
	p.animation, err = source.LoadAnimation(assets.Files, name)
	if err != nil {
		return nil, err
	}
	p.hands = name == "opening"
	p.white = ebiten.NewImage(1, 1)
	p.white.Fill(color.White)
	p.batch.Options.FillRule = ebiten.FillRuleEvenOdd
	return p, nil
}

func (p *Preview) Update() error { p.frame++; return nil }

func (p *Preview) Draw(dst *ebiten.Image) {
	if p.opening != nil {
		p.opening.SetTick((p.frame%source.OpeningTicks + source.OpeningTicks) % source.OpeningTicks)
		p.opening.Draw(dst)
		return
	}
	if p.pattern != nil {
		p.pattern.SetTick((p.frame%source.PatternTicks + source.PatternTicks) % source.PatternTicks)
		p.pattern.Draw(dst)
		return
	}
	dst.Fill(source.RGB12(0x102))
	if p.picture != nil {
		var op ebiten.DrawImageOptions
		if p.hires {
			op.GeoM.Scale(0.5, 1)
			op.GeoM.Translate(16, 0)
		}
		dst.DrawImage(p.picture, &op)
		return
	}
	frame := p.animation.Frames[(p.frame%len(p.animation.Frames)+len(p.animation.Frames))%len(p.animation.Frames)]
	// Draw both original bitplane masks independently, retaining contour parity.
	for plane := 0; plane < 2; plane++ {
		p.batch.Begin(dst, p.white)
		word := uint16(0x347)
		if plane == 1 {
			word = 0x68a
		}
		vertex := func(point source.Point) ebiten.Vertex {
			x, y := int(point.X)*351/256, int(point.Y)*289/204
			if p.hands {
				x -= 15
				y = min(289, int(point.Y)*350/204+8)
			}
			return render.Vertex(float64(x), float64(y), 0, 0, source.RGB12(word))
		}
		for _, polygon := range frame {
			if polygon.Mask>>uint(plane)&1 == 0 || len(polygon.Points) < 3 {
				continue
			}
			first := vertex(polygon.Points[0])
			for i := 1; i+1 < len(polygon.Points); i++ {
				p.batch.Triangle(first, vertex(polygon.Points[i]), vertex(polygon.Points[i+1]))
			}
		}
		p.batch.Flush()
	}
}

func (*Preview) Layout(int, int) (int, int) { return Width, Height }

func (p *Preview) Close() {
	if p.opening != nil {
		_ = p.opening.Close()
	}
	if p.pattern != nil {
		_ = p.pattern.Close()
	}
	if p.white != nil {
		p.white.Deallocate()
	}
	if p.picture != nil {
		p.picture.Deallocate()
	}
}

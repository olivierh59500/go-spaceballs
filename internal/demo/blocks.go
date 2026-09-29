package demo

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

// BlocksEffect shares the contour/ring-mask renderer with the neighboring
// scenes, while retaining the authored corner palettes and two-plane shadow.
type BlocksEffect struct {
	model      source.BlockModel
	clock      *source.BlocksClock
	animation  source.Animation
	masks      *maskRing
	palette    *ebiten.Image
	paletteSub *ebiten.Image
	pixels     []byte
	shader     *ebiten.Shader
	options    ebiten.DrawRectShaderOptions
	closed     bool
	dirty      bool
}

func NewBlocksEffect() (*BlocksEffect, error) {
	bytes, err := assets.Files.ReadFile("raw/second-effects.bin")
	if err != nil {
		return nil, err
	}
	model, err := source.ReadBlockModel(bytes)
	if err != nil {
		return nil, err
	}
	animation, err := source.LoadAnimation(assets.Files, "blocks")
	if err != nil {
		return nil, err
	}
	shader, err := ebiten.NewShader([]byte(blocksShader))
	if err != nil {
		return nil, err
	}
	e := &BlocksEffect{model: model, clock: source.NewBlocksClock(model), animation: animation,
		shader: shader, pixels: make([]byte, 30*17*4), dirty: true}
	e.masks = newMaskRing()
	e.palette = render.NewSurface(Width, Height)
	e.paletteSub = e.palette.SubImage(image.Rect(0, 0, 30, 17)).(*ebiten.Image)
	e.options.Images[2] = e.palette
	return e, nil
}

func (e *BlocksEffect) Update(f kit.Frame) error {
	e.SetTick(int(math.Round(f.Time * FPS)))
	return nil
}

func (e *BlocksEffect) SetTick(tick int) {
	tick = max(0, min(tick, source.BlocksTicks-1))
	if e.clock.Tick > tick+1 {
		e.clock = source.NewBlocksClock(e.model)
		e.masks.clear()
	}
	for e.clock.Tick <= tick {
		if e.clock.Step() {
			e.dirty = true
			e.masks.prepare(e.clock.Current, e.animation.Frames[e.clock.Frame], false)
		}
	}
}

func (e *BlocksEffect) Draw(dst *ebiten.Image) {
	if e.dirty {
		for row := 0; row < 17; row++ {
			for col := 0; col < 15; col++ {
				for plane, word := range [2]uint16{e.clock.Background[row][col], e.clock.Shadow[row][col]} {
					at := (row*30 + col + plane*15) * 4
					c := source.RGB12(word)
					e.pixels[at], e.pixels[at+1], e.pixels[at+2], e.pixels[at+3] = c.R, c.G, c.B, 255
				}
			}
		}
		// Only 2,040 bytes of palette data are uploaded; geometry stays on the GPU.
		e.paletteSub.WritePixels(e.pixels)
		e.dirty = false
	}
	e.options.Images[0], e.options.Images[1] = e.masks.planes[e.clock.Display[0]], e.masks.planes[e.clock.Display[1]]
	dst.DrawRectShader(Width, Height, e.shader, &e.options)
}

func (e *BlocksEffect) Close() error {
	if !e.closed {
		e.masks.close()
		e.palette.Deallocate()
		e.shader.Deallocate()
		e.closed = true
	}
	return nil
}

const blocksShader = `//kage:unit pixels
package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	p := texCoord - imageSrc0Origin()
	body := step(0.5, imageSrc0At(texCoord).a)
	shadow := step(0.5, imageSrc1At(texCoord).a)
	if body > 0.5 {
		return vec4(0, 0, 0, 1)
	}
	row := 0.0
	if p.y >= 36 {
		row = 1 + floor((p.y - 36)/16)
		if row >= 13 {
			row += 1
		}
		row = min(16, row)
	}
	col := min(14, max(0, floor((p.x + 8)/24)))
	return imageSrc2At(imageSrc0Origin() + vec2(col + shadow*15, row) + vec2(0.5))
}
`

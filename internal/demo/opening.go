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

// OpeningEffect keeps the original held-plane and triple-buffer composition.
// It samples an absolute time; seeking backward resets its retained masks.
type OpeningEffect struct {
	clock   *source.OpeningClock
	banks   map[string]source.Animation
	planes  [3][2]*ebiten.Image
	white   *ebiten.Image
	shader  *ebiten.Shader
	batch   *render.Batch
	palette []float32
	options ebiten.DrawRectShaderOptions
	closed  bool
}

func NewOpeningEffect() (*OpeningEffect, error) {
	e := &OpeningEffect{clock: source.NewOpeningClock(), banks: make(map[string]source.Animation),
		batch: render.NewBatch(4096), palette: make([]float32, 12)}
	for _, name := range []string{"hands", "opening", "first"} {
		bank, err := source.LoadAnimation(assets.Files, name)
		if err != nil {
			return nil, err
		}
		e.banks[name] = bank
	}
	shader, err := ebiten.NewShader([]byte(openingShader))
	if err != nil {
		return nil, err
	}
	e.shader = shader
	for i := range e.planes {
		for j := range e.planes[i] {
			e.planes[i][j] = render.NewSurface(Width, Height)
		}
	}
	e.white = ebiten.NewImage(1, 1)
	e.white.Fill(color.White)
	e.batch.Options.FillRule = ebiten.FillRuleEvenOdd
	e.batch.Options.Blend = ebiten.BlendXor
	e.options.Uniforms = map[string]any{"Palette": e.palette}
	return e, nil
}

func (e *OpeningEffect) Update(f kit.Frame) error {
	e.SetTick(int(math.Round(f.Time * FPS)))
	return nil
}

func (e *OpeningEffect) SetTick(tick int) {
	tick = max(0, min(tick, source.OpeningTicks-1))
	if e.clock.Tick > tick+1 {
		e.clock = source.NewOpeningClock()
		for _, planes := range e.planes {
			for _, plane := range planes {
				plane.Clear()
			}
		}
	}
	for e.clock.Tick <= tick {
		if e.clock.Step() {
			e.prepare()
		}
	}
}

func (e *OpeningEffect) prepare() {
	c := e.clock
	planes := e.planes[c.Current]
	if c.SecondOnly {
		planes[1].Clear()
	} else {
		planes[0].Clear()
		if c.FillSecond {
			planes[1].Clear()
		}
	}
	if !c.Draw {
		return
	}
	frame := e.banks[c.Bank].Frames[c.Frame]
	for plane := 0; plane < 2; plane++ {
		bit := plane
		if c.SecondOnly {
			bit--
		}
		if bit < 0 {
			continue
		}
		e.batch.Begin(planes[plane], e.white)
		for _, polygon := range frame {
			if polygon.Mask>>uint(bit)&1 == 0 || len(polygon.Points) < 3 {
				continue
			}
			point := func(p source.Point) ebiten.Vertex {
				x, y := int(p.X)*351/256, int(p.Y)*289/204
				if c.SecondOnly {
					x -= 15
					y = min(289, int(p.Y)*350/204+8)
				}
				return render.Vertex(float64(x), float64(y), 0, 0, color.White)
			}
			if plane == 1 && !c.FillSecond {
				// Once the source stops clearing/filling this plane, it adds edges
				// to the retained image rather than XORing another filled silhouette.
				for i, p := range polygon.Points {
					a, b := point(p), point(polygon.Points[(i+1)%len(polygon.Points)])
					dx, dy := float64(b.DstX-a.DstX), float64(b.DstY-a.DstY)
					length := math.Hypot(dx, dy)
					if length == 0 {
						continue
					}
					nx, ny := float32(-dy/length/2), float32(dx/length/2)
					q := [4]ebiten.Vertex{a, b, b, a}
					q[0].DstX += nx
					q[0].DstY += ny
					q[1].DstX += nx
					q[1].DstY += ny
					q[2].DstX -= nx
					q[2].DstY -= ny
					q[3].DstX -= nx
					q[3].DstY -= ny
					e.batch.Quad(q)
				}
				continue
			}
			first := point(polygon.Points[0])
			for i := 1; i+1 < len(polygon.Points); i++ {
				e.batch.Triangle(first, point(polygon.Points[i]), point(polygon.Points[i+1]))
			}
		}
		e.batch.Flush()
	}
}

func (e *OpeningEffect) Draw(dst *ebiten.Image) {
	for i, word := range e.clock.Palette {
		c := source.RGB12(word)
		e.palette[i*3], e.palette[i*3+1], e.palette[i*3+2] = float32(c.R)/255, float32(c.G)/255, float32(c.B)/255
	}
	e.options.Images[0] = e.planes[e.clock.Display[0]][0]
	e.options.Images[1] = e.planes[e.clock.Display[1]][1]
	dst.DrawRectShader(Width, Height, e.shader, &e.options)
}

func (e *OpeningEffect) Close() error {
	if !e.closed {
		for _, planes := range e.planes {
			for _, plane := range planes {
				plane.Deallocate()
			}
		}
		e.white.Deallocate()
		e.shader.Deallocate()
		e.closed = true
	}
	return nil
}

const openingShader = `//kage:unit pixels
package main

var Palette [4]vec3

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
	p := texCoord - imageSrc0Origin()
	one := int(step(0.5, imageSrc0At(imageSrc0Origin() + p).a))
	two := int(step(0.5, imageSrc1At(imageSrc1Origin() + p).a))
	return vec4(Palette[one + two*2], 1)
}
`

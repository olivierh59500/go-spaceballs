package demo

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
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
	bank     *composite.ContourBank
	lookup   *composite.BitplanePalette
	palette  [8]color.NRGBA
	planes   [3]*ebiten.Image
	offsets  [3][2]float32
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
	e, err := newPatternRenderer(material)
	if err != nil {
		return nil, err
	}
	e.frames = frames
	e.SetTick(0)
	return e, nil
}

func newPatternRenderer(material source.PatternMaterial) (*PatternEffect, error) {
	var colors [8]color.NRGBA
	lookup, err := composite.NewBitplanePalette(composite.BitplanePaletteConfig{
		Width: 640, Height: 512, Planes: 3, Palette: colors[:],
		Channels: []composite.BitplaneChannel{composite.BitplaneAlpha, composite.BitplaneRed, composite.BitplaneRed},
	})
	if err != nil {
		return nil, err
	}
	e := &PatternEffect{material: material, lookup: lookup,
		bank: newContourBank(640, 512, 1, 1, ebiten.Blend{}), canvas: render.NewSurface(342, Height), texture: ebiten.NewImageFromImage(material.Image()),
		offsets: [3][2]float32{{32, 0}}}
	e.mask = e.bank.Image(0, 0)
	e.planes = [3]*ebiten.Image{e.mask, e.texture, e.texture}
	return e, nil
}

func (e *PatternEffect) Update(f kit.Frame) error {
	e.SetTick(int(math.Round(f.Time * FPS)))
	return nil
}

func (e *PatternEffect) SetTick(tick int) { e.state = e.material.State(tick) }

func (e *PatternEffect) Draw(dst *ebiten.Image) {
	paintContours(e.bank, 0, 0, true, func(batch *render.Batch) {
		if e.state.Frame < 0 || e.state.Frame >= len(e.frames) {
			return
		}
		for _, polygon := range e.frames[e.state.Frame] {
			filledContour(batch, polygon.Points, false)
		}
	})
	for i, word := range e.state.Palette {
		c := source.RGB12(word)
		e.palette[i] = color.NRGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	}
	// Display fetch begins 32 low-resolution pixels before the visible window.
	// The first material plane also keeps the source's fine-scroll delay.
	e.offsets[1] = [2]float32{float32(e.state.X + 17), float32(e.state.Y)}
	e.offsets[2] = [2]float32{32, float32(e.state.SecondY)}
	if err := e.lookup.SetPalette(e.palette[:]); err != nil {
		panic(err)
	}
	if err := e.lookup.DrawOffsets(e.canvas, e.planes[:], e.offsets[:]); err != nil {
		panic(err)
	}
	var op ebiten.DrawImageOptions
	// DIWSTRT $1c71 and DIWSTOP $3ec7 expose 342 low-resolution pixels.
	op.GeoM.Scale(float64(Width)/float64(e.canvas.Bounds().Dx()), 1)
	dst.DrawImage(e.canvas, &op)
}

func (e *PatternEffect) Close() error {
	e.bank.Close()
	e.canvas.Deallocate()
	e.texture.Deallocate()
	return e.lookup.Close()
}

package demo

import (
	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
	"math"
)

type FinaleEffect struct {
	model  source.FinaleModel
	clock  *source.FinaleClock
	raster *PatternEffect
	base   map[string]int
	show   bool
}

func NewFinaleEffect() (*FinaleEffect, error) {
	controller, err := assets.Files.ReadFile("raw/final-effects.bin")
	if err != nil {
		return nil, err
	}
	packed, err := assets.Files.ReadFile("raw/later-animation.bin")
	if err != nil {
		return nil, err
	}
	model, err := source.ReadFinaleModel(controller, packed)
	if err != nil {
		return nil, err
	}
	raster, err := newPatternRenderer(model.Material)
	if err != nil {
		return nil, err
	}
	raster.canvas.Deallocate()
	raster.canvas = render.NewSurface(Width, Height)
	raster.maskShift[0] = 0
	e := &FinaleEffect{model: model, clock: source.NewFinaleClock(model), raster: raster, base: make(map[string]int), show: true}
	for _, name := range []string{"final-a", "final-b", "final-c"} {
		bank, err := source.LoadAnimation(assets.Files, name)
		if err != nil {
			raster.Close()
			return nil, err
		}
		e.base[name] = len(raster.frames)
		raster.frames = append(raster.frames, bank.Frames...)
	}
	return e, nil
}
func (e *FinaleEffect) Update(f kit.Frame) error {
	tick := max(0, min(source.FinaleTicks-1, int(math.Round(f.Time*FPS))))
	if e.clock.Tick > tick+1 {
		e.clock = source.NewFinaleClock(e.model)
	}
	for e.clock.Tick <= tick {
		e.clock.Step()
	}
	e.raster.state = e.model.Material.Motion(tick)
	e.raster.state.Palette = e.clock.Palette
	if e.clock.Visible.Frame >= 0 {
		e.raster.state.Frame = e.base[e.clock.Visible.Bank] + e.clock.Visible.Frame
	}
	e.show = e.clock.Show
	return nil
}
func (e *FinaleEffect) Draw(dst *ebiten.Image) {
	if !e.show {
		dst.Fill(source.RGB12(e.clock.Palette[0]))
		return
	}
	e.raster.Draw(dst)
}
func (e *FinaleEffect) Close() error { return e.raster.Close() }

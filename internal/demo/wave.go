package demo

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

// WaveEffect reuses the first musical scene's material and GPU composition.
// Its source controller supplies different motion tables, pose cues and colors.
type WaveEffect struct {
	raster *PatternEffect
	clock  *source.WaveClock
	model  source.WaveModel
	base   map[string]int
}

func NewWaveEffect() (*WaveEffect, error) {
	controller, err := assets.Files.ReadFile("raw/fourth-effects.bin")
	if err != nil {
		return nil, err
	}
	packed, err := assets.Files.ReadFile("raw/later-animation.bin")
	if err != nil {
		return nil, err
	}
	animation, err := assets.Files.ReadFile("raw/fourth-animation.bin")
	if err != nil {
		return nil, err
	}
	model, err := source.ReadWaveModel(controller, packed, animation)
	if err != nil {
		return nil, err
	}
	raster, err := newPatternRenderer(model.Material)
	if err != nil {
		return nil, err
	}
	e := &WaveEffect{raster: raster, model: model, clock: source.NewWaveClock(model), base: make(map[string]int)}
	// The source addresses its mask four bytes before the working buffer,
	// compensating the 32-pixel display-fetch lead used by both materials.
	raster.offsets[0][0] = 0
	for _, name := range []string{"small-turn", "small-dance", "small-final"} {
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

func (e *WaveEffect) Update(f kit.Frame) error {
	e.SetTick(int(math.Round(f.Time * FPS)))
	return nil
}

func (e *WaveEffect) SetTick(tick int) {
	tick = max(0, min(tick, source.WaveTicks-1))
	if e.clock.Tick > tick+1 {
		e.clock = source.NewWaveClock(e.model)
	}
	for e.clock.Tick <= tick {
		e.clock.Step()
	}
	e.raster.state = e.model.Material.Motion(tick)
	e.raster.state.Palette = e.clock.Palette
	if e.clock.Visible.Frame >= 0 {
		e.raster.state.Frame = e.base[e.clock.Visible.Bank] + e.clock.Visible.Frame
	}
}

func (e *WaveEffect) Draw(dst *ebiten.Image) { e.raster.Draw(dst) }
func (e *WaveEffect) Close() error           { return e.raster.Close() }

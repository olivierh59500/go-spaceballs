package demo

import (
	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
	"math"
)

type OutlineEffect struct {
	model  source.OutlineModel
	clock  *source.OutlineClock
	raster *BlocksEffect
	banks  map[string]source.Animation
}

func NewOutlineEffect() (*OutlineEffect, error) {
	data, err := assets.Files.ReadFile("raw/ninth-effects.bin")
	if err != nil {
		return nil, err
	}
	model, err := source.ReadOutlineModel(data)
	if err != nil {
		return nil, err
	}
	raster, err := newBlocksRenderer()
	if err != nil {
		return nil, err
	}
	e := &OutlineEffect{model: model, clock: source.NewOutlineClock(model), raster: raster, banks: make(map[string]source.Animation)}
	raster.material.ControlOffset[1] = 1
	for _, name := range []string{"late-a", "late-b", "late-c", "late-d", "late-e", "late-f"} {
		e.banks[name], err = source.LoadAnimation(assets.Files, name)
		if err != nil {
			raster.Close()
			return nil, err
		}
	}
	return e, nil
}
func (e *OutlineEffect) Update(f kit.Frame) error {
	tick := max(0, min(source.OutlineTicks-1, int(math.Round(f.Time*FPS))))
	if e.clock.Tick > tick+1 {
		e.clock = source.NewOutlineClock(e.model)
		e.raster.masks.clear()
	}
	for e.clock.Tick <= tick {
		if e.clock.Step() {
			e.raster.masks.prepareOutline(e.clock.Current, e.banks[e.clock.Bank].Frames[e.clock.Frame], e.clock.Mirror)
			e.raster.dirty = true
		}
	}
	return nil
}
func (e *OutlineEffect) Draw(dst *ebiten.Image) {
	e.raster.drawGrid(dst, e.clock.Grid, [17][15]uint16{}, [2]int{e.clock.Display, e.clock.Display})
}
func (e *OutlineEffect) Close() error { return e.raster.Close() }

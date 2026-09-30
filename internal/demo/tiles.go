package demo

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/effects"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

// TileEffect restores the short loading composition's three 224-pixel bands.
// DCK selects palettes over one indexed image per tile without readback.
type TileEffect struct {
	model    source.TileModel
	clock    *source.TileClock
	images   []*ebiten.Image
	renderer *effects.IndexedImageBank
	closed   bool
}

func NewTileEffect() (*TileEffect, error) {
	return newTileEffect("fourth-loader")
}

func newTileEffect(name string) (*TileEffect, error) {
	controller, err := assets.Files.ReadFile("raw/" + name + ".bin")
	if err != nil {
		return nil, err
	}
	packed, err := assets.Files.ReadFile("raw/second-animation.bin")
	if err != nil {
		return nil, err
	}
	model, err := source.ReadTileModel(controller, packed)
	if err != nil {
		return nil, err
	}
	e := &TileEffect{model: model, clock: source.NewTileClock(model)}
	// A temporary decoding palette stores each three-bit index as red=index*17;
	// the authored model and its clock retain their original colors and cues.
	indexed := model
	for index := range indexed.Palette[0] {
		indexed.Palette[0][index] = uint16(index << 8)
	}
	for tile := range indexed.Tiles {
		pixels, err := indexed.Image(tile, 0)
		if err != nil {
			e.Close()
			return nil, err
		}
		e.images = append(e.images, ebiten.NewImageFromImage(pixels))
	}
	palettes := make([]effects.IndexedImagePalette, len(model.Palette))
	for index, words := range model.Palette {
		palettes[index].Colors = make([]color.NRGBA, len(words))
		for entry, word := range words {
			palettes[index].Colors[entry] = source.RGB12(word)
		}
	}
	slots := make([]effects.IndexedImageSlot, len(e.clock.Tiles))
	for band := range slots {
		slots[band].Hidden = true
		slots[band].Options.GeoM.Translate(64, float64(68+band*46))
	}
	e.renderer, err = effects.NewIndexedImageBank(effects.IndexedImageBankConfig{
		Images: e.images, Palettes: palettes, Slots: slots, MaxSlots: len(slots),
		Channel: composite.BitplaneRed, Scale: 15, Select: e.selectTiles,
	})
	if err != nil {
		e.Close()
		return nil, err
	}
	if err := e.renderer.Update(kit.Frame{}); err != nil {
		e.Close()
		return nil, err
	}
	return e, nil
}

func (e *TileEffect) selectTiles(_ kit.Frame, slots []effects.IndexedImageSlot) error {
	for band, tile := range e.clock.Tiles {
		slots[band].Image, slots[band].Palette = tile, e.clock.Palette
		slots[band].Hidden = tile < 0 || e.clock.Palette < 0
	}
	return nil
}

func (e *TileEffect) Update(f kit.Frame) error {
	tick := max(0, min(e.model.Ticks-1, int(math.Round(f.Time*FPS))))
	if e.clock.Tick > tick+1 {
		e.clock = source.NewTileClock(e.model)
	}
	for e.clock.Tick <= tick {
		e.clock.Step()
	}
	return e.renderer.Update(f)
}

func (e *TileEffect) Draw(dst *ebiten.Image) {
	dst.Fill(color.Black)
	if e.clock.Palette < 0 {
		return
	}
	dst.SubImage(image.Rect(0, 0, Width, 280)).(*ebiten.Image).Fill(source.RGB12(e.model.Palette[e.clock.Palette][0]))
	e.renderer.Draw(dst)
	if err := e.renderer.Err(); err != nil {
		panic(err)
	}
}

func (e *TileEffect) Close() error {
	if !e.closed {
		if e.renderer != nil {
			e.renderer.Close()
		}
		for _, image := range e.images {
			image.Deallocate()
		}
		e.closed = true
	}
	return nil
}

package demo

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

// TileEffect restores the short loading composition's three 224-pixel bands.
// The native renderer selects predecoded tile/palette images without readback.
type TileEffect struct {
	model  source.TileModel
	clock  *source.TileClock
	images [2][]*ebiten.Image
	closed bool
}

func NewTileEffect() (*TileEffect, error) {
	controller, err := assets.Files.ReadFile("raw/fourth-loader.bin")
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
	for palette := range e.images {
		for tile := range model.Tiles {
			pixels, err := model.Image(tile, palette)
			if err != nil {
				e.Close()
				return nil, err
			}
			e.images[palette] = append(e.images[palette], ebiten.NewImageFromImage(pixels))
		}
	}
	return e, nil
}

func (e *TileEffect) Update(f kit.Frame) error {
	tick := max(0, min(source.TileTicks-1, int(math.Round(f.Time*FPS))))
	if e.clock.Tick > tick+1 {
		e.clock = source.NewTileClock(e.model)
	}
	for e.clock.Tick <= tick {
		e.clock.Step()
	}
	return nil
}

func (e *TileEffect) Draw(dst *ebiten.Image) {
	dst.Fill(color.Black)
	dst.SubImage(image.Rect(0, 0, Width, 280)).(*ebiten.Image).Fill(source.RGB12(e.model.Palette[e.clock.Palette][0]))
	for band, tile := range e.clock.Tiles {
		var op ebiten.DrawImageOptions
		op.GeoM.Translate(64, float64(68+band*46))
		dst.DrawImage(e.images[e.clock.Palette][tile], &op)
	}
}

func (e *TileEffect) Close() error {
	if !e.closed {
		for _, images := range e.images {
			for _, image := range images {
				image.Deallocate()
			}
		}
		e.closed = true
	}
	return nil
}

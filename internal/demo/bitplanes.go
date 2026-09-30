package demo

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/composite"
)

// Palette and image bounds are established by each scene's constructor.
// A failed draw identifies a broken renderer invariant, like Ebitengine's
// own checks for invalid shader source dimensions.
func drawBitplanes(lookup *composite.BitplanePalette, dst *ebiten.Image, planes []*ebiten.Image, colors []color.NRGBA) {
	if err := lookup.SetPalette(colors); err != nil {
		panic(err)
	}
	if err := lookup.Draw(dst, planes); err != nil {
		panic(err)
	}
}

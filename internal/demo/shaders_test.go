package demo

import (
	"image"
	"image/color"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

var shaderProbe *shaderCheck
var shaderRunError error

// macOS requires Ebitengine's window creation on the process main thread.
func TestMain(m *testing.M) {
	const size, span = 64, 64 * 64 / 8
	bits := make([]byte, 6*span)
	for i := range bits {
		bits[i] = byte(i*53 ^ i/7 ^ i>>3)
	}
	colors := make([]uint16, 64)
	for i := 0; i < 32; i++ {
		colors[i] = uint16(i*73) & 0xfff
		colors[i+32] = (colors[i] & 0xeee) >> 1
	}
	shaderProbe = &shaderCheck{bits: bits, palette: colors}
	ebiten.SetWindowSize(size, size)
	ebiten.SetWindowTitle("Native material verification")
	ebiten.SetVsyncEnabled(false)
	shaderRunError = ebiten.RunGame(shaderProbe)
	os.Exit(m.Run())
}

// Compare the six-plane GPU composition with the independent CPU planar
// decoder, including a material image placed on Ebitengine's shared atlas.
func TestNoiseShaderMatchesPlanarPixelsAcrossImageOrigins(t *testing.T) {
	const size, span = 64, 64 * 64 / 8
	g := shaderProbe
	if shaderRunError != nil {
		t.Fatal(shaderRunError)
	}
	expected, err := source.DecodePlanar(g.bits, source.Planar{Width: size, Height: size, RowStride: size / 8,
		PlaneOffsets: []int{0, span, span * 2, span * 3, span * 4, span * 5}, Palette: g.palette})
	if err != nil {
		t.Fatal(err)
	}
	if g.err != nil {
		t.Fatal(g.err)
	}
	if len(g.pixels) != len(expected.Pix) {
		t.Fatal("GPU verification did not finish")
	}
	for i, value := range expected.Pix {
		if g.pixels[i] != value {
			t.Fatalf("GPU material differs from decoded bitplanes at byte %d: %d != %d", i, g.pixels[i], value)
		}
	}
}

type shaderCheck struct {
	bits    []byte
	palette []uint16
	pixels  []byte
	done    bool
	err     error
}

func (*shaderCheck) Layout(int, int) (int, int) { return 64, 64 }
func (g *shaderCheck) Update() error {
	if g.done {
		return ebiten.Termination
	}
	return nil
}
func (g *shaderCheck) Draw(dst *ebiten.Image) {
	if g.done {
		return
	}
	const size, span = 64, 64 * 64 / 8
	colors := make([]color.NRGBA, len(g.palette))
	for i, word := range g.palette {
		colors[i] = source.RGB12(word)
	}
	lookup, err := composite.NewBitplanePalette(composite.BitplanePaletteConfig{
		Width: size, Height: size, Planes: 6, Palette: colors,
		Channels: []composite.BitplaneChannel{composite.BitplaneAlpha, composite.BitplaneAlpha,
			composite.BitplaneAlpha, composite.BitplaneAlpha, composite.BitplaneAlpha, composite.BitplaneRed},
	})
	if err != nil {
		g.err, g.done = err, true
		return
	}
	defer lookup.Close()
	var planes [6]*ebiten.Image
	for i := range planes {
		pixels := make([]byte, size*size*4)
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				at := (y*size + x) * 4
				if i == 5 {
					pixels[at+3] = 255
				}
				if g.bits[i*span+y*(size/8)+x/8]>>uint(7-x%8)&1 != 0 {
					pixels[at], pixels[at+1], pixels[at+2], pixels[at+3] = 255, 255, 255, 255
				}
			}
		}
		if i == 5 {
			// Atlas placement differs from every unmanaged mask surface.
			planes[i] = ebiten.NewImageFromImage(&image.NRGBA{Pix: pixels, Stride: size * 4, Rect: image.Rect(0, 0, size, size)})
		} else {
			planes[i] = render.NewSurface(size, size)
			planes[i].WritePixels(pixels)
		}
		defer planes[i].Deallocate()
	}
	if err := lookup.Draw(dst, planes[:]); err != nil {
		g.err, g.done = err, true
		return
	}
	g.pixels = make([]byte, size*size*4)
	dst.ReadPixels(g.pixels)
	g.done = true
}

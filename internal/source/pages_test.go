package source

import (
	"github.com/olivierh59500/go-spaceballs/assets"
	"testing"
)

// The indexed adapters must reproduce the independently recovered RGB pages
// before any transition is applied by the native palette shader.
func TestIndexedPagesRetainRecoveredArtworkAndPalettes(t *testing.T) {
	read := func(name string) []byte { b, _ := assets.Files.ReadFile("raw/" + name + ".bin"); return b }
	packed, director := read("pattern"), read("director")
	for _, offset := range []int{0, 0x14c8, 0x1fc8, 0x220c} {
		p, err := TitlePage(packed, director, offset)
		if err != nil {
			t.Fatal(err)
		}
		original, err := TitlePicture(packed, director, offset)
		if err != nil {
			t.Fatal(err)
		}
		for y := 0; y < p.Pixels.Bounds().Dy(); y++ {
			for x := 0; x < p.Pixels.Bounds().Dx(); x++ {
				index := int(p.Pixels.NRGBAAt(x, y).R) / 17
				if RGB12(p.Palette[index]) != original.NRGBAAt(x, y) {
					t.Fatal("title page index/palette changed", offset, x, y)
				}
			}
		}
	}
	p, err := DragonPage(read("credits-packed"), read("second-director"))
	if err != nil {
		t.Fatal(err)
	}
	original, err := Dragon(read("credits-packed"), read("second-director"))
	if err != nil {
		t.Fatal(err)
	}
	for y := 0; y < p.Pixels.Bounds().Dy(); y++ {
		for x := 0; x < p.Pixels.Bounds().Dx(); x++ {
			index := int(p.Pixels.NRGBAAt(x, y).R) / 17
			if RGB12(p.Palette[index]) != original.NRGBAAt(x, y) {
				t.Fatal("dragon page index/palette changed", x, y)
			}
		}
	}
	credit, err := CreditPage(director)
	if err != nil || len(credit.Palette) != 2 || credit.Palette[0] != 0 || credit.Palette[1] != 0xdef {
		t.Fatal("credit color register order changed", err)
	}
	close, err := ClosingPage(read("final-effects"))
	if err != nil {
		t.Fatal(err)
	}
	final, err := ClosingPicture(read("final-effects"))
	if err != nil {
		t.Fatal(err)
	}
	for y := 0; y < 280; y++ {
		for x := 0; x < 352; x++ {
			index := int(close.Pixels.NRGBAAt(x, y).R) / 17
			if RGB12(close.Palette[index]) != final.NRGBAAt(x, y) {
				t.Fatal("closing page index/palette changed", x, y)
			}
		}
	}
}

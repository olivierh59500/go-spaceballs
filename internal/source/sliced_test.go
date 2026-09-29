package source

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/olivierh59500/go-spaceballs/assets"
)

func TestSlicedCuesRewindsZoomAndPalettesMatchOriginalCPU(t *testing.T) {
	read := func(name string) []byte { b, _ := assets.Files.ReadFile("raw/" + name + ".bin"); return b }
	previous, err := ReadAngularModel(read("fifth-effects"), read("fifth-animation"), read("third-animation"))
	if err != nil {
		t.Fatal(err)
	}
	model, err := ReadSlicedModel(read("sixth-effects"), previous.Material)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open("testdata/original-sliced-clock.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil || len(rows) != SlicedTicks {
		t.Fatal("incomplete sliced fixture", len(rows), err)
	}
	c := NewSlicedClock(model)
	bases := map[string]int{"outline-a": 0x410ce, "outline-b": 0x44874}
	for tick, row := range rows {
		if len(row) != 46 {
			t.Fatal("truncated sliced state", tick)
		}
		c.Step()
		lookup := 0x7d2e4
		if c.Zoom {
			lookup = 0x7d924
		}
		got := []int{tick, c.Software, c.ShapeTick, bases[c.Bank], c.Frame, int(uint16(c.remaining)),
			0x5a000 + c.Current*12760, 0x6cb10 + c.Texture*12760, lookup}
		for _, mask := range c.Display {
			got = append(got, 0x5a000+mask*12760)
		}
		for _, word := range c.Palette[:32] {
			got = append(got, int(word))
		}
		for i, value := range got {
			base := 10
			if i == 3 || i >= 6 && i <= 8 {
				base = 16
			}
			want, err := strconv.ParseInt(row[i], base, 32)
			if err != nil || int(want) != value {
				t.Fatalf("tick %d, field %d: %x != original %s (%v)", tick, i, value, row[i], err)
			}
		}
	}
}

package source

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/olivierh59500/go-spaceballs/assets"
)

func TestRibbonPlaybackAndPaletteFlashesMatchOriginalCPU(t *testing.T) {
	data, _ := assets.Files.ReadFile("raw/seventh-effects.bin")
	model, err := ReadRibbonModel(data)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open("testdata/original-ribbons-clock.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil || len(rows) != RibbonTicks {
		t.Fatal("incomplete ribbon fixture", len(rows), err)
	}
	c := NewRibbonClock(model)
	addresses := [4]int{0x5a000, 0x603b0, 0x66760, 0x51bea}
	for tick, row := range rows {
		if len(row) != 26 {
			t.Fatal("truncated ribbon state", tick)
		}
		c.Step()
		lookup := 0x7e66c
		if c.Mirror == 1 {
			lookup = 0x79e28
		} else if c.mode == 2 {
			lookup = 0x7ecac
		}
		got := []int{tick, c.Software, c.Frame, int(uint16(c.remaining)), c.direction, c.mode,
			addresses[c.Current], lookup, addresses[c.Display[0]], addresses[c.Display[1]]}
		for _, word := range c.Palette {
			got = append(got, int(word))
		}
		for i, value := range got {
			base := 10
			if i == 6 || i == 7 {
				base = 16
			}
			want, err := strconv.ParseInt(row[i], base, 64)
			if err != nil || int(want) != value {
				t.Fatalf("tick %d, field %d: %x != original %s (%v)", tick, i, value, row[i], err)
			}
		}
	}
}

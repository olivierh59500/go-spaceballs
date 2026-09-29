package source

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/olivierh59500/go-spaceballs/assets"
)

func TestNoiseCompositionMatchesOriginalCPU(t *testing.T) {
	controller, _ := assets.Files.ReadFile("raw/third-effects.bin")
	material, _ := assets.Files.ReadFile("raw/third-animation.bin")
	model, err := ReadNoiseModel(controller, material)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open("testdata/original-noise-clock.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil || len(rows) != NoiseTicks {
		t.Fatal("incomplete noise controller fixture", len(rows), err)
	}
	c := NewNoiseClock()
	for c.Prime() {
	}
	for tick, row := range rows {
		if len(row) != 43 {
			t.Fatal("truncated noise state", tick, len(row))
		}
		values := make([]int, len(row))
		for i, text := range row {
			base := 10
			if i == 4 || i == 5 {
				base = 16
			}
			v, err := strconv.ParseInt(text, base, 32)
			if err != nil {
				t.Fatal(err)
			}
			values[i] = int(v)
		}
		c.Step()
		got := []int{tick, c.Software, c.Frame, int(uint16(c.remaining)), 0x5a000 + c.Current*12760, 0x6cb10 + c.Texture*12760}
		for _, plane := range c.Display {
			got = append(got, 0x5a000+plane*12760)
		}
		for _, word := range model.Palette[:32] {
			got = append(got, int(word))
		}
		for i, value := range got {
			if value != values[i] {
				t.Fatalf("tick %d, field %d: %x != original %x", tick, i, value, values[i])
			}
		}
	}
}

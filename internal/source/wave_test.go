package source

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/olivierh59500/go-spaceballs/assets"
)

func TestWaveMotionCuesPalettesAndDisplayedPosesMatchOriginalCPU(t *testing.T) {
	controller, _ := assets.Files.ReadFile("raw/fourth-effects.bin")
	packed, _ := assets.Files.ReadFile("raw/later-animation.bin")
	animation, _ := assets.Files.ReadFile("raw/fourth-animation.bin")
	model, err := ReadWaveModel(controller, packed, animation)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open("testdata/original-wave-clock.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil || len(rows) != WaveTicks {
		t.Fatal("incomplete wave fixture", len(rows), err)
	}
	c := NewWaveClock(model)
	bases := map[string]int{"small-turn": 0xb4080, "small-dance": 0xb4bc6, "small-final": 0xb649a}
	addresses := [3]int{0x5a000, 0x603b0, 0x66760}
	for tick, row := range rows {
		if len(row) != 23 {
			t.Fatal("truncated wave state", tick, len(row))
		}
		values := make([]int, len(row))
		for i, text := range row {
			base := 10
			if i == 4 || i >= 7 && i <= 9 {
				base = 16
			}
			v, err := strconv.ParseInt(text, base, 32)
			if err != nil {
				t.Fatal(err)
			}
			values[i] = int(v)
		}
		c.Step()
		motion := model.Material.Motion(tick)
		got := []int{tick, c.Tick, c.Software, c.ShapeTick, bases[c.Bank], c.Frame, int(uint16(c.remaining)),
			addresses[c.current], addresses[c.display] - 4, bases[c.Visible.Bank], c.Visible.Frame,
			motion.X, motion.Y, motion.SecondY, int(motion.Scroll)}
		for _, word := range c.Palette {
			got = append(got, int(word))
		}
		for i, value := range got {
			if value != values[i] {
				t.Fatalf("tick %d, field %d: %x != original %x", tick, i, value, values[i])
			}
		}
	}
}

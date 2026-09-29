package source

import (
	"encoding/csv"
	"github.com/olivierh59500/go-spaceballs/assets"
	"os"
	"strconv"
	"testing"
)

func TestFinaleVariableCadenceCuesColorsAndVisiblePosesMatchOriginalCPU(t *testing.T) {
	controller, _ := assets.Files.ReadFile("raw/final-effects.bin")
	packed, _ := assets.Files.ReadFile("raw/later-animation.bin")
	model, err := ReadFinaleModel(controller, packed)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open("testdata/original-finale-clock.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil || len(rows) != FinaleTicks {
		t.Fatal("incomplete finale fixture", len(rows), err)
	}
	c := NewFinaleClock(model)
	addresses := [3]int{0x5a000, 0x5faa0, 0x65540}
	bases := map[string]int{"final-a": 0x43000, "final-b": 0x43888, "final-c": 0x476fa}
	for tick, row := range rows {
		if len(row) != 24 {
			t.Fatal("truncated finale state", tick)
		}
		c.Step()
		motion := model.Material.Motion(tick)
		displayMode := 0x200
		if c.Show {
			displayMode = 0x3200
		}
		got := []int{tick, c.Tick, c.Software, bases[c.Bank], c.Frame, int(uint16(c.remaining)), c.period, displayMode, addresses[c.current], addresses[c.display] - 4, bases[c.Visible.Bank], c.Visible.Frame, motion.X, motion.Y, motion.SecondY, int(motion.Scroll)}
		for _, word := range c.Palette {
			got = append(got, int(word))
		}
		for i, value := range got {
			base := 10
			if i == 3 || i >= 8 && i <= 10 {
				base = 16
			}
			want, err := strconv.ParseInt(row[i], base, 32)
			if err != nil || value != int(want) {
				t.Fatalf("tick %d field %d: %x != source %s (%v)", tick, i, value, row[i], err)
			}
		}
	}
}

package source

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func TestOpeningClockMatchesOriginalCPU(t *testing.T) {
	f, err := os.Open("testdata/original-opening-clock.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil || len(rows) != OpeningTicks {
		t.Fatal("invalid independent opening fixture", len(rows), err)
	}
	c := NewOpeningClock()
	addresses := [3]int{0x5a000, 0x603b0, 0x66760}
	for tick, row := range rows {
		if len(row) != 13 {
			t.Fatal("incomplete source state", tick)
		}
		values := make([]int, len(row))
		for i, text := range row {
			base := 10
			if i == 3 || i >= 6 && i <= 8 {
				base = 16
			}
			v, err := strconv.ParseInt(text, base, 32)
			if err != nil {
				t.Fatal(err)
			}
			values[i] = int(v)
		}
		c.Step()
		base, remaining := 0x3fec6, 65-c.ShapeTick
		if c.Bank == "opening" {
			base, remaining = 0x80000, 198-c.ShapeTick
		}
		if c.Bank == "first" {
			base, remaining = 0x83dfe, max(-1, 478-c.ShapeTick)
		}
		second := 0
		if c.SecondOnly {
			second = 1
		}
		got := []int{tick, c.Tick, c.ShapeTick, base, int(uint16(remaining)), second, addresses[c.Current],
			addresses[c.Display[0]], addresses[c.Display[1]] + 0x31d8}
		for _, word := range c.Palette {
			got = append(got, int(word))
		}
		for i, v := range got {
			if v != values[i] {
				t.Fatalf("tick %d, field %d: %x != original %x", tick, i, v, values[i])
			}
		}
	}
	before := *c
	if c.Step() || *c != before {
		t.Fatal("opening continued past its source exit")
	}
}

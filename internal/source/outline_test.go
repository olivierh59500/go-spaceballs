package source

import (
	"encoding/binary"
	"encoding/csv"
	"github.com/olivierh59500/go-spaceballs/assets"
	"os"
	"strconv"
	"testing"
)

func TestOutlineCuesMirrorsAndEveryBlockGradientMatchOriginalCPU(t *testing.T) {
	data, _ := assets.Files.ReadFile("raw/ninth-effects.bin")
	model, err := ReadOutlineModel(data)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open("testdata/original-outline-clock.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil || len(rows) != OutlineTicks {
		t.Fatal("incomplete outline fixture", len(rows), err)
	}
	grids, err := os.ReadFile("testdata/original-outline-grids.bin")
	if err != nil || len(grids) != (OutlineTicks/2+1)*17*15*3*2 {
		t.Fatal("incomplete outline grids", len(grids), err)
	}
	c := NewOutlineClock(model)
	check := func(index int) {
		t.Helper()
		for row := 0; row < 17; row++ {
			for col := 0; col < 15; col++ {
				for color, word := range [3]uint16{c.Grid[row][col], 0, 0} {
					at := ((index*17*15+row*15+col)*3 + color) * 2
					want := binary.BigEndian.Uint16(grids[at:])
					if word != want {
						t.Fatalf("grid %d, %d,%d, color %d: %03x != source %03x", index, row, col, color, word, want)
					}
				}
			}
		}
	}
	check(0)
	bases := map[string]int{"late-a": 0x52182, "late-b": 0x53c50, "late-c": 0x5586e, "late-d": 0x56f52, "late-e": 0x57c72, "late-f": 0x58b00}
	for tick, row := range rows {
		if len(row) != 14 {
			t.Fatal("truncated outline state", tick)
		}
		if c.Step() {
			check(c.Software)
		}
		lookup := 0x7d2e4
		if c.Mirror {
			lookup = 0x7d924
		}
		got := []int{tick, c.ShapeTick, bases[c.Bank], c.Frame, int(uint16(c.remaining)), c.cue, 0x5a000 + c.Current*12760, 0x5a000 + c.Display*12760, lookup}
		for _, color := range c.colors.color {
			got = append(got, int(color))
		}
		got = append(got, c.colors.step)
		for i, value := range got {
			base := 10
			if i == 2 || i >= 6 && i <= 8 {
				base = 16
			}
			want, err := strconv.ParseInt(row[i], base, 32)
			if err != nil || value != int(want) {
				t.Fatalf("tick %d, field %d: %x != source %s (%v)", tick, i, value, row[i], err)
			}
		}
	}
}

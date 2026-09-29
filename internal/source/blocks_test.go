package source

import (
	"encoding/binary"
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/olivierh59500/go-spaceballs/assets"
)

func TestBlockClockAndEveryOriginalGradient(t *testing.T) {
	data, _ := assets.Files.ReadFile("raw/second-effects.bin")
	model, err := ReadBlockModel(data)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open("testdata/original-blocks-clock.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows, err := csv.NewReader(file).ReadAll()
	if err != nil || len(rows) != BlocksTicks {
		t.Fatal("incomplete independent block clock", len(rows), err)
	}
	grids, err := os.ReadFile("testdata/original-blocks-grids.bin")
	if err != nil || len(grids) != (BlocksTicks/2+1)*17*15*3*2 {
		t.Fatal("incomplete executed block gradients", len(grids), err)
	}
	c := NewBlocksClock(model)
	checkGrid := func(index int) {
		t.Helper()
		for row := 0; row < 17; row++ {
			for col := 0; col < 15; col++ {
				for plane, word := range [3]uint16{c.Background[row][col], 0, c.Shadow[row][col]} {
					at := ((index*17*15+row*15+col)*3 + plane) * 2
					want := binary.BigEndian.Uint16(grids[at:])
					if word != want {
						t.Fatalf("shape tick %d, grid %d,%d, color %d: %03x != original %03x", index, row, col, plane, word, want)
					}
				}
			}
		}
	}
	checkGrid(0)
	for tick, row := range rows {
		if len(row) != 21 {
			t.Fatal("truncated block state", tick)
		}
		values := make([]int, len(row))
		for i, text := range row {
			base := 10
			if i == 4 {
				base = 16
			}
			value, err := strconv.ParseInt(text, base, 32)
			if err != nil {
				t.Fatal(err)
			}
			values[i] = int(value)
		}
		if c.Step() {
			checkGrid(c.Software)
		}
		done := 0
		if c.Done {
			done = 1
		}
		got := []int{tick, c.Software, int(uint16(c.remaining)), done, 0x5a000 + c.Current*0x31d8,
			0x5a000 + c.Display[0]*0x31d8, 0x5a000 + c.Display[1]*0x31d8}
		for _, corner := range c.primary.color {
			got = append(got, int(corner))
		}
		got = append(got, 0, 0, 0, 0)
		for _, corner := range c.secondary.color {
			got = append(got, int(corner))
		}
		got = append(got, c.primary.step, c.secondary.step)
		for i, value := range got {
			if value != values[i] {
				t.Fatalf("tick %d, field %d: %x != original %x", tick, i, value, values[i])
			}
		}
	}
}

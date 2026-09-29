package source

import (
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/olivierh59500/go-spaceballs/assets"
)

func TestPatternClockMatchesOriginalCPU(t *testing.T) {
	bytes, err := assets.Files.ReadFile("raw/later-animation.bin")
	if err != nil {
		t.Fatal(err)
	}
	p, err := ReadPatternMaterial(bytes)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open("testdata/original-pattern-clock.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows, err := csv.NewReader(file).ReadAll()
	if err != nil || len(rows) != PatternTicks {
		t.Fatalf("invalid independent CPU fixture: %d rows, %v", len(rows), err)
	}
	for tick, row := range rows {
		if len(row) != 17 {
			t.Fatal("truncated source state", tick)
		}
		values := make([]int, len(row))
		for i, text := range row {
			base := 10
			if i == 3 {
				base = 16
			}
			value, err := strconv.ParseInt(text, base, 32)
			if err != nil {
				t.Fatal(err)
			}
			values[i] = int(value)
		}
		got := p.State(tick)
		fields := []int{tick, got.VBL, got.ShapeTick, got.X, got.Y, got.SecondY, int(got.Scroll)}
		want := []int{values[0], values[1], values[2], values[5], values[6], values[7], values[8]}
		for i := range fields {
			if fields[i] != want[i] {
				t.Fatalf("tick %d, field %d: %d != original %d", tick, i, fields[i], want[i])
			}
		}
		for i, word := range got.Palette {
			if int(word) != values[9+i] {
				t.Fatalf("tick %d, color %d: %03x != original %03x", tick, i, word, values[9+i])
			}
		}
		// Check the bank switch and consumed frame count independently as well.
		base, remaining := 0xb4080, 25-got.ShapeTick
		if got.ShapeTick > 26 {
			base, remaining = 0xb4230, max(-1, 209-(got.ShapeTick-26))
		}
		if base != values[3] || uint16(remaining) != uint16(values[4]) {
			t.Fatalf("authored bank cursor differs at %d: %x/%d", tick, base, remaining)
		}
	}
}

func TestPatternMaterialDecodeAndBounds(t *testing.T) {
	bytes, _ := assets.Files.ReadFile("raw/later-animation.bin")
	p, err := ReadPatternMaterial(bytes)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.X) != 500 || len(p.Y) != 477 {
		t.Fatal("authored movement table lengths changed")
	}
	if fmt.Sprintf("%x", sha256.Sum256(p.Bits)) != "fd622c231f209116e7680ac183e0346416c88a4c9681936f0e458ba6ba66fa73" {
		t.Fatal("authored pattern bitmap changed")
	}
	if _, err := ReadPatternMaterial(bytes[:0x80e0]); err == nil {
		t.Fatal("truncated material was accepted")
	}
}

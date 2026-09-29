package source

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/olivierh59500/go-spaceballs/assets"
)

func TestOutlineDuetAndWhiteExitMatchOriginalCPU(t *testing.T) {
	read := func(name string) []byte { b, _ := assets.Files.ReadFile("raw/" + name + ".bin"); return b }
	previous, err := ReadAngularModel(read("fifth-effects"), read("fifth-animation"), read("third-animation"))
	if err != nil {
		t.Fatal(err)
	}
	model, err := ReadDuetModel(read("eighth-effects"), previous.Material)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open("testdata/original-duet-clock.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil || len(rows) != DuetTicks {
		t.Fatal("incomplete duet fixture", len(rows), err)
	}
	c := NewDuetClock(model)
	bases := map[string]int{"two-body-a": 0x41174, "two-body-b": 0x4552a}
	for tick, row := range rows {
		if len(row) != 49 {
			t.Fatal("truncated duet state", tick)
		}
		c.Step()
		endingPointer := 0
		if c.EndingTexture >= 0 {
			endingPointer = 0x6cb10 + c.EndingTexture*12760
		}
		got := []int{tick, c.Software, c.ShapeTick, bases[c.Bank], c.Frame, int(uint16(c.remaining)), c.cue,
			0x5a000 + c.Current*12760, 0x6cb10 + c.Texture*12760, endingPointer, int(c.EndingColors[0]), int(c.EndingColors[1])}
		for _, plane := range c.Display {
			got = append(got, 0x5a000+plane*12760)
		}
		for _, word := range c.Palette[:32] {
			got = append(got, int(word))
		}
		for i, value := range got {
			base := 10
			if i == 3 || i >= 7 && i <= 9 {
				base = 16
			}
			want, err := strconv.ParseInt(row[i], base, 32)
			if err != nil || value != int(want) {
				t.Fatalf("tick %d, field %d: %x != original %s (%v)", tick, i, value, row[i], err)
			}
		}
	}
}

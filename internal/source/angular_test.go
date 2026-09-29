package source

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/olivierh59500/go-spaceballs/assets"
)

func TestAngularCuesMasksAndPaletteMatchOriginalCPU(t *testing.T) {
	controller, _ := assets.Files.ReadFile("raw/fifth-effects.bin")
	material, _ := assets.Files.ReadFile("raw/fifth-animation.bin")
	preceding, _ := assets.Files.ReadFile("raw/third-animation.bin")
	model, err := ReadAngularModel(controller, material, preceding)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open("testdata/original-angular-clock.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil || len(rows) != AngularTicks {
		t.Fatal("incomplete angular fixture", len(rows), err)
	}
	c := NewAngularClock(model)
	bases := map[string]int{"zoom-a": 0x50dc4, "zoom-b": 0x515b8, "zoom-c": 0x5271a, "zoom-d": 0x53952, "zoom-e": 0x546f2}
	for tick, row := range rows {
		if len(row) != 46 {
			t.Fatal("truncated angular state", tick)
		}
		c.Step()
		got := []int{tick, c.Software, c.ShapeTick, bases[c.Bank], c.Frame, int(uint16(c.remaining)), c.age,
			0x5a000 + c.Current*12760, 0x6cb10 + c.Texture*12760}
		for _, mask := range c.Display {
			got = append(got, 0x5a000+mask*12760)
		}
		for _, word := range model.Material.Palette[:32] {
			got = append(got, int(word))
		}
		for i, value := range got {
			base := 10
			if i == 3 || i == 7 || i == 8 {
				base = 16
			}
			want, err := strconv.ParseInt(row[i], base, 32)
			if err != nil || value != int(want) {
				t.Fatalf("tick %d, field %d: %x != original %s (%v)", tick, i, value, row[i], err)
			}
		}
	}
}

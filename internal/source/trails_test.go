package source

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/olivierh59500/go-spaceballs/assets"
)

func TestTrailsClockMatchesOriginalCPU(t *testing.T) {
	data, _ := assets.Files.ReadFile("raw/first-effects.bin")
	model, err := ReadTrailModel(data)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open("testdata/original-trails-clock.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil || len(rows) != TrailsTicks {
		t.Fatal("invalid independent trails fixture", len(rows), err)
	}
	c := NewTrailsClock(model)
	bases := map[string]int{"trails-a": 0x45e62, "trails-b": 0x47212, "trails-c": 0x47a9a, "trails-d": 0x48300}
	for tick, row := range rows {
		if len(row) != 46 {
			t.Fatal("truncated trail state", tick)
		}
		values := make([]int, len(row))
		for i, text := range row {
			base := 10
			if i == 4 || i == 7 || i == 8 {
				base = 16
			}
			v, err := strconv.ParseInt(text, base, 32)
			if err != nil {
				t.Fatal(err)
			}
			values[i] = int(v)
		}
		c.Step()
		lookup := 0x7d2e4
		if c.Zoom {
			lookup = 0x7d924
		}
		got := []int{tick, c.Tick, c.Software, c.ShapeTick, bases[c.Bank], c.Frame, int(uint16(c.remaining)), 0x5a000 + c.Current*0x31d8, lookup}
		for _, buffer := range c.Display {
			got = append(got, 0x5a000+buffer*0x31d8)
		}
		for _, color := range c.Palette {
			got = append(got, int(color))
		}
		for i, value := range got {
			if value != values[i] {
				t.Fatalf("tick %d, field %d: %x != original %x", tick, i, value, values[i])
			}
		}
	}
}

func TestTrailMorphsMatchOriginalPreparedCoordinates(t *testing.T) {
	f, err := os.Open("testdata/original-trails-morphs.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	rows, err := r.ReadAll()
	if err != nil || len(rows) == 0 {
		t.Fatal("missing original trail morphs", err)
	}
	banks := make(map[string]Animation)
	for address, name := range map[string]string{"45e62": "trails-a", "47212": "trails-b", "47a9a": "trails-c", "48300": "trails-d"} {
		banks[address], err = LoadAnimation(assets.Files, name)
		if err != nil {
			t.Fatal(err)
		}
	}
	integer := func(s string) int {
		v, err := strconv.Atoi(s)
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	for _, row := range rows {
		frame, index := integer(row[1]), integer(row[2])
		var shapes []Polygon
		for _, shape := range banks[row[0]].Frames[frame] {
			if shape.Morphed {
				shapes = append(shapes, shape)
			}
		}
		if index >= len(shapes) {
			t.Fatal("missing original trail morph", row[:4])
		}
		p := shapes[index]
		if int(p.Mask) != integer(row[3]) || len(p.Points)*2 != len(row)-4 {
			t.Fatal("trail contour vertex count differs", row[:4])
		}
		for i, point := range p.Points {
			if int(point.X) != integer(row[4+i*2]) || int(point.Y) != integer(row[5+i*2]) {
				t.Fatalf("%s frame %d morph %d point %d: Go %+v, original %s,%s", row[0], frame, index, i, point, row[4+i*2], row[5+i*2])
			}
		}
	}
}

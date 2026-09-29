package source

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/olivierh59500/go-spaceballs/assets"
)

func TestCompleteRecoveredAnimationInventoryHasBoundedContours(t *testing.T) {
	total := 0
	for _, spec := range AnimationBanks {
		bank, err := LoadAnimation(assets.Files, spec.Name)
		if err != nil {
			t.Fatal(spec.Name, err)
		}
		total += len(bank.Frames)
		for frame, polygons := range bank.Frames {
			for _, polygon := range polygons {
				if polygon.Mask > 15 {
					t.Fatal("plane mask escapes source range", spec.Name, frame)
				}
				for _, p := range polygon.Points {
					if p.X < -512 || p.X > 512 || p.Y < -512 || p.Y > 512 {
						t.Fatal("decoded contour escapes byte-coordinate domain", spec.Name, frame, p)
					}
				}
			}
		}
	}
	if len(AnimationBanks) != 32 || total < 3000 {
		t.Fatal("incomplete recovered inventory", len(AnimationBanks), total)
	}
}

// Fixtures execute the original morph reader, including unequal vertex counts.
func TestMorphContoursMatchAllOriginalCPUCoordinates(t *testing.T) {
	banks := make(map[string]Animation)
	for base, name := range map[string]string{"80000": "intro", "83dfe": "first-animation"} {
		bytes, err := assets.Files.ReadFile("raw/" + name + ".bin")
		if err != nil {
			t.Fatal(err)
		}
		bank, err := ReadAnimation(bytes)
		if err != nil {
			t.Fatal(err)
		}
		banks[base] = bank
	}
	f, err := os.Open("testdata/original-morphs.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	rows, err := r.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	integer := func(s string) int {
		i, err := strconv.Atoi(s)
		if err != nil {
			t.Fatal(err)
		}
		return i
	}
	for _, row := range rows {
		frame, index := integer(row[1]), integer(row[2])
		var shapes []Polygon
		for _, p := range banks[row[0]].Frames[frame] {
			if p.Morphed {
				shapes = append(shapes, p)
			}
		}
		if index >= len(shapes) {
			t.Fatal("missing original morph", row[:4])
		}
		p := shapes[index]
		if int(p.Mask) != integer(row[3]) || len(p.Points)*2 != len(row)-4 {
			t.Fatal("morph metadata differs", row[:4], len(p.Points))
		}
		for i, point := range p.Points {
			if int(point.X) != integer(row[4+i*2]) || int(point.Y) != integer(row[5+i*2]) {
				t.Fatalf("%s frame %d morph %d point %d: Go %+v, CPU %s,%s", row[0], frame, index, i, point, row[4+i*2], row[5+i*2])
			}
		}
	}
}

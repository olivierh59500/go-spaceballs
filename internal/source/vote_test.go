package source

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/olivierh59500/go-spaceballs/assets"
)

func TestVoteCardsAndEarlyTransportMatchOriginalCPU(t *testing.T) {
	controller, _ := assets.Files.ReadFile("raw/main-effects.bin")
	packed, _ := assets.Files.ReadFile("raw/second-animation.bin")
	model, err := ReadTileModel(controller, packed)
	if err != nil {
		t.Fatal(err)
	}
	decoded := bytes.Join(model.Tiles, nil)
	if fmt.Sprintf("%x", sha256.Sum256(decoded)) != "c149c46640bd0f4021eac48e59086fa87ddfa27c67a43279fb54ddef263e939e" {
		t.Fatal("21 decoded tiles differ from the original unpacker output")
	}
	f, err := os.Open("testdata/original-vote-clock.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil || len(rows) != 150 {
		t.Fatal("incomplete original tile controller fixture", len(rows), err)
	}
	c := NewTileClock(model)
	for tick, row := range rows {
		if len(row) != 22 {
			t.Fatal("truncated tile controller state", tick)
		}
		c.Step()
		got := []int{tick, c.Tick, c.Tick % 5, c.Updates % 2, 0x5a000 + c.Current*11592}
		got = append(got, c.positions[:]...)
		for _, word := range model.Palette[c.Palette] {
			got = append(got, int(word))
		}
		for i, value := range got {
			base := 10
			if i == 4 {
				base = 16
			}
			want, err := strconv.ParseInt(row[i], base, 32)
			if err != nil || value != int(want) {
				t.Fatalf("tick %d, field %d: %x != original %s (%v)", tick, i, value, row[i], err)
			}
		}
		for band, tile := range c.Tiles {
			if tile < 0 {
				if fmt.Sprintf("%x", sha256.Sum256(make([]byte, 3864))) != row[17+band*2] {
					t.Fatal("uncopied startup band differs", tick, band)
				}
				continue
			}
			address, _ := strconv.ParseInt(row[16+band*2], 16, 32)
			if int(address) != 0x5a000+c.Display*11592+band*3864 {
				t.Fatal("tile display address differs", tick, band)
			}
			if fmt.Sprintf("%x", sha256.Sum256(model.Tiles[tile])) != row[17+band*2] {
				t.Fatal("displayed tile differs from original copied bitmap", tick, band, tile)
			}
		}
	}
}

package source

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/olivierh59500/go-spaceballs/assets"
)

// The independent hash comes from executing the original 0x5273c CPU routine.
func TestPackedPictureMatchesOriginal68000Output(t *testing.T) {
	data, err := assets.Files.ReadFile("raw/pattern.bin")
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := Unpack(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(decoded) != 37356 || fmt.Sprintf("%x", sha256.Sum256(decoded)) != "376babf81457c628169a8500dc30d6fcff0db56d3d38d5ca6fa4c150cc526b7f" {
		t.Fatal("decoded picture differs from the original CPU output")
	}
}

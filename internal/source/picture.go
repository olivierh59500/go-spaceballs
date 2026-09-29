package source

import (
	"encoding/binary"
	"fmt"
	"image"
)

// TitlePicture retains the three prepared planes and authored copper palette.
// Offset selects one of the four large title/credit pictures in the bank.
func TitlePicture(packed, director []byte, offset int) (*image.NRGBA, error) {
	if offset < 0 || offset >= len(packed) {
		return nil, fmt.Errorf("source: packed picture offset outside its bank")
	}
	data, err := Unpack(packed[offset:])
	if err != nil {
		return nil, err
	}
	const pal = 0x52968 - 0x52000
	if len(director) < pal+30 {
		return nil, fmt.Errorf("source: missing title-picture palette")
	}
	colors := make([]uint16, 8)
	for i := range colors {
		colors[i] = binary.BigEndian.Uint16(director[pal+i*4:])
	}
	return DecodePlanar(data, Planar{Width: 352, Height: 283, RowStride: 44,
		PlaneOffsets: []int{0, 0x30a4, 0x6148}, Palette: colors})
}

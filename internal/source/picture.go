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

// Dragon keeps the hires 640-pixel planes selected by the second director.
// Its pointers start four bytes before the unpacked output in cleared memory.
func Dragon(packed, director []byte) (*image.NRGBA, error) {
	data, err := Unpack(packed)
	if err != nil {
		return nil, err
	}
	const pal = 0x3e5f2 - 0x3e000
	if len(data) != 81920 || len(director) < pal+32 {
		return nil, fmt.Errorf("source: missing hires picture data/palette")
	}
	colors := make([]uint16, 16)
	for i := range colors {
		colors[i] = binary.BigEndian.Uint16(director[pal+i*2:])
	}
	memory := make([]byte, len(data)+4)
	copy(memory[4:], data)
	return DecodePlanar(memory, Planar{Width: 640, Height: 256, RowStride: 80,
		PlaneOffsets: []int{0, 0x5000, 0xa000, 0xf000}, Palette: colors})
}

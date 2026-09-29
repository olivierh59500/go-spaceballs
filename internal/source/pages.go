package source

import (
	"encoding/binary"
	"fmt"
	"image"
)

// IndexedPage retains palette indices separately from the authored RGB12 colors.
type IndexedPage struct {
	Pixels  *image.NRGBA
	Palette []uint16
	Hires   bool
}

func TitlePage(packed, director []byte, offset int) (IndexedPage, error) {
	var p IndexedPage
	if offset < 0 || offset >= len(packed) || len(director) < 0x986 {
		return p, fmt.Errorf("source: missing title page")
	}
	data, err := Unpack(packed[offset:])
	if err != nil {
		return p, err
	}
	p.Palette = make([]uint16, 8)
	indices := make([]uint16, 8)
	for i := range p.Palette {
		p.Palette[i] = binary.BigEndian.Uint16(director[0x968+i*4:])
		indices[i] = uint16(i) * 0x111
	}
	p.Pixels, err = DecodePlanar(data, Planar{Width: 352, Height: 283, RowStride: 44, PlaneOffsets: []int{0, 0x30a4, 0x6148}, Palette: indices})
	return p, err
}

func CreditPage(director []byte) (IndexedPage, error) {
	var p IndexedPage
	if len(director) <= 0x996 {
		return p, fmt.Errorf("source: missing credit page")
	}
	data, err := Unpack(director[0x996:])
	if err != nil {
		return p, err
	}
	p.Palette = []uint16{0, 0xdef}
	p.Pixels, err = DecodePlanar(data, Planar{Width: 352, Height: 290, RowStride: 44, PlaneOffsets: []int{0}, Palette: []uint16{0, 0x111}})
	return p, err
}

func DragonPage(packed, director []byte) (IndexedPage, error) {
	var p IndexedPage
	data, err := Unpack(packed)
	if err != nil {
		return p, err
	}
	if len(director) < 0x612 || len(data) != 81920 {
		return p, fmt.Errorf("source: missing dragon page")
	}
	p.Palette = make([]uint16, 16)
	indices := make([]uint16, 16)
	for i := range p.Palette {
		p.Palette[i] = binary.BigEndian.Uint16(director[0x5f2+i*2:])
		indices[i] = uint16(i) * 0x111
	}
	memory := make([]byte, len(data)+4)
	copy(memory[4:], data)
	p.Hires = true
	p.Pixels, err = DecodePlanar(memory, Planar{Width: 640, Height: 256, RowStride: 80, PlaneOffsets: []int{0, 0x5000, 0xa000, 0xf000}, Palette: indices})
	return p, err
}

func ClosingPage(controller []byte) (IndexedPage, error) {
	var p IndexedPage
	data, err := Unpack(controller[0x1a62:])
	if err != nil {
		return p, err
	}
	p.Palette = make([]uint16, 4)
	for at := 0x113c; at+3 < len(controller); at += 4 {
		reg := binary.BigEndian.Uint16(controller[at:])
		if reg == 0xffff {
			break
		}
		if reg >= 0x180 && reg <= 0x186 {
			p.Palette[(reg-0x180)/2] = binary.BigEndian.Uint16(controller[at+2:])
		}
	}
	p.Pixels, err = DecodePlanar(data, Planar{Width: 352, Height: 280, RowStride: 44, PlaneOffsets: []int{0, 0x3020}, Palette: []uint16{0, 0x111, 0x222, 0x333}})
	return p, err
}

package source

import (
	"fmt"
	"image"
	"image/color"
)

// Planar describes the original bitplanes, including interleaved font rows.
type Planar struct {
	Width, Height, RowStride int
	PlaneOffsets             []int
	Palette                  []uint16
	TransparentZero          bool
}

func RGB12(value uint16) color.NRGBA {
	return color.NRGBA{R: uint8(value>>8&15) * 17, G: uint8(value>>4&15) * 17, B: uint8(value&15) * 17, A: 255}
}

// DecodePlanar keeps original bit order and palette indices; no scaling occurs.
func DecodePlanar(data []byte, spec Planar) (*image.NRGBA, error) {
	if spec.Width < 1 || spec.Height < 1 || spec.Width > 8192 || spec.Height > 8192 ||
		spec.RowStride < (spec.Width+7)/8 || len(spec.PlaneOffsets) < 1 || len(spec.PlaneOffsets) > 8 ||
		len(spec.Palette) < 1<<len(spec.PlaneOffsets) {
		return nil, fmt.Errorf("source: invalid planar image")
	}
	for _, offset := range spec.PlaneOffsets {
		if offset < 0 || offset+(spec.Height-1)*spec.RowStride+(spec.Width+7)/8 > len(data) {
			return nil, fmt.Errorf("source: truncated bitplane")
		}
	}
	result := image.NewNRGBA(image.Rect(0, 0, spec.Width, spec.Height))
	for y := 0; y < spec.Height; y++ {
		for x := 0; x < spec.Width; x++ {
			index := 0
			for plane, offset := range spec.PlaneOffsets {
				bit := data[offset+y*spec.RowStride+x/8] >> uint(7-x%8) & 1
				index |= int(bit) << plane
			}
			if index == 0 && spec.TransparentZero {
				continue
			}
			result.SetNRGBA(x, y, RGB12(spec.Palette[index]))
		}
	}
	return result, nil
}

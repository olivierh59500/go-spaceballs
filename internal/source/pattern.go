package source

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
)

const PatternTicks = 480

// PatternMaterial is the authored one-bit, 640 by 512 spirograph bitmap and
// its two circular movement tables. Two independent views provide two planes.
type PatternMaterial struct {
	Bits []byte
	X, Y []uint16
}

func ReadPatternMaterial(data []byte) (PatternMaterial, error) {
	var p PatternMaterial
	if len(data) < 0x8882 {
		return p, fmt.Errorf("source: truncated pattern material")
	}
	var err error
	p.Bits, err = Unpack(data)
	if err != nil {
		return p, err
	}
	if len(p.Bits) != 640*512/8 {
		return p, fmt.Errorf("source: unexpected pattern bitmap size")
	}
	words := func(bytes []byte) []uint16 {
		out := make([]uint16, len(bytes)/2)
		for i := range out {
			out[i] = binary.BigEndian.Uint16(bytes[i*2:])
		}
		return out
	}
	p.X, p.Y = words(data[0x80e0:0x84c8]), words(data[0x84c8:0x8882])
	for _, x := range p.X {
		if x > 288 {
			return p, fmt.Errorf("source: pattern X escapes its bitmap")
		}
	}
	for _, y := range p.Y {
		if y > 222 {
			return p, fmt.Errorf("source: pattern Y escapes its bitmap")
		}
	}
	return p, nil
}

func (p PatternMaterial) Image() *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, 640, 512))
	for y := 0; y < 512; y++ {
		for x := 0; x < 640; x++ {
			value := uint8(0)
			if p.Bits[y*80+x/8]>>uint(7-x%8)&1 != 0 {
				value = 255
			}
			dst.SetNRGBA(x, y, color.NRGBA{value, value, value, 255})
		}
	}
	return dst
}

// PatternState preserves the PAL VBL clock and the half-rate contour clock.
// Frame is -1 until a completed buffer becomes visible, or after the bank ends.
type PatternState struct {
	VBL, ShapeTick, Frame int
	X, Y, SecondY         int
	Scroll                uint16
	Palette               [8]uint16
}

func (p PatternMaterial) State(tick int) PatternState {
	tick = max(0, min(tick, PatternTicks-1))
	s := p.Motion(tick)
	s.ShapeTick, s.Frame = tick/2+1, tick/2-1
	if s.Frame >= 236 {
		s.Frame = -1
	}
	if s.ShapeTick >= 3 {
		s.Palette = [8]uint16{0, 0, 0x170, 0, 0x707, 0, 0x711, 0}
	}
	if tick >= 20 {
		step := min(tick-20, 16)
		s.Palette[3] = BlendRGB12(0, 0xd00, step, 16)
		s.Palette[5] = BlendRGB12(0, 0xd70, step, 16)
		s.Palette[7] = BlendRGB12(0, 0xbb0, step, 16)
	}
	if tick >= 460 {
		step := min(tick-460, 16)
		for i, c := range s.Palette {
			s.Palette[i] = BlendRGB12(c, 0, step, 16)
		}
	}
	return s
}

// Motion is shared by both spirograph passages. It leaves animation/palette
// timing to their controllers and preserves the two circular material views.
func (p PatternMaterial) Motion(tick int) PatternState {
	s := PatternState{VBL: tick + 1, Frame: -1}
	// The initialization runs the same pointer step once before the first VBL.
	s.X = int(p.X[(3*(tick+2))%len(p.X)])
	s.Y = int(p.Y[(2*(tick+2))%len(p.Y)])
	s.SecondY = int(p.Y[(3*(tick+2))%len(p.Y)])
	s.Scroll = uint16((^s.X & 15) << 4)
	return s
}

// BlendRGB12 follows the original signed, per-nibble integer interpolation.
func BlendRGB12(from, to uint16, numerator, denominator int) uint16 {
	var out uint16
	for shift := 0; shift < 12; shift += 4 {
		a, b := int(from>>shift&15), int(to>>shift&15)
		out |= uint16(a+(b-a)*numerator/denominator) << shift
	}
	return out
}

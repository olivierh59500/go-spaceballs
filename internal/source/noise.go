package source

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
)

const NoiseTicks = 350

type NoiseModel struct {
	Palette [64]uint16
	Bits    [][]byte
}

func ReadNoiseModel(controller, material []byte) (NoiseModel, error) {
	var out NoiseModel
	if len(controller) < 0xc1c || len(material) < 4*12760 {
		return out, fmt.Errorf("source: truncated noise composition")
	}
	word := func(at int) uint16 { return binary.BigEndian.Uint16(controller[at:]) }
	background, first, second, third := word(0xae8), word(0xaea), word(0xaec), word(0xaee)
	out.Palette[0] = background
	for i := 0; i < 31; i++ {
		index := int(word(0x426 + i*2))
		if index > 31 {
			return out, fmt.Errorf("source: noise palette permutation outside its table")
		}
		switch {
		case i < 10:
			out.Palette[index] = BlendRGB12(first, second, i, 11)
		case i < 21:
			out.Palette[index] = BlendRGB12(second, third, i-10, 11)
		default:
			out.Palette[index] = BlendRGB12(second, third, 2, 11)
		}
	}
	for i := 0; i < 32; i++ {
		out.Palette[i+32] = (out.Palette[i] & 0xeee) >> 1
	}
	for i := 0; i < 4; i++ {
		out.Bits = append(out.Bits, append([]byte(nil), material[i*12760:(i+1)*12760]...))
	}
	return out, nil
}

func MonochromeImage(bits []byte) (*image.NRGBA, error) {
	if len(bits) != 12760 {
		return nil, fmt.Errorf("source: wrong 352 by 290 bitplane size")
	}
	dst := image.NewNRGBA(image.Rect(0, 0, 352, 290))
	for y := 0; y < 290; y++ {
		for x := 0; x < 352; x++ {
			value := uint8(0)
			if bits[y*44+x/8]>>uint(7-x%8)&1 != 0 {
				value = 255
			}
			dst.SetNRGBA(x, y, color.NRGBA{value, value, value, 255})
		}
	}
	return dst, nil
}

type NoiseClock struct {
	Tick, Software, Frame int
	Current, Texture      int
	Display               [5]int
	remaining, next       int
	phase, primed         int
}

func NewNoiseClock() *NoiseClock {
	return &NoiseClock{Frame: -1, remaining: 14, next: 7, phase: 1}
}

// Prime performs the seven source startup draws before interrupts are enabled.
func (c *NoiseClock) Prime() bool {
	if c.primed >= 7 {
		return false
	}
	c.primed++
	c.prepare()
	return true
}

func (c *NoiseClock) Step() bool {
	if c.Tick >= NoiseTicks {
		return false
	}
	c.Tick++
	if c.Tick&1 != 0 {
		return false
	}
	c.Software++
	c.phase = (c.phase + 1) % 8
	c.Texture = c.phase / 2
	c.prepare()
	return true
}

func (c *NoiseClock) prepare() {
	c.Display[0] = c.Current
	for i := 1; i < 5; i++ {
		c.Display[i] = (c.Current + i + 1) % 6
	}
	c.Current = (c.Current + 1) % 6
	if c.remaining == -1 {
		c.remaining, c.next = 21, 0
	}
	c.remaining--
	c.Frame, c.next = c.next, c.next+1
}

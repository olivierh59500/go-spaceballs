package source

import (
	"encoding/binary"
	"fmt"
)

const SlicedTicks = 1618

type SlicedModel struct {
	Material NoiseModel
	Cues     []TrailCue
}

func ReadSlicedModel(controller []byte, material NoiseModel) (SlicedModel, error) {
	out := SlicedModel{Material: material}
	if len(controller) < 0x10ce {
		return out, fmt.Errorf("source: truncated sliced-dancer controller")
	}
	word := func(at int) uint16 { return binary.BigEndian.Uint16(controller[at:]) }
	colors := [4]uint16{word(0xf9a), word(0xf9c), word(0xf9e), word(0xfa0)}
	out.Material.Palette[0] = colors[0]
	for i := 0; i < 31; i++ {
		index := int(word(0x796 + i*2))
		if index > 31 {
			return out, fmt.Errorf("source: sliced palette index escapes its table")
		}
		switch {
		case i < 10:
			out.Material.Palette[index] = BlendRGB12(colors[1], colors[2], i+1, 11)
		case i < 21:
			out.Material.Palette[index] = BlendRGB12(colors[2], colors[3], i-9, 11)
		default:
			out.Material.Palette[index] = BlendRGB12(colors[2], colors[3], 2, 11)
		}
	}
	for i := 0; i < 32; i++ {
		out.Material.Palette[i+32] = (out.Material.Palette[i] & 0xeee) >> 1
	}
	for at := 0x452; at < 0x68a; at += 8 {
		base := binary.BigEndian.Uint32(controller[at:])
		bank := "outline-a"
		if base == 0x44874 {
			bank = "outline-b"
		} else if base != 0x410ce {
			return out, fmt.Errorf("source: invalid sliced cue bank")
		}
		out.Cues = append(out.Cues, TrailCue{bank, int(binary.BigEndian.Uint32(controller[at+4:]))})
	}
	return out, nil
}

type SlicedClock struct {
	Tick, Software, ShapeTick int
	Current, Texture          int
	Display                   [5]int
	Bank                      string
	Frame                     int
	Zoom, Draw                bool
	Palette                   [64]uint16
	remaining, next, cue      int
	phase                     int
	model                     SlicedModel
	flags                     [4]bool
	colorToggle               bool
}

func NewSlicedClock(model SlicedModel) *SlicedClock {
	c := &SlicedClock{model: model, Palette: model.Material.Palette, Bank: "outline-a", Frame: -1, remaining: 229, phase: 1}
	for range 7 {
		c.display()
	}
	return c
}

func (c *SlicedClock) Step() bool {
	if c.Tick >= SlicedTicks {
		return false
	}
	if c.Tick >= 1608 && c.Tick&1 != 0 {
		c.Palette[0] = (c.Palette[0] - 0x100) & 0xfff
		c.Palette[32] = (c.Palette[0] & 0xeee) >> 1
	}
	c.Tick++
	if c.Tick&1 != 0 {
		return false
	}
	c.Software++
	if c.Software < 4 {
		return false
	}
	c.ShapeTick++
	phase := c.ShapeTick % 11
	if phase == 0 {
		entry := c.cue
		if entry == 33 {
			c.ShapeTick -= 3
		}
		cue := c.model.Cues[entry]
		c.cue = (c.cue + 1) % len(c.model.Cues)
		c.Bank, c.next = cue.Bank, cue.Frame
		length := 232
		if c.Bank == "outline-b" {
			length = 210
		}
		c.remaining = length - cue.Frame - 2
	} else if phase == 10 {
		for i, special := range []int{9, 17, 44, 56} {
			if c.cue == special && !c.flags[i] {
				amount := 3
				if i == 2 {
					amount = 2
				}
				c.ShapeTick -= amount
				c.flags[i] = true
			}
		}
	}
	c.phase = (c.phase + 1) % 8
	c.Texture = c.phase / 2
	c.display()
	c.Draw = c.remaining >= 0
	if c.Draw {
		c.remaining--
		phase := c.ShapeTick % 44
		if phase == 0 {
			c.Zoom = false
		} else if phase == 22 {
			c.Zoom = true
		}
		c.Frame, c.next = c.next, c.next+1
	}
	return true
}

func (c *SlicedClock) display() {
	for i := range c.Display {
		c.Display[i] = c.Current
	}
	phase := int(uint32(int32(c.ShapeTick-1)) % 22)
	count := min(4, phase)
	for i := range count {
		c.Display[5-count+i] = (c.Current + 5 - i) % 6
	}
	if phase == 0 {
		for i := 1; i < 32; i++ {
			c.Palette[i] = ^c.Palette[i] & 0xfff
		}
		c.colorToggle = !c.colorToggle
		c.Palette[0] = 0x500
		if c.colorToggle {
			c.Palette[0] = 0x104
		}
		for i := 0; i < 32; i++ {
			c.Palette[i+32] = (c.Palette[i] & 0xeee) >> 1
		}
	}
	c.Current = (c.Current + 1) % 6
}

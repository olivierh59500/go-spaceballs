package source

import (
	"encoding/binary"
	"fmt"
)

const DuetTicks = 1041

type DuetModel struct {
	Material NoiseModel
	Cues     []TrailCue
	Palette  [2][64]uint16
}

func ReadDuetModel(controller []byte, material NoiseModel) (DuetModel, error) {
	out := DuetModel{Material: material}
	if len(controller) < 0x1174 {
		return out, fmt.Errorf("source: truncated duet controller")
	}
	word := func(at int) uint16 { return binary.BigEndian.Uint16(controller[at:]) }
	for set := range out.Palette {
		colors := [4]uint16{word(0xfb8 + set*8), word(0xfba + set*8), word(0xfbc + set*8), word(0xfbe + set*8)}
		out.Palette[set][0] = colors[0]
		for i := 0; i < 31; i++ {
			index := int(word(0x6c2 + i*2))
			if index > 31 {
				return out, fmt.Errorf("source: duet palette index outside its table")
			}
			switch {
			case i < 10:
				out.Palette[set][index] = BlendRGB12(colors[1], colors[2], i, 11)
			case i < 21:
				out.Palette[set][index] = BlendRGB12(colors[2], colors[3], i-10, 11)
			default:
				out.Palette[set][index] = BlendRGB12(colors[2], colors[3], 2, 11)
			}
		}
		for i := 0; i < 32; i++ {
			out.Palette[set][i+32] = (out.Palette[set][i] & 0xeee) >> 1
		}
	}
	out.Material.Palette = out.Palette[0]
	for at := 0x3a2; at < 0x53a; at += 8 {
		bank := "two-body-a"
		address := binary.BigEndian.Uint32(controller[at:])
		if address == 0x4552a {
			bank = "two-body-b"
		} else if address != 0x41174 {
			return out, fmt.Errorf("source: invalid duet cue bank")
		}
		out.Cues = append(out.Cues, TrailCue{bank, int(binary.BigEndian.Uint32(controller[at+4:]))})
	}
	return out, nil
}

type DuetClock struct {
	Tick, Software, ShapeTick int
	Current, Texture          int
	Display                   [5]int
	Bank                      string
	Frame                     int
	Draw, Ending              bool
	Palette                   [64]uint16
	EndingColors              [2]uint16
	EndingTexture             int
	remaining, next, cue      int
	phase, set, fade          int
	model                     DuetModel
}

func NewDuetClock(model DuetModel) *DuetClock {
	c := &DuetClock{model: model, Bank: "two-body-a", Frame: 1, remaining: 252, next: 2, ShapeTick: 44, phase: 1,
		EndingColors: [2]uint16{0x201, 0x502}, EndingTexture: -1}
	for range 7 {
		c.display()
	}
	c.Palette = model.Palette[0]
	return c
}

func (c *DuetClock) Step() bool {
	if c.Tick >= DuetTicks {
		return false
	}
	c.Tick++
	if c.Tick&1 != 0 {
		return false
	}
	c.Software++
	c.Ending = c.ShapeTick >= 552
	c.phase = (c.phase + 1) % 8
	if !c.Ending {
		c.Texture = c.phase / 2
	} else {
		if c.phase == 0 {
			c.EndingTexture = 0
		} else {
			c.Texture = c.phase / 2
		}
		c.fade++
		c.EndingColors[0] = BlendRGB12(0x201, 0xfff, c.fade, 11)
		c.EndingColors[1] = BlendRGB12(0x502, 0xfff, c.fade, 11)
	}
	if c.Software < 3 {
		return false
	}
	c.ShapeTick++
	if c.ShapeTick%11 == 0 {
		cue := c.model.Cues[c.cue]
		c.cue = (c.cue + 1) % len(c.model.Cues)
		c.Bank, c.next = cue.Bank, cue.Frame
		c.remaining = 255 - cue.Frame
	}
	c.display()
	if c.remaining == -1 {
		if c.Bank != "two-body-b" {
			c.Bank, c.remaining, c.next = "two-body-b", 254, 0
		} else {
			c.remaining++
			c.next--
		}
	}
	c.remaining--
	c.Frame, c.next, c.Draw = c.next, c.next+1, true
	return true
}

func (c *DuetClock) display() {
	for i := range c.Display {
		c.Display[i] = c.Current
	}
	phase := int(uint32(int32(c.ShapeTick-1)) % 11)
	if c.ShapeTick >= 50 && (c.cue < 25 || c.cue > 34) {
		phase = 65535
	}
	count := min(4, phase)
	for i := range count {
		c.Display[5-count+i] = (c.Current + 5 - i) % 6
	}
	if phase == 0 && c.cue >= 25 {
		c.set = 1 - c.set
		c.Palette = c.model.Palette[c.set]
	}
	c.Current = (c.Current + 1) % 6
}

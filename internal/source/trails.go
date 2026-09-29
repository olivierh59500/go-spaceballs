package source

import (
	"encoding/binary"
	"fmt"
)

const TrailsTicks = 566

type TrailCue struct {
	Bank  string
	Frame int
}

type TrailModel struct {
	Cues    []TrailCue
	Palette [32]uint16
}

func ReadTrailModel(data []byte) (TrailModel, error) {
	var out TrailModel
	if len(data) < 0xe62 {
		return out, fmt.Errorf("source: truncated trail controller")
	}
	banks := map[uint32]string{0x45e62: "trails-a", 0x47212: "trails-b", 0x47a9a: "trails-c", 0x48300: "trails-d"}
	for at := 0x3aa; at < 0x482; at += 8 {
		address := binary.BigEndian.Uint32(data[at:])
		bank, ok := banks[address]
		if !ok {
			return out, fmt.Errorf("source: unknown trail cue table %x", address)
		}
		frame := int(binary.BigEndian.Uint32(data[at+4:]))
		length := map[string]int{"trails-a": 69, "trails-b": 50, "trails-c": 37, "trails-d": 49}[bank]
		if frame < 0 || frame >= length-1 {
			return out, fmt.Errorf("source: trail cue frame escapes its bank")
		}
		out.Cues = append(out.Cues, TrailCue{bank, frame})
	}
	word := func(at int) uint16 { return binary.BigEndian.Uint16(data[at:]) }
	background, first, second, third := word(0xd2a), word(0xd2c), word(0xd2e), word(0xd30)
	out.Palette[0] = background
	for i := 0; i < 31; i++ {
		index := int(word(0x554 + i*2))
		if index > 31 {
			return out, fmt.Errorf("source: trail palette permutation outside its table")
		}
		switch {
		case i < 10:
			out.Palette[index] = BlendRGB12(first, second, i+1, 11)
		case i < 21:
			out.Palette[index] = BlendRGB12(second, third, i-9, 11)
		default:
			out.Palette[index] = BlendRGB12(second, third, 2, 11)
		}
	}
	return out, nil
}

// TrailsClock retains the six masks and five independently addressed display
// planes. The authored cue loop changes the pose every eleven shape ticks.
type TrailsClock struct {
	Tick, Software, ShapeTick int
	Current                   int
	Display                   [5]int
	Palette                   [32]uint16
	Bank                      string
	Frame, remaining, next    int
	Zoom, Draw                bool
	model                     TrailModel
	cue                       int
}

func NewTrailsClock(model TrailModel) *TrailsClock {
	c := &TrailsClock{model: model, Bank: "trails-a", Frame: -1, remaining: 66, ShapeTick: 40}
	// Startup clears seven consecutively selected masks and installs the palette.
	for range 7 {
		c.display()
	}
	c.Palette = model.Palette
	return c
}

func (c *TrailsClock) Step() bool {
	if c.Tick >= TrailsTicks {
		return false
	}
	c.Tick++
	if c.Tick&1 != 0 {
		return false
	}
	c.Software++
	if c.Software >= 4 {
		c.ShapeTick++
		if c.ShapeTick >= 52 && c.ShapeTick%11 == 0 {
			cue := c.model.Cues[c.cue]
			c.cue = (c.cue + 1) % len(c.model.Cues)
			c.Bank, c.next = cue.Bank, cue.Frame
			length := 69
			switch c.Bank {
			case "trails-b":
				length = 50
			case "trails-c":
				length = 37
			case "trails-d":
				length = 49
			}
			c.remaining = length - cue.Frame - 2
		}
	}
	c.display()
	c.Draw = c.Software >= 4 && c.remaining >= 0
	if c.Draw {
		c.remaining--
		c.Frame = c.next
		c.next++
		phase := c.ShapeTick % 88
		c.Zoom = c.Bank == "trails-b" && phase >= 22 && phase < 44
	}
	return true
}

func (c *TrailsClock) display() {
	for i := range c.Display {
		c.Display[i] = c.Current
	}
	phase := c.ShapeTick - 1
	if c.ShapeTick <= 50 {
		phase += 4
	}
	phase %= 22
	if phase == 0 {
		if c.ShapeTick >= 55 {
			for i, word := range c.Palette {
				c.Palette[i] = ^word & 0xfff
			}
		}
	} else {
		count := min(phase, 4)
		for i := range count {
			c.Display[5-count+i] = (c.Current + 5 - i) % 6
		}
	}
	c.Current = (c.Current + 1) % 6
}

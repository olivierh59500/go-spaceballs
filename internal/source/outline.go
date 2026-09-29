package source

import (
	"encoding/binary"
	"fmt"
)

const OutlineTicks = 1385

type OutlineModel struct {
	Colors [][4]uint16
	Cues   []TrailCue
}

func ReadOutlineModel(data []byte) (OutlineModel, error) {
	var out OutlineModel
	if len(data) < 0x3182 {
		return out, fmt.Errorf("source: truncated outline controller")
	}
	for at := 0xfb0; at < 0x1048; at += 8 {
		var values [4]uint16
		for i := range values {
			values[i] = binary.BigEndian.Uint16(data[at+i*2:])
		}
		out.Colors = append(out.Colors, values)
	}
	banks := map[uint32]string{0x52182: "late-a", 0x53c50: "late-b", 0x5586e: "late-c", 0x56f52: "late-d", 0x57c72: "late-e", 0x58b00: "late-f"}
	for at := 0x2118; at < 0x2350; at += 8 {
		bank, ok := banks[binary.BigEndian.Uint32(data[at:])]
		if !ok {
			return out, fmt.Errorf("source: unknown outline cue bank")
		}
		out.Cues = append(out.Cues, TrailCue{bank, int(binary.BigEndian.Uint32(data[at+4:]))})
	}
	return out, nil
}

type OutlineClock struct {
	Tick, Software, ShapeTick int
	Current, Display          int
	Bank                      string
	Frame                     int
	Mirror                    bool
	Grid                      [17][15]uint16
	colors                    blockColors
	remaining, next, cue      int
	model                     OutlineModel
}

func NewOutlineClock(model OutlineModel) *OutlineClock {
	c := &OutlineClock{model: model, Bank: "late-a", Frame: -1, remaining: 74, Current: 1, colors: blockColors{values: model.Colors, step: 15}}
	c.colors.advance()
	c.Grid = BlockGrid(c.colors.color)
	return c
}

func (c *OutlineClock) Step() bool {
	if c.Tick >= OutlineTicks {
		return false
	}
	c.Tick++
	if c.Tick&1 != 0 {
		return false
	}
	c.Software++
	c.ShapeTick++
	c.Display = c.Current
	c.Current = (c.Current + 1) % 6
	c.Grid = BlockGrid(c.colors.color)
	if c.ShapeTick%11 == 0 {
		entry := c.cue
		if entry == 33 {
			c.ShapeTick -= 11
		}
		cue := c.model.Cues[entry]
		c.cue = (c.cue + 1) % len(c.model.Cues)
		c.Bank, c.next = cue.Bank, cue.Frame
		length := map[string]int{"late-a": 76, "late-b": 76, "late-c": 50, "late-d": 45, "late-e": 47, "late-f": 37}[cue.Bank]
		c.remaining = length - cue.Frame - 2
	}
	c.colors.advance()
	c.Mirror = c.cue >= 38 || c.Bank == "late-c" || c.Bank == "late-f"
	if c.remaining == -1 {
		c.Bank, c.remaining, c.next = "late-a", 74, 0
	}
	c.remaining--
	c.Frame, c.next = c.next, c.next+1
	return true
}

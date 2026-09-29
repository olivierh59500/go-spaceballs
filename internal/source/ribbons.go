package source

import (
	"encoding/binary"
	"fmt"
)

const RibbonTicks = 816

type RibbonModel struct{ Colors [3][16]uint16 }

func ReadRibbonModel(data []byte) (RibbonModel, error) {
	var out RibbonModel
	if len(data) < 0x149c {
		return out, fmt.Errorf("source: truncated ribbon controller")
	}
	for set := range out.Colors {
		for i := range out.Colors[set] {
			out.Colors[set][i] = binary.BigEndian.Uint16(data[0x142c+set*32+i*2:])
		}
	}
	return out, nil
}

type RibbonPose struct {
	Frame, Mirror int
}

type RibbonClock struct {
	Tick, Software      int
	Current             int
	Display             [2]int
	Pose                [4]RibbonPose
	Palette             [16]uint16
	Frame, Mirror       int
	remaining, next     int
	direction, mode     int
	fadeA, fadeB, fadeC int
	model               RibbonModel
}

func NewRibbonClock(model RibbonModel) *RibbonClock {
	c := &RibbonClock{model: model, Frame: -1, remaining: 143, direction: 1, Mirror: 0, Display: [2]int{1, 2}}
	for i := range c.Pose {
		c.Pose[i].Frame = -1
	}
	return c
}

func (c *RibbonClock) Step() bool {
	if c.Tick >= RibbonTicks {
		return false
	}
	c.Tick++
	n := c.Tick
	blend := func(from, to [16]uint16, step, total int) {
		for i := range c.Palette {
			c.Palette[i] = BlendRGB12(from[i], to[i], step, total)
		}
	}
	var white [16]uint16
	for i := range white {
		white[i] = 0xfff
	}
	if n <= 10 {
		blend([16]uint16{}, white, c.fadeA, 10)
		c.fadeA++
	} else if n <= 21 {
		if n == 11 {
			c.fadeA = 0
		}
		blend(white, c.model.Colors[0], c.fadeA, 10)
		c.fadeA++
	}
	if n >= 281 && n <= 291 {
		blend(c.model.Colors[0], white, c.fadeB, 10)
		c.fadeB++
	} else if n >= 292 && n <= 301 {
		blend(white, c.model.Colors[1], c.fadeC, 9)
		c.fadeC++
	}
	if n == 350 {
		c.fadeB, c.fadeC = 0, 0
	}
	if n >= 548 && n <= 558 {
		blend(c.model.Colors[1], white, c.fadeB, 10)
		c.fadeB++
	} else if n >= 559 && n <= 568 {
		blend(white, c.model.Colors[2], c.fadeC, 9)
		c.fadeC++
	}
	if n&1 != 0 {
		return false
	}
	c.Software++
	c.Display = [2]int{(c.Current + 1) % 4, (c.Current + 2) % 4}
	draw := c.Software >= 2
	if draw {
		if c.remaining == -1 {
			c.Mirror = 1
			if c.mode == 1 {
				c.Mirror = 0
			}
			c.remaining = 143
			c.direction = -c.direction
			if c.mode != 1 {
				c.remaining -= 5
				c.next += c.direction * 5
			}
			c.mode++
		}
		c.remaining--
		c.Frame, c.next = c.next, c.next+c.direction
		c.Pose[c.Current] = RibbonPose{c.Frame, c.Mirror}
	}
	c.Current = (c.Current + 3) % 4
	c.Pose[c.Current] = RibbonPose{Frame: -1}
	return draw
}

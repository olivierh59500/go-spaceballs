package source

import (
	"encoding/binary"
	"fmt"
)

// WaveTicks includes the source's complete programmed fade. The director's
// earlier minimum shape wait remains separate from loading handoff timing.
const WaveTicks, WaveMinimumTicks = 846, 755

type WaveModel struct {
	Material PatternMaterial
	Cues     []TrailCue
	Palette  [8]uint16
}

func ReadWaveModel(controller, packed, animation []byte) (WaveModel, error) {
	var out WaveModel
	if len(controller) < 0xee2 || len(animation) < 0x3442 {
		return out, fmt.Errorf("source: truncated wave composition")
	}
	var err error
	out.Material.Bits, err = Unpack(packed)
	if err != nil {
		return out, err
	}
	if len(out.Material.Bits) != 40960 {
		return out, fmt.Errorf("source: wrong wave material dimensions")
	}
	words := func(start, end int) []uint16 {
		values := make([]uint16, (end-start)/2)
		for i := range values {
			values[i] = binary.BigEndian.Uint16(animation[start+i*2:])
		}
		return values
	}
	out.Material.X, out.Material.Y = words(0x2ca0, 0x3088), words(0x3088, 0x3442)
	banks := map[uint32]string{0xb4080: "small-turn", 0xb4bc6: "small-dance", 0xb649a: "small-final"}
	for at := 0x64a; at < 0x7aa; at += 8 {
		bank, ok := banks[binary.BigEndian.Uint32(controller[at:])]
		if !ok {
			return out, fmt.Errorf("source: unknown wave cue bank")
		}
		out.Cues = append(out.Cues, TrailCue{bank, int(binary.BigEndian.Uint32(controller[at+4:]))})
	}
	for i := range out.Palette {
		out.Palette[i] = binary.BigEndian.Uint16(controller[0xece+i*2:])
	}
	return out, nil
}

type WavePose struct {
	Bank  string
	Frame int
}

type WaveClock struct {
	Tick, Software, ShapeTick int
	Bank                      string
	Frame                     int
	Visible                   WavePose
	Palette                   [8]uint16
	model                     WaveModel
	remaining, next, cue      int
	current, display          int
	poses                     [3]WavePose
	fade                      int
}

func NewWaveClock(model WaveModel) *WaveClock {
	c := &WaveClock{model: model, Bank: "small-turn", Frame: -1, remaining: 36, display: 1}
	for i := range c.poses {
		c.poses[i].Frame = -1
	}
	c.Visible.Frame = -1
	return c
}

func (c *WaveClock) Step() {
	if c.Tick >= WaveTicks {
		return
	}
	if c.Tick >= 828 && c.Tick <= 844 {
		for _, i := range []int{1, 2, 4, 6} {
			c.Palette[i] = BlendRGB12(c.model.Palette[i], 0, c.fade, 16)
		}
		c.fade++
	}
	c.Tick++
	if c.Tick&1 == 0 {
		return
	}
	c.Software++
	if c.Software < 3 {
		return
	}
	c.ShapeTick++
	c.display = (c.current + 1) % 3
	c.Visible = c.poses[c.display]
	phase := c.ShapeTick % 11
	if phase == 0 {
		cue := c.model.Cues[c.cue]
		c.cue = (c.cue + 1) % len(c.model.Cues)
		c.Bank, c.next = cue.Bank, cue.Frame
		count := map[string]int{"small-turn": 37, "small-dance": 101, "small-final": 59}[cue.Bank]
		c.remaining = count - cue.Frame - 2
	} else if phase == 1 {
		for i, word := range c.Palette {
			c.Palette[i] = -word & 0xfff
		}
	}
	if c.remaining >= 0 {
		c.remaining--
		c.Frame = c.next
		c.next++
		c.poses[c.current] = WavePose{c.Bank, c.Frame}
	}
	c.current = (c.current + 2) % 3
	c.poses[c.current] = WavePose{Frame: -1}
	if c.ShapeTick == 3 {
		c.Palette = c.model.Palette
	}
}

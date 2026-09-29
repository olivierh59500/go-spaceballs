package source

import (
	"encoding/binary"
	"fmt"
	"image"
)

const FinaleTicks = 1716

type FinaleModel struct {
	Material PatternMaterial
	Cues     []TrailCue
	Colors   [4][8]uint16
}

func ReadFinaleModel(controller, packed []byte) (FinaleModel, error) {
	var out FinaleModel
	if len(controller) < 0x1a6e {
		return out, fmt.Errorf("source: truncated finale controller")
	}
	var err error
	out.Material.Bits, err = Unpack(packed)
	if err != nil {
		return out, err
	}
	if len(out.Material.Bits) != 40960 {
		return out, fmt.Errorf("source: invalid final material dimensions")
	}
	words := func(start, end int) []uint16 {
		v := make([]uint16, (end-start)/2)
		for i := range v {
			v[i] = binary.BigEndian.Uint16(controller[start+i*2:])
		}
		return v
	}
	out.Material.X, out.Material.Y = words(0x12c0, 0x16a8), words(0x16a8, 0x1a62)
	banks := map[uint32]string{0x43000: "final-a", 0x43888: "final-b", 0x476fa: "final-c"}
	for at := 0x7a4; at < 0xa94; at += 8 {
		bank, ok := banks[binary.BigEndian.Uint32(controller[at:])]
		if !ok {
			return out, fmt.Errorf("source: invalid finale cue")
		}
		out.Cues = append(out.Cues, TrailCue{bank, int(binary.BigEndian.Uint32(controller[at+4:]))})
	}
	for set := range out.Colors {
		for i := range out.Colors[set] {
			index := 2*i + 1
			if i >= 4 {
				index = 2 * (i - 4)
			}
			out.Colors[set][index] = binary.BigEndian.Uint16(controller[0x764+set*16+i*2:])
		}
	}
	return out, nil
}

func ClosingPicture(controller []byte) (*image.NRGBA, error) {
	if len(controller) <= 0x1a62 {
		return nil, fmt.Errorf("source: missing closing picture")
	}
	data, err := Unpack(controller[0x1a62:])
	if err != nil {
		return nil, err
	}
	colors := make([]uint16, 4)
	for at := 0x113c; at+3 < len(controller); at += 4 {
		reg := binary.BigEndian.Uint16(controller[at:])
		if reg == 0xffff {
			break
		}
		if reg >= 0x180 && reg <= 0x186 {
			colors[(reg-0x180)/2] = binary.BigEndian.Uint16(controller[at+2:])
		}
	}
	return DecodePlanar(data, Planar{Width: 352, Height: 280, RowStride: 44, PlaneOffsets: []int{0, 0x3020}, Palette: colors})
}

type FinaleClock struct {
	Tick, Software                             int
	Bank                                       string
	Frame                                      int
	Visible                                    WavePose
	Palette                                    [8]uint16
	Show                                       bool
	period, phase, next, remaining, cue, color int
	current, display                           int
	poses                                      [3]WavePose
	fade                                       int
	saved                                      [8]uint16
	model                                      FinaleModel
}

func NewFinaleClock(model FinaleModel) *FinaleClock {
	c := &FinaleClock{model: model, Bank: "final-b", Frame: -1, remaining: 169, period: 11, display: 1, Show: true}
	for i := range c.poses {
		c.poses[i].Frame = -1
	}
	c.Visible.Frame = -1
	return c
}
func (c *FinaleClock) Step() {
	if c.Tick >= FinaleTicks {
		return
	}
	c.Tick++
	n := c.Tick
	c.Show = true
	if n >= 720 {
		c.Palette[2], c.Palette[4], c.Palette[6] = 0xfff, 0xfff, 0xfff
		c.Show = n&2 != 0
	}
	for _, change := range [][2]int{{720, 8}, {750, 8}, {800, 7}, {950, 6}, {1100, 5}, {1250, 4}, {1400, 3}, {1500, 2}} {
		if n == change[0] {
			c.period = change[1]
		}
	}
	if n >= 1701 {
		if c.fade == 0 {
			c.saved = c.Palette
		}
		for _, i := range []int{0, 1, 3, 5, 7} {
			c.Palette[i] = BlendRGB12(c.saved[i], 0xfff, c.fade, 15)
		}
		c.fade++
	}
	c.phase++
	interval := 2
	if n >= 720 {
		interval = 4
	}
	if c.phase != interval {
		return
	}
	c.phase = 0
	c.Software++
	c.display = (c.current + 1) % 3
	c.Visible = c.poses[c.display]
	c.nextBufferClear()
	phase := c.Software % c.period
	if phase == 0 {
		cue := c.model.Cues[c.cue]
		c.cue = (c.cue + 1) % len(c.model.Cues)
		c.Bank, c.next = cue.Bank, cue.Frame
		count := map[string]int{"final-a": 50, "final-b": 170, "final-c": 130}[cue.Bank]
		c.remaining = count - cue.Frame - 2
	} else if phase == 1 {
		if n <= 720 {
			c.Palette = c.model.Colors[c.color]
			c.color = (c.color + 1) % 4
		} else {
			c.Palette[2], c.Palette[4], c.Palette[6] = 0xfff, 0xfff, 0xfff
			if c.fade == 0 {
				for i, add := range []uint16{0x102, 0x231, 0x322, 0x413} {
					idx := 2*i + 1
					c.Palette[idx] = (c.Palette[idx] + add) & 0x777
				}
			}
		}
	}
	if c.remaining >= 0 {
		c.remaining--
		c.Frame, c.next = c.next, c.next+1
		c.poses[c.current] = WavePose{c.Bank, c.Frame}
	}
	c.current = (c.current + 2) % 3
}
func (c *FinaleClock) nextBufferClear() { c.poses[(c.current+2)%3] = WavePose{Frame: -1} }

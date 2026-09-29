package source

import (
	"encoding/binary"
	"fmt"
)

const AngularTicks = 715

type AngularModel struct {
	Material NoiseModel
	Cues     []TrailCue
}

func ReadAngularModel(controller, material, preceding []byte) (AngularModel, error) {
	var out AngularModel
	if len(controller) < 0xdc4 || len(material) != 50688 || len(preceding) < 51040 {
		return out, fmt.Errorf("source: truncated angular composition")
	}
	// The nine-track overlay leaves the last 352 material bytes resident.
	bits := append([]byte(nil), preceding[:51040]...)
	copy(bits, material)
	for i := 0; i < 4; i++ {
		out.Material.Bits = append(out.Material.Bits, bits[i*12760:(i+1)*12760])
	}
	word := func(at int) uint16 { return binary.BigEndian.Uint16(controller[at:]) }
	colors := [4]uint16{word(0xc90), word(0xc92), word(0xc94), word(0xc96)}
	out.Material.Palette[0] = colors[0]
	for i := 0; i < 31; i++ {
		index := int(word(0x582 + i*2))
		if index > 31 {
			return out, fmt.Errorf("source: angular palette index outside its table")
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
	banks := map[uint32]string{0x50dc4: "zoom-a", 0x515b8: "zoom-b", 0x5271a: "zoom-c", 0x53952: "zoom-d", 0x546f2: "zoom-e"}
	for at := 0x376; at < 0x476; at += 8 {
		bank, ok := banks[binary.BigEndian.Uint32(controller[at:])]
		if !ok {
			return out, fmt.Errorf("source: unknown angular cue bank")
		}
		out.Cues = append(out.Cues, TrailCue{bank, int(binary.BigEndian.Uint32(controller[at+4:]))})
	}
	return out, nil
}

type AngularClock struct {
	Tick, Software, ShapeTick int
	Current, Texture          int
	Display                   [5]int
	Bank                      string
	Frame                     int
	Draw                      bool
	remaining, next, cue      int
	age, phase                int
	oldBank                   string
	model                     AngularModel
}

func NewAngularClock(model AngularModel) *AngularClock {
	c := &AngularClock{model: model, Bank: "zoom-b", oldBank: "zoom-b", Frame: -1, remaining: 50, phase: 1}
	for range 7 {
		c.display()
	}
	return c
}

func (c *AngularClock) Step() bool {
	if c.Tick >= AngularTicks {
		return false
	}
	c.Tick++
	if c.Tick&1 != 0 {
		return false
	}
	c.Software++
	c.phase = (c.phase + 1) % 8
	c.Texture = c.phase / 2
	if c.Software < 6 {
		return false
	}
	c.ShapeTick++
	if c.ShapeTick%11 == 0 {
		cue := c.model.Cues[c.cue]
		c.cue = (c.cue + 1) % len(c.model.Cues)
		c.Bank, c.next = cue.Bank, cue.Frame
		length := map[string]int{"zoom-a": 34, "zoom-b": 53, "zoom-c": 52, "zoom-d": 33, "zoom-e": 18}[cue.Bank]
		c.remaining = length - cue.Frame - 2
	}
	c.display()
	c.age++
	if c.Bank != c.oldBank {
		c.age, c.oldBank = 0, c.Bank
	}
	c.Draw = c.remaining >= 0
	if c.Draw {
		c.remaining--
		c.Frame = c.next
		c.next++
	}
	return true
}

func (c *AngularClock) display() {
	for i := range c.Display {
		c.Display[i] = c.Current
	}
	count := min(4, c.age)
	for i := range count {
		c.Display[5-count+i] = (c.Current + 5 - i) % 6
	}
	c.Current = (c.Current + 1) % 6
}

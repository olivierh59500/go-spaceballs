package source

const OpeningTicks = 961

// OpeningClock retains the resident's two independently addressed display
// planes. During the middle passage, the first plane is held while only the
// second plane rotates through the three working buffers.
type OpeningClock struct {
	Tick, Copper int
	Current      int
	Display      [2]int
	Palette      [4]uint16
	Bank         string
	Frame        int
	SecondOnly   bool
	FillSecond   bool
	Draw         bool
	anchor       uint16
	fade         int
}

func NewOpeningClock() *OpeningClock {
	return &OpeningClock{Current: 2, Display: [2]int{1, 1}, Palette: [4]uint16{0xfff, 0xfff, 0xfff, 0xfff},
		Bank: "hands", Frame: -1, anchor: 0xfff}
}

// Step returns whether a copper tick prepares a new buffer. All fades use the
// old VBL counter, as the source updates colors before incrementing it.
func (c *OpeningClock) Step() bool {
	if c.Tick >= OpeningTicks {
		return false
	}
	t := c.Tick
	if t >= 60 && t <= 92 {
		color := BlendRGB12(0xfff, 0x214, c.fade, 31)
		c.Palette[1], c.Palette[3], c.anchor = color, color, color
		c.fade++
	}
	if t == 115 {
		c.fade = 0
	}
	if t >= 116 && t <= 122 {
		c.Palette[3] = BlendRGB12(c.anchor, 0x103, c.fade, 6)
		c.fade++
	}
	if t == 130 {
		c.fade, c.anchor = 0, c.Palette[3]
	}
	if t >= 389 && t <= 399 {
		c.Palette[3] = BlendRGB12(c.anchor, 0x214, c.fade, 10)
		c.Palette[1] = BlendRGB12(0x214, 0x214, c.fade, 10)
		c.fade++
	}
	if t == 420 {
		c.fade = 0
	}
	if t >= 579 && t <= 599 {
		c.Palette[1] = BlendRGB12(0x214, 0x214, c.fade, 20)
		c.Palette[2] = BlendRGB12(0x214, 0x214, c.fade, 20)
		c.Palette[3] = BlendRGB12(0x214, 0x214, c.fade, 20)
		c.fade++
	}
	c.Tick++
	if t&1 != 0 {
		return false
	}
	c.Copper++
	c.SecondOnly = c.Copper >= 58 && c.Copper <= 199
	// Address order: $5a000, $603b0, $66760. The working-buffer cycle descends.
	visible := (c.Current + 1) % 3
	c.Current = (c.Current + 2) % 3
	c.Display[1] = visible
	if !c.SecondOnly {
		c.Display[0] = visible
	}
	c.FillSecond = c.Tick <= 580 || c.SecondOnly
	switch {
	case c.Copper < 58:
		c.Bank, c.Frame = "hands", c.Copper-1
	case c.Copper <= 199:
		c.Bank, c.Frame = "opening", c.Copper-58
	default:
		c.Bank, c.Frame = "first", 55+c.Copper-200
	}
	c.Draw = c.Frame < mapFrameCount(c.Bank)
	if c.Copper == 3 {
		c.Palette[0], c.Palette[2] = 0x102, 0x102
	}
	if c.Copper == 202 {
		c.Palette[2] = c.Palette[1]
	}
	return true
}

func mapFrameCount(bank string) int {
	switch bank {
	case "hands":
		return 66
	case "opening":
		return 142
	default:
		return 335
	}
}

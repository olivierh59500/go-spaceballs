package source

import (
	"encoding/binary"
	"fmt"
)

const BlocksTicks = 422

type BlockModel struct {
	Background, Shadow [][4]uint16
}

func ReadBlockModel(data []byte) (BlockModel, error) {
	var out BlockModel
	if len(data) < 0x29f0 {
		return out, fmt.Errorf("source: truncated block controller")
	}
	read := func(start, end int) [][4]uint16 {
		var values [][4]uint16
		for at := start; at < end; at += 8 {
			var corners [4]uint16
			for i := range corners {
				corners[i] = binary.BigEndian.Uint16(data[at+i*2:])
			}
			values = append(values, corners)
		}
		return values
	}
	// The zero corner quartet after each loop is its last interpolation target.
	out.Background, out.Shadow = read(0xa72, 0xb12), read(0xf18, 0xfc0)
	return out, nil
}

// BlockBlend follows the generated 16-step lookup, which floors the complete
// weighted component. It differs from the source's signed RGB12 fade routine.
func BlockBlend(a, b uint16, step int) uint16 {
	var color uint16
	for shift := 0; shift < 12; shift += 4 {
		from, to := int(a>>shift&15), int(b>>shift&15)
		color |= uint16((from*(16-step)+to*step)/16) << shift
	}
	return color
}

// BlockGrid keeps the duplicated twelfth/ thirteenth band in the source's
// copper lists. The middle palette entry and both-mask entry are black.
func BlockGrid(corners [4]uint16) [17][15]uint16 {
	var grid [17][15]uint16
	for row := range grid {
		vertical := row + 1
		if row >= 13 {
			vertical--
		}
		left := BlockBlend(corners[0], corners[2], vertical)
		right := BlockBlend(corners[1], corners[3], vertical)
		for col := range grid[row] {
			grid[row][col] = BlockBlend(left, right, col+1)
		}
	}
	return grid
}

type blockColors struct {
	values [][4]uint16
	key    int
	step   int
	color  [4]uint16
}

func (c *blockColors) advance() {
	if c.step == 15 {
		c.step = 0
		if c.key >= len(c.values)-1 {
			c.key = 0
		}
		c.key++
	}
	for i := range c.color {
		c.color[i] = BlockBlend(c.values[c.key-1][i], c.values[c.key][i], c.step+1)
	}
	c.step++
}

type BlocksClock struct {
	Tick, Software, Frame int
	Current               int
	Display               [2]int
	Done                  bool
	Background, Shadow    [17][15]uint16
	primary, secondary    blockColors
	remaining, next       int
}

func NewBlocksClock(model BlockModel) *BlocksClock {
	c := &BlocksClock{Frame: -1, remaining: 209, primary: blockColors{values: model.Background, step: 15},
		secondary: blockColors{values: model.Shadow, step: 15}}
	c.display()
	c.primary.advance()
	c.secondary.color = c.primary.color
	c.Background, c.Shadow = BlockGrid(c.primary.color), BlockGrid(c.secondary.color)
	return c
}

func (c *BlocksClock) Step() bool {
	if c.Tick >= BlocksTicks {
		return false
	}
	c.Tick++
	if c.Tick&1 != 0 {
		return false
	}
	c.Software++
	c.display()
	c.Background, c.Shadow = BlockGrid(c.primary.color), BlockGrid(c.secondary.color)
	c.primary.advance()
	c.secondary.advance()
	if c.remaining == -1 {
		c.remaining, c.next, c.Done = 209, 0, true
	}
	c.remaining--
	c.Frame, c.next = c.next, c.next+1
	return true
}

func (c *BlocksClock) display() {
	c.Display[0], c.Display[1] = c.Current, (c.Current+3)%6
	c.Current = (c.Current + 1) % 6
}

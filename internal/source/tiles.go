package source

import (
	"encoding/binary"
	"fmt"
	"image"
)

const TileTicks = 143

type TileModel struct {
	Tiles        [][]byte
	Cues         [][3]int
	Palette      [2][8]uint16
	Prime, Ticks int
}

func ReadTileModel(controller, packed []byte) (TileModel, error) {
	var out TileModel
	out.Prime, out.Ticks = 4, TileTicks
	out.Cues = make([][3]int, 9)
	cueBase, paletteBase := 0xb9e, 0x1052
	if len(controller) > 0x208 && binary.BigEndian.Uint32(controller[0x1fc:]) == 0x368 {
		out.Prime, out.Ticks = 0, 150
		out.Cues = make([][3]int, 10)
		cueBase, paletteBase = 0xb68, 0x1024
	}
	if len(controller) < 0x1092 {
		return out, fmt.Errorf("source: truncated tile controller")
	}
	for i := 0; i < 21; i++ {
		at := int(binary.BigEndian.Uint32(controller[i*16+2:])) - 0xc0000
		if at < 0 || at >= len(packed) {
			return out, fmt.Errorf("source: tile packer pointer escapes its payload")
		}
		data, err := Unpack(packed[at:])
		if err != nil {
			return out, err
		}
		if len(data) != 3864 {
			return out, fmt.Errorf("source: unexpected tile bitmap dimensions")
		}
		out.Tiles = append(out.Tiles, data)
	}
	for row := range out.Cues {
		for band := range out.Cues[row] {
			address := int(binary.BigEndian.Uint32(controller[cueBase+(row*3+band)*4:]))
			if address < 0xe8000 || (address-0xe8000)%3864 != 0 || (address-0xe8000)/3864 >= len(out.Tiles) {
				return out, fmt.Errorf("source: tile cue pointer escapes decoded materials")
			}
			out.Cues[row][band] = (address - 0xe8000) / 3864
		}
	}
	for palette := range out.Palette {
		for i := range out.Palette[palette] {
			out.Palette[palette][i] = binary.BigEndian.Uint16(controller[paletteBase+palette*32+i*2:])
		}
	}
	return out, nil
}

func (m TileModel) Image(tile, palette int) (*image.NRGBA, error) {
	if tile < 0 || tile >= len(m.Tiles) || palette < 0 || palette > 1 {
		return nil, fmt.Errorf("source: invalid tile/palette selection")
	}
	return DecodePlanar(m.Tiles[tile], Planar{Width: 224, Height: 46, RowStride: 28,
		PlaneOffsets: []int{0, 1288, 2576}, Palette: m.Palette[palette][:]})
}

type TileClock struct {
	Tick, Updates, Palette int
	Tiles                  [3]int
	Current, Display       int
	positions              [3]int
	model                  TileModel
	groups                 [3][3]int
}

func NewTileClock(model TileModel) *TileClock {
	c := &TileClock{model: model, Palette: -1}
	for i := range c.groups {
		for j := range c.groups[i] {
			c.groups[i][j] = -1
		}
	}
	for range model.Prime {
		c.advance()
	}
	return c
}

func (c *TileClock) Step() {
	if c.Tick >= c.model.Ticks {
		return
	}
	c.Tick++
	if c.model.Prime == 0 && c.Tick <= 5 {
		c.advance()
	}
	if c.Tick%5 == 0 {
		c.advance()
	}
}

func (c *TileClock) advance() {
	c.Updates++
	c.Palette = (c.Updates + 1) % 2
	c.positions[0]++
	length := len(c.model.Cues)
	if c.positions[0] >= length {
		c.positions[0] -= length
		c.positions[1]++
		c.positions[2] += 2
	}
	c.positions[1] = (c.positions[1] + 1) % length
	c.positions[2] = (c.positions[2] + 1) % length
	// The source initially advances twice, then retains its last working group.
	// Preserve that transport rather than inventing a modulo-three rotation.
	c.Display = c.Current
	c.Current = min(2, c.Current+1)
	for i := range c.Tiles {
		c.groups[c.Current][i] = c.model.Cues[c.positions[i]][i]
	}
	c.Tiles = c.groups[c.Display]
}

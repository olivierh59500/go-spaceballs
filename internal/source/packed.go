package source

import (
	"encoding/binary"
	"fmt"
)

// Unpack retains the source's reverse-word bit transport and backward copies.
// Its 12-byte header stores four offset widths, output size and packed size.
func Unpack(data []byte) ([]byte, error) {
	if len(data) < 16 {
		return nil, fmt.Errorf("source: truncated packed header")
	}
	length, packed := int(binary.BigEndian.Uint32(data[4:])), int(binary.BigEndian.Uint32(data[8:]))
	if length < 1 || length > 4<<20 || packed < 4 || packed%4 != 0 || packed > len(data)-12 {
		return nil, fmt.Errorf("source: invalid packed dimensions %d/%d", length, packed)
	}
	r := reverseBits{data: data[12 : 12+packed], at: packed}
	r.word = r.pull()
	out := make([]byte, length)
	p := length
	for p > 0 && r.err == nil {
		if r.bits(1) != 0 {
			count := 0
			for {
				n := int(r.bits(3))
				count += n
				if n != 7 || r.err != nil {
					break
				}
			}
			if count < 1 || count > p {
				return nil, fmt.Errorf("source: literal run escapes output at %d with count %d", p, count)
			}
			for range count {
				var value byte
				for bit := 0; bit < 8; bit++ {
					value |= byte(r.bit()) << bit
				}
				p--
				out[p] = value
			}
		} else {
			kind := int(r.bits(2))
			// Header entries are DBF counters, so offsets consume one more bit.
			width, count := int(data[kind])+1, kind+2
			if width < 1 || width > 16 {
				return nil, fmt.Errorf("source: invalid packed offset width")
			}
			if kind == 3 {
				bits, limit := 3, 7
				if r.bits(1) == 0 {
					bits, limit = 7, 127
					count += 14
				}
				for {
					n := int(r.bits(bits))
					count += n
					if n != limit || r.err != nil {
						break
					}
				}
			}
			offset := int(r.bits(width))
			if count > p || p+offset >= len(out) {
				return nil, fmt.Errorf("source: packed match escapes output")
			}
			for range count {
				value := out[p+offset]
				p--
				out[p] = value
			}
		}
	}
	if r.err != nil {
		return nil, r.err
	}
	return out, nil
}

type reverseBits struct {
	data []byte
	at   int
	word uint32
	err  error
}

func (r *reverseBits) pull() uint32 {
	if r.at < 4 {
		r.err = fmt.Errorf("source: packed bitstream exhausted")
		return 0
	}
	r.at -= 4
	return binary.BigEndian.Uint32(r.data[r.at:])
}

func (r *reverseBits) bit() uint32 {
	bit := r.word & 1
	r.word >>= 1
	if r.word == 0 {
		next := r.pull()
		r.word = next>>1 | bit<<31
		bit = next & 1
	}
	return bit
}

func (r *reverseBits) bits(count int) uint32 {
	var result uint32
	for range count {
		result = result<<1 | r.bit()
	}
	return result
}

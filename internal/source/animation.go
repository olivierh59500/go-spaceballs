package source

import (
	"encoding/binary"
	"fmt"
)

type Point struct{ X, Y int16 }
type Polygon struct {
	Mask    byte
	Morphed bool
	Points  []Point
}
type Animation struct{ Frames [][]Polygon }

// ReadAnimation keeps the source's byte contours and referenced morph records.
// The word frame header is a DBF count, followed by relative longword pointers.
func ReadAnimation(data []byte) (Animation, error) {
	var result Animation
	if len(data) < 6 {
		return result, fmt.Errorf("source: truncated animation")
	}
	count := int(binary.BigEndian.Uint16(data)) + 1
	if count > 20000 || 2+count*4 > len(data) {
		return result, fmt.Errorf("source: invalid animation pointer count")
	}
	for frame := 0; frame < count; frame++ {
		at := int(binary.BigEndian.Uint32(data[2+frame*4:]))
		if at < 0 || at >= len(data) {
			return result, fmt.Errorf("source: frame %d pointer escapes its bank", frame)
		}
		polygons := int(data[at])
		at++
		var shapes []Polygon
		for range polygons {
			if at >= len(data) {
				return result, fmt.Errorf("source: truncated polygon command")
			}
			command := data[at]
			at++
			if command < 0xe6 {
				points, err := readBytePoints(data, at)
				if err != nil {
					return result, err
				}
				shapes = append(shapes, Polygon{Mask: (command - 0xd0) >> 1, Points: points})
				at += 1 + len(points)*2
			} else {
				if at+6 > len(data) {
					return result, fmt.Errorf("source: truncated morph record")
				}
				back, forward := int(binary.BigEndian.Uint16(data[at:])), int(binary.BigEndian.Uint16(data[at+2:]))
				first, err := readBytePoints(data, at-back+1)
				if err != nil {
					return result, err
				}
				second, err := readBytePoints(data, at+forward+1)
				if err != nil {
					return result, err
				}
				points, err := morphPoints(first, second, int(data[at+4]), int(data[at+5]))
				if err != nil {
					return result, err
				}
				if len(points) > 0 {
					shapes = append(shapes, Polygon{Mask: ((command | 1) - 0xe4) >> 1, Points: points, Morphed: true})
				}
				at += 6
			}
		}
		result.Frames = append(result.Frames, shapes)
	}
	return result, nil
}

func readBytePoints(data []byte, at int) ([]Point, error) {
	if at < 0 || at >= len(data) || at+1+int(data[at])*2 > len(data) {
		return nil, fmt.Errorf("source: byte contour escapes its bank")
	}
	count := int(data[at])
	points := make([]Point, count)
	for i := range points {
		points[i] = Point{X: int16(data[at+2+i*2]), Y: int16(data[at+1+i*2])}
	}
	return points, nil
}

// Morphs resample the shorter contour with the source's Q3 edge accumulators,
// then use signed integer division for its explicit numerator/denominator.
func morphPoints(first, second []Point, numerator, denominator int) ([]Point, error) {
	if len(first) <= 1 || len(second) <= 1 {
		return nil, nil
	}
	if denominator == 0 {
		return nil, fmt.Errorf("source: zero morph denominator")
	}
	// A0 initially visits the forward reference, while A1 visits the back one.
	first, second = second, first
	if len(second) > len(first) {
		first, second = second, first
		numerator = denominator - numerator
	}
	output := make([]Point, len(first))
	long, short := len(first)-1, len(second)-1
	segment := 0
	for i, target := range first {
		for segment+1 < short && (segment+1)*long/short <= i {
			segment++
		}
		start, end := segment*long/short, (segment+1)*long/short
		p, next := second[segment], second[segment+1]
		x, y := int(p.X), int(p.Y)
		if i == long {
			x, y = int(second[short].X), int(second[short].Y)
		} else {
			x = (x*8 + (i-start)*((int(next.X)-x)*8/(end-start))) >> 3
			y = (y*8 + (i-start)*((int(next.Y)-y)*8/(end-start))) >> 3
		}
		output[i] = Point{X: int16(x + (int(target.X)-x)*numerator/denominator),
			Y: int16(y + (int(target.Y)-y)*numerator/denominator)}
	}
	return output, nil
}

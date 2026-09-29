package source

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

const DiskBytes, TrackBytes = 901120, 5632

type Region struct {
	Name          string
	Track, Tracks int
	Address       uint32
}

// Regions retain the original two loader stages and subsequent transfer order.
var Regions = []Region{
	{"resident", 0, 2, 0x3ec00}, {"intro", 2, 3, 0x80000}, {"loader-music", 5, 7, 0x409e6},
	{"first-animation", 12, 6, 0x83dfe}, {"pattern", 18, 3, 0xa0000}, {"director", 21, 1, 0x52000},
	{"picture", 22, 2, 0x4d000}, {"main-music-head", 24, 36, 0x1000},
	{"main-music-tail", 60, 6, 0x32800}, {"credits-packed", 66, 4, 0xc0000},
	{"second-director", 70, 1, 0x3e000}, {"effects-director", 71, 1, 0x50000},
	{"later-animation", 72, 7, 0xb0000}, {"main-effects", 79, 1, 0x40000},
	{"second-animation", 80, 5, 0xc0000},
	{"first-effects", 85, 3, 0x45000}, {"third-animation", 91, 10, 0x6cb10},
	{"second-effects", 101, 2, 0x50000}, {"background-patterns", 103, 5, 0xa0000},
	{"third-effects", 108, 1, 0x40000}, {"fourth-loader", 109, 1, 0x50000},
	{"fourth-effects", 110, 1, 0x45000}, {"fourth-animation", 111, 3, 0xb4080},
	{"fifth-effects", 114, 4, 0x50000}, {"fifth-animation", 118, 9, 0x6cb10},
	{"sixth-effects", 127, 7, 0x40000}, {"seventh-effects", 134, 5, 0x4c000},
	{"eighth-effects", 139, 7, 0x40000}, {"ninth-effects", 146, 8, 0x4f000},
	{"final-animation", 154, 6, 0x43000}, {"final-effects", 88, 2, 0x40000},
}

func ValidateDisk(data []byte) error {
	if len(data) != DiskBytes || fmt.Sprintf("%x", sha256.Sum256(data)) != "6c60b98bf206921ceea76c7b4be25ad95c8fd81088a493916921ef70ab0cab9f" {
		return fmt.Errorf("source: disk differs from the verified supplied image")
	}
	return nil
}

func ReadRegion(data []byte, region Region) ([]byte, error) {
	start, end := region.Track*TrackBytes, (region.Track+region.Tracks)*TrackBytes
	if region.Tracks < 1 || start < 0 || end > len(data) {
		return nil, fmt.Errorf("source: disk region %s escapes the image", region.Name)
	}
	return append([]byte(nil), data[start:end]...), nil
}

// ReadModule verifies a four-channel MOD header before computing its byte span.
// This rejects the incidental M.K. sequence inside the loader's sample data.
func ReadModule(data []byte, start int) ([]byte, error) {
	if start < 0 || len(data)-start < 1084 {
		return nil, fmt.Errorf("source: truncated module header")
	}
	hdr := data[start : start+1084]
	if string(hdr[1080:]) != "M.K." || hdr[950] < 1 || hdr[950] > 128 {
		return nil, fmt.Errorf("source: invalid module signature/order count")
	}
	patterns, samples := 0, 0
	for _, order := range hdr[952 : 952+int(hdr[950])] {
		if order > 127 {
			return nil, fmt.Errorf("source: invalid module order")
		}
		patterns = max(patterns, int(order)+1)
	}
	for i := 0; i < 31; i++ {
		s := hdr[20+i*30 : 50+i*30]
		if s[24] > 15 || s[25] > 64 {
			return nil, fmt.Errorf("source: invalid sample tuning/volume")
		}
		samples += int(binary.BigEndian.Uint16(s[22:])) * 2
	}
	length := 1084 + patterns*1024 + samples
	if length > len(data)-start {
		return nil, fmt.Errorf("source: truncated module patterns/samples")
	}
	return append([]byte(nil), data[start:start+length]...), nil
}

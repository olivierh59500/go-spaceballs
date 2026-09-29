package source

import (
	"fmt"
	"io/fs"
)

type AnimationBank struct {
	Name, File     string
	Offset, Frames int
}

// AnimationBanks identifies every complete pointer table in the loaded streams.
// Their ordering is an inventory, not a substitute for the source choreography.
var AnimationBanks = []AnimationBank{
	{"opening", "intro", 0, 142}, {"first", "first-animation", 0, 335},
	{"turn", "later-animation", 0x4080, 26}, {"long-turn", "later-animation", 0x4230, 210},
	{"blocks", "background-patterns", 0, 211},
	{"small-turn", "fourth-animation", 0, 37}, {"small-dance", "fourth-animation", 0xb46, 101}, {"small-final", "fourth-animation", 0x241a, 59},
	{"trails-a", "first-effects", 0xe62, 69}, {"trails-b", "first-effects", 0x2212, 50}, {"trails-c", "first-effects", 0x2a9a, 37}, {"trails-d", "first-effects", 0x3300, 49},
	{"middle-turn", "third-effects", 0xc1c, 24},
	{"zoom-a", "fifth-effects", 0xdc4, 34}, {"zoom-b", "fifth-effects", 0x15b8, 53}, {"zoom-c", "fifth-effects", 0x271a, 52},
	{"zoom-d", "fifth-effects", 0x3952, 33}, {"zoom-e", "fifth-effects", 0x46f2, 18},
	{"outline-a", "sixth-effects", 0x10ce, 232}, {"outline-b", "sixth-effects", 0x4874, 210},
	{"close-up", "seventh-effects", 0x149c, 145},
	{"two-body-a", "eighth-effects", 0x1174, 256}, {"two-body-b", "eighth-effects", 0x552a, 256},
	{"late-a", "ninth-effects", 0x3182, 76}, {"late-b", "ninth-effects", 0x4c50, 76}, {"late-c", "ninth-effects", 0x686e, 50},
	{"late-d", "ninth-effects", 0x7f52, 45}, {"late-e", "ninth-effects", 0x8c72, 47}, {"late-f", "ninth-effects", 0x9b00, 37},
	{"final-a", "final-animation", 0, 50}, {"final-b", "final-animation", 0x888, 170}, {"final-c", "final-animation", 0x46fa, 130},
}

func LoadAnimation(files fs.FS, name string) (Animation, error) {
	for _, bank := range AnimationBanks {
		if bank.Name != name {
			continue
		}
		bytes, err := fs.ReadFile(files, "raw/"+bank.File+".bin")
		if err != nil {
			return Animation{}, err
		}
		if bank.Offset < 0 || bank.Offset >= len(bytes) {
			return Animation{}, fmt.Errorf("source: animation bank escapes its payload")
		}
		animation, err := ReadAnimation(bytes[bank.Offset:])
		if err != nil {
			return animation, err
		}
		if len(animation.Frames) != bank.Frames {
			return animation, fmt.Errorf("source: authored animation frame count changed")
		}
		return animation, nil
	}
	return Animation{}, fmt.Errorf("source: unknown animation bank %q", name)
}

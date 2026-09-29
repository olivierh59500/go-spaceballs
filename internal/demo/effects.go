package demo

import (
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

// NewEffect constructs a converted unit through DCK's common effect interface.
// A nil effect identifies a raw asset name handled by the diagnostic viewer.
func NewEffect(name string) (kit.Effect, int, error) {
	var effect kit.Effect
	var ticks int
	var err error
	switch name {
	case "opening-effect":
		effect, err = NewOpeningEffect()
		ticks = source.OpeningTicks
	case "pattern":
		effect, err = NewPatternEffect()
		ticks = source.PatternTicks
	case "trails":
		effect, err = NewTrailsEffect()
		ticks = source.TrailsTicks
	case "blocks-effect":
		effect, err = NewBlocksEffect()
		ticks = source.BlocksTicks
	case "noise":
		effect, err = NewNoiseEffect()
		ticks = source.NoiseTicks
	case "wave":
		effect, err = NewWaveEffect()
		ticks = source.WaveTicks
	}
	if err != nil {
		return nil, 0, err
	}
	return effect, ticks, nil
}

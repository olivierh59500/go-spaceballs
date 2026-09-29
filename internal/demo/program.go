package demo

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/sound"
	playback "github.com/olivierh59500/democonstructionkit/sound/ebiten"
	"github.com/olivierh59500/democonstructionkit/timeline"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/source"
	"image/color"
	"math"
	"time"
)

type Segment struct {
	Name      string
	Ticks     int
	Fade      string
	FadeStart int
	HoldTick  int
}

// ProgramSegments includes both original loading compositions. The early
// decode/drive holds are anchored to the supplied reference recording.
var ProgramSegments = []Segment{
	{"boot", 60, "", 0, 0}, {"opening-effect", 961, "", 0, 0}, {"opening-hold", 66, "", 0, 960},
	{"state", 23, "out", 0, 0}, {"black", 46, "", 0, 0}, {"of", 23, "out", 0, 0}, {"black", 40, "", 0, 0},
	{"the", 23, "out", 0, 0}, {"black", 50, "", 0, 0}, {"art", 23, "out", 0, 0}, {"black", 30, "", 0, 0},
	{"credits", 287, "credits", 0, 0}, {"dragon", 21, "dragon-in", 0, 0}, {"dragon", 361, "", 0, 0}, {"dragon", 17, "dragon-out", 0, 0},
	{"black", 158, "", 0, 0}, {"pattern", source.PatternTicks, "", 0, 0}, {"vote", 150, "", 0, 0},
	{"trails", source.TrailsTicks, "", 0, 0}, {"blocks-effect", source.BlocksTicks, "", 0, 0}, {"noise", source.NoiseTicks, "", 0, 0},
	{"tiles", source.TileTicks, "", 0, 0}, {"wave", source.WaveMinimumTicks, "", 0, 0}, {"angular", source.AngularTicks, "", 0, 0},
	{"sliced", source.SlicedTicks, "", 0, 0}, {"ribbons", source.RibbonTicks, "", 0, 0}, {"duet", source.DuetTicks, "", 0, 0},
	{"outline", source.OutlineTicks, "", 0, 0}, {"finale", source.FinaleTicks, "", 0, 0}, {"closing", 17, "closing-in", 0, 0},
}

type Program struct {
	timing    *timeline.Sequence
	active    kit.Effect
	index     int
	current   Segment
	tick      int
	audio     bool
	music     *playback.Player
	musicName string
	mainStart int
	closing   kit.Effect
	closed    bool
}

func NewProgram(audio bool) (*Program, error) {
	var durations []time.Duration
	at, main := 0, -1
	for _, s := range ProgramSegments {
		durations = append(durations, time.Duration(s.Ticks)*time.Second/FPS)
		if s.Name == "pattern" {
			main = at
		}
		at += s.Ticks
	}
	timing, err := timeline.NewSequence(durations, false)
	if err != nil {
		return nil, err
	}
	return &Program{timing: timing, index: -1, audio: audio, mainStart: main}, nil
}
func (p *Program) Update(f kit.Frame) error {
	if p.closed {
		return nil
	}
	p.tick = max(0, int(math.Round(f.Time*FPS)))
	i, local, ok := p.timing.At(time.Duration(p.tick) * time.Second / FPS)
	if !ok {
		i = len(ProgramSegments) - 1
		local = time.Duration(16) * time.Second / FPS
	}
	if i != p.index {
		if p.active != nil {
			kit.Close(p.active)
		}
		p.index = i
		p.current = ProgramSegments[i]
		effect, err := p.makeEffect(p.current)
		if err != nil {
			return err
		}
		p.active = effect
	}
	if p.audio {
		if err := p.updateMusic(); err != nil {
			return err
		}
	}
	if p.active != nil {
		sample := local.Seconds()
		if p.current.HoldTick > 0 {
			sample = float64(p.current.HoldTick) / FPS
		}
		f.Time = sample
		return p.active.Update(f)
	}
	return nil
}
func (p *Program) makeEffect(s Segment) (kit.Effect, error) {
	if s.Name == "boot" || s.Name == "black" {
		return nil, nil
	}
	if s.Name == "opening-hold" {
		return NewOpeningEffect()
	}
	if s.Name == "state" || s.Name == "of" || s.Name == "the" || s.Name == "art" || s.Name == "credits" || s.Name == "dragon" || s.Name == "closing" {
		e, err := NewPageEffect(s.Name)
		if err != nil {
			return nil, err
		}
		base := e.model.Palette
		zero := make([]uint16, len(base))
		white := make([]uint16, len(base))
		for i := range white {
			white[i] = 0xfff
		}
		switch s.Fade {
		case "out":
			e.Fade(base, zero, s.FadeStart, 23)
		case "credits":
			e.Fade([]uint16{0, 0}, base, 0, 8)
		case "dragon-in":
			from := make([]uint16, len(base))
			for i := range from {
				from[i] = 0xbde
			}
			e.Fade(from, base, 0, 21)
		case "dragon-out":
			to := append([]uint16(nil), base...)
			for i := 1; i < len(to); i++ {
				to[i] = 0
			}
			e.Fade(base, to, 0, 17)
		case "closing-in":
			e.Fade(white, base, 0, 17)
		}
		return e, nil
	}
	e, _, err := NewEffect(s.Name)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, fmt.Errorf("demo: missing production segment %q", s.Name)
	}
	return e, nil
}
func (p *Program) updateMusic() error {
	name := "loader.mod"
	start := time.Duration(292) * time.Second / FPS
	// The recovered module's envelope aligns at 44.12 s in the recording;
	// the visual director starts the first material at 43.78 s.
	if p.tick >= 2206 {
		name = "state-of-the-art.mod"
		start = time.Duration(2206) * time.Second / FPS
	}
	if time.Duration(p.tick)*time.Second/FPS < start {
		return nil
	}
	if p.musicName == name {
		return nil
	}
	if p.music != nil {
		p.music.Close()
	}
	data, err := assets.Files.ReadFile("raw/" + name)
	if err != nil {
		return err
	}
	p.music, err = playback.Open(nil, name, data, sound.Options{Loop: name == "loader.mod", Interpolation: true})
	if err != nil {
		return err
	}
	position := time.Duration(p.tick)*time.Second/FPS - start
	if position > 0 {
		if err := p.music.Seek(position); err != nil {
			return err
		}
	}
	p.music.Play()
	p.musicName = name
	return nil
}
func (p *Program) Draw(dst *ebiten.Image) {
	dst.Fill(color.Black)
	if p.current.Name == "boot" {
		dst.Fill(color.White)
	}
	if p.active != nil {
		p.active.Draw(dst)
	}
}
func (p *Program) Close() error {
	if !p.closed {
		if p.active != nil {
			kit.Close(p.active)
		}
		if p.music != nil {
			p.music.Close()
		}
		p.closed = true
	}
	return nil
}
func (p *Program) Duration() time.Duration { return p.timing.Duration() }

type ProductionGame struct {
	program            *Program
	tick, start, limit int
}

func NewProductionGame(audio bool, start, limit int) (*ProductionGame, error) {
	p, err := NewProgram(audio)
	if err != nil {
		return nil, err
	}
	g := &ProductionGame{program: p, start: max(0, start), limit: max(0, limit)}
	if err := p.Update(kit.Frame{Time: float64(g.start) / FPS}); err != nil {
		p.Close()
		return nil, err
	}
	return g, nil
}
func (g *ProductionGame) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if g.limit > 0 && g.tick >= g.limit {
		return ebiten.Termination
	}
	g.tick++
	return g.program.Update(kit.Frame{Time: float64(g.start+g.tick) / FPS})
}
func (g *ProductionGame) Draw(dst *ebiten.Image)   { g.program.Draw(dst) }
func (*ProductionGame) Layout(int, int) (int, int) { return Width, Height }
func (g *ProductionGame) Close()                   { g.program.Close() }

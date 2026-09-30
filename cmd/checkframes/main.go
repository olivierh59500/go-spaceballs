// Command checkframes fingerprints complete scene frames for renderer migrations.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	kit "github.com/olivierh59500/democonstructionkit"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/go-spaceballs/internal/demo"
)

var banks = [...]string{"opening-effect", "pattern", "trails", "blocks-effect", "noise", "tiles", "vote", "wave", "angular", "sliced", "ribbons", "duet", "outline", "finale"}

type sample struct {
	Bank   string `json:"bank"`
	Tick   int    `json:"tick"`
	SHA256 string `json:"sha256"`
}

type report struct {
	Width   int      `json:"width"`
	Height  int      `json:"height"`
	Rate    int      `json:"ticks_per_second"`
	Frames  int      `json:"drawn_frames"`
	Samples []sample `json:"samples"`
	Mode    string   `json:"mode,omitempty"`
}

type probe struct {
	bank, tick, duration int
	effect               kit.Effect
	surface              *ebiten.Image
	pixels               []byte
	report               report
	done                 bool
	err                  error
	production           *demo.ProductionGame
	full                 bool
	lastUnit             string
	unitStart            int
	segmentStarts        map[int]bool
}

func (*probe) Layout(int, int) (int, int) { return demo.Width, demo.Height }

func (p *probe) Update() error {
	if p.done {
		return ebiten.Termination
	}
	return nil
}

func (p *probe) Draw(dst *ebiten.Image) {
	if p.done {
		return
	}
	if p.surface == nil {
		p.surface = render.NewSurface(demo.Width, demo.Height)
		p.pixels = make([]byte, demo.Width*demo.Height*4)
		p.report.Width, p.report.Height, p.report.Rate = demo.Width, demo.Height, demo.FPS
	}
	if p.full {
		p.drawProduction(dst)
		return
	}
	// Every source update is followed by a draw, including unsampled frames.
	// This keeps retained-mask history identical to normal playback.
	for step := 0; step < 12 && !p.done; step++ {
		if p.effect == nil {
			p.effect, p.duration, p.err = demo.NewEffect(banks[p.bank])
			if p.err != nil || p.effect == nil {
				if p.err == nil {
					p.err = fmt.Errorf("missing scene %q", banks[p.bank])
				}
				p.done = true
				return
			}
		}
		frame := kit.Frame{Tick: uint64(p.tick), Time: float64(p.tick) / demo.FPS, Delta: 1.0 / demo.FPS}
		if p.err = p.effect.Update(frame); p.err != nil {
			p.done = true
			return
		}
		p.surface.Fill(color.Black)
		p.effect.Draw(p.surface)
		p.report.Frames++
		if p.tick < 12 || p.tick%23 == 0 || p.tick >= p.duration-3 {
			p.surface.ReadPixels(p.pixels)
			digest := sha256.Sum256(p.pixels)
			p.report.Samples = append(p.report.Samples, sample{Bank: banks[p.bank], Tick: p.tick, SHA256: hex.EncodeToString(digest[:])})
		}
		p.tick++
		if p.tick == p.duration {
			p.err = kit.Close(p.effect)
			p.effect = nil
			p.bank++
			p.tick = 0
			p.done = p.err != nil || p.bank == len(banks)
		}
	}
	dst.DrawImage(p.surface, nil)
}

func (p *probe) drawProduction(dst *ebiten.Image) {
	if p.production == nil {
		p.production, p.err = demo.NewProductionGame(false, 0, 0)
		if p.err != nil {
			p.done = true
			return
		}
		p.report.Mode = "production"
	}
	for step := 0; step < 12 && !p.done; step++ {
		if p.report.Frames > 0 {
			if p.err = p.production.Update(); p.err != nil {
				p.done = true
				return
			}
		}
		unit, tick := p.production.Position()
		if unit != p.lastUnit || p.segmentStarts[tick] {
			p.lastUnit, p.unitStart = unit, tick
		}
		p.production.Draw(p.surface)
		p.report.Frames++
		if tick-p.unitStart < 24 || tick%23 == 0 || tick >= p.duration-3 {
			p.surface.ReadPixels(p.pixels)
			digest := sha256.Sum256(p.pixels)
			p.report.Samples = append(p.report.Samples, sample{Bank: unit, Tick: tick, SHA256: hex.EncodeToString(digest[:])})
		}
		p.done = tick >= p.duration
	}
	dst.DrawImage(p.surface, nil)
}

func (p *probe) Close() {
	if p.production != nil {
		p.production.Close()
	}
	if p.effect != nil {
		kit.Close(p.effect)
	}
	if p.surface != nil {
		p.surface.Deallocate()
	}
}

func main() {
	output := flag.String("output", "", "write frame fingerprint report as JSON")
	production := flag.Bool("production", false, "fingerprint the complete director including all pages, fades and handoffs")
	flag.Parse()
	if *output == "" {
		log.Fatal("-output is required")
	}
	p := &probe{full: *production}
	if p.full {
		p.segmentStarts = make(map[int]bool, len(demo.ProgramSegments))
		for _, segment := range demo.ProgramSegments {
			p.segmentStarts[p.duration] = true
			p.duration += segment.Ticks
		}
		p.duration += demo.FPS // Include a stationary closing hold.
	}
	defer p.Close()
	ebiten.SetWindowSize(demo.Width, demo.Height)
	ebiten.SetWindowTitle("Spaceballs / frame verification")
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetVsyncEnabled(false)
	if err := ebiten.RunGame(p); err != nil {
		log.Fatal(err)
	}
	if p.err != nil {
		log.Fatal(p.err)
	}
	data, err := json.MarshalIndent(p.report, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*output, append(data, '\n'), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d full frames drawn; %d fingerprints recorded\n", p.report.Frames, len(p.report.Samples))
}

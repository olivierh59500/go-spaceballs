// Command video exports the full native production with both recovered modules.
package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/video"
	"github.com/olivierh59500/go-spaceballs/internal/demo"
)

type recordingHost struct{ *demo.ProductionGame }

// RecordingChapter supplies composition boundaries to DCK's video exporter.
func (h recordingHost) RecordingChapter() string {
	name, _ := h.Position()
	return strings.ReplaceAll(name, "-", " ")
}

func main() {
	config := video.Config{
		Output: "recordings/spaceballs-state-of-the-art.mp4",
		Title:  "Spaceballs: State of the Art Go",
		Width:  demo.Width * 2, Height: demo.Height * 2,
		FPS: demo.FPS, TPS: demo.FPS, SampleRate: 48000,
		Duration: 250 * time.Second, PosterAt: 229760 * time.Millisecond,
	}
	config.Flags(flag.CommandLine)
	flag.Parse()
	if config.Duration <= 0 {
		log.Fatal("The stationary closing page requires a positive recording duration")
	}
	if err := video.Run(config, func() (ebiten.Game, error) {
		game, err := demo.NewProductionGame(true, 0, 0)
		return recordingHost{game}, err
	}); err != nil {
		log.Fatal(err)
	}
}

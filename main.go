package main

import (
	"flag"
	"github.com/hajimehoshi/ebiten/v2"
	capture "github.com/olivierh59500/democonstructionkit/fidelity/ebiten"
	"github.com/olivierh59500/go-spaceballs/internal/demo"
	"log"
)

func main() {
	start := flag.Int("start", 0, "start at this 50 Hz production tick")
	limit := flag.Int("ticks", 0, "optional update limit for validation")
	silent := flag.Bool("silent", false, "disable device music")
	directory := flag.String("capture", "", "capture one native production frame without device audio")
	flag.Parse()
	if *directory != "" {
		if err := capture.Run(capture.Config{Directory: *directory, Frames: []int{0}, Width: demo.Width, Height: demo.Height}, func() (ebiten.Game, error) { return demo.NewProductionGame(false, *start, 0) }); err != nil {
			log.Fatal(err)
		}
		return
	}
	game, err := demo.NewProductionGame(!*silent, *start, *limit)
	if err != nil {
		log.Fatal(err)
	}
	defer game.Close()
	ebiten.SetTPS(demo.FPS)
	ebiten.SetWindowSize(demo.Width*3, demo.Height*3)
	ebiten.SetWindowTitle("Spaceballs: State of the Art Go")
	ebiten.SetRunnableOnUnfocused(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

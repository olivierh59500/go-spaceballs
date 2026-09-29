// Command preview validates original animation and picture assets through DCK.
package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	capture "github.com/olivierh59500/democonstructionkit/fidelity/ebiten"
	"github.com/olivierh59500/go-spaceballs/internal/demo"
)

func main() {
	bank := flag.String("bank", "intro", "source bank or composed effect: intro, first-animation, picture, dragon, pattern")
	frame := flag.Int("frame", 0, "initial source frame")
	offset := flag.Int("picture-offset", 0, "packed picture offset within its source bank")
	directory := flag.String("capture", "", "write one deterministic native asset screenshot")
	flag.Parse()
	if *directory != "" {
		err := capture.Run(capture.Config{Directory: *directory, Frames: []int{0}, Width: demo.Width, Height: demo.Height},
			func() (ebiten.Game, error) { return demo.NewPreview(*bank, *frame, *offset) })
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	game, err := demo.NewPreview(*bank, *frame, *offset)
	if err != nil {
		log.Fatal(err)
	}
	defer game.Close()
	ebiten.SetTPS(demo.FPS)
	ebiten.SetWindowSize(demo.Width*3, demo.Height*3)
	ebiten.SetWindowTitle("Spaceballs / source asset preview")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

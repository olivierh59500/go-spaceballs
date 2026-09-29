// Command preview validates original animation and picture assets through DCK.
package main

import (
	"flag"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	capture "github.com/olivierh59500/democonstructionkit/fidelity/ebiten"
	"github.com/olivierh59500/democonstructionkit/sound"
	playback "github.com/olivierh59500/democonstructionkit/sound/ebiten"
	"github.com/olivierh59500/go-spaceballs/assets"
	"github.com/olivierh59500/go-spaceballs/internal/demo"
)

func main() {
	bank := flag.String("bank", "opening-effect", "source bank or composed effect: opening-effect, hands, intro, first-animation, picture, dragon, pattern, trails, blocks-effect, noise, tiles, wave, angular, sliced")
	frame := flag.Int("frame", 0, "initial source frame")
	offset := flag.Int("picture-offset", 0, "packed picture offset within its source bank")
	directory := flag.String("capture", "", "write one deterministic native asset screenshot")
	music := flag.Bool("music", false, "play the original module through DCK during an interactive preview")
	limit := flag.Int("ticks", 0, "stop an interactive validation run after this many updates; zero keeps running")
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
	game.SetUpdateLimit(*limit)
	if *music {
		name := "state-of-the-art.mod"
		if *bank == "opening-effect" || *bank == "hands" || *bank == "intro" || *bank == "first-animation" || *bank == "picture" || *bank == "dragon" {
			name = "loader.mod"
		}
		data, err := assets.Files.ReadFile("raw/" + name)
		if err != nil {
			log.Fatal(err)
		}
		player, err := playback.Open(nil, name, data, sound.Options{Loop: true, Interpolation: true})
		if err != nil {
			log.Fatal(err)
		}
		defer player.Close()
		position := time.Duration(max(0, *frame)) * time.Second / demo.FPS
		if *bank == "trails" {
			position += 12600 * time.Millisecond
		}
		if *bank == "blocks-effect" {
			position += 23920 * time.Millisecond
		}
		if *bank == "noise" {
			position += 32360 * time.Millisecond
		}
		if *bank == "wave" {
			position += 42220 * time.Millisecond
		}
		if *bank == "tiles" {
			position += 39360 * time.Millisecond
		}
		if err := player.Seek(position); err != nil {
			log.Fatal(err)
		}
		player.Play()
	}
	ebiten.SetTPS(demo.FPS)
	ebiten.SetWindowSize(demo.Width*3, demo.Height*3)
	ebiten.SetWindowTitle("Spaceballs / source asset preview")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

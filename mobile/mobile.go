// Package mobile hosts the same native production on Android.
package mobile

import (
	"log"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	enginemobile "github.com/hajimehoshi/ebiten/v2/mobile"
	"github.com/olivierh59500/go-spaceballs/internal/demo"
)

type startRequest struct {
	tick   int
	verify bool
}

type host struct {
	game         *demo.ProductionGame
	requests     chan startRequest
	verification bool
	lastUnit     string
	lastReport   time.Time
}

var gameHost = &host{requests: make(chan startRequest, 1)}

func init() {
	ebiten.SetTPS(demo.FPS)
	enginemobile.SetGame(gameHost)
}

func (h *host) Update() error {
	request, reset := startRequest{}, false
	select {
	case request = <-h.requests:
		reset = true
	default:
	}
	if h.game == nil || reset {
		if h.game != nil {
			h.game.Close()
		}
		// Create audio and graphics only after Android has installed its view.
		var err error
		h.game, err = demo.NewProductionGame(!request.verify || request.tick == 0, request.tick, 0)
		if err != nil {
			return err
		}
		h.verification = request.verify
		h.lastUnit, h.lastReport = "", time.Time{}
	}
	if err := h.game.Update(); err != nil {
		return err
	}
	unit, tick := h.game.Position()
	if unit != h.lastUnit {
		log.Printf("spaceballs_scene name=%s tick=%d", unit, tick)
		h.lastUnit = unit
	}
	if h.verification && time.Since(h.lastReport) >= 5*time.Second {
		var memory runtime.MemStats
		runtime.ReadMemStats(&memory)
		log.Printf("spaceballs_verify tick=%d tps=%.1f fps=%.1f heap_mib=%.1f", tick, ebiten.ActualTPS(), ebiten.ActualFPS(), float64(memory.HeapAlloc)/(1<<20))
		h.lastReport = time.Now()
	}
	return nil
}

func (h *host) Draw(dst *ebiten.Image) {
	if h.game != nil {
		h.game.Draw(dst)
	}
}

func (*host) Layout(int, int) (int, int) { return demo.Width, demo.Height }

// ConfigureVerification starts a checkpoint without changing its 50 Hz clock.
// Tick zero plays both soundtracks; positive checkpoints use a muted seek.
func ConfigureVerification(tick int) bool {
	if tick < 0 || tick > 13000 {
		return false
	}
	select {
	case gameHost.requests <- startRequest{tick: tick, verify: true}:
		return true
	default:
		return false
	}
}

// Restart starts a fresh production when Android recreates the activity.
func Restart() {
	select {
	case <-gameHost.requests:
	default:
	}
	gameHost.requests <- startRequest{}
}

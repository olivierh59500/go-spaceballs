package demo

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

// maskRing is shared by productions that retain six monochrome source poses.
// A controller chooses the working/display masks; this renderer owns geometry.
type maskRing struct {
	planes [6]*ebiten.Image
	white  *ebiten.Image
	batch  *render.Batch
}

func newMaskRing() *maskRing {
	r := &maskRing{white: ebiten.NewImage(1, 1), batch: render.NewBatch(4096)}
	r.white.Fill(color.White)
	for i := range r.planes {
		r.planes[i] = render.NewSurface(Width, Height)
	}
	r.batch.Options.FillRule = ebiten.FillRuleEvenOdd
	return r
}

func (r *maskRing) clear() {
	for _, plane := range r.planes {
		plane.Clear()
	}
}

func (r *maskRing) prepare(index int, frame []source.Polygon, zoom bool) {
	plane := r.planes[index]
	plane.Clear()
	r.batch.Begin(plane, r.white)
	for _, polygon := range frame {
		filledContour(r.batch, polygon.Points, zoom)
	}
	r.batch.Flush()
}

func (r *maskRing) close() {
	for _, plane := range r.planes {
		plane.Deallocate()
	}
	r.white.Deallocate()
}

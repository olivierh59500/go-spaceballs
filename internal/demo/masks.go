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

// prepareSliced preserves the source's edge-endpoint Y exchange in its second
// bank. Ray parity spans let those modified edges fill without assuming a
// conventional closed silhouette after the exchange.
func (r *maskRing) prepareSliced(index int, frame []source.Polygon, zoom, exchangeY bool) {
	plane := r.planes[index]
	plane.Clear()
	r.batch.Begin(plane, r.white)
	point := func(p source.Point) ebiten.Vertex {
		x, y := int(p.X)*351/256, int(p.Y)*289/204
		if zoom {
			x = max(0, min(351, int(p.X)*680/256-130))
			y = max(0, min(289, int(p.Y)*550/204-135))
		}
		return render.Vertex(float64(x), float64(y), 0, 0, color.White)
	}
	for _, polygon := range frame {
		if len(polygon.Points) < 3 {
			continue
		}
		for i, p := range polygon.Points {
			a, b := point(p), point(polygon.Points[(i+1)%len(polygon.Points)])
			if exchangeY {
				a.DstY, b.DstY = b.DstY, a.DstY
			}
			if a.DstY == b.DstY {
				continue
			}
			leftA, leftB := a, b
			leftA.DstX, leftB.DstX = -1, -1
			r.batch.Quad([4]ebiten.Vertex{a, b, leftB, leftA})
		}
	}
	r.batch.Flush()
}

package demo

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/go-spaceballs/internal/source"
)

// contourVertex applies the original integer lookup before GPU parity filling.
func contourVertex(p source.Point, zoom bool) ebiten.Vertex {
	x, y := int(p.X)*351/256, int(p.Y)*289/204
	if zoom {
		x = max(0, min(351, int(p.X)*580/256-30))
		y = max(0, min(289, int(p.Y)*450/204-100))
	}
	return render.Vertex(float64(x), float64(y), 0, 0, color.White)
}

func filledContour(batch *render.Batch, points []source.Point, zoom bool) {
	batch.Fan(len(points), func(i int) ebiten.Vertex { return contourVertex(points[i], zoom) })
}

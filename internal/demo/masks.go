package demo

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/democonstructionkit/composite"
	"github.com/olivierh59500/democonstructionkit/render"
	"github.com/olivierh59500/go-spaceballs/internal/source"
	"image/color"
)

type maskRing struct {
	bank   *composite.ContourBank
	planes [6]*ebiten.Image
}

func newContourBank(width, height, slots, layers int, blend ebiten.Blend) *composite.ContourBank {
	bank, err := composite.NewContourBank(composite.ContourBankConfig{Width: width, Height: height, Slots: slots, Layers: layers, FillRule: ebiten.FillRuleEvenOdd, Blend: blend, BatchTriangles: 4096})
	if err != nil {
		panic(err)
	}
	return bank
}

func paintContours(bank *composite.ContourBank, slot, layer int, clear bool, draw func(*render.Batch)) {
	if err := bank.Paint(slot, layer, clear, draw); err != nil {
		panic(err)
	}
}

func newMaskRing() *maskRing {
	r := &maskRing{bank: newContourBank(Width, Height, 6, 1, ebiten.Blend{})}
	for i := range r.planes {
		r.planes[i] = r.bank.Image(i, 0)
	}
	return r
}
func (r *maskRing) clear() { r.bank.Clear() }
func (r *maskRing) close() { r.bank.Close() }

func (r *maskRing) prepare(index int, frame []source.Polygon, zoom bool) {
	paintContours(r.bank, index, 0, true, func(batch *render.Batch) {
		for _, p := range frame {
			filledContour(batch, p.Points, zoom)
		}
	})
}

func (r *maskRing) prepareOutline(index int, frame []source.Polygon, mirror bool) {
	point := func(p source.Point) ebiten.Vertex {
		x := int(p.X)
		if mirror {
			x = 256 - x
		}
		return render.Vertex(float64(x*351/256), float64(int(p.Y)*289/204), 0, 0, color.White)
	}
	paintContours(r.bank, index, 0, true, func(batch *render.Batch) {
		for _, p := range frame {
			batch.StrokeContour(len(p.Points), func(i int) ebiten.Vertex { return point(p.Points[i]) }, render.ContourStroke{SkipHorizontal: true})
		}
	})
}

func (r *maskRing) prepareSliced(index int, frame []source.Polygon, zoom, exchangeY bool) {
	point := func(p source.Point) ebiten.Vertex {
		x, y := int(p.X)*351/256, int(p.Y)*289/204
		if zoom {
			x = max(0, min(351, int(p.X)*680/256-130))
			y = max(0, min(289, int(p.Y)*550/204-135))
		}
		return render.Vertex(float64(x), float64(y), 0, 0, color.White)
	}
	paintContours(r.bank, index, 0, true, func(batch *render.Batch) {
		for _, p := range frame {
			batch.ParityContour(len(p.Points), func(i int) ebiten.Vertex { return point(p.Points[i]) }, render.ParityContour{RayX: -1, SwapY: exchangeY})
		}
	})
}

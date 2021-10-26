package imageErrorHeap

import (
	"container/heap"
	"image"
	"image/color"
	"image/draw"
	"math"
)

type ImageRectangleError struct {
	// bb of sub-image
	Rect image.Rectangle
	// AvgColor denotes the average color for a sub-image
	AvgColor color.RGBA64
	// AvgError denotes mean square error of color compared to average color
	AvgError float64
}

func newImageRectangleError(img *image.RGBA, rect image.Rectangle) ImageRectangleError {
	trimmedImg := moveTopLeft(img.SubImage(rect), rect)
	avgColor := averageColor(trimmedImg)
	mse := calculateColorMSE(trimmedImg, avgColor)

	return ImageRectangleError{
		Rect:     rect,
		AvgColor: avgColor,
		AvgError: mse,
	}
}

type ImageErrorHeap struct{
	img *image.RGBA
	imgRectError []ImageRectangleError
}

func (h ImageErrorHeap) Tmp() []ImageRectangleError{
	return h.imgRectError
}

func NewImageErrorHeap(img *image.RGBA) *ImageErrorHeap{
	imgErrorHeap := &ImageErrorHeap{
		img:          img,
		imgRectError: make([]ImageRectangleError, 0),
	}
	heap.Init(imgErrorHeap)
	heap.Push(imgErrorHeap, newImageRectangleError(img, img.Bounds()))
	return imgErrorHeap
}

func (h ImageErrorHeap) Len() int { return len(h.imgRectError) }

// Less - sorts desc AvgError, so that we pop biggest error first
func (h ImageErrorHeap) Less(i, j int) bool {
	return h.imgRectError[i].AvgError > h.imgRectError[j].AvgError
}

func (h ImageErrorHeap) Swap(i, j int) {
	h.imgRectError[i], h.imgRectError[j] = h.imgRectError[j], h.imgRectError[i]
}

func (h *ImageErrorHeap) Push(x interface{}) {
	(*h).imgRectError = append((*h).imgRectError, x.(ImageRectangleError))
}

func (h *ImageErrorHeap) Pop() interface{} {
	old := (*h).imgRectError
	n := len(old)
	x := old[n-1]
	(*h).imgRectError = old[0 : n-1]
	return x
}

// PushHelp - QoL method to push new record
func (h *ImageErrorHeap) PushHelp(rect image.Rectangle) ImageRectangleError{
	ire := newImageRectangleError(h.img, rect)
	heap.Push(h, ire)
	return ire
}

// PopHelp - QoL method to pop new record
func (h *ImageErrorHeap) PopHelp() ImageRectangleError {
	return heap.Pop(h).(ImageRectangleError)
}

func moveTopLeft(img image.Image, rect image.Rectangle) *image.RGBA {
	newImg := image.NewRGBA(rect.Sub(rect.Min))
	draw.Draw(newImg, newImg.Rect, img, rect.Min, draw.Src)
	return newImg
}

func averageColor(i image.Image) color.RGBA64 {
	var (
		area = uint64(i.Bounds().Dx() * i.Bounds().Dy())
		cumR, cumG, cumB, cumA uint64 = 0, 0, 0, 0
	)

	for y:=i.Bounds().Min.Y; y < i.Bounds().Max.Y; y++ {
		for x:=i.Bounds().Min.X; x < i.Bounds().Max.X; x++ {
			r, g, b, a := i.At(x, y).RGBA()
			cumR += uint64(r)
			cumG += uint64(g)
			cumB += uint64(b)
			cumA += uint64(a)
		}
	}

	return color.RGBA64{
		R: uint16(cumR / area),
		G: uint16(cumG / area),
		B: uint16(cumB / area),
		A: uint16(cumA / area),
	}
}

func calculateColorMSE(i image.Image, c color.RGBA64) float64 {
	// TODO: change to use CIELAB distance calculation instead, for now use rgb MSE
	se := func (c1, c2 color.Color) float64 {
		r1, g1, b1, _ := c1.RGBA()
		r2, g2, b2, _ := c2.RGBA()
		rErr, gErr, bErr := r1-r2, g1-g2, b1-b2
		return float64(rErr * rErr) + float64(gErr * gErr) + float64(bErr * bErr)
	}

	s := i.Bounds().Size()
	area := float64(s.X * s.Y)
	mse := float64(0)
	for y:=i.Bounds().Min.Y; y < i.Bounds().Max.Y; y++ {
		for x:=i.Bounds().Min.X; x < i.Bounds().Max.X; x++ {
			mse += se(i.At(x, y), c) / area
		}
	}
	return math.Sqrt(mse) * math.Log2(area) // weighted average by size
}
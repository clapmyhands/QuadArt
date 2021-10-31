package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"log"
	"os"
	imageErrorHeap "quadart/heap"
	"runtime"
)

const filename = "E137odiVkAQjj7T.jpg"
//const filename = "tmp.jpg"

func main() {
	reader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func(reader *os.File) {
		_ = reader.Close()
	}(reader)

	originalImg, err := jpeg.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	copyImg := image.NewRGBA(originalImg.Bounds())
	draw.Draw(copyImg, originalImg.Bounds(), originalImg, originalImg.Bounds().Min, draw.Src)

	dc := gg.NewContext(originalImg.Bounds().Dx(), originalImg.Bounds().Dy())

	imgErrHeap := imageErrorHeap.NewImageErrorHeap(copyImg)
	var crosses []image.Rectangle
	var images []image.Image
	for i := 0; i < 1000; i++ {
		ire := imgErrHeap.PopHelp()
		for _, rect := range split4Rectangle(ire.Rect) {
			tmp := imgErrHeap.PushHelp(rect)
			dc.DrawRectangle(
				float64(tmp.Rect.Min.X),
				float64(tmp.Rect.Min.Y),
				float64(tmp.Rect.Dx()),
				float64(tmp.Rect.Dy()),
			)
			dc.SetColor(tmp.AvgColor)
			dc.Fill()
		}
		crosses = append(crosses, ire.Rect)

		if i%5 == 0 {
			// copy and fill stroke before printing image
			// this is done such that stroke is always done on the last step as 1 action
			// avoids issue of multiple line stacking causing deeper line
			copyDC := gg.NewContextForImage(dc.Image())
			for _, cross := range crosses {
				drawCross(copyDC, cross, 0.5)
			}
			copyDC.SetRGBA(0, 0, 0, 0.4)
			copyDC.Stroke()

			_ = copyDC.SavePNG(fmt.Sprintf("./out/%d.png", i))
		}
		images = append(images, dc.Image())

		if i%10 == 0 {
			logMemoryUsage()
		}
	}

	//imagesToGIF(images, 0)
	//logMemoryUsage()
	//var (
	//	images  []image.Image
	//	rects   = []image.Rectangle{copyImg.Bounds()}
	//	crosses []image.Rectangle
	//)
	//for i:=0; i<8; i++ {
	//	var newRects = make([]image.Rectangle, 0, len(rects)*4)
	//	for _, rect := range rects {
	//		// append next loop's rectangles
	//		newRects = append(newRects, split4Rectangle(rect)...)
	//
	//		trimmedImg := moveTopLeft(copyImg.SubImage(rect), rect)
	//		avgColor := averageColor(trimmedImg)
	//
	//		dc.SetColor(avgColor)
	//		dc.DrawRectangle(float64(rect.Min.X), float64(rect.Min.Y), float64(rect.Max.X), float64(rect.Max.Y))
	//		dc.Fill()
	//	}
	//
	//	// draw previous split's crosses
	//	fmt.Printf("Crosses: %d\n", len(crosses))
	//	for _, cross := range crosses {
	//		drawCross(dc, cross, 0.8)
	//	}
	//	crosses = append(crosses, rects...)
	//
	//	rects = newRects
	//
	//	_ = dc.SavePNG(fmt.Sprintf("./out/%d.png", i))
	//	images = append(images, dc.Image())
	//}
}

func logMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf(`Alloc = %v MiB	TotalAlloc = %v MiB	Sys = %v MiB	NumGC = %v`,
		bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func imagesToGIF(imgs []image.Image, delay int) {
	var (
		images = make([]*image.Paletted, 0, 5)
		delays = make([]int, 0, 5)
	)
	for _, img := range imgs {
		images = append(images, imageToPaletted(img))
		delays = append(delays, 0)
	}
	f, err := os.OpenFile("pepe.gif", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_ = gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}

func split4Rectangle(rect image.Rectangle) []image.Rectangle {
	var (
		minX = rect.Min.X
		minY = rect.Min.Y
		midX = rect.Min.X + rect.Dx()/2
		midY = rect.Min.Y + rect.Dy()/2
		maxX = rect.Max.X
		maxY = rect.Max.Y
	)

	return []image.Rectangle{
		image.Rect(minX, minY, midX, midY),
		image.Rect(midX, minY, maxX, midY),
		image.Rect(minX, midY, midX, maxY),
		image.Rect(midX, midY, maxX, maxY),
	}
}

func trimmedSubImage(img *image.RGBA, rect image.Rectangle) *image.RGBA {
	return moveTopLeft(img.SubImage(rect), rect)
}

func moveTopLeft(img image.Image, rect image.Rectangle) *image.RGBA {
	newImg := image.NewRGBA(rect.Sub(rect.Min))
	draw.Draw(newImg, newImg.Rect, img, rect.Min, draw.Src)
	return newImg
}

func drawCross(dc *gg.Context, r image.Rectangle, lineWidth float64) {
	midX := r.Min.X + r.Dx()/2
	midY := r.Min.Y + r.Dy()/2
	dc.DrawLine(float64(midX), float64(r.Min.Y), float64(midX), float64(r.Max.Y))
	dc.DrawLine(float64(r.Min.X), float64(midY), float64(r.Max.X), float64(midY))
	dc.SetLineWidth(lineWidth)
}

func imageToPaletted(img image.Image) *image.Paletted {
	palettedImage := image.NewPaletted(img.Bounds(), palette.WebSafe)
	draw.Draw(palettedImage, palettedImage.Rect, img, img.Bounds().Min, draw.Src)
	return palettedImage
}

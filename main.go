package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"log"
	"os"
	imageErrorHeap "quadtree/heap"
)

//const filename = "E137odiVkAQjj7T.jpg"
const filename = "tmp.jpg"

// test draw image
// func main() {
// 	reader, err := os.Open(filename)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer reader.Close()
//
// 	img, err := jpeg.Decode(reader)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	bounds := img.Bounds()
//
// 	log.Println(bounds)
//
// 	dimg := image.NewRGBA(image.Rectangle{bounds.Min, bounds.Max})
// 	draw.Draw(
// 		dimg,
// 		img.Bounds(),
// 		img,
// 		img.Bounds().Min,
// 		draw.Src,
// 	)
//
// 	r := image.Rectangle{dimg.Bounds().Min, image.Point{
// 		X: dimg.Bounds().Max.X / 2,
// 		Y: dimg.Bounds().Max.Y / 2,
// 	}}
//
// 	draw.Draw(dimg, r, img, image.Point{
// 		X: dimg.Bounds().Max.X / 2,
// 		Y: dimg.Bounds().Max.Y / 2,
// 	}, draw.Src)
//
// 	dc := gg.NewContextForImage(dimg)
// 	dc.DrawImage(dimg, 0, 0)
// 	dc.SavePNG("out.png")
//

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
	//var crosses []image.Rectangle
	var images  []image.Image
	for i:=0; i<1000; i++ {
		ire := imgErrHeap.PopHelp()
		//for _, tmp := range imgErrHeap.Tmp() {
		//	fmt.Println(tmp.AvgError)
		//}
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

		////draw previous split's crosses
		//fmt.Printf("Crosses: %d\n", len(crosses))
		//crosses = append(crosses, ire.Rect)
		//for _, cross := range crosses {
		//	drawCross(dc, cross, 0.8)
		//}

		if i % 5 == 0 {
			fmt.Println("test")
			_ = dc.SavePNG(fmt.Sprintf("./out/%d.png", i))
		}
		images = append(images, dc.Image())
	}

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

func imagesToGIF(imgs []image.Image) {
	var (
		images = make([]*image.Paletted, 0, 5)
		delays = make([]int, 0, 5)
	)
	for _, img := range imgs {
		images = append(images, imageToPaletted(img))
		delays = append(delays, 100)
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

func split4(img *image.RGBA) []*image.RGBA {
	var (
		b = img.Bounds()
		minX = b.Min.X
		minY = b.Min.Y
		midX = b.Min.X + b.Dx()/2
		midY = b.Min.Y + b.Dy()/2
		maxX = b.Max.X
		maxY = b.Max.Y
	)

	topLeftImg := trimmedSubImage(img, image.Rect(minX, minY, midX, midY))
	topRightImg := trimmedSubImage(img, image.Rect(midX, minY, maxX, midY))
	botLeftImg := trimmedSubImage(img, image.Rect(minX, midY, midX, maxY))
	botRightImg := trimmedSubImage(img, image.Rect(midX, midY, maxX, maxY))
	return []*image.RGBA{topLeftImg, topRightImg, botLeftImg, botRightImg}
}

func trimmedSubImage(img *image.RGBA, rect image.Rectangle) *image.RGBA {
	return moveTopLeft(img.SubImage(rect), rect)
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

	if area <= 1 {
		fmt.Println(area)
		r, g, b, a := i.At(i.Bounds().Min.X, i.Bounds().Min.Y).RGBA()
		return color.RGBA64{
			R: uint16(r),
			G: uint16(g),
			B: uint16(b),
			A: uint16(a),
		}
	}

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

func drawCross(dc *gg.Context, r image.Rectangle, lineWidth float64) {
	midX := r.Min.X + r.Dx()/2
	midY := r.Min.Y + r.Dy()/2
	dc.DrawLine(float64(midX), float64(r.Min.Y), float64(midX), float64(r.Max.Y))
	dc.DrawLine(float64(r.Min.X), float64(midY), float64(r.Max.X), float64(midY))
	dc.SetLineWidth(lineWidth)
	dc.SetRGBA(0, 0, 0, 0.5)
	dc.Stroke()
}

func imageToPaletted(img image.Image) *image.Paletted {
	palettedImage := image.NewPaletted(img.Bounds(), palette.Plan9)
	draw.Draw(palettedImage, palettedImage.Rect, img,  img.Bounds().Min, draw.Over)
	return palettedImage
}
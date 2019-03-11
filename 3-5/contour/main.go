package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"gocv.io/x/gocv"
)

func matToImage(fileExt gocv.FileExt, mat gocv.Mat) image.Image {
	srcb, _ := gocv.IMEncode(fileExt, mat)
	src, _, _ := image.Decode(bytes.NewReader(srcb))
	return src
}

func getDst() draw.Image {
	s, _ := png.Decode(os.Stdin)
	dst, _ := s.(draw.Image)
	return dst
}

func whiteToTransparent(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var col color.RGBA
			r, _, _, _ := src.At(x, y).RGBA()
			if r == 65535 {
				col = color.RGBA{0, 0, 0, 255}
			}
			dst.Set(x, y, col)
		}
	}
	return dst
}

func main() {
	// グレイスケール化
	cvtSrc := gocv.IMRead("./src/go.jpeg", gocv.IMReadColor)
	gray := gocv.NewMatWithSize(460, 460, gocv.MatTypeCV64F)
	gocv.CvtColor(cvtSrc, &gray, gocv.ColorBGRToGray)

	// 二値化
	thresholdDst := gocv.NewMatWithSize(460, 460, gocv.MatTypeCV64F)
	gocv.Threshold(gray, &thresholdDst, 127, 255, gocv.ThresholdBinaryInv)

	// 輪郭抽出
	points := gocv.FindContours(thresholdDst, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	// 輪郭の内側を白にする
	gocv.DrawContours(&thresholdDst, points, -1, color.RGBA{255, 255, 255, 0}, -1)

	src := matToImage(gocv.PNGFileExt, cvtSrc)
	r := src.Bounds()

	maskSrc := matToImage(gocv.PNGFileExt, thresholdDst)
	mask := whiteToTransparent(maskSrc)
	png.Encode(os.Stdout, mask)

	dst := getDst()

	draw.DrawMask(dst, r, src, image.Pt(0, 0), mask, image.Pt(0, 0), draw.Over)
	png.Encode(os.Stdout, dst)
}

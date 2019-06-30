package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
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
	s, _, _ := image.Decode(os.Stdin)
	dst, _ := s.(draw.Image)
	return dst
}

func white2mask(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var col color.RGBA
			c := color.Gray16Model.Convert(src.At(x, y))
			gray, _ := c.(color.Gray16)
			if gray != color.Black {
				col = color.RGBA{0, 0, 0, 255}
			}
			dst.Set(x, y, col)
		}
	}
	return dst
}

func main() {
	// グレイスケール化
	cvtSrc := gocv.IMRead("./src/fuku_gopher.png", gocv.IMReadColor)
	gray := gocv.NewMatWithSize(460, 460, gocv.MatTypeCV64F)
	gocv.CvtColor(cvtSrc, &gray, gocv.ColorBGRToGray)

	// 二値化
	thresholdDst := gocv.NewMatWithSize(460, 460, gocv.MatTypeCV64F)
	gocv.Threshold(gray, &thresholdDst, 150, 150, gocv.ThresholdBinaryInv)

	// 輪郭抽出
	points := gocv.FindContours(thresholdDst, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	// 輪郭の内側を白にする
	gocv.DrawContours(&thresholdDst, points, -1, color.RGBA{255, 255, 255, 0}, -1)

	// image.Image に変換
	src := matToImage(gocv.PNGFileExt, cvtSrc)
	maskSrc := matToImage(gocv.PNGFileExt, thresholdDst)

	// draw.Draw の 引数準備
	dst := getDst()
	r := src.Bounds()
	mask := white2mask(maskSrc)

	// draw.Draw 実行
	draw.DrawMask(dst, r, src, image.Pt(0, 0), mask, image.Pt(0, 0), draw.Over)
	png.Encode(os.Stdout, dst)
}

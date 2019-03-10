package main

import (
	"fmt"
	"image/color"

	"gocv.io/x/gocv"
)

func main() {

	// classifier の初期化
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()
	if !classifier.Load("./src/haarcascade_frontalface_alt.xml") {
		fmt.Println("Error reading cascade file")
		return
	}

	// 画像をMat形式に
	img := gocv.IMRead("./src/nogi.jpg", gocv.IMReadColor)
	if img.Empty() {
		fmt.Println("Error reading image")
		return
	}

	// 顔検知
	rects := classifier.DetectMultiScale(img)
	fmt.Printf("found %d faces\n", len(rects))

	// 囲むための色
	blue := color.RGBA{0, 0, 255, 0}

	// 認識した顔の数だけ四角で囲む
	for _, r := range rects {
		// 引数は (Mat形式のデータ、範囲、色、線の太さ)
		gocv.Rectangle(&img, r, blue, 3)
	}

	// 結果を画像に書き出し
	gocv.IMWrite("result.png", img)
}

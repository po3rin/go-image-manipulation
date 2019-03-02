package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	img, _ := png.Decode(os.Stdin)
	bounds := img.Bounds() // (0,0)-(460,460)

	// 受け取った画像と同じ大きさのカラーモードGray16の画像を生成.
	// この時点では460*460の真っ黒の画像です。
	dst := image.NewGray16(bounds)

	// 1ピクセルずつ処理します。
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 元画像の(x,y)ピクセルのカラーをGray16に変換。
			c := color.Gray16Model.Convert(img.At(x, y))
			gray, _ := c.(color.Gray16)

			// 先ほど作ったdstに先ほど変換したカラーをセット。
			dst.Set(x, y, gray)
		}
	}

	png.Encode(os.Stdout, dst)
}

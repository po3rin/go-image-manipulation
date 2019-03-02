package main

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
)

func main() {
	img, _ := png.Decode(os.Stdin)

	// カラーモードがNRGBAか確認
	fmt.Println(img.ColorModel() == color.NRGBAModel) // output : true

	// 画像の境界を所得。単位はpx
	fmt.Println(img.Bounds()) // output : (0,0)-(460,460)

	// 指定したピクセルのRGBAを返します。
	fmt.Println(img.At(0, 0)) // output : {255 255 255 255}
}

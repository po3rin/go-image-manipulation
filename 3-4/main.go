package main

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
	"github.com/po3rin/go-image-manipulation/3-3/resize/lerp"
	"golang.org/x/image/font"
)

//Resize は与えられた画像を線形補間法を使用して画像を拡大・縮小する。
func Resize(img image.Image, xRatio, yRatio float64) image.Image {
	// 拡大後のサイズを計算
	width := int(float64(img.Bounds().Size().X) * xRatio)
	height := int(float64(img.Bounds().Size().Y) * yRatio)

	// 元となる新しい拡大画像を生成する
	newRect := image.Rect(0, 0, width, height)
	dst := image.NewRGBA64(newRect)
	bounds := dst.Bounds()

	// 1画素ずつ線形補正法を使ってカラーを計算してdstにセット
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x, y, LerpEffect(img, xRatio, yRatio, x, y))
		}
	}
	return dst
}

// getLerpParam 1軸に対してLerpで使うパラメータと近傍点の座標所得
func getLerpParam(dstPos int, ratio float64) (int, int, float64) {
	// 拡大前の座標の所得 (拡大後の座標 / リサイズ比率)
	v1float := float64(dstPos) / ratio

	// 拡大前の座標から最も近い2つの整数値を所得
	v1 := int(v1float)
	v2 := v1 + 1

	// (拡大前の座標の浮動小数点数 - 拡大前の座標の整数値)
	v3 := v1float - float64(v1)
	return v1, v2, v3
}

// initGetOneColorFunc RGBAいずれかを返す(x,y)座標依存関数を返す。
func initGetOneColorFunc(src image.Image, colorName string) lerp.PosDependFunc {
	return func(x, y int) float64 {
		var c uint32
		switch colorName {
		case "R":
			c, _, _, _ = src.At(x, y).RGBA()
		case "G":
			_, c, _, _ = src.At(x, y).RGBA()
		case "B":
			_, _, c, _ = src.At(x, y).RGBA()
		case "A":
			_, _, _, c = src.At(x, y).RGBA()
		}
		return float64(c)
	}
}

// LerpEffect (x,y)に対してLerpを行った結果のカラーを返す
func LerpEffect(src image.Image, xRatio, yRatio float64, x, y int) color.RGBA64 {
	//	元画像の近傍４画素の座標と、alpha、betaを所得
	x1, x2, alpha := getLerpParam(x, xRatio)
	y1, y2, beta := getLerpParam(y, yRatio)

	// 元画像の近傍４画素の座標を格納
	ps := lerp.Points{
		{X: x1, Y: y1},
		{X: x2, Y: y1},
		{X: x1, Y: y2},
		{X: x2, Y: y2},
	}

	// RGBAそれぞれの値に対してLerpを適用
	r := lerp.Lerp(initGetOneColorFunc(src, "R"), alpha, beta, ps)
	g := lerp.Lerp(initGetOneColorFunc(src, "G"), alpha, beta, ps)
	b := lerp.Lerp(initGetOneColorFunc(src, "B"), alpha, beta, ps)
	a := lerp.Lerp(initGetOneColorFunc(src, "A"), alpha, beta, ps)

	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

// NewRect 指定した大きさの画像を指定色で塗りつぶした画像を生成
func NewRect(r image.Rectangle, c color.Color) draw.Image {
	dst := image.NewRGBA(r)
	rect := dst.Rect
	for h := rect.Min.Y; h < rect.Max.Y; h++ {
		for v := rect.Min.X; v < rect.Max.X; v++ {
			dst.Set(v, h, c)
		}
	}
	return dst
}

// GetSrc 合成する画像を読み込んで横幅300の画像にリサイズした結果を消す
func GetSrc() image.Image {
	src, _, _ := image.Decode(os.Stdin)
	rate := float64(300) / float64(src.Bounds().Max.X)
	src = Resize(src, rate, rate)
	return src
}

// GetCover カバー画像を読み込み、OGP画像(横幅1200px)に合うようにリサイズした結果を消す
func GetCover() image.Image {
	f, _ := os.Open("./cover.jpg")
	defer f.Close()
	cover, _, _ := image.Decode(f)
	rate := float64(1200) / float64(cover.Bounds().Max.X)
	cover = Resize(cover, rate, rate)
	return cover
}

// DrawText テキストの合成
func DrawText(img draw.Image, text string) image.Image {
	// フォントファイルを読み込んでfreetype.Fontにパース
	file, _ := os.Open("./mplus-1c-regular.ttf")
	fontBytes, _ := ioutil.ReadAll(file)
	f, _ := freetype.ParseFont(fontBytes)

	// freetypeの機能で画像に文字を合成
	if f != nil {
		c := freetype.NewContext()
		c.SetFont(f)
		c.SetFontSize(38)
		c.SetClip(img.Bounds())
		c.SetDst(img)
		c.SetSrc(NewRect(img.Bounds(), color.RGBA{255, 255, 255, 255}))
		c.SetHinting(font.HintingNone)
		pt := freetype.Pt(300, 500)
		_, _ = c.DrawString(text, pt)
	}
	return img
}

func main() {
	r := image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{1200, 630}}

	dst := NewRect(r, color.RGBA{0, 0, 0, 250})
	src := GetSrc()
	cover := GetCover()
	mask := NewRect(r, color.RGBA{0, 0, 0, 140})

	// coverをsrcとしてdstに合成
	draw.DrawMask(
		dst, r,
		cover, image.Pt(0, 0),
		mask, image.Pt(0, 0),
		draw.Over,
	)
	// srcを真ん中に合成。
	draw.Draw(
		dst, r,
		src, image.Pt(
			-r.Bounds().Max.X/2+src.Bounds().Max.X/2,
			-r.Bounds().Max.Y/2+src.Bounds().Max.Y/2,
		),
		draw.Over,
	)

	text := "Goではじめる画像処理、画像解析"
	d := DrawText(dst, text)

	// 書き出し
	png.Encode(os.Stdout, d)
}

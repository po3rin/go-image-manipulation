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
	"github.com/po3rin/resize"
	"golang.org/x/image/font"
)

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
	src = resize.Resize(src, rate, rate)
	return src
}

// GetCover カバー画像を読み込み、OGP画像(横幅1200px)に合うようにリサイズした結果を消す
func GetCover() image.Image {
	f, _ := os.Open("./src/cover.jpg")
	defer f.Close()
	cover, _, _ := image.Decode(f)
	rate := float64(1200) / float64(cover.Bounds().Max.X)
	cover = resize.Resize(cover, rate, rate)
	return cover
}

// DrawText テキストの合成
func DrawText(img draw.Image, text string) image.Image {
	// フォントファイルを読み込んでfreetype.Fontにパース
	file, _ := os.Open("./src/mplus-1c-regular.ttf")
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
	mask := NewRect(r, color.RGBA{0, 0, 0, 60})

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

# go-image-manipulation

golang.tokyoが技術書典6で配布する技術書の@po3rinが執筆する第n章(未定)のコードです。

## コンテンツ

|ディレクトリ|内容|
|:-----------|:------------|
|3-2|imageパッケージの基本を抑えよう|
|3-3|線形補正法を実装して画像リサイズの仕組みを学ぼう|
|3-4|画像やテキストを合成してOGP画像を生成してみよう|
|3-5|OpenCVを使って画像解析をやってみよう|

## OpenCV環境の動かし方

3-5節ではOpenCVを使う為、OpenCVをGoで動かす為のDockerfileを用意しています。コマンドを渡せば実行できるMakefileを用意しています。

```console
$ make run CMD="go run 3-5/hellogocv/main.go"
```

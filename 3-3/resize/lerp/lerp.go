package lerp

// Point lerpで使う座標
type Point struct {
	X int
	Y int
}

// Points lerpで使う近傍4座標
type Points [4]Point

// Lerper lerpの為のパラメータ
type Lerper struct {
	A  float64
	B  float64
	Ps Points
}

// NewLerper Lerperの初期化
func NewLerper(a, b float64, ps Points) Lerper {
	return Lerper{
		A:  a,
		B:  b,
		Ps: ps,
	}
}

// PosDependFunc 座標依存関数
type PosDependFunc func(x, y int) float64

// Lerp calicurate relp
func (l Lerper) Lerp(f PosDependFunc) float64 {
	n := (1.0-l.B)*(1.0-l.A)*f(l.Ps[0].X, l.Ps[0].Y) +
		l.A*(1.0-l.B)*f(l.Ps[1].X, l.Ps[0].Y) +
		l.B*(1-l.A)*f(l.Ps[0].X, l.Ps[1].Y) +
		l.A*l.B*f(l.Ps[1].X, l.Ps[1].Y)
	return n
}

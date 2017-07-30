package math

// 点
type Point struct {
	X, Y float32
}

// 大小
type Size struct {
	Width, Height float32
}

// 矩形
type Rect struct {
	Min, Max Point
}

func (rect *Rect) Dx() float32{
	return rect.Max.X - rect.Min.X
}

func (rect *Rect) Dy() float32 {
	return rect.Max.Y - rect.Min.Y
}

func (rect *Rect) Center() Point {
	return Point{(rect.Max.X - rect.Min.X) / 2, (rect.Max.Y - rect.Min.X) / 2}
}








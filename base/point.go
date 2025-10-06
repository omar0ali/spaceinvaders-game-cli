package base

type Point struct {
	X, Y int
}

type PointFloat struct {
	X, Y float64
}

type PointInterface interface {
	GetX() float64
	GetY() float64
	SetX(float64)
	SetY(float64)
	AppendX(float64)
	AppendY(float64)
}

func (p Point) GetX() float64 {
	return float64(p.X)
}

func (p Point) GetY() float64 {
	return float64(p.Y)
}

func (p *Point) SetX(x float64) {
	p.X = int(x)
}

func (p *Point) SetY(y float64) {
	p.Y = int(y)
}

func (p *Point) AppendX(x float64) {
	p.X += int(x)
}

func (p *Point) AppendY(y float64) {
	p.Y += int(y)
}

func (p PointFloat) GetX() float64 {
	return p.X
}

func (p PointFloat) GetY() float64 {
	return p.Y
}

func (p *PointFloat) AppendX(x float64) {
	p.X += x
}

func (p *PointFloat) AppendY(y float64) {
	p.Y += y
}

package rayTracer

import (
	"fmt"
)

const PIOVER180 float64 = 0.017453292519943295769236907684886

type Point struct {
	X, Y, Z float64
}

type Vector struct {
	X, Y, Z float64
}

func (v *Vector) AddEq(v2 Vector) Vector {
	v.X += v2.X
	v.Y += v2.Y
	v.Z += v2.Z
	return *v
}

func (p *Point) Add(v Vector) Vector {
	p2 := Vector{p.X + v.X, p.Y + v.Y, p.Z + v.Z}
	return p2
}

func (p *Point) Sub(v Vector) Vector {
	p2 := Vector{p.X - v.X, p.Y - v.Y, p.Z - v.Z}
	return p2
}

func (p *Point) AddP(v Point) Vector {
	p2 := Vector{p.X + v.X, p.Y + v.Y, p.Z + v.Z}
	return p2
}

func (p *Point) SubP(v Point) Vector {
	p2 := Vector{p.X - v.X, p.Y - v.Y, p.Z - v.Z}
	return p2
}

func (p *Point) String() string {
	s := ""
	s += fmt.Sprintf("(%f, %f, %f)", p.X, p.Y, p.Z)
	return s
}

func (v *Vector) AddV(v2 Vector) Vector {
	v1 := Vector{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
	return v1
}

func (v *Vector) AddP(p1, p2 Point) Vector {
	v1 := Vector{p1.X + p2.X, p1.Y + p2.Y, p1.Z + p2.Z}
	return v1
}
func (v *Vector) SubP(p1, p2 Point) Vector {
	v1 := Vector{p1.X - p2.X, p1.Y - p2.Y, p1.Z - p2.Z}
	return v1
}

func (v *Vector) MulC(c float64) Vector {
	v1 := Vector{v.X * c, v.Y * c, v.Z * c}
	return v1
}

func (v1 *Vector) Sub(v2 Vector) Vector {
	v := Vector{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z}
	return v
}

func (v1 *Vector) MulV(v2 Vector) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

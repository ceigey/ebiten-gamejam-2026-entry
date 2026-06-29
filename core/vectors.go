package core

import "math"

type Vec2 struct {
	X float64
	Y float64
}

func (vec *Vec2) Normalize() Vec2 {
	length := vec.Magnitude()
	if length == 0 {
		return Vec2{0, 0}
	}
	return Vec2{vec.X / length, vec.Y / length}
}

func (vec *Vec2) Magnitude() float64 {
	return math.Hypot(vec.X, vec.Y)
}

func (vec *Vec2) Coords() (float64, float64) {
	return vec.X, vec.Y
}

func (vec *Vec2) Reverse() Vec2 {
	return Vec2{
		X: -vec.X,
		Y: -vec.Y,
	}
}

func (vec *Vec2) Multiply(n float64) Vec2 {
	return Vec2{
		X: vec.X * n,
		Y: vec.Y * n,
	}
}

package core

import "math"

type Vec2 struct {
	X float64
	Y float64
}

func (vec Vec2) Normalize() Vec2 {
	length := math.Hypot(vec.X, vec.Y)
	if length == 0 {
		return Vec2{0, 0}
	}
	return Vec2{vec.X / length, vec.Y / length}
}

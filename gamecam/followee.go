package gamecam

import "github.com/g3n/engine/math32"

type Followee interface {
	Height() float64
	Position() math32.Vector3
	Theta() float64
}

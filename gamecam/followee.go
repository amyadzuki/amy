package gamecam

import "github.com/g3n/engine/math32"

type Followee interface {
	Facing() (float64, float64)
	HeightToEye() float64
	Position() math32.Vector3
}

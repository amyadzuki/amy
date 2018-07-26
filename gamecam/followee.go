package gamecam

import "github.com/g3n/engine/math32"

type Followee interface {
	FacingRadiansCcwOfEast() float64
	HeightToEye() float64
	Position() math32.Vector3
}

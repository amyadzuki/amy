package gamecam

import "github.com/g3n/engine/math32"

type Followee interface {
	Base() float64
	FacingNormalized() (float64, float64)
	FrontOfEye() float64
	HeightToCap() float64
	HeightToEye() float64
	Position() math32.Vector3
}

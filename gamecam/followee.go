package gamecam

import "github.com/g3n/engine/math32"

type Followee interface {
	FacingNormalized() (float64, float64)
	YAtEye() float64
	ZAtBot() float64
	ZAtCap() float64
	ZAtEye() float64
	Position() math32.Vector3
}

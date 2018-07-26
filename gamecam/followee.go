package gamecam

import "github.com/g3n/engine/math32"

type Followee interface {
	Height() float32
	Position() math32.Vector3
}

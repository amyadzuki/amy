package gamecam

import (
	"github.com/amyadzuki/amystuff/bitfield"
)

type CamMode struct {
	bitfield.Uint8
}

var ImplScreen func(CamMode) bool = nil

func (cm *CamMode) Init(mask uint8) {
	cm.Uint8 = bitfield.Uint8(mask)
}

func (cm CamMode) Screen() bool {
	if ImplScreen != nil {
		return ImplScreen(cm)
	} else {
		return cm.Any(ScreenReasons) && !cm.Any(WorldOverrides)
	}
}

func (cm CamMode) World() bool {
	return !cm.Screen()
}

const (
	DefaultToScreen uint8 = iota << 1
	ScreenButtonHeld
	ScreenToggleOn
	MiddleMouseHeld
	FirstUnusedBit
)

var (
	ScreenReasons  = DefaultToScreen | ScreenButtonHeld | ScreenToggleOn
	WorldOverrides = MiddleMouseHeld
)

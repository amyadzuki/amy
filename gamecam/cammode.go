package gamecam

import (
	"github.com/amyadzuki/amygolib/bitfield"
)

type CamMode struct {
	bitfield.Uint8
}

var ImplScreen func(CamMode) bool = nil
var ForceScreen, ForceWorld bool

func (cm *CamMode) Init(mask uint8) {
	cm.Uint8 = bitfield.Uint8(mask)
}

func (cm CamMode) Screen() bool {
	switch {
	case ImplScreen != nil:
		return ImplScreen(cm)
	case ForceScreen:
		return true
	case ForceWorld:
		return false
	default:
		return cm.Any(ScreenReasons) && !cm.Any(WorldOverrides)
	}
}

func (cm CamMode) World() bool {
	return !cm.Screen()
}

const (
	DefaultToScreen uint8 = 1 << iota
	ScreenButtonHeld
	ScreenToggleOn
	MiddleMouseHeld
	FirstUnusedBit
)

var (
	ScreenReasons  = DefaultToScreen | ScreenButtonHeld | ScreenToggleOn
	WorldOverrides = MiddleMouseHeld
)

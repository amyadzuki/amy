package styles

import (
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
)

var CloseButton gui.ButtonStyles
var ClosingButton gui.ButtonStyles

func init() {
	s := gui.StyleDefault()

	CloseButton = s.Button
	CloseButton.Over.BgColor = math32.Color4{0.75, 0, 0, 1}

	ClosingButton = s.Button
	ClosingButton.Normal.BgColor = CloseButton.Over.BgColor
	ClosingButton.Over.BgColor = math32.Color4{1, 0, 0, 1}
}

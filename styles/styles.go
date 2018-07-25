package styles

import (
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
)

var AmyDark gui.Style
var AmyDarkCloseButton gui.ButtonStyles
var AmyDarkClosingButton gui.ButtonStyles
var AmyDarkHelpButton gui.ButtonStyles
var AmyDarkHelpingButton gui.ButtonStyles

func init() {
	AmyDark = *gui.StyleDefault()
	AmyDark.Button.Normal.BgColor = math32.Color4{0, 0, 0, 0.125}
	AmyDark.Button.Over.BgColor = math32.Color4{0.25, 0.125, 0.375, 0.25}
	AmyDarkCloseButton = AmyDark.Button
	AmyDarkCloseButton.Over.BgColor = math32.Color4{0.75, 0, 0, 1}
	AmyDarkClosingButton = AmyDarkCloseButton
	AmyDarkClosingButton.Normal.BgColor = AmyDarkCloseButton.Over.BgColor
	AmyDarkClosingButton.Over.BgColor = math32.Color4{1, 0, 0, 1}
	AmyDarkHelpButton = AmyDark.Button
	AmyDarkHelpingButton = AmyDarkHelpButton
}

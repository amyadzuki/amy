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
var AmyDarkWindowContent math32.Color4
var InvisibleButton gui.ButtonStyles

func init() {
	AmyDark = *gui.StyleDefault()
	AmyDark.Button.Normal.BgColor = math32.Color4{0, 0, 0, 0.25}
	AmyDark.Button.Over.BgColor = math32.Color4{0.25, 0.125, 0.375, 0.5}
	AmyDark.Button.Focus.BgColor = math32.Color4{0.25, 0.125, 0.375, 0.5}
	AmyDark.Button.Pressed.BgColor = math32.Color4{0.25, 0.125, 0.375, 1}
	AmyDark.Button.Disabled.BgColor = math32.Color4{0, 0, 0, 0}

	AmyDark.Window.Normal.TitleStyle.BgColor = math32.Color4{0.1294117647, 0.58823529411, 0.95294117647, 0.25}
	AmyDark.Window.Over.TitleStyle.BgColor = AmyDark.Window.Normal.TitleStyle.BgColor
	AmyDark.Window.Over.TitleStyle.BgColor.A = 1
	AmyDark.Window.Focus.TitleStyle.BgColor = AmyDark.Window.Normal.TitleStyle.BgColor
	AmyDark.Window.Disabled.TitleStyle.BgColor = AmyDark.Window.Normal.TitleStyle.BgColor

	// 0x21/255 0x96/255 0xF3/255
	// 0.1294117647, 0.58823529411, 0.95294117647
	AmyDarkWindowContent = math32.Color4{0, 0, 0, 0.25}

	AmyDarkCloseButton = AmyDark.Button
	AmyDarkCloseButton.Over.BgColor = math32.Color4{0.75, 0, 0, 0.5}
	AmyDarkCloseButton.Focus.BgColor = math32.Color4{0.75, 0, 0, 0.5}
	AmyDarkCloseButton.Pressed.BgColor = math32.Color4{1, 0, 0, 1}
	AmyDarkClosingButton = AmyDarkCloseButton
	AmyDarkClosingButton.Normal.BgColor = math32.Color4{1, 0, 0, 0.5}
	AmyDarkClosingButton.Over.BgColor = AmyDarkCloseButton.Pressed.BgColor
	AmyDarkClosingButton.Focus.BgColor = AmyDarkCloseButton.Pressed.BgColor
	AmyDarkClosingButton.Pressed.BgColor = AmyDarkCloseButton.Pressed.BgColor

	AmyDarkHelpButton = AmyDark.Button
	AmyDarkHelpingButton = AmyDarkHelpButton

	InvisibleButton.Normal.BorderColor = math32.Color4{0, 0, 0, 0}
	InvisibleButton.Over.BorderColor = math32.Color4{0, 0, 0, 0}
	InvisibleButton.Focus.BorderColor = math32.Color4{0, 0, 0, 0}
	InvisibleButton.Pressed.BorderColor = math32.Color4{0, 0, 0, 0}
	InvisibleButton.Disabled.BorderColor = math32.Color4{0, 0, 0, 0}
	InvisibleButton.Normal.BgColor = math32.Color4{0, 0, 0, 0}
	InvisibleButton.Over.BgColor = math32.Color4{0, 0, 0, 0}
	InvisibleButton.Focus.BgColor = math32.Color4{0, 0, 0, 0}
	InvisibleButton.Pressed.BgColor = math32.Color4{0, 0, 0, 0}
	InvisibleButton.Disabled.BgColor = math32.Color4{0, 0, 0, 0}
	InvisibleButton.Normal.FgColor = math32.Color4{0, 0, 0, 0}
	InvisibleButton.Over.FgColor = math32.Color4{0, 0, 0, 0}
	InvisibleButton.Focus.FgColor = math32.Color4{0, 0, 0, 0}
	InvisibleButton.Pressed.FgColor = math32.Color4{0, 0, 0, 0}
	InvisibleButton.Disabled.FgColor = math32.Color4{0, 0, 0, 0}
}

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
	AmyDark.Button.Normal.BgColor = math32.Color4{0, 0, 0, 0.25}
	AmyDark.Button.Over.BgColor = math32.Color4{0.25, 0.125, 0.375, 0.5}
	AmyDark.Button.Focus.BgColor = math32.Color4{0.25, 0.125, 0.375, 0.5}
	AmyDark.Button.Pressed.BgColor = math32.Color4{0.25, 0.125, 0.375, 1}
	AmyDark.Button.Disabled.BgColor = math32.Color4{0, 0, 0, 0}
	AmyDark.Window.Normal.TitleBgColor = math32.Color4{0.1294117647, 0.58823529411, 0.95294117647, 0.25}
	AmyDark.Window.Over.TitleBgColor = AmyDark.Window.Normal.TitleBgColor
	AmyDark.Window.Over.TitleBgColor.A = 1
	AmyDark.Window.Focus.TitleBgColor = AmyDark.Window.Normal.TitleBgColor
	AmyDark.Window.Disabled.TitleBgColor = AmyDark.Window.Normal.TitleBgColor
	// 0x21/255 0x96/255 0xF3/255
	// 0.1294117647, 0.58823529411, 0.95294117647

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
}

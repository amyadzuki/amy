package account

import "github.com/g3n/engine/gui"

type Chara interface {
	Account
	StyleWidgetChanger(*gui.Button) // style the button to change charas
}

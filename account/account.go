package account

import "fmt"

type User interface {
	HintLeft() string // Race and class, for example
	HintRight() string // Level, for example
	LogIn() error // Log in with the server
	Peek() error // Show the account in the GUI preview pane
	Register() error // Register with the server

	fmt.Stringer // The name of the account for display
}

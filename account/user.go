package account

type User interface {
	Account
	Charas() ([]Chara, int, error) // the int is the index of the last one played
}

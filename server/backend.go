package server

type Backend interface {
	Body(...interface{}) []byte
	Error(int, string, ...interface{})
	Http(int, interface{}) error
	Https(string, string, string, interface{}) error
	Path(...interface{}) string
	Write(string, ...interface{})
}

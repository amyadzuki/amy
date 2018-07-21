package server

import "github.com/amyadzuki/amystuff/onfail"

type Server interface {
	Api(Backend, uint32, ...interface{})
	Http(Backend, ...interface{})
	Https(Backend, ...interface{})
	Serve(Backend, string, string, string, string, onfail.OnFail)
}

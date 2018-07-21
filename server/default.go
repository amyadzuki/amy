package server

import "github.com/amyadzuki/amystuff/onfail"

func Api(backend Backend, version uint32, args ...interface{}) {
	Default.Api(backend, version, args...)
}

func Http(backend Backend, args ...interface{}) {
	Default.Http(backend, args...)
}

func Https(backend Backend, args ...interface{}) {
	Default.Https(backend, args...)
}

func Serve(backend Backend, httpAddr, httpsAddr, certPath, keyPath string, onFail onfail.OnFail) {
	Default.Serve(backend, httpAddr, httpsAddr, certPath, keyPath, onFail)
}

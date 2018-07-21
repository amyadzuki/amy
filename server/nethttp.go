package server

import (
	"fmt"
	"ioutil"

	"net/http"
)

type NetHttpBackend struct {
	ImplBody  func(...interface{}) ([]byte, error)
	ImplError func(int, string, ...interface{})
	ImplHttp  func(string, interface{}) error
	ImplHttps func(string, string, string, interface{}) error
	ImplPath  func(...interface{}) string
	ImplWrite func(string, ...interface{})
}

func (impl *NetHttpBackend) Body(args ...interface{}) (body []byte, err error) {
	if impl.ImplBody != nil {
		body, err = impl.ImplBody(args...)
	} else {
		body, err = ioutil.ReadAll(args[1].(*http.Request).Body)
	}
	return
}

func (impl *NetHttpBackend) Error(status int, reason string, args ...interface{}) {
	if impl.ImplError != nil {
		impl.ImplError(status, reason, args...)
	} else {
		http.Error(args[0].(http.ResponseWriter), reason, status)
	}
}

func (impl *NetHttpBackend) Http(addr string, handler interface{}) error {
	if impl.ImplHttp != nil {
		return impl.ImplHttp(addr, handler)
	} else {
		return http.ListenAndServe(addr, handler.(http.Handler))
	}
}

func (impl *NetHttpBackend) Https(addr, certPath, keyPath string, handler interface{}) error {
	if impl.ImplHttps != nil {
		return impl.ImplHttps(addr, certPath, keyPath, handler)
	} else {
		return http.ListenAndServeTLS(addr, certPath, keyPath, handler.(http.Handler))
	}
}

func (impl *NetHttpBackend) Path(args ...interface{}) string {
	if impl.ImplPath != nil {
		return impl.ImplPath(args...)
	} else {
		return args[1].(*http.Request).URL.Path
	}
}

func (impl *NetHttpBackend) Write(body string, args ...interface{}) {
	if impl.ImplWrite != nil {
		impl.ImplWrite(body, args...)
	} else {
		fmt.Fprintln(args[0].(http.ResponseWriter), body)
	}
}

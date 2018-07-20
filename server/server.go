package server

import (
	"../str"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
)

var FastHttp FastHttpBackend
var NetHttp NetHttpBackend
var Default Basic

type Backend interface {
	Body(...interface{}) []byte
	Error(int, string, ...interface{})
	Http(int, interface{}) error
	Https(string, string, string, interface{}) error
	Path(...interface{}) string
	Write(string, ...interface{})
}

// github.com/valyala/fasthttp backend:

type FastHttpBackend struct {
	ImplBody  func(...interface{}) []byte
	ImplError func(int, string, ...interface{})
	ImplHttp  func(int, interface{}) error
	ImplHttps func(string, string, string, interface{}) error
	ImplPath  func(...interface{}) string
	ImplWrite func(string, ...interface{})
}

func (impl *FastHttpBackend) Body(args ...interface{}) []byte {
	if impl.ImplBody != nil {
		return impl.ImplBody(args...)
	} else {
		return args[0].(*fasthttp.RequestCtx).PostBody()
	}
}

func (impl *FastHttpBackend) Error(status int, reason string, args ...interface{}) {
	if impl.ImplError != nil {
		impl.ImplError(reason, status, args...)
	} else {
		args[0].(*fasthttp.RequestCtx).Error(reason, status)
	}
}

func (impl *FastHttpBackend) Http(addr string, handler interface{}) error {
	if impl.ImplHttp != nil {
		return impl.ImplHttp(addr, handler)
	} else {
		return fasthttp.ListenAndServe(addr, handler.(fasthtttp.RequestHandler))
	}
}

func (impl *FastHttpBackend) Https(addr, certPath, keyPath string, handler interface{}) error {
	if impl.ImplHttps != nil {
		return impl.ImplHttps(addr, certPath, keyPath, handler)
	} else {
		return fasthttp.ListenAndServeTLS(addr, certPath, keyPath, handler.(fasthtttp.RequestHandler))
	}
}

func (impl *FastHttpBackend) Path(args ...interface{}) string {
	if impl.ImplPath != nil {
		return impl.ImplPath(args...)
	} else {
		return args[0].(*fasthttp.RequestCtx).Path()
	}
}

func (impl *FastHttpBackend) Write(body string, args ...interface{}) {
	if impl.ImplWrite != nil {
		impl.ImplWrite(body, args...)
	} else {
		fmt.Fprintln(args[0].(*fasthttp.RequestCtx), body)
	}
}

// net/http backend:

type NetHttpBackend struct {
	ImplBody  func(...interface{}) []byte
	ImplError func(int, string, ...interface{})
	ImplPath  func(...interface{}) string
	ImplWrite func(string, ...interface{})
}

func (impl *NetHttpBackend) Body(args ...interface{}) []byte {
	if impl.ImplBody != nil {
		return impl.ImplBody(args...)
	} else {
		return args[1].(*http.Request).Body
	}
}

func (impl *NetHttpBackend) Error(status int, reason string, args ...interface{}) {
	if impl.ImplError != nil {
		impl.ImplError(reason, status, args...)
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

type Server interface {
	Api(Backend, uint32, ...interface{})
	Http(Backend, ...interface{})
	Https(Backend, ...interface{})
	Serve(Backend, string, string, string, onfail.Func)
}

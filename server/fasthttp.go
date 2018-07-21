package server

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type FastHttpBackend struct {
	ImplBody  func(...interface{}) []byte
	ImplError func(int, string, ...interface{})
	ImplHttp  func(string, interface{}) error
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
		impl.ImplError(status, reason, args...)
	} else {
		args[0].(*fasthttp.RequestCtx).Error(reason, status)
	}
}

func (impl *FastHttpBackend) Http(addr string, handler interface{}) error {
	if impl.ImplHttp != nil {
		return impl.ImplHttp(addr, handler)
	} else {
		return fasthttp.ListenAndServe(addr, handler.(fasthttp.RequestHandler))
	}
}

func (impl *FastHttpBackend) Https(addr, certPath, keyPath string, handler interface{}) error {
	if impl.ImplHttps != nil {
		return impl.ImplHttps(addr, certPath, keyPath, handler)
	} else {
		return fasthttp.ListenAndServeTLS(addr, certPath, keyPath, handler.(fasthttp.RequestHandler))
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

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

// Any Server:

type Server interface {
	Api(Backend, uint32, ...interface{})
	Http(Backend, ...interface{})
	Https(Backend, ...interface{})
	Serve(Backend, string, string, string, onfail.Func)
}

// A Basic Server:

type Basic struct {
	ImplApi   func(Backend, interface{}, uint32)
	ImplHttp  func(Backend, interface{})
	ImplHttps func(Backend, interface{})
	ImplServe func(Backend, string, string, string, onfail.Func)
}

func (server *Basic) Api(backend Backend, version uint32, args ...interface{}) {
	if server.ImplApi != nil {
		server.ImplApi(backend, version, args...)
	} else {
		basicApi(server, backend, version, args...)
	}
}

func (server *Basic) Http(backend Backend, args ...interface{}) {
	if server.ImplHttp != nil {
		server.ImplHttp(backend, args...)
	} else {
		basicHttp(server, backend, args...)
	}
}

func (server *Basic) Https(backend Backend, args ...interface{}) {
	if server.ImplHttps != nil {
		server.ImplHttps(backend, args...)
	} else {
		basicHttps(server, backend, args...)
	}
}

func (server *Basic) Serve(backend Backend, addr, certPath, keyPath string, onFail onfail.Func) {
	if server.ImplServe != nil {
		server.ImplServe(backend, addr, certPath, keyPath, onFail)
	} else {
		basicServe(server, backend, addr, certPath, keyPath, onFail)
	}
}

// Default method implementations for the Basic Server

func basicApi(server Server, backend Backend, version uint32, args ...interface{}) {
	backend.Error(501, "API not implemented")
}

func basicHttp(server Server, backend Backend, args ...interface{}) {
	basicHttpx(server, backend, false, args...)
}

func basicHttps(server Server, backend Backend, args ...interface{}) {
	basicHttpx(server, backend, true, args...)
}

func basicHttpx(server Server, backend Backend, secure bool, args ...interface{}) {
	path := backend.Path()
	endedInHtml, lessHtml := str.CaseHasSuffix(path, ".html")
	if endedInHtml {
		path = lessHtml
	} else {
		endedInHtm, lessHtm := str.CaseHasSuffix(path, ".htm")
		if endedInHtm {
			path = lessHtm
		}
	}
	endedInIndex, lessIndex := str.CaseHasSuffix(path, "index")
	if endedInIndex {
		path = lessIndex
	}
	endedInSlash, lessSlash := str.CaseHasSuffix(path, "/")
	if endedInSlash {
		path = lessSlash
	}

	if !secure {
		beganWithApiV, lessApiV := str.CaseHasPrefix(path, "/api/v")
		if beganWithApiV {
			fieldsSlash := strings.Split(lessApiV, "/")
			if len(fieldsSlash) < 2 {
				backend.Error(404, "")
				return
			}
			fieldsVersion := strings.Split(fieldsSlash[0], ".")
			if len(fieldsVersion) < 1 || len(fieldsVersion) > 4 {
				backend.Error(404, "")
				return
			}
			var version uint32
			for idx := 0; idx < 4; idx++ {
				version <<= 8
				if idx < len(fieldsVersion) {
					part, err := strconv.ParseInt(fieldsVersion[idx], 10, 8)
					if err != nil {
						backend.Error(404, "Invalid API version number")
						return
					}
					version |= uint32(part)
				}
			}
			server.Api(backend, version, args...)
			return
		}
	}

}

func basicServe(server Server, backend Backend, addr, certPath, keyPath string, onFail onfail.Func) {
	go func() {
		onFail("HTTP", backend.Http(addr, server.Http))
	}()
	go func() {
		onFail("HTTPS", backend.Https(addr, certPath, keyPath, server.Https))
	}()
	return
}

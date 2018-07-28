package server

import (
	"strconv"
	"strings"

	"github.com/amyadzuki/amygolib/onfail"
	"github.com/amyadzuki/amygolib/str"

	"github.com/pkg/errors"
)

type Basic struct {
	ImplApi   func(Backend, uint32, ...interface{})
	ImplHttp  func(Backend, ...interface{})
	ImplHttps func(Backend, ...interface{})
	ImplServe func(Backend, string, string, string, string, onfail.OnFail)
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

func (server *Basic) Serve(backend Backend, httpAddr, httpsAddr, certPath, keyPath string, onFail onfail.OnFail) {
	if server.ImplServe != nil {
		server.ImplServe(backend, httpAddr, httpsAddr, certPath, keyPath, onFail)
	} else {
		basicServe(server, backend, httpAddr, httpsAddr, certPath, keyPath, onFail)
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

func basicServe(server Server, backend Backend, httpAddr, httpsAddr, certPath, keyPath string, onFail onfail.OnFail) {
	go func() {
		err := backend.Http(httpAddr, server.Http)
		onFail.Fail(errors.Wrap(err, "HTTP"))
	}()
	go func() {
		err := backend.Https(httpsAddr, certPath, keyPath, server.Https)
		onFail.Fail(errors.Wrap(err, "HTTPS"))
	}()
	return
}

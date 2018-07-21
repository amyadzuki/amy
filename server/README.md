# server
The `server` package contains a hackable HTTP/HTTPS server framework.

## Imports
```go
import "github.com/amyadzuki/amystuff/server"
```

## Usage Examples
```go
package main

import (
	"flag"
	"github.com/amyadzuki/amystuff/onfail"
	"github.com/amyadzuki/amystuff/server"
)

func main() {
	http := flag.Int("http", 80, "http port")
	https := flag.Int("https", 443, "https port")
	cert := flag.String("cert", "cert.pem", "Path to TLS certificate file")
	key := flag.String("key", "key.pem", "Path to TLS key file")
	flag.Parse()
	server.Default.ImplApi = api
	server.Serve(
		server.FastHttp,
		":"+strconv.Itoa(*http),
		":"+strconv.Itoa(*https),
		*cert,
		*key,
		onfail.Panic,
	)
}

func api(backend server.Backend, version uint32, args ...interface{}) {
	switch version {
	case 0x01000000:
		backend.Error(501, "Not implemented", args...)
	default:
		backend.Error(404, "Not found", args...)
	}
}
```

# str
The `str` package contains some mid-level string operations.

## Imports
```go
import "github.com/amyadzuki/amy/str"
```

## Usage Examples
```go
package main

import (
	"fmt"
	"github.com/amyadzuki/amy/str"
)

func main() {
	fmt.Printf("%v: \"%s\"\n", str.CaseHasPrefix("Hello, world!", "heLLo"), str.Simp("@Amy.Adzuki.1234!"))
}
```

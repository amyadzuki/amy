# query

The `query` package is designed to save time when using databases.  It has only been tested with sqlite.

## Imports

```go
import "github.com/amyadzuki/amy/query"
```

## Usage Example

```go
package example

import (
	"database/sql"
	"github.com/amyadzuki/amy/query"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func MySQLQuery(key interface{}) (values []string, err error) {
	q := query.NewByValue(db)
	// or q := query.New(db)
	// or q := query.Query{DB: db}
	// or var q query.Query; q.Init(db)
	q.SQL = `SELECT "value" FROM "Table" WHERE "key" = ?`
	q.Query(key)
	for q.NextOrClose() {
		var value string
		q.ScanKeepOpen(&value)
		if !q.Ok() {
			err = q.Error
			return
		}
		values = append(values, value)
	}
	return
}
```



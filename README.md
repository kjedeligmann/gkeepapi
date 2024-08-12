Inspired by [gkeepapi](https://github.com/kiwiz/gkeepapi) Python library.

This is a stub of the (mobile) Google Keep API client package, so it has limited functionality and is under construction. For now, you can use it to view all of your Google Keep notes as follows:

```go
package main

import (
	"fmt"
	"log"

	"github.com/kjedeligmann/gkeepapi"
)

var email, masterToken, gaid string // your credentials

func main() {
	var keep gkeepapi.Keep
	keep.Authenticate(email, gaid, masterToken)
	notes, err := keep.List()
	if err != nil {
		log.Fatal(err)
	}
	for _, note := range notes {
		fmt.Printf("%s\n%s\n\n", note.Title, note.Text)
	}
}
```

You can get `masterToken` using [gpsoauth](https://github.com/kjedeligmann/gpsoauth) package.

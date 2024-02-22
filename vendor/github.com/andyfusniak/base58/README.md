# base58

Package base58 provides a cryptographically secure random base58 string generator.

Base58 has an alphabet of 58 easily readable characters. It excludes letters that might look ambiguous when printed (0 – zero, I – capital i, O – capital o and l – lower-case L). Unlike Base64 it does not contain any URI reserved characters so is suitable for use in URL query parameters.

Example usage:

```go
package main

import (
	"fmt"
	"log"

	"github.com/andyfusniak/base58"
)

func main() {
	s, err := base58.RandString(8)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("base58 string %s\n", s)
}
```

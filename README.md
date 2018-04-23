# go-lsd-parametrized

Generate function of calculate Levestein Distance by specific parameters written in Go.

## Usage

```go
package main

import (
    "fmt"

    . "github.com/deltam/go-lsd-parametrized"
)

func main() {
    a, b := "kitten", "shitting"
    fmt.Printf("compare string: %s, %s\n", a, b)

    // normal lsd
    fmt.Printf("normal lsd = %f\n", Lsd(a, b))

    // custom lsd
    lsdf := GenerateLsdFunc(LevenshteinParam{Insert: 1, Delete: 0.1, Replace: 1})
    fmt.Printf("custom lsd = %f\n", lsdf(a, b))
}
```

```sh
$ go run main.go
compare string: kitten, shitting
normal lsd = 4.000000
custom lsd = 2.200000
```

## Use Case

- Clastering error messages
    - Decrease `Replace` parameter

## License

MIT License

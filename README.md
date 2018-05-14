# go-lsd-parametrized

Calculate Leveshtein Distance by specific parameters written in Go.

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
    params := LevenshteinParam{Insert: 0.1, Delete: 1, Replace: 0.01}
    fmt.Printf("custom lsd = %f\n", params.Distance(a, b))
}
```

```sh
$ go run main.go
compare string: kitten, shitting
normal lsd = 4
custom lsd = 0.220000
```

## Use Case

- Clastering error messages

## License

MIT License

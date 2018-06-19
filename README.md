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

    // standard
    fmt.Printf("standard = %d\n", Lsd(a, b))

    // weighted
    wd := Weights{Insert: 0.1, Delete: 1, Replace: 0.01}
    fmt.Printf("weighted = %f\n", wd.Distance(a, b))

    // weighted and normalized
    nd := Normalized(wd)
    fmt.Printf("normalized = %f\n", nd.Distance(a, b))
}
```

```sh
$ go run main.go
compare string: kitten, shitting
standard = 4
weighted = 0.220000
normalized = 0.027500
```

## Use Case

- Clastering error messages

## License

MIT License

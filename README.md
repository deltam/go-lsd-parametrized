# go-lsd-parametrized

Calculate Leveshtein Distance by specific parameters written in Go.

[godoc](https://godoc.org/github.com/deltam/go-lsd-parametrized)

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

    // weighted by rune
    wr := ByRune(&Weights{1, 1, 1}).
        Insert("g", 0.1).
        Insert("h", 0.01).
        Replace("k", "s", 0.001).
        Replace("e", "i", 0.0001)
    fmt.Printf("rune weight = %f\n", wr.Distance(a, b))
}
```

```sh
$ go run main.go
compare string: kitten, shitting
standard = 4
weighted = 0.220000
normalized = 0.027500
rune weight = 0.111100
```

## Custom Distance

```go
type LengthDiff struct{}

func (_ LengthDiff) Distance(a, b string) float64 {
    d := utf8.RuneCountInString(a) - utf8.RuneCountInString(b)
    return math.Abs(float64(d))
}

func main() {
    d := LengthDiff{}
    fmt.Println(d.Distance("kitten", "shitting"))
    // Output:
    // 2

    group := []string{"", "a", "ab", "abc"}
    s, dist := lsdp.Nearest(d, "xx", group)
    fmt.Println(s, dist)
    // Output:
    // ab 0
}
```

## Use Case

- Clustering error messages

## License

MIT License

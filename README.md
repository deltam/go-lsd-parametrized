# go-lsd-parametrized

Weighted Leveshtein Distance and its extended interfaces written in Go.

[godoc](https://godoc.org/github.com/deltam/go-lsd-parametrized)

## Installation

```sh
go get -u github.com/deltam/go-lsd-parametrized
```

## Usage

```go
package main

import (
    "fmt"

    "github.com/deltam/go-lsd-parametrized"
)

func main() {
    a, b := "kitten", "shitting"

    // standard
    fmt.Println(lsdp.Lsd(a, b))
    // Output:
    // 4

    // weighted
    wd := lsdp.Weights{Insert: 0.1, Delete: 1, Replace: 0.01}
    fmt.Println(wd.Distance(a, b))
    // Output:
    // 0.22

    // weighted and normalized
    nd := lsdp.Normalized(wd)
    fmt.Println(nd.Distance(a, b))
    // Output:
    // 0.0275

    // weighted by rune
    wr := lsdp.ByRune(&lsdp.Weights{1, 1, 1}).
        Insert("g", 0.1).
        Insert("h", 0.01).
        Replace("k", "s", 0.001).
        Replace("e", "i", 0.0001)
    fmt.Println(wr.Distance(a, b))
    // Output:
    // 0.1111
}
```

## Operators

```go
func main() {
    std := lsdp.Weights{1, 1, 1}
    fruits := []string{"apple", "orange", "lemon", "water melon"}

    // find nearest string
    s, d := lsdp.Nearest(std, "aple", fruits)
    fmt.Println(s, d)
    // Output:
    // apple 2

    // calculate distance of each strings
    ds := lsdp.DistanceAll(std, "aple", fruits)
    fmt.Println(ds)
    // Output:
    // [1 4 5 9]
}
```

## Custom Distance

```go
func lenDiff(a, b string) float64 {
    d := utf8.RuneCountInString(a) - utf8.RuneCountInString(b)
    return math.Abs(float64(d))
}

func main() {
    var d lsdp.DistanceFunc = lenDiff
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

Composite two Distances

```go
func Far(dm1, dm2 lsdp.DistanceMeasurer) lsdp.DistanceMeasurer {
    return &far{dm1: dm1, dm2: dm2}
}

type far struct {
    dm1, dm2 lsdp.DistanceMeasurer
}

func (f *far) Distance(a, b string) float64 {
    d1 := f.dm1.Distance(a, b)
    d2 := f.dm2.Distance(a, b)
    if d1 > d2 {
        return d1
    }
    return d2
}

func main() {
    a, b := "kitten", "shitting"

    std := lsdp.Weights{Insert: 1, Delete: 1, Replace: 1}
    fmt.Println(std.Distance(a, b))
    // Output:
    // 4

    wd := lsdp.Weights{Insert: 10, Delete: 1, Replace: 0.1}
    fmt.Println(wd.Distance(a, b))
    // Output:
    // 20.2

    fd := Far(std, wd)
    fmt.Println(fd.Distance(a, b))
    // Output:
    // 20.2
}
```

## Use Case

- Clustering error messages

## License

MIT License

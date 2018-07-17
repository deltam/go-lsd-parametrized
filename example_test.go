package lsdp_test

import (
	"fmt"

	"github.com/deltam/go-lsd-parametrized"
)

func ExampleWeights_Distance() {
	wd := lsdp.Weights{Insert: 0.1, Delete: 1, Replace: 0.01}
	fmt.Println(wd.Distance("kitten", "shitting"))
	// Output: 0.22
}

func ExampleByRune() {
	wr := lsdp.ByRune(&lsdp.Weights{1, 1, 1}).
		Insert("a", 0.1).
		Delete("b", 0.01).
		Replace("c", "d", 0.001)
	fmt.Println(wr.Distance("bc", "ad"))
	// Output: 0.111
}

func ExampleNearest() {
	std := lsdp.Weights{1, 1, 1}
	group := []string{"apple", "orange", "lemon", "water melon"}
	fmt.Println(lsdp.Nearest(std, "mon", group))
	// Output:
	// lemon 2
}

func ExampleDistanceAll() {
	std := lsdp.Weights{1, 1, 1}
	group := []string{"apple", "orange", "lemon", "water melon"}
	fmt.Println(lsdp.DistanceAll(std, "mon", group))
	// Output:
	// [5 5 2 8]
}

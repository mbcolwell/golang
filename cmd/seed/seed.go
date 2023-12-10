package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mbcolwell/golang/internal/seed"
)

func main() {
	start := time.Now()
	input := flag.String("input-file", "data/seed/input.txt", "Input file containing seeds and maps")
	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Unable to parse flags %v\n", err)
		os.Exit(1)
	}

	seeds, maps := seed.ReadInput(*input)
	fmt.Printf("Took %s to read input\n", time.Since(start))
	start = time.Now()

	fmt.Printf("Lowest location number given the seed inputs: %d\n", seed.MinSeed(seeds, maps))
	fmt.Printf("Took %s to do part 1\n", time.Since(start))
	start = time.Now()
	fmt.Printf("Lowest location number given the range seed inputs: %d\n", seed.MinSeedRange(seeds, maps))
	fmt.Printf("Took %s to do part 2\n", time.Since(start))
}

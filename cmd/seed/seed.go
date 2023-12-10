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

	fmt.Printf("Lowest location number given the seed inputs: %d\n", seed.MinSeed(seeds, maps))

	fmt.Printf("Seeds took %s\n", time.Since(start))
}

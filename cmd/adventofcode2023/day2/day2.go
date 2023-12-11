package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mbcolwell/golang/internal/adventofcode2023/day2"
)

func main() {
	start := time.Now()
	input := flag.String(
		"input-file", "internal/adventofcode2023/day2/input.txt", "Input file containing seeds and maps",
	)
	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Unable to parse flags %v\n", err)
		os.Exit(1)
	}

	part1_values := day2.Part1Games(*input)
	fmt.Printf("Part 1: sum of values is %d\n", day2.SumNumbers(part1_values))
	part2_values := day2.Part2Powers(*input)
	fmt.Printf("Part 2: sum of values is %d\n", day2.SumNumbers(part2_values))

	fmt.Printf("Took %s\n", time.Since(start))
}

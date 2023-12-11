package day2

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func Part1Games(filepath string) []int {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var lineTxt string
	var games []int
	var g int

	maxColor := map[string]int{
		"red":   12,
		"green": 13,
		"blue":  14,
	}

	scanner := bufio.NewScanner(file)

out:
	for scanner.Scan() {
		lineTxt = scanner.Text()[5:] // Remove "Game "

		lsplit := strings.Split(lineTxt, ": ")
		g, err = strconv.Atoi(lsplit[0])
		if err != nil {
			log.Fatalf("Unable to convert game to integer: %v\n", err)
		}

		for _, gameString := range strings.Split(lsplit[1], "; ") {
			for _, cube := range strings.Split(gameString, ", ") {
				cubeData := strings.Split(cube, " ")
				n, err := strconv.Atoi(cubeData[0])
				if err != nil {
					log.Fatalf("Unable to convert cube to integer: %v\n", err)
				}
				if n > maxColor[cubeData[1]] {
					continue out
				}
			}
		}

		games = append(games, g)
	}
	return games
}

func Part2Powers(filepath string) []int {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var lineTxt string
	var powers []int
	var r, g, b int

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		r, g, b = 0, 0, 0
		lineTxt = strings.Split(scanner.Text(), ": ")[1]

		for _, gameString := range strings.Split(lineTxt, "; ") {
			for _, cube := range strings.Split(gameString, ", ") {
				cubeData := strings.Split(cube, " ")
				n, err := strconv.Atoi(cubeData[0])
				if err != nil {
					log.Fatalf("Unable to convert cube to integer: %v\n", err)
				}
				if cubeData[1] == "red" && n > r {
					r = n
				}
				if cubeData[1] == "green" && n > g {
					g = n
				}
				if cubeData[1] == "blue" && n > b {
					b = n
				}
			}
		}

		powers = append(powers, r*g*b)
	}

	return powers
}

func SumNumbers(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}
	return numbers[0] + SumNumbers(numbers[1:])
}

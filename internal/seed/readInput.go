package seed

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type SeedMap struct {
	Source      []int
	Destination []int
	Range       []int
}

func ReadInput(filepath string) (seeds []int, maps []SeedMap) {
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
	seeds = []int{}
	maps = []SeedMap{}
	newMap := SeedMap{}

	scanner := bufio.NewScanner(file)

	// Get the seeds row
	scanner.Scan()
	lineTxt = scanner.Text()
	for _, s := range strings.Split(lineTxt, " ")[1:] {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("Unable to convert see to integer: %v\n", err)
		}
		seeds = append(seeds, i)
	}
	scanner.Scan() // Skips the blank line

	// Get all the maps
	for scanner.Scan() {
		lineTxt = scanner.Text()

		if len(lineTxt) == 0 {
			maps = append(maps, newMap)
			newMap = SeedMap{}
			continue
		}
		if lineTxt[len(lineTxt)-4:] == "map:" {
			continue
		}

		for n, s := range strings.Split(lineTxt, " ") {
			i, err := strconv.Atoi(s)
			if err != nil {
				log.Fatalf("Unable to convert see to integer: %v\n", err)
			}
			if n == 0 {
				newMap.Destination = append(newMap.Destination, i)
			} else if n == 1 {
				newMap.Source = append(newMap.Source, i)
			} else {
				newMap.Range = append(newMap.Range, i)
			}
		}

	}
	maps = append(maps, newMap)

	return seeds, maps
}

package seed

import "slices"

func seedToLocation(seed int, maps []SeedMap) int {
	for _, m := range maps {
		for i := 0; i < len(m.Source); i++ {
			if seed >= m.Source[i] && seed < m.Source[i]+m.Range[i] {
				seed += m.Destination[i] - m.Source[i]
				break
			}
		}
	}
	return seed
}

func MinSeed(seeds []int, maps []SeedMap) int {
	locations := make([]int, len(seeds))
	for i, s := range seeds {
		locations[i] = seedToLocation(s, maps)
	}

	return slices.Min(locations)
}

package seed

import "slices"

func seedToLocationSlow(seed int, maps []SeedMap) int {
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
		locations[i] = seedToLocationSlow(s, maps)
	}

	return slices.Min(locations)
}

type seedRange struct {
	starts  []int
	lengths []int
}

func mapRange(initialisation seedRange, maps []SeedMap) seedRange {
	var startRange, endRange int
	var output seedRange

	for n, m := range maps {
		if n != 0 {
			initialisation = output
		}
		output = seedRange{}
		for sr := 0; sr < len(initialisation.starts); sr++ {
			startRange = initialisation.starts[sr]
			endRange = initialisation.starts[sr] + initialisation.lengths[sr]
			for r := 0; r < len(m.Source); r++ {
				left := startRange >= m.Source[r] && startRange < m.Source[r]+m.Range[r]
				right := endRange >= m.Source[r] && endRange < m.Source[r]+m.Range[r]

				conversion := m.Destination[r] - m.Source[r]

				if left && right { // Input is a subset of the output
					output.starts = append(output.starts, startRange+conversion)
					output.lengths = append(output.lengths, initialisation.lengths[sr])
				} else if left { // Left half of input overlaps with right half of output
					newStart := startRange + conversion
					output.starts = append(output.starts, newStart)
					output.lengths = append(output.lengths, m.Destination[r]+m.Range[r]-newStart)
				} else if right { // Right half of input overlaps with left half of output
					output.starts = append(output.starts, m.Destination[r])
					output.lengths = append(
						output.lengths, startRange+initialisation.lengths[sr]+conversion-m.Destination[r],
					)
				} else if startRange < m.Source[r] && endRange >= m.Source[r]+m.Range[r] { // Input is a superset
					output.starts = append(output.starts, m.Destination[r])
					output.lengths = append(output.lengths, m.Range[r])
				}
			}
		}
	}
	return output
}

func MinSeedRange(seeds []int, maps []SeedMap) int {
	sr := seedRange{}
	for i, s := range seeds {
		if i%2 == 0 {
			sr.starts = append(sr.starts, s)
		} else {
			sr.lengths = append(sr.lengths, s)
		}
	}
	output := mapRange(sr, maps)
	return slices.Min(output.starts)
}

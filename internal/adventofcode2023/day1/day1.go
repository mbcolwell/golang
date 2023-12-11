package day1

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func ParseInput(filepath string, part2 bool) []int {
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
	var newInt int
	values := []int{}

	intLiterals := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		newInt = 0
		lineTxt = scanner.Text()

	out1:
		for i := 0; i < len(lineTxt); i++ {
			if part2 {
				for n, il := range intLiterals {
					if i+len(il) > len(lineTxt) {
						continue
					}
					if lineTxt[i:i+len(il)] == il {
						newInt += (n + 1) * 10
						break out1
					}
				}
			}
			n, err := strconv.Atoi(lineTxt[i : i+1])
			if err == nil {
				newInt += n * 10
				break
			}
		}

	out2:
		for i := len(lineTxt); i > 0; i-- {
			if part2 {
				for n, il := range intLiterals {
					if i-len(il) < 0 {
						continue
					}
					if lineTxt[i-len(il):i] == il {
						newInt += n + 1
						break out2
					}
				}
			}
			n, err := strconv.Atoi(lineTxt[i-1 : i])
			if err == nil {
				newInt += n
				break
			}
		}
		values = append(values, newInt)
	}

	return values
}

func SumNumbers(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}
	return numbers[0] + SumNumbers(numbers[1:])
}

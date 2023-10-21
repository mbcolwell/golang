package quiz

import (
	"encoding/csv"
	"fmt"
	"os"
)

type problem struct {
	q string
	a string
}

func exit(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func GetCsvProblems(filepath string) []problem {
	fd, err := os.Open(filepath)
	defer fd.Close()

	exit(err, fmt.Sprintf("Unable to open %s\n", filepath))

	fileReader := csv.NewReader(fd)
	records, err := fileReader.ReadAll()

	exit(err, fmt.Sprintf("Unable to read %s\n", filepath))

	ret := make([]problem, len(records))
	for i, record := range records {
		ret[i] = problem{
			q: record[0],
			a: record[1],
		}
	}
	return ret
}

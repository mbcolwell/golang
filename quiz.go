package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
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

func Timer(sleepTime time.Duration, done chan bool) {
	time.Sleep(sleepTime * time.Second)
	fmt.Println("\n\nOut of time!!")
	done <- true
}

func QuizMaster(problems []problem, n *int, done chan bool) {
	var resp string
	for i, p := range problems {
		fmt.Printf("Problem %d: %s ", i+1, p.q)
		fmt.Scanf("%s\n", &resp)
		if resp == p.a {
			*n++
		}
	}
	done <- true
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file of format 'question,answer'")
	t := flag.Int("timer", 30, "time limit to answers all quiz questions in seconds")
	flag.Parse()

	problems := GetCsvProblems(*csvFilename)

	n := 0

	done := make(chan bool)
	go QuizMaster(problems, &n, done)
	go Timer(time.Duration(*t), done)
	<-done

	fmt.Printf("You scored %d out of %d\n", n, len(problems))
}

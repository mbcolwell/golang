package quiz

import (
	"flag"
	"fmt"
	"time"
)

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

func Quiz() {
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

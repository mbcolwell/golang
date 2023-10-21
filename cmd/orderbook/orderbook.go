package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/mbcolwell/golang/internal/orderbook"
)

func main() {
	outputFile := flag.String("log-filepath", "output.log", "File to store the output of the program in")
	console_out := flag.Bool("print-con", false, "Whether to print the output to the console")
	err := flag.CommandLine.Parse(os.Args[2:])

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Value for n is not an integer")
		os.Exit(1)
	}

	fmt.Println(*console_out)
	fmt.Println(n)
	fmt.Println(*outputFile)

	orderbook.ReadStream(os.Stdin)
}

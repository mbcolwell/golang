package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mbcolwell/golang/internal/orderbook"
)

func main() {
	start := time.Now()
	outputFile := flag.String("log-filepath", "output.log", "File to store the output of the program in")
	console_out := flag.Bool("print-con", false, "Whether to print the output to the console")
	err := flag.CommandLine.Parse(os.Args[2:])

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Value for n is not an integer: %v\n", err)
		os.Exit(1)
	}

	var wrt io.Writer
	out, err := os.OpenFile(*outputFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()
	if *console_out {
		wrt = io.MultiWriter(os.Stdout, out)
	} else {
		wrt = io.Writer(out)
	}
	log.SetOutput(wrt)
	log.SetFlags(0)

	book := map[string]orderbook.Ladder{}

	reader := bufio.NewReader(os.Stdin)
	var msg orderbook.Message
	EOF := 0

	for {
		msg, EOF = orderbook.ReadMessage(reader)
		if EOF == 1 {
			break
		}

		if orderbook.ProcessMessage(n, msg, &book) {
			ticker := string(msg.Order.Symbol[:])
			log.Println(orderbook.FormatLadder(
				n, ticker, msg.Header.SeqNo, *book[ticker+"B"].Depth, *book[ticker+"S"].Depth))
		}
	}
	log.Printf("Orderbook took %s\n", time.Since(start))
}

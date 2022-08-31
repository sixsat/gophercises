package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	csvFile := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	quizzes, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	score := 0

outerloop:
	for i, quiz := range quizzes {
		fmt.Printf("Problem #%d: %s = ", i+1, quiz[0])

		var input string
		ansCh := make(chan string)
		go func() {
			fmt.Scanln(&input)
			ansCh <- input
		}()

		select {
		case <-ansCh:
			if input == strings.TrimSpace(quiz[1]) {
				score++
			}
		case <-time.After(time.Duration(*timeLimit) * time.Second):
			fmt.Println()
			break outerloop
		}
	}

	fmt.Printf("You scored %d out of %d.\n", score, len(quizzes))
}

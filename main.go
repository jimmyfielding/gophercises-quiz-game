package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

type problem struct {
	Question string
	Answer   string
}

func main() {
	var fileName string
	var timeLimit int
	var shuffle bool
	flag.StringVar(&fileName, "f", "problems.csv", "-f test.csv, filename of csv containing quiz")
	flag.IntVar(&timeLimit, "t", 30, "-t 150, timelimit of the quiz timer")
	flag.BoolVar(&shuffle, "s", false, "-s shuffle problem order")
	flag.Parse()
	problems, err := readProblems(fileName)
	if err != nil {
		log.Fatal(err)
	}

	askProblems(problems, timeLimit, shuffle)
}

func askProblems(problems []problem, timeLimit int, shuffle bool) {
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	correct := 0
	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
	}

problemloop:
	for i, p := range problems {
		fmt.Printf("#%d: %s\n", i, p.Question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scan(&answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemloop
		case answer := <-answerCh:
			if answer == p.Answer {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func readProblems(filePath string) ([]problem, error) {
	csvfile, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	r := csv.NewReader(csvfile)

	var ps []problem
	var p problem
	for {
		s, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		p = problem{
			Question: s[0],
			Answer:   s[1],
		}

		ps = append(ps, p)
	}

	return ps, nil
}

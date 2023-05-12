package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var colorReset = "\033[0m"
var colorRed = "\033[31m"
var colorGreen = "\033[32m"

func main() {
	csvFile := flag.String("csv", "problems.csv", "a csv file in 'question', 'answer' format")
	timeLimit := flag.Int("limit", 60, "time limit in seconds")
	flag.Parse()
	fmt.Printf("\n🏆 Welome to the QUIZ Game, Time Limit (%d) seconds 🏆\n\n", *timeLimit)

	questionsAndAnswers := readCsv(*csvFile)

	reader := bufio.NewReader(os.Stdin)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for question, correctAnswer := range questionsAndAnswers {
		fmt.Printf("%s = ", question)
		answerCh := make(chan string)
		go func() {
			userResponse, err := reader.ReadString('\n') // can also use Scanf, it removes spaces
			if err != nil {
				log.Fatal(err)
			}
			userResponse = strings.TrimSpace(userResponse)
			answerCh <- userResponse
		}()
		select {
		case <-timer.C:
			fmt.Printf("\n⏰ Time is up!\n")
			fmt.Printf("Score = %d/%d\n", correct, len(questionsAndAnswers))
			return
		case answer := <-answerCh:

			if correctAnswer == answer { // can also use strings.Compare(response, answer) == 0
				fmt.Print(string(colorGreen))
				fmt.Println("✅ Correct", string(colorReset))
				correct++
			} else {
				fmt.Print(string(colorRed))
				fmt.Println("❌ Incorrect", string(colorReset))
			}
		}

	}

	fmt.Printf("Score = %d/%d\n", correct, len(questionsAndAnswers))
}

func readCsv(filename string) map[string]string {
	m := make(map[string]string)
	f, err := os.Open(filename)
	if err != nil {
		fmt.Print(string(colorRed))
		fmt.Printf("Failed top open the CSV file: %q", filename)
		fmt.Println(string(colorReset))
		os.Exit(1)
	}
	fmt.Print(string(colorGreen))
	fmt.Println("Successfully opened CSV file", string(colorReset))
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range lines {
		m[line[0]] = strings.TrimSpace(line[1])
	}
	return m
}

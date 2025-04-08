package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const ARG_VALUE_MISSING = 1
const ARG_VALUE_MISSING_OR_INVALID = 2
const CSV_READ_ERROR = 3
const INTERNAL_ERROR = 10

func readFile(path string) (*[][]string, error) {
	file, fileOpenErr := os.Open(path)

	if fileOpenErr != nil {
		return nil, fmt.Errorf("the file at path %s could not be opened", path)
	}

	reader := csv.NewReader(file)

	questions, fileReadError := reader.ReadAll()

	if fileReadError != nil {
		return nil, errors.New("there was a problem reading the CSV file provided")
	}

	return &questions, nil
}

func main() {
	filePath := flag.String("csv", "./problems.csv", "The CSV file to read. Default is ./problems.csv")
	timeArg := flag.Int("time", 30, "The timer for the quiz")

	flag.Parse()

	questions, questionParseErr := readFile(*filePath)

	if questionParseErr != nil {
		fmt.Println(questionParseErr.Error())
		os.Exit(CSV_READ_ERROR)
	}

	playGame(questions, *timeArg)
}

func playGame(questions *[][]string, seconds int) {
	fmt.Println("Press the enter key to start the quiz!")

	// Wait for user to input to start quiz
	fmt.Scanln()

	timer := time.NewTimer(time.Duration(seconds) * time.Second)

	numCorrect := 0
	numWrong := 0

	go func() {
		<-timer.C
		fmt.Println("\nGame over! You ran out of time!")
		calculateAndPrintScore(numCorrect, len(*questions))
		os.Exit(0)
	}()

	for i, item := range *questions {
		fmt.Printf("%s ", item[0])
		var ans string
		_, scanErr := fmt.Scanln(&ans)

		if scanErr != nil {
			fmt.Println("A problem occurred reading your last answer. Try again")
			i--
			continue
		}

		if strings.Trim(ans, " ") == strings.Trim(item[1], " ") {
			numCorrect++
		} else {
			numWrong++
		}
	}

	timer.Stop()
	calculateAndPrintScore(numCorrect, len(*questions))
}

func calculateAndPrintScore(correct int, numQuestions int) {
	percent := (float32(correct) / float32(numQuestions)) * 100.00

	fmt.Printf("You got %d correct out of %d questions. Your score is %2.2f%%\n", correct, numQuestions, percent)
}

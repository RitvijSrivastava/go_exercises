package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

var filePath string = "./quiz_game/problems.csv"

type problem struct {
	question string
	answer   string
}

func extractQuestionAndAnswer(question string) problem {
	lastIndex := strings.LastIndex(question, ",")
	if lastIndex == -1 {
		panic("Invalid Question!")
	}

	return problem{
		question: strings.TrimSpace(question[:lastIndex]),
		answer:   strings.TrimSpace(question[lastIndex+1:]),
	}
}

func main() {

	// Open a file
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	csvFileReader := csv.NewReader(file)

	var problems []problem // To store all the questions and answers

	// Extract and parse problems
	for {

		if statement, _ := csvFileReader.Read(); statement == nil {
			break
		} else {
			question := statement[0]
			answer := statement[1]

			question = strings.TrimSpace(question)
			answer = strings.TrimSpace(answer)

			newProblem := problem{
				question: question,
				answer:   answer,
			}
			problems = append(problems, newProblem)
		}

	}

	// Run quiz
	totalQuestions := len(problems)
	correctAnswers := 0

	consoleReader := bufio.NewReader(os.Stdin)

	for _, statement := range problems {
		fmt.Print(statement.question, "? ")

		// Take user input
		userInput, _ := consoleReader.ReadString('\n')

		// Clean the user input
		userInput = strings.Replace(userInput, "\n", "", -1)
		userInput = strings.TrimSpace(userInput)

		if strings.Compare(userInput, statement.answer) == 0 {
			correctAnswers++
		}
	}

	// Print the result
	fmt.Printf("\nYou Scored: %v/%v\n", correctAnswers, totalQuestions)

}

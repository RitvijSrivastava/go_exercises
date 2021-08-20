package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const defaultFilePath = "./problems.csv"

type problem struct {
	question string
	answer   string
}

func main() {

	// Setup command line flags
	var quizFilePath string
	flag.StringVar(&quizFilePath, "file", defaultFilePath, "fully legal path to the csv file containing the quiz")

	var shouldRandomizeQuestions bool
	flag.BoolVar(&shouldRandomizeQuestions, "randomize", false, "set if questions should be randomized")

	var quizTime int
	flag.IntVar(&quizTime, "time", 30, "time required to complete the quiz (in seconds)")

	flag.Parse()

	// Open a file
	file, err := os.Open(quizFilePath)
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

	// Randomize Question order if the "randomize" flag is set.
	if shouldRandomizeQuestions {
		randomizeQuestions(problems)
	}

	// Run quiz
	runQuiz(problems, quizTime)
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

func randomizeQuestions(problems []problem) {
	rand.Seed(time.Now().UnixMilli())
	rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
}

func runQuiz(problems []problem, quizTime int) {

	quizTimer := time.NewTimer(time.Duration(quizTime) * time.Second)

	totalQuestions := len(problems)
	correctlyAnswered := 0
	questionsAttempted := 0

	consoleReader := bufio.NewReader(os.Stdin)

	// Start quiz after the user presses "Enter"
	fmt.Printf("Press Enter to start the quiz")
	consoleReader.ReadString('\n')

quiz:
	for _, statement := range problems {
		fmt.Print(statement.question, "? ")
		userAnswerChannel := make(chan string)

		go func() {
			// Take user input
			userInput, _ := consoleReader.ReadString('\n')

			// Clean the user input
			userInput = strings.Replace(userInput, "\n", "", -1)
			userInput = strings.TrimSpace(userInput)

			userAnswerChannel <- userInput
		}()

		select {
		case <-quizTimer.C:
			break quiz
		case userInput := <-userAnswerChannel:
			if strings.Compare(userInput, statement.answer) == 0 {
				correctlyAnswered++
			}
			questionsAttempted++
		}
	}

	// Print result
	printResult(correctlyAnswered, questionsAttempted, totalQuestions)
}

func printResult(correctlyAnswered, questionsAttempted, totalQuestions int) {

	fmt.Printf("\nResult:\n")
	fmt.Println("Questions Attempted:", questionsAttempted)
	fmt.Println("Correct Answers:", correctlyAnswered)
	fmt.Printf("\nYour score: %v/%v\n", correctlyAnswered, totalQuestions)
}
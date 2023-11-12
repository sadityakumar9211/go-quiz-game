package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {

	var path string
	var duration int
	var isShuffle bool

	flag.StringVar(&path, "file", "problems.csv", "Specify a CSV file containing records")
	flag.IntVar(&duration, "limit", 30, "Specify a quiz duration")
	flag.BoolVar(&isShuffle, "shuffle", false, "Specify shuffle option")

	flag.Parse()

	records, err := ReadCSV(path)

	if isShuffle {
		shuffle(records)
	}

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	fmt.Println("Successfully! Parsed the Questions file")
	for i := 3; i > 0; i -= 1 {
		fmt.Printf("Starting the Quiz in %v...\n", i)
		time.Sleep(time.Second)
	}
	fmt.Print("\n\n\n")

	var correctCount int
	var incorrectCount int
	doneC := make(chan struct{})

	// kind of a timer which sends of a signal to doneC channel to
	// which another go routinte is listening for stopping the quiz
	go func() {
		time.Sleep(time.Duration(duration * int(time.Second)))
		doneC <- struct{}{}
	}()

	for i, record := range records {
		select {
		case <-doneC:
			goto end
		default:
			fmt.Printf("Problem #%v: %v - ", i+1, record[0])
			var userInput string
			fmt.Scan(&userInput)
			userInput = strings.Trim(userInput, " \n")
			if strings.EqualFold(userInput, record[0]) {
				correctCount += 1
			} else {
				incorrectCount += 1
			}
		}
	}
end:
	fmt.Println("You have successfully completed the quiz game.")
	fmt.Printf("Your Score: %v out of %v", correctCount, correctCount+incorrectCount)

}

// Reads the CSV file and prints its content line by line
func ReadCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	// Checks for the error
	if err != nil {
		return nil, err
	}

	// Closes the file
	defer file.Close()

	// The csv.NewReader() function is called in
	// which the object os.File passed as its parameter
	// and this creates a new csv.Reader that reads
	// from the file
	reader := csv.NewReader(file)

	// ReadAll reads all the records from the CSV file
	// and Returns them as slice of slices of string
	// and an error if any
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// This is an acceptable implementation
func shuffle(records [][]string) {
	// Using the Durstenfeld Shuffle (Modern Fisher Yates Shuffle)
	lastUnshuffledIdx := len(records) - 1
	for i := 0; i < len(records); i++ {
		j := rand.Intn(len(records) - i)
		records[j], records[lastUnshuffledIdx] = records[lastUnshuffledIdx], records[j]
		lastUnshuffledIdx--
	}
}

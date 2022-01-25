package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/vidhanio/wordle"
)

func main() {
	validWords := make([]string, 370103)

	// Load the words into memory
	file, err := os.Open("words.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		validWords = append(validWords, scanner.Text())
	}

	// Close the file
	file.Close()

	commonWords := make([]string, 10000)

	// Load the words into memory
	file, err = os.Open("commonwords.txt")
	if err != nil {
		panic(err)
	}

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		commonWords = append(commonWords, scanner.Text())
	}

	// Close the file
	file.Close()

	// Create a new game
	w, err := wordle.New(5, 100000, commonWords, validWords)
	if err != nil {
		panic(err)
	}

	// ask user for input
	reader := bufio.NewReader(os.Stdin)
	for !w.Won() {
		fmt.Print("Enter a guess: ")
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]
		charGuesses, err := w.Guess(text)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(charGuesses)
		}
	}
}

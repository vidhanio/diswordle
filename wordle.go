package wordle

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type GuessType int

const (
	GuessTypeCorrect GuessType = iota
	GuessTypeWrongPosition
	GuessTypeWrong
)

type void struct{}

type Wordle struct {
	word []rune // Correct word

	commonWords map[string]void // Common words to use as correct words
	validWords  map[string]void // All valid words

	guessesAllowed int           // Number of guesses allowed
	wordLength     int           // Length of word
	guesses        [][]rune      // Guesses
	guessTypes     [][]GuessType // Character guesses

	cancelled bool
}

func New(wordLength, guessesAllowed int, commonWords, validWords []string) (*Wordle, error) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	// Get random word with the given length

	randomWord := make([]rune, wordLength)
	wordsWithLength := make([][]rune, 0)

	for _, w := range commonWords {
		if len(w) == wordLength {
			wordsWithLength = append(wordsWithLength, []rune(w))
		}
	}

	if len(wordsWithLength) == 0 {
		return nil, fmt.Errorf("invalid word length: no words with length %d in common words", wordLength)
	}

	randomWord = wordsWithLength[r.Intn(len(wordsWithLength))]

	// Initialize Wordle
	w := &Wordle{
		word:           randomWord,
		validWords:     make(map[string]void, len(validWords)),
		guessesAllowed: guessesAllowed,
		wordLength:     wordLength,
		guesses:        make([][]rune, 0),
	}

	// Convert valid words to runes
	for _, word := range validWords {
		w.validWords[word] = void{}
	}

	return w, nil
}

// Guess a word
func (w *Wordle) Guess(guess string) ([]GuessType, error) {
	if w.Done() {
		return nil, fmt.Errorf("game is done")
	}

	guess = strings.ToLower(guess)
	guessRunes := []rune(guess)

	if guessRunes == nil || len(guessRunes) != w.wordLength {
		return nil, fmt.Errorf("invalid word length: wanted %d, got %d", w.wordLength, len(guessRunes))
	}

	if _, ok := w.validWords[guess]; !ok {
		return nil, fmt.Errorf("invalid word: %s was not found in the dictionary", guess)
	}

	charGuesses := make([]GuessType, w.wordLength)

	for i, g := range guessRunes {
		if g < 'a' || g > 'z' {
			return nil, fmt.Errorf("invalid word: %c is not a letter", g)
		}

		if g == w.word[i] {
			charGuesses[i] = GuessTypeCorrect
		} else if contains(w.word, g) {
			charGuesses[i] = GuessTypeWrongPosition
		} else {
			charGuesses[i] = GuessTypeWrong
		}
	}

	w.guesses = append(w.guesses, guessRunes)
	w.guessTypes = append(w.guessTypes, charGuesses)

	return charGuesses, nil
}

func (w *Wordle) Word() string {
	if w.Done() {
		return string(w.word)
	}

	return ""
}

func (w *Wordle) WordLength() int {
	return w.wordLength
}

func (w *Wordle) Guesses() []string {
	guesses := make([]string, len(w.guesses))

	for i, g := range w.guesses {
		guesses[i] = string(g)
	}

	return guesses
}

func (w *Wordle) GuessTypes() [][]GuessType {
	return w.guessTypes
}

func (w *Wordle) GuessesLeft() int {
	return w.guessesAllowed - len(w.guesses)
}

func (w *Wordle) Cancel() {
	w.cancelled = true
}

func (w *Wordle) Cancelled() bool {
	return w.cancelled
}

func (w *Wordle) Won() bool {
	if len(w.guesses) == 0 {
		return false
	}

	return equal(w.guesses[len(w.guesses)-1], w.word)
}

func (w *Wordle) Lost() bool {
	return len(w.guesses) >= w.guessesAllowed
}

func (w *Wordle) Done() bool {
	return w.Won() || w.Lost() || w.Cancelled()
}

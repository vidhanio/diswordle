package wordle

import "errors"

var (
	ErrNoWordWithLength = errors.New("No word with length")
	ErrInvalidWord      = errors.New("Invalid word")
	ErrTooManyGuesses   = errors.New("Too many guesses")
	ErrInvalidGuess     = errors.New("Invalid guess")
)

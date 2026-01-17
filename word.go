package wordlegameengine

import (
	"fmt"
	"sort"
)

const WordLength = 5

type Word [WordLength]byte

var ErrInvalidLength = fmt.Errorf("word must be %d letters", WordLength)
var ErrInvalidCharacter = fmt.Errorf("word must contain only lowercase a-z")
var ErrNotInWordlist = fmt.Errorf("word not in allowed guesses")

func NewWord(s string) (Word, error) {
	if len(s) != WordLength {
		return Word{}, ErrInvalidLength
	}
	for i := 0; i < len(s); i++ {
		if s[i] < 'a' || s[i] > 'z' {
			return Word{}, ErrInvalidCharacter
		}
	}
	var w Word
	copy(w[:], s)
	return w, nil
}

func (w Word) String() string {
	return string(w[:])
}

func (w *Word) Validate() error {
	for i := 0; i < WordLength; i++ {
		if w[i] < 'a' || w[i] > 'z' {
			return ErrInvalidCharacter
		}
	}

	s := w.String()
	idx := sort.Search(len(AllowedGuesses), func(i int) bool {
		return AllowedGuesses[i].String() >= s
	})

	if idx < len(AllowedGuesses) && AllowedGuesses[idx].String() == s {
		return nil
	}
	return ErrNotInWordlist
}

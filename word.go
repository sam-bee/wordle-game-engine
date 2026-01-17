package wordlegameengine

import (
	"fmt"
	"sort"
)

const WordLength = 5

type Word [WordLength]byte

func errInvalidLength(s string) error {
	return fmt.Errorf("%q must be %d letters", s, WordLength)
}

func errInvalidCharacter(s string) error {
	return fmt.Errorf("%q must contain only lowercase a-z", s)
}

func errNotInWordlist(s string) error {
	return fmt.Errorf("%q not in allowed guesses", s)
}

func NewWord(s string) (Word, error) {
	var w Word
	if err := parseWord(s, w[:]); err != nil {
		return Word{}, err
	}
	return w, nil
}

func (w Word) String() string {
	return string(w[:])
}

func (w *Word) Validate() error {
	s := w.String()
	if err := validateCharacters(s); err != nil {
		return err
	}
	if !isInWordlist(s, AllowedGuesses) {
		return errNotInWordlist(s)
	}
	return nil
}

func parseWord(s string, dest []byte) error {
	if len(s) != WordLength {
		return errInvalidLength(s)
	}
	for i := 0; i < len(s); i++ {
		if s[i] < 'a' || s[i] > 'z' {
			return errInvalidCharacter(s)
		}
	}
	copy(dest, s)
	return nil
}

func validateCharacters(s string) error {
	for i := 0; i < len(s); i++ {
		if s[i] < 'a' || s[i] > 'z' {
			return errInvalidCharacter(s)
		}
	}
	return nil
}

func isInWordlist(s string, wordlist []Word) bool {
	idx := sort.Search(len(wordlist), func(i int) bool {
		return wordlist[i].String() >= s
	})
	return idx < len(wordlist) && wordlist[idx].String() == s
}

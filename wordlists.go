package wordlegameengine

import (
	"bufio"
	"fmt"
	"os"
)

var AllowedGuesses []Word
var AllowedSolutions []Word

func LoadWordlists(dataDir string) error {
	var err error

	AllowedGuesses, err = loadWordlist(dataDir + "/allowed-guesses.txt")
	if err != nil {
		return err
	}

	AllowedSolutions, err = loadWordlist(dataDir + "/allowed-solutions.txt")
	if err != nil {
		return err
	}

	return nil
}

func loadWordlist(filepath string) ([]Word, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []Word
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			word, err := NewWord(line)
			if err != nil {
				return nil, fmt.Errorf("invalid word %q in %s: %w", line, filepath, err)
			}
			words = append(words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

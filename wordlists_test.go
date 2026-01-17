package wordlegameengine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWordlists(t *testing.T) {
	// Create temp directory with test files
	tmpDir := t.TempDir()

	t.Run("successful load", func(t *testing.T) {
		// Reset global state
		AllowedGuesses = nil
		AllowedSolutions = nil

		guessesFile := filepath.Join(tmpDir, "allowed-guesses.txt")
		solutionsFile := filepath.Join(tmpDir, "allowed-solutions.txt")

		err := os.WriteFile(guessesFile, []byte("apple\nberry\ncrane\n"), 0644)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(solutionsFile, []byte("delta\neager\n"), 0644)
		if err != nil {
			t.Fatal(err)
		}

		err = LoadWordlists(tmpDir)
		if err != nil {
			t.Errorf("LoadWordlists() error = %v, want nil", err)
		}

		if len(AllowedGuesses) != 3 {
			t.Errorf("len(AllowedGuesses) = %d, want 3", len(AllowedGuesses))
		}
		if len(AllowedSolutions) != 2 {
			t.Errorf("len(AllowedSolutions) = %d, want 2", len(AllowedSolutions))
		}

		if AllowedGuesses[0].String() != "apple" {
			t.Errorf("AllowedGuesses[0] = %q, want %q", AllowedGuesses[0].String(), "apple")
		}
		if AllowedSolutions[1].String() != "eager" {
			t.Errorf("AllowedSolutions[1] = %q, want %q", AllowedSolutions[1].String(), "eager")
		}
	})

	t.Run("missing guesses file", func(t *testing.T) {
		AllowedGuesses = nil
		AllowedSolutions = nil

		emptyDir := t.TempDir()
		err := LoadWordlists(emptyDir)
		if err == nil {
			t.Error("LoadWordlists() error = nil, want error for missing file")
		}
	})

	t.Run("missing solutions file", func(t *testing.T) {
		AllowedGuesses = nil
		AllowedSolutions = nil

		partialDir := t.TempDir()
		guessesFile := filepath.Join(partialDir, "allowed-guesses.txt")
		err := os.WriteFile(guessesFile, []byte("apple\n"), 0644)
		if err != nil {
			t.Fatal(err)
		}

		err = LoadWordlists(partialDir)
		if err == nil {
			t.Error("LoadWordlists() error = nil, want error for missing solutions file")
		}
	})

	t.Run("invalid word in guesses", func(t *testing.T) {
		AllowedGuesses = nil
		AllowedSolutions = nil

		invalidDir := t.TempDir()
		guessesFile := filepath.Join(invalidDir, "allowed-guesses.txt")
		solutionsFile := filepath.Join(invalidDir, "allowed-solutions.txt")

		err := os.WriteFile(guessesFile, []byte("apple\nTOOLONG\n"), 0644)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(solutionsFile, []byte("delta\n"), 0644)
		if err != nil {
			t.Fatal(err)
		}

		err = LoadWordlists(invalidDir)
		if err == nil {
			t.Error("LoadWordlists() error = nil, want error for invalid word")
		}
	})

	t.Run("empty files", func(t *testing.T) {
		AllowedGuesses = nil
		AllowedSolutions = nil

		emptyFilesDir := t.TempDir()
		guessesFile := filepath.Join(emptyFilesDir, "allowed-guesses.txt")
		solutionsFile := filepath.Join(emptyFilesDir, "allowed-solutions.txt")

		err := os.WriteFile(guessesFile, []byte(""), 0644)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(solutionsFile, []byte(""), 0644)
		if err != nil {
			t.Fatal(err)
		}

		err = LoadWordlists(emptyFilesDir)
		if err != nil {
			t.Errorf("LoadWordlists() error = %v, want nil for empty files", err)
		}

		if len(AllowedGuesses) != 0 {
			t.Errorf("len(AllowedGuesses) = %d, want 0", len(AllowedGuesses))
		}
		if len(AllowedSolutions) != 0 {
			t.Errorf("len(AllowedSolutions) = %d, want 0", len(AllowedSolutions))
		}
	})
}

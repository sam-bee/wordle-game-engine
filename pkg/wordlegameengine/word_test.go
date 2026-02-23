package wordlegameengine

import (
	"strings"
	"testing"
)

func TestNewWord(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		errContains string
	}{
		{"valid word", "hello", false, ""},
		{"valid word all a", "aaaaa", false, ""},
		{"valid word all z", "zzzzz", false, ""},
		{"too short", "hell", true, "must be 5 letters"},
		{"too long", "helloo", true, "must be 5 letters"},
		{"empty string", "", true, "must be 5 letters"},
		{"uppercase letter", "Hello", true, "must contain only lowercase a-z"},
		{"contains number", "hell0", true, "must contain only lowercase a-z"},
		{"contains space", "hell ", true, "must contain only lowercase a-z"},
		{"contains hyphen", "he-lo", true, "must contain only lowercase a-z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			word, err := NewWord(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWord(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err != nil {
				if !strings.Contains(err.Error(), tt.input) {
					t.Errorf("NewWord(%q) error = %v, want error containing input", tt.input, err)
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("NewWord(%q) error = %v, want error containing %q", tt.input, err, tt.errContains)
				}
			}
			if err == nil && word.String() != tt.input {
				t.Errorf("NewWord(%q).String() = %q, want %q", tt.input, word.String(), tt.input)
			}
		})
	}
}

func TestWord_String(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"hello", "hello"},
		{"world", "world"},
		{"aaaaa", "aaaaa"},
		{"zzzzz", "zzzzz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			word, _ := NewWord(tt.input)
			if got := word.String(); got != tt.input {
				t.Errorf("Word.String() = %q, want %q", got, tt.input)
			}
		})
	}
}

func TestWord_Validate(t *testing.T) {
	// Set up test wordlist
	oldGuesses := AllowedGuesses
	defer func() { AllowedGuesses = oldGuesses }()

	AllowedGuesses = []Word{}
	for _, s := range []string{"apple", "berry", "crane", "delta", "eager"} {
		w, _ := NewWord(s)
		AllowedGuesses = append(AllowedGuesses, w)
	}

	tests := []struct {
		name        string
		word        Word
		wantErr     bool
		errContains string
	}{
		{"valid word in list (first)", mustNewWord("apple"), false, ""},
		{"valid word in list (middle)", mustNewWord("crane"), false, ""},
		{"valid word in list (last)", mustNewWord("eager"), false, ""},
		{"valid word not in list", mustNewWord("zebra"), true, "not in allowed guesses"},
		{"invalid uppercase", Word{'H', 'e', 'l', 'l', 'o'}, true, "must contain only lowercase a-z"},
		{"invalid null byte", Word{0, 'e', 'l', 'l', 'o'}, true, `"\x00ello" must contain only lowercase a-z`},
		{"invalid number", Word{'h', 'e', 'l', 'l', '0'}, true, "must contain only lowercase a-z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.word.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Word.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Word.Validate() error = %v, want error containing %q", err, tt.errContains)
			}
		})
	}
}

func mustNewWord(s string) Word {
	w, err := NewWord(s)
	if err != nil {
		panic(err)
	}
	return w
}

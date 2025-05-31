package service

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTitleFromString(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		maxWords   int
		maxLen     int
		wantExact  string
		wantRegexp string
	}{
		{
			name:       "EmptyInput",
			input:      "",
			maxWords:   5,
			maxLen:     10,
			wantRegexp: `^untitled-[[:xdigit:]]{8}$`,
		},
		{
			name:      "ShortInputWithinLimits",
			input:     "hello world",
			maxWords:  5,
			maxLen:    20,
			wantExact: "hello world",
		},
		{
			name:      "ExceedsMaxWords",
			input:     "one two three",
			maxWords:  2,
			maxLen:    100,
			wantExact: "one two",
		},
		{
			name:      "ExceedsMaxLenByCharacter",
			input:     "a b c",
			maxWords:  10,
			maxLen:    2,
			wantExact: "a",
		},
		{
			name:       "SingleGiantWord",
			input:      "abcdefghijk",
			maxWords:   5,
			maxLen:     5,
			wantRegexp: `^abcde-[[:xdigit:]]{8}$`,
		},
		{
			name:      "FirstWordExactMaxLen",
			input:     "abcde fgh",
			maxWords:  5,
			maxLen:    5,
			wantExact: "abcde",
		},
		{
			name:      "WordsAndSpacesTrimmed",
			input:     "   leading and trailing   ",
			maxWords:  3,
			maxLen:    20,
			wantExact: "leading and trailing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := TitleFromString(tt.input, tt.maxWords, tt.maxLen)

			if tt.wantExact != "" {
				require.Equal(t, tt.wantExact, out)
			}
			if tt.wantRegexp != "" {
				require.Regexp(t, regexp.MustCompile(tt.wantRegexp), out)
			}
		})
	}
}

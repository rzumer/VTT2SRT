package parsing

import (
	"os"
	"bufio"
	"strings"
	"github.com/rzumer/vtt2srt/util"
)

type Parser struct {
	input *bufio.Reader
}

func NewParser(inputPath string) (*Parser, error) {
	inputFile, err := os.Open(inputPath)
	
	return &Parser {
		input: bufio.NewReader(inputFile),
	}, err
}

// Parses the first line of the file to ensure that it is a valid WebVTT file.
func (parser *Parser) Valid() bool {
	line, err := parser.readLine()
	
	if err != nil {
		return false
	}
	
	if len(line) < 6 {
		return false
	}
	
	if !strings.HasPrefix(line, "WEBVTT") {
		return false
	}
	
	if len(line) > 6 {
		separators := []byte { ' ', '\t', '\n' }
		
		// The content is valid if WEBVTT is followed by one of the separators allowed.
		if !util.Contains(separators, line[6]) {
			return false
		}
	}
	
	return true
}

func (parser *Parser) ParseHeader() []string {
	return nil
}

func (parser *Parser) ParseCue() []string {
	return nil
}

func (parser *Parser) readLine() (string, error) {
	line, err := parser.input.ReadString('\n')
	return replaceInvalidCharacters(line), err 
}

func replaceInvalidCharacters(input string) string {
	output := input
	
	// Replace all U+0000 NULL characters by U+FFFD REPLACEMENT CHARACTERs.
	output = strings.Replace(output, "\u0000", "\uFFFD", -1)
	
	// Replace each U+000D CARRIAGE RETURN U+000A LINE FEED (CRLF) character pair
	// by a single U+000A LINE FEED (LF) character.
	output = strings.Replace(output, "\u000D\u000A", "\u000A", -1)
	
	// Replace all remaining U+000D CARRIAGE RETURN characters by U+000A LINE FEED (LF) characters.
	output = strings.Replace(output, "\u000D", "\u000A", -1)
	
	return output
}

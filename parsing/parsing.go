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
	line, _ := parser.readLine()
	
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

func (parser *Parser) ParseAllCues() []string {
	output := make([]string, 0)
	
	// Skip the header.
	parser.ParseHeader()
	
	// Parse each cue until the end of the file.
	for {
		
		// Skip empty lines between blocks.
		for char, err := parser.input.Peek(1); err == nil && char[0] == '\n'; {
			parser.readLine()
		}
		
		// Parse the next cue.
		cue, _ := parser.ParseCue()
		output = append(output, cue...)
		
		// Stop when the reader reaches the end of the file.
		// Otherwise, add an empty line to separate cues.
		if peeked, _ := parser.input.Peek(1); len(peeked) == 0 {
			break
		} else {
			output = append(output, "")
		}
	}
	
	return output
}

func (parser *Parser) ParseHeader() ([]string, error) {
	if char, _ := parser.input.Peek(1); char[0] == '\n' {
		return nil, nil
	}
	
	return parser.collectBlock(true)
}

func (parser *Parser) ParseCue() ([]string, error) {
	return parser.collectBlock(false)
}

func (parser *Parser) collectBlock(inHeader bool) ([]string, error) {
	var err error
	cueText := make([]string, 0)
	lineCount := 0
	seenEOF := false
	seenArrow := false
	
	for {
		line, err := parser.readLine()
		line = strings.TrimSpace(line)
		lineCount++
		
		if err != nil {
			seenEOF = true
		}
		
		if strings.Contains(line, "-->") {
			if !inHeader && (lineCount == 1 || (lineCount == 2 && !seenArrow)) {
				seenArrow = true
				cueText = append(cueText, line)
			} else {
				break
			}
		} else if line == "" {
			break
		} else {
			if !inHeader && lineCount == 2 {
				// Process stylesheet or region (currently unsupported).
				if strings.HasPrefix(line, "STYLE") {
					break
				} else if strings.HasPrefix(line, "REGION") {
					break
				}
			}
			
			// Append a line of content.
			cueText = append(cueText, line)
		}
		
		if seenEOF {
			break
		}
	}
	
	if len(cueText) > 0 {
		return cueText, err
	}
	
	return nil, err
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

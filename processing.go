package main

import(
	"os"
	"path/filepath"
	"bufio"
	"errors"
	"strconv"
	"strings"
	"regexp"
	"fmt"
)

func validate(inputPath string, outputPath string) (bool, error) {
	// Ensure that the input file and the output path exist.
	_, inError := os.Stat(inputPath)
	_, outError := os.Stat(filepath.Dir(outputPath))
	
	if os.IsNotExist(inError) {
		return false, errors.New("input file not found")
	}
	
	if os.IsNotExist(outError) {
		return false, errors.New("output path not found")
	}
	
	if inError != nil || outError != nil {
		return false, errors.New("unknown I/O error")
	}

	// Ensure that the input file is readable.
	inputFile, readError := os.Open(inputPath)
	
	if readError != nil {
		return false, errors.New("input file access error")
	}
	
	defer inputFile.Close()
	
	// Ensure that the first line of the file is "WEBVTT".
	scanner := bufio.NewScanner(inputFile)
	
	if !scanner.Scan() || scanner.Text() != "WEBVTT" {
		return false, errors.New("input file is not a valid VTT file")
	}
	
	return true, nil
}

func convert(inputPath string) []string {
	inputFile, _ := os.Open(inputPath)
	
	defer inputFile.Close()
	
	scanner := bufio.NewScanner(inputFile)
	counter := 1
	output := make([]string, 0)
	
	for scanner.Scan() {
		text := scanner.Text()
		
		// Skip header content and consecutive empty lines.
		if(text == "WEBVTT" || text == "") {
			continue
		}
		
		// Output the subtitle number.
		output = append(output, strconv.Itoa(counter))
		
		// Split the timestamps expected to be on the line.
		timestamps := strings.Split(text, " --> ")
		
		// If two timestamps are not found on the line, keep searching.
		if len(timestamps) < 2 {
			continue
		}
		
		// Convert the timestamps to the hh:mm:ss,zzz SRT format.
		convertedLine := convertTimestamp(timestamps[0]) + " --> " + convertTimestamp(timestamps[1])
		
		output = append(output, convertedLine)
		
		// Output each line for the current subtitle.
		for scanner.Scan() {
			text = scanner.Text()
			output = append(output, text)
			
			if(text == "") {
				break;
			}
		}
		
		counter++
	}
	
	return output
}

func save(output []string, outputPath string) {
	outputFile, _ := os.Create(outputPath)
	
	writer := bufio.NewWriter(outputFile)
	
	for _, outputLine := range output {
		writer.WriteString(outputLine + "\n")
	}
	
	writer.Flush()
}

func convertTimestamp(timestamp string) string {
	vttTimestampRegexp := `(?:(\d{2}):)?(\d{2}):(\d{2}).(\d{3})`
	matcher, _ := regexp.Compile(vttTimestampRegexp)
	
	submatches := matcher.FindAllStringSubmatch(timestamp, -1)
	
	if submatches == nil || len(submatches) == 0 || len(submatches[0]) < 5 {
		return timestamp
	}
	
	return fmt.Sprintf("%02s:%02s:%02s,%03s", submatches[0][1], submatches[0][2], submatches[0][3], submatches[0][4])
}

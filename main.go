package main

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"
)

func main() {
	// Retrieve command line arguments.
	args := os.Args[1:]
	
	// If not enough arguments are provided, display usage information.
	if len(args) < 1 || len(args) > 2 {
		fmt.Printf("Usage: VTT2SRT input [output]\n" +
			"\n" +
			"input\tThe input file path.\n" +
			"output\tThe output file path.\n")
		
		return
	}
	
	// Set the input path received as argument.
	inputPath := args[0]
	
	// Set the default output path.
	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".srt"
	
	// If the output path is received as a parameter, set it to the argument.
	if len(args) == 2 {
		outputPath = args[1]
	}
	
	// Ensure that the file paths do not collide.
	if inputPath == outputPath {
		outputExt := filepath.Ext(outputPath)
		outputPath = strings.TrimSuffix(outputPath, outputExt) + "_out" + outputExt
	}
	
	// Ensure that the input and output paths are valid.
	valid, err := validate(inputPath, outputPath)
	if !valid {
		fmt.Printf("Validation error: " + err.Error() + ".\n")
		return
	}

	// Process the file and save its result to the output file.
	save(convert(parseInput(inputPath)), outputPath)
	fmt.Printf("Done.\n")
	return
}

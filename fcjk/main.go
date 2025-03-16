// main.go
package main

import (
	"bufio"
	"fcjk/formatter"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func main() {
	// Define command line flags
	inputFile := flag.String("i", "", "Path to input file")
	flag.Parse()

	if *inputFile == "" {
		// Handle stdin input mode
		handleStdinInput()
	} else {
		// Handle file input mode
		handleFileInput(*inputFile)
	}
}

// main.go - handleStdinInput 函数修改
func handleStdinInput() {
	// 检测操作系统并提供适当的提示
	eofHint := "Ctrl+D"
	if runtime.GOOS == "windows" {
		eofHint = "Ctrl+Z followed by Enter"
	}

	fmt.Printf("Enter text to format (Press %s to finish):\n", eofHint)

	scanner := bufio.NewScanner(os.Stdin)
	var inputText strings.Builder

	for scanner.Scan() {
		inputText.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	// Format the text
	formattedText := formatter.FormatText(inputText.String())

	// Output to stdout
	fmt.Println("\nFormatted text:")
	fmt.Println(formattedText)
}

func handleFileInput(inputFile string) {
	// Check if file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: File '%s' does not exist\n", inputFile)
		os.Exit(1)
	}

	// Read input file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Format the text
	formattedText := formatter.FormatText(string(content))

	// Create output file
	outputFile := strings.TrimSuffix(inputFile, ".txt") + "_output.txt"
	if !strings.HasSuffix(inputFile, ".txt") {
		outputFile = inputFile + "_output.txt"
	}

	// Write to output file
	err = os.WriteFile(outputFile, []byte(formattedText), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Formatted text has been saved to '%s'\n", outputFile)
}

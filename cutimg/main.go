package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Define command line flags
	dirPath := flag.String("d", "", "Directory path containing PNG files to process")
	flag.Parse()

	// Check if directory path is provided
	if *dirPath == "" {
		fmt.Println("Error: Please provide a directory path using the -d flag")
		flag.Usage()
		os.Exit(1)
	}

	// Check if the directory exists
	dirInfo, err := os.Stat(*dirPath)
	if err != nil {
		fmt.Printf("Error: Could not access directory %s: %v\n", *dirPath, err)
		os.Exit(1)
	}

	if !dirInfo.IsDir() {
		fmt.Printf("Error: %s is not a directory\n", *dirPath)
		os.Exit(1)
	}

	// Walk through the directory
	err = filepath.Walk(*dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Process only PNG files
		if strings.ToLower(filepath.Ext(path)) == ".png" {
			err := processPNG(path)
			if err != nil {
				fmt.Printf("Error processing %s: %v\n", path, err)
			} else {
				fmt.Printf("Successfully processed: %s\n", path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking through directory: %v\n", err)
		os.Exit(1)
	}
}

func processPNG(filePath string) error {
	baseDir := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)
	fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	outputPath := filepath.Join(baseDir, fileNameWithoutExt+"_modify.png")

	// Construct the ffmpeg command for center-crop to 3:4 aspect ratio (portrait)
	// We keep the height and adjust width accordingly to get 3:4 aspect ratio
	cmd := exec.Command(
		"ffmpeg",
		"-i", filePath,
		"-vf", "crop=ih*3/4:ih:iw/2-ih*3/8:0", // width=height*3/4, center horizontally
		"-c:v", "png", // PNG codec
		"-pred", "mixed", // Use mixed prediction for lossless output
		"-lossless", "1", // Ensure lossless output
		outputPath,
	)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v\nOutput: %s", err, string(output))
	}

	return nil
}

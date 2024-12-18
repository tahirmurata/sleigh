package elf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GetInputFile retrieves and opens the input file for the specified year and day.
// If the input file doesn't exist locally, it downloads it from adventofcode.com.
func (e *Elf) GetInputFile(year int, day int) (*os.File, error) {
	filename := fmt.Sprintf("input.txt")
	file, err := os.Open(filename)
	if err == nil {
		return file, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		// handle the case where the file doesn't exist

		// Create file
		fileCreated, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		defer fileCreated.Close()

		// Prepare request
		url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		request.AddCookie(&http.Cookie{Name: "session", Value: e.token})

		// Send request
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Write response to file
		_, err = io.Copy(fileCreated, resp.Body)
		if err != nil {
			return nil, err
		}

		// Open file
		file, err = os.Open(filename)
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	return nil, err
}

// LinesFromFile reads lines from a file and returns them as a slice of strings.
func LinesFromFile(file *os.File) ([]string, error) {
	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return lines, nil
}

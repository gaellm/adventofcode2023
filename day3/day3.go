package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"unicode"
)

// read each line of file and add to a slice
func readLines(filename string) ([]string, error) {
	lines := make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		return lines, errors.New("fail to read lines from " + filename + " due to error " + err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

// findEngineSymbols takes a slice of strings (lines) and creates a map of maps
// to store non-alphanumeric characters found in the lines.
// The outer map uses line numbers as keys, and the inner map uses character indexes
// as keys, storing the non-alphanumeric characters found in the respective positions.
func findEngineSymbols(lines []string) map[int]map[int]string {

	matrix := make(map[int]map[int]string)

	for lineNb, line := range lines {
		for charIndex, char := range line {
			if !unicode.IsLetter(char) && !unicode.IsNumber(char) && string(char) != "." {
				_, ok := matrix[lineNb]
				if ok {
					matrix[lineNb][charIndex] = string(char)
				} else {
					matrix[lineNb] = make(map[int]string)
					matrix[lineNb][charIndex] = string(char)
				}
			}
		}
	}
	return matrix
}

func main() {
	fmt.Println("Hello _")
}

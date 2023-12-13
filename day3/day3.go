package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

type engineNumber struct {
	numberStr    string
	number       int
	startIndice  int
	endIndice    int
	isPartNumber bool
	line         int
}

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
			if !unicode.IsNumber(char) && string(char) != "." {
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

// findEngineNumbers takes a slice of strings (lines) and identifies engine numbers
// within those lines. It returns a map where the key is the line number, and the
// value is a slice of engineNumber structures representing the engine numbers found
// in that line.
func findEngineNumbers(lines []string) ([]engineNumber, error) {

	engineNumbers := []engineNumber{}

	for lineNb, line := range lines {
		var currentNumber engineNumber

		for i, char := range line {
			if unicode.IsDigit(char) {
				if currentNumber.numberStr == "" {
					currentNumber.line = lineNb
					currentNumber.startIndice = i
					currentNumber.endIndice = i
					currentNumber.numberStr = string(char)
				} else {
					currentNumber.numberStr += string(char)
					currentNumber.endIndice = i
				}
			} else {
				if currentNumber.numberStr != "" {

					nb, err := strconv.Atoi(currentNumber.numberStr)
					if err != nil {
						return nil, errors.New("failed to convert count to integer: " + err.Error())
					}
					currentNumber.number = nb

					engineNumbers = append(engineNumbers, currentNumber)
					currentNumber = engineNumber{}
				}
			}
		}

		//process end of line
		if currentNumber.numberStr != "" {
			nb, err := strconv.Atoi(currentNumber.numberStr)
			if err != nil {
				return nil, errors.New("failed to convert count to integer: " + err.Error())
			}
			currentNumber.number = nb

			engineNumbers = append(engineNumbers, currentNumber)
			currentNumber = engineNumber{}
		}
	}
	return engineNumbers, nil
}

// function to check if a number is in an interval
func isNumberInInterval(nb int, interval []int) bool {
	if nb >= interval[0] && nb <= interval[1] {
		return true
	}
	return false
}

func isPartNumber(eNbr engineNumber, symbolsMatrix map[int]map[int]string) bool {

	//check lines
	possibleSymbolLines := []int{eNbr.line - 1, eNbr.line, eNbr.line + 1}
	var effectiveSymbolsLines = []int{}
	for _, line := range possibleSymbolLines {
		if _, ok := symbolsMatrix[line]; ok {
			effectiveSymbolsLines = append(effectiveSymbolsLines, line)
		}
	}
	if len(effectiveSymbolsLines) == 0 {
		return false
	}

	//check indices
	interval := []int{eNbr.startIndice - 1, eNbr.endIndice + 1}
	for _, line := range effectiveSymbolsLines {
		for indice := range symbolsMatrix[line] {
			if isNumberInInterval(indice, interval) {
				return true
			}
		}
	}

	return false
}

func getEnginPartNumbers(lines []string) ([]int, error) {
	var partNumbers []int
	symbols := findEngineSymbols(lines)
	numbers, err := findEngineNumbers(lines)
	if err != nil {
		return []int{}, err
	}

	for _, nb := range numbers {
		if isPartNumber(nb, symbols) {
			partNumbers = append(partNumbers, nb.number)
		}

	}

	return partNumbers, nil
}

// iterate an ints slice and sum it
func sumInts(ints []int) int {
	sum := 0
	for _, i := range ints {
		sum += i
	}
	return sum
}

func main() {

	//read input file
	lines, err := readLines("./input.txt")
	if err != nil {
		panic(err.Error())
	}

	partNumbers, err := getEnginPartNumbers(lines)
	if err != nil {
		panic(err.Error())
	}

	//fmt.Println("Part numbers: ", partNumbers)
	fmt.Println("Sum of part numbers: ", sumInts(partNumbers))

}

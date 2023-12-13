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

type gear struct {
	numbers []int
	ratio   int
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

func getEnginPartNumbers(lines []string, symbols map[int]map[int]string, numbers []engineNumber) ([]int, error) {
	var partNumbers []int

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

func filterPositiveNumbers(arr []int) []int {
	var result []int

	for _, num := range arr {
		if num >= 0 {
			result = append(result, num)
		}
	}

	return result
}

func isIntInArray(target int, arr []int) bool {
	for _, num := range arr {
		if num == target {
			return true
		}
	}
	return false
}

func findEngineNumbersInInterval(numbers []engineNumber, possibleLines []int, possiblesIndiceInterval []int) []int {

	possibleLines = filterPositiveNumbers(possibleLines)
	var result []int

	for _, number := range numbers {

		if isIntInArray(number.line, possibleLines) {

			if isNumberInInterval(number.startIndice, possiblesIndiceInterval) {
				result = append(result, number.number)
				continue
			}

			if isNumberInInterval(number.endIndice, possiblesIndiceInterval) {
				result = append(result, number.number)
			}
		}
	}
	return result
}

func getGears(symbolsMatrix map[int]map[int]string, engineNumbers []engineNumber) []gear {

	var gears []gear

	// Iterate through the outer map
	for symbolLine, indiceMap := range symbolsMatrix {
		possibleLines := []int{symbolLine - 1, symbolLine, symbolLine + 1}
		// Iterate through the inner map
		for indice, symbol := range indiceMap {
			if symbol != "*" {
				continue
			}
			possiblesIndiceInterval := []int{indice - 1, indice + 1}
			allNumbersInInterval := findEngineNumbersInInterval(engineNumbers, possibleLines, possiblesIndiceInterval)
			if len(allNumbersInInterval) == 2 {
				gears = append(gears, gear{
					numbers: allNumbersInInterval,
					ratio:   allNumbersInInterval[0] * allNumbersInInterval[1],
				})
			}
		}
	}

	return gears
}

func main() {

	//read input file
	lines, err := readLines("./input.txt")
	if err != nil {
		panic(err.Error())
	}

	symbols := findEngineSymbols(lines)
	numbers, err := findEngineNumbers(lines)
	if err != nil {
		panic(err.Error())
	}

	partNumbers, err := getEnginPartNumbers(lines, symbols, numbers)
	if err != nil {
		panic(err.Error())
	}

	gears := getGears(symbols, numbers)
	var sumGearRatio int
	for _, gear := range gears {
		sumGearRatio += gear.ratio
	}

	fmt.Println("Part 1 - sum of part numbers: ", sumInts(partNumbers))
	fmt.Println("Part 2 - sum of gears ratio: ", sumGearRatio)
}

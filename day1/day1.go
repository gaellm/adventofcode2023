package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var numberMap = map[string]string{
	"one":   "1",
	"two":   "2",
	"three": "3",
	"four":  "4",
	"five":  "5",
	"six":   "6",
	"seven": "7",
	"eight": "8",
	"nine":  "9",
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

func reverseStr(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// replaceRegularSpelledOutNumber replaces spelled-out numbers in a given string
// with their corresponding digits using a predefined map. From left to right.
// For example, if the input string contains "eightwo", the function will replace
// it with "8wo" based on the predefined numberMap.
func replaceRegularSpelledOutNumber(str string) string {

	// Regular expression to match all number words in the map
	var numberWords []string
	for word := range numberMap {
		numberWords = append(numberWords, word)
	}

	// Replace function that uses the map to replace words with digits
	replaceFunc := func(match string) string {

		lowerMatch := strings.ToLower(match) // Ensure case-insensitivity
		if digit, exists := numberMap[lowerMatch]; exists {
			return digit
		}
		return match
	}

	pattern := `(` + strings.Join(numberWords, "|") + `)`
	re := regexp.MustCompile(pattern)

	return re.ReplaceAllStringFunc(str, replaceFunc)
}

// replaceReverseSpelledOutNumber replaces spelled-out numbers in a given string
// with their corresponding digits using a predefined map. From right to left.
// For example, if the input string contains "eightwo", the function will replace
// it with "eigh2" based on the predefined numberMap.
func replaceReverseSpelledOutNumber(str string) string {

	// Regular expression to match all number words in the map
	var numberWords []string
	for word := range numberMap {
		numberWords = append(numberWords, reverseStr(word))
	}

	// Replace function that uses the map to replace words with digits
	replaceFunc := func(match string) string {

		lowerMatch := strings.ToLower(match) // Ensure case-insensitivity
		if digit, exists := numberMap[reverseStr(lowerMatch)]; exists {
			return digit
		}
		return match
	}

	pattern := `(` + strings.Join(numberWords, "|") + `)`
	re := regexp.MustCompile(pattern)

	return reverseStr(re.ReplaceAllStringFunc(reverseStr(str), replaceFunc))
}

// get first string digit
func findFirstDigit(input string) string {
	for _, char := range input {
		if unicode.IsDigit(char) {
			return string(char)
		}
	}
	return ""
}

// keep only the first and the last digit character of string. Taking in account that
// some of the digits are actually spelled out with letters
func keepFirstAndLast(line string) string {

	first := findFirstDigit(replaceRegularSpelledOutNumber(line))
	last := findFirstDigit(reverseStr(replaceReverseSpelledOutNumber(line)))

	return first + last
}

// iterate an ints slice and sum it
func sumInts(ints []int) int {
	sum := 0
	for _, i := range ints {
		sum += i
	}
	return sum
}

func lines2Ints(lines []string) ([]int, error) {

	var ints []int

	for _, line := range lines {

		//get first and last number
		numberStr := keepFirstAndLast(line)
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return ints, errors.New("fail to transform keeped digits " + numberStr + " to integer with error " + err.Error())
		}
		ints = append(ints, number)
	}

	return ints, nil

}

func main() {

	//read input file
	lines, err := readLines("./input.txt")
	if err != nil {
		panic(err.Error())
	}

	//get the calibration value
	ints, err := lines2Ints(lines)
	if err != nil {
		panic(err.Error())
	}

	//sum all the calibration values
	fmt.Println(sumInts(ints))
}

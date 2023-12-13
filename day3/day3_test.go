package main

import (
	"os"
	"reflect"
	"testing"
)

func TestReadLines(t *testing.T) {
	// Create a temporary file with sample content
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	content := "Line 1\nLine 2\nLine 3"
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}

	// Close the file before testing
	tmpfile.Close()

	// Test reading lines from the temporary file
	lines, _ := readLines(tmpfile.Name())

	// Verify the expected content
	expected := []string{"Line 1", "Line 2", "Line 3"}
	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Expected '%s', but got '%s'", expected[i], line)
		}
	}

	// Test reading from a non-existent file
	_, err = readLines("nonexistentfile.txt")
	if err == nil {
		t.Error("Expected an error for a non-existent file, but got none")
	}
}

func TestFindEngineSymbols(t *testing.T) {
	// Example input lines
	lines := []string{
		".....#123,",
		".........!",
	}

	// Call the function to find engine symbols
	symbolMatrix := findEngineSymbols(lines)

	// Expected result based on the example input
	expectedMatrix := map[int]map[int]string{
		0: {5: "#", 9: ","},
		1: {9: "!"},
	}

	// Check if the result matches the expected result
	if !reflect.DeepEqual(symbolMatrix, expectedMatrix) {
		t.Errorf("Result does not match the expected matrix. Got: %v, Expected: %v", symbolMatrix, expectedMatrix)
	}
}

func TestFindEngineNumbers(t *testing.T) {
	// Example input lines
	lines := []string{
		"467..114..",
		"...*......",
		"..35..633.",
		".......982",
		".......9..",
	}

	// Call the function to find engine numbers
	result, _ := findEngineNumbers(lines)

	// Expected result based on the example input
	expectedResult := []engineNumber{
		{line: 0, numberStr: "467", startIndice: 0, endIndice: 2, number: 467, isPartNumber: false}, {line: 0, numberStr: "114", startIndice: 5, endIndice: 7, number: 114, isPartNumber: false},
		{line: 2, numberStr: "35", startIndice: 2, endIndice: 3, number: 35, isPartNumber: false}, {line: 2, numberStr: "633", startIndice: 6, endIndice: 8, number: 633, isPartNumber: false},
		{line: 3, numberStr: "982", startIndice: 7, endIndice: 9, number: 982, isPartNumber: false}, {line: 4, numberStr: "9", startIndice: 7, endIndice: 7, number: 9, isPartNumber: false}}

	// Check if the result matches the expected result
	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Result does not match the expected result. Got: %v, Expected: %v", result, expectedResult)
	}
}

func TestIsNumberInInterval(t *testing.T) {
	// Test cases
	testCases := []struct {
		number   int
		interval []int
		expected bool
	}{
		{5, []int{0, 10}, true},
		{15, []int{0, 10}, false},
		{0, []int{0, 10}, true},
		{10, []int{0, 10}, true},
		{-5, []int{-10, 5}, true},
		{-15, []int{-10, 5}, false},
	}

	// Iterate through test cases
	for _, testCase := range testCases {
		result := isNumberInInterval(testCase.number, testCase.interval)

		// Check if the result matches the expected result
		if result != testCase.expected {
			t.Errorf("For number %d and interval %v, expected %t, but got %t", testCase.number, testCase.interval, testCase.expected, result)
		}
	}
}

func TestIsPartNumber(t *testing.T) {

	symbolMatrix := map[int]map[int]string{
		1: {3: "*"},
		3: {6: "#"},
		4: {3: "*"},
	}
	// Test cases
	testCases := []struct {
		engineNumber   engineNumber
		symbolsMatrix  map[int]map[int]string
		expectedResult bool
	}{
		{
			engineNumber{"467", 467, 0, 2, false, 0},
			symbolMatrix,
			true,
		},
		{
			engineNumber{"114", 114, 5, 7, false, 0},
			symbolMatrix,
			false,
		},
		{
			engineNumber{"633", 633, 6, 8, false, 2},
			symbolMatrix,
			true,
		},
	}

	// Iterate through test cases
	for _, testCase := range testCases {
		result := isPartNumber(testCase.engineNumber, testCase.symbolsMatrix)

		// Check if the result matches the expected result
		if result != testCase.expectedResult {
			t.Errorf("For engine number %v and symbols matrix %v, expected %t, but got %t", testCase.engineNumber, testCase.symbolsMatrix, testCase.expectedResult, result)
		}
	}
}

func TestGetEnginPartNumbers(t *testing.T) {
	// Test cases
	testCases := []struct {
		lines          []string
		expectedResult []int
		expectedError  error
	}{
		{
			[]string{
				"467..114..",
				"...*......",
				"..35..633.",
				"......#...",
				"617%......",
				".....+.58.",
				"..592.....",
				"......755.",
				"...$./....",
				".664.598..",
			},
			[]int{467, 35, 633, 617, 592, 755, 664, 598},
			nil,
		},
		{
			[]string{
				".......@...",
				"........982",
				".370.......",
				"...*.......",
			},
			[]int{982, 370},
			nil,
		},
		{
			[]string{
				"...........",
				"........982",
				".370.......",
				"...%.......",
			},
			[]int{370},
			nil,
		},
		{
			[]string{
				"...........",
				"....982....",
				"...........",
			},
			[]int{},
			nil,
		},
		{
			[]string{
				"...........",
				"......9....",
				"+.....)...*",
			},
			[]int{9},
			nil,
		},
	}

	// Iterate through test cases
	for _, testCase := range testCases {
		result, err := getEnginPartNumbers(testCase.lines)

		if len(testCase.expectedResult) > 0 {
			// Check if the result matches the expected result
			if !reflect.DeepEqual(result, testCase.expectedResult) {
				t.Errorf("For lines %v, expected %v, but got %v", testCase.lines, testCase.expectedResult, result)
			}
		} else {
			if len(result) != len(testCase.expectedResult) {
				t.Errorf("For lines %v, expected %v, but got %v", len(testCase.lines), len(testCase.expectedResult), result)
			}
		}

		// Check if the error matches the expected error
		if (err != nil && testCase.expectedError == nil) || (err == nil && testCase.expectedError != nil) {
			t.Errorf("For lines %v, expected error %v, but got %v", testCase.lines, testCase.expectedError, err)
		}
	}
}

func TestSumInts(t *testing.T) {
	testCases := []struct {
		ints           []int
		expectedResult int
	}{
		{[]int{467, 35, 633, 617, 592, 755, 664, 598}, 4361},
	}

	result := sumInts(testCases[0].ints)
	if result != testCases[0].expectedResult {
		t.Errorf("Expected %v, but got %v", testCases[0].expectedResult, result)
	}
}

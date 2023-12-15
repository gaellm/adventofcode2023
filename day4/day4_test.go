package main

import (
	"fmt"
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

func TestParseNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
	}{
		{"41 48 83 86 17", []int{41, 48, 83, 86, 17}},
		{"1 2   3", []int{1, 2, 3}},
		{"10 20 30", []int{10, 20, 30}},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := parseNumbers(test.input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		input    string
		expected card
	}{
		{
			"Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53",
			card{
				cardNumber:  1,
				winningNbrs: []int{41, 48, 83, 86, 17},
				numbers:     []int{83, 86, 6, 31, 17, 9, 48, 53},
			},
		},
		{
			"Card 23: 13 32 20 16 61 | 61 30 68 82 17 32 24 19",
			card{
				cardNumber:  23,
				winningNbrs: []int{13, 32, 20, 16, 61},
				numbers:     []int{61, 30, 68, 82, 17, 32, 24, 19},
			},
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := parseLine(test.input)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestCountNumbersInSlice(t *testing.T) {
	tests := []struct {
		slice1   []int
		slice2   []int
		expected int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{3, 4, 5, 6, 7, 8, 0}, 3},
		{[]int{10, 20, 30}, []int{20, 30, 40, 50}, 2},
		{[]int{1, 2, 3}, []int{4, 5, 6}, 0},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v and %v", test.slice1, test.slice2), func(t *testing.T) {
			result := countNumbersInSlice(test.slice1, test.slice2)

			if result != test.expected {
				t.Errorf("Expected %d, got %d", test.expected, result)
			}
		})
	}
}

func TestDoubleOnEachMatch(t *testing.T) {
	tests := []struct {
		matchNbr int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 8},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("matchNbr=%d", test.matchNbr), func(t *testing.T) {
			result := doubleOnEachMatch(test.matchNbr)

			if result != test.expected {
				t.Errorf("Expected %d, got %d", test.expected, result)
			}
		})
	}
}

func TestProcessCardsPoints(t *testing.T) {
	cards := []card{
		{1, []int{41, 48, 83, 86, 17}, []int{83, 86, 6, 31, 17, 9, 48, 53}, 0, 0},
		{2, []int{10, 20, 30}, []int{20, 30, 40, 50}, 0, 0},
		{3, []int{1, 2, 3}, []int{4, 5, 6}, 0, 0},
	}

	expected := []card{
		{1, []int{41, 48, 83, 86, 17}, []int{83, 86, 6, 31, 17, 9, 48, 53}, 8, 4},
		{2, []int{10, 20, 30}, []int{20, 30, 40, 50}, 2, 2},
		{3, []int{1, 2, 3}, []int{4, 5, 6}, 0, 0},
	}

	result := processCardsPoints(cards)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

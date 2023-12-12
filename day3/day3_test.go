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
		"Hello#123,",
		"...Worlcd!",
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

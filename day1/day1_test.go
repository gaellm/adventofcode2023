package main

import (
	"errors"
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

func TestReplaceRegularSpelledOutNumber(t *testing.T) {

	input := "eightwo"
	expected := "8wo"

	actual := replaceRegularSpelledOutNumber(input)
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

}

func TestReplaceReverseSpelledOutNumber(t *testing.T) {

	input := "eightwo"
	expected := "eigh2"

	actual := replaceReverseSpelledOutNumber(input)
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

}

func TestReverseStr(t *testing.T) {

	input := "eightwo"
	expected := "owthgie"

	actual := reverseStr(input)
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

}

func TestFindFirstDigit(t *testing.T) {

	input := "dfdsf8two"
	expected := "8"

	actual := findFirstDigit(input)
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

}

func TestKeepFirstAndLast(t *testing.T) {

	input := "eightwo"
	expected := "82"

	actual := keepFirstAndLast(input)
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
	}

}

func TestLines2Ints(t *testing.T) {
	testCases := []struct {
		inputLines  []string
		expected    []int
		expectedErr error
	}{
		{[]string{"fdsf1efdsf2fdsf", "456def789two"}, []int{12, 42}, nil},
		{[]string{"invalid", "12invalid34"}, nil, errors.New("fail to transform keeped digits  to integer with error strconv.Atoi: parsing \"\": invalid syntax")},
		{nil, nil, nil},
	}

	for _, testCase := range testCases {
		result, err := lines2Ints(testCase.inputLines)

		if !reflect.DeepEqual(err, testCase.expectedErr) {
			t.Errorf("For input lines %v, expected error %v, but got %v", testCase.inputLines, testCase.expectedErr, err)
		}

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("For input lines %v, expected %v, but got %v", testCase.inputLines, testCase.expected, result)
		}
	}
}

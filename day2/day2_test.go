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

func TestGetGameID(t *testing.T) {
	testCases := []struct {
		input       string
		expectedID  int
		expectedErr bool
	}{
		{"Game 1: 3 blue,", 1, false},
		{"Game 42: 5 red,", 42, false},
		{"No game here", 0, true},
		{"Game abc: 2 green,", 0, true},
		{"Game : 1 yellow,", 0, true},
	}

	for _, testCase := range testCases {
		result, err := getGameID(testCase.input)

		if testCase.expectedErr && err == nil {
			t.Errorf("Expected an error for input: %s", testCase.input)
		}

		if !testCase.expectedErr && err != nil {
			t.Errorf("Unexpected error for input: %s", testCase.input)
		}

		if result != testCase.expectedID {
			t.Errorf("For input '%s', expected %d, but got %d", testCase.input, testCase.expectedID, result)
		}
	}
}

func TestParseGameSets(t *testing.T) {
	testCases := []struct {
		input       string
		expected    []cubeSet
		expectedErr bool
	}{
		{
			"Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green",
			[]cubeSet{
				{blue: 3, red: 4},
				{blue: 6, red: 1, green: 2},
				{green: 2},
			},
			false,
		},
		{
			"Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red",
			[]cubeSet{
				{green: 1, blue: 6, red: 3},
				{red: 6, green: 3},
			},
			false,
		},
	}

	for _, testCase := range testCases {
		result, err := parseGameSets(testCase.input)

		if testCase.expectedErr && err == nil {
			t.Errorf("Expected an error for input: %s", testCase.input)
		}

		if !testCase.expectedErr && err != nil {
			t.Errorf("Unexpected error for input: %s", testCase.input)
		}

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("For input '%s', expected %v, but got %v", testCase.input, testCase.expected, result)
		}
	}
}

func TestGetMaxColorGameSet(t *testing.T) {
	testCases := []struct {
		input    []cubeSet
		expected cubeSet
	}{
		{
			[]cubeSet{
				{3, 5, 2},
				{1, 8, 4},
				{7, 3, 6},
			},
			cubeSet{7, 8, 6},
		},
		{
			[]cubeSet{
				{2, 4, 7},
				{5, 2, 1},
				{3, 6, 9},
			},
			cubeSet{5, 6, 9},
		},
		{
			[]cubeSet{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			cubeSet{7, 8, 9},
		},
	}

	for _, testCase := range testCases {
		result := getMaxColorGameSet(testCase.input)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("For input %v, expected %v, but got %v", testCase.input, testCase.expected, result)
		}
	}
}

func TestIsGameSetsPossible(t *testing.T) {
	// Assuming theBag is a global variable or defined elsewhere in your code
	theBag := cubeSet{10, 10, 10}

	testCases := []struct {
		input    []cubeSet
		expected bool
	}{
		{
			[]cubeSet{
				{3, 5, 2},
				{1, 8, 4},
				{7, 3, 6},
			},
			true,
		},
		{
			[]cubeSet{
				{2, 4, 7},
				{5, 2, 1},
				{3, 6, 9},
			},
			true,
		},
		{
			[]cubeSet{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 19},
			},
			false,
		},
	}

	for _, testCase := range testCases {
		result := isGameSetsPossible(testCase.input)

		if result != testCase.expected {
			t.Errorf("For input %v and theBag %v, expected %v, but got %v", testCase.input, theBag, testCase.expected, result)
		}
	}
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type cubeSet struct {
	red   int
	blue  int
	green int
}

var theBag = cubeSet{
	red:   12,
	blue:  14,
	green: 13,
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

func getGameID(input string) (int, error) {
	// Define a regular expression pattern to match the game ID
	pattern := `Game (\d+):`
	re := regexp.MustCompile(pattern)

	// Find the first match in the input string
	match := re.FindStringSubmatch(input)

	if len(match) < 2 {
		return 0, errors.New("no game ID found in the string")
	}

	// Convert the matched value to an integer
	id, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, errors.New("failed to convert game ID to integer: " + err.Error())
	}

	return id, nil
}

func parseGameSets(line string) ([]cubeSet, error) {
	// Split the line into sets using semicolons as separators
	setStrings := strings.Split(line, ";")

	// Initialize a slice to store the parsed cubeSets
	var cubeSets []cubeSet

	// Define a regular expression pattern to extract numbers and colors from each set
	pattern := `(\d+) (blue|red|green)`
	re := regexp.MustCompile(pattern)

	// Iterate over each set string
	for _, setString := range setStrings {
		// Find all matches in the set string
		matches := re.FindAllStringSubmatch(setString, -1)

		// Initialize cubeSet for the current set
		set := cubeSet{}

		// Process each match and update the cubeSet
		for _, match := range matches {
			count, err := strconv.Atoi(match[1])
			if err != nil {
				return nil, errors.New("failed to convert count to integer: " + err.Error())
			}

			color := match[2]

			switch color {
			case "blue":
				set.blue += count
			case "red":
				set.red += count
			case "green":
				set.green += count
			default:
				return nil, errors.New("invalid color found: " + color)
			}
		}

		// Append the cubeSet to the result slice
		cubeSets = append(cubeSets, set)
	}

	return cubeSets, nil
}

func getMaxColorGameSet(gameSets []cubeSet) cubeSet {

	var maxGameSet cubeSet

	for _, gameSet := range gameSets {

		if gameSet.blue > maxGameSet.blue {
			maxGameSet.blue = gameSet.blue
		}
		if gameSet.green > maxGameSet.green {
			maxGameSet.green = gameSet.green
		}
		if gameSet.red > maxGameSet.red {
			maxGameSet.red = gameSet.red
		}

	}

	return maxGameSet
}

func isGameSetsPossible(gameSets []cubeSet) bool {

	for _, game := range gameSets {

		if game.blue > theBag.blue {
			return false
		}

		if game.red > theBag.red {
			return false
		}

		if game.green > theBag.green {
			return false
		}
	}

	return true
}

func main() {

	var possibleIdsSum int
	var gameSetPowerSum int

	//read input file
	games, err := readLines("./input.txt")
	if err != nil {
		panic(err.Error())
	}

	for _, game := range games {

		gameSets, err := parseGameSets(game)
		if err != nil {
			panic(err.Error())
		}

		maxGameSet := getMaxColorGameSet(gameSets)
		gameSetPower := maxGameSet.blue * maxGameSet.green * maxGameSet.red
		gameSetPowerSum += gameSetPower

		if isGameSetsPossible(gameSets) {
			id, err := getGameID(game)
			if err != nil {
				panic(err.Error())
			}

			possibleIdsSum += id
		}
	}

	fmt.Println("possible IDs sum: ", possibleIdsSum)
	fmt.Println("power sets sum: ", gameSetPowerSum)
}

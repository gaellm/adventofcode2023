package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type card struct {
	cardNumber   int
	winningNbrs  []int
	numbers      []int
	points       int
	winningTimes int
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

func parseNumbers(s string) ([]int, error) {
	var result []int

	numStrs := strings.Fields(s)
	for _, numStr := range numStrs {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number: %v", err)
		}
		result = append(result, num)
	}

	return result, nil
}

func parseLine(line string) (card, error) {
	var result card

	// Split the line into parts using ":" and "|"
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return result, errors.New("invalid input format")
	}

	// Extract card number
	cardNumber, err := strconv.Atoi(strings.TrimSpace(parts[0][len("Card "):]))
	if err != nil {
		return result, fmt.Errorf("failed to parse card number: %v", err)
	}
	result.cardNumber = cardNumber

	// Extract numbers and winning numbers
	numbersPart := strings.TrimSpace(parts[1])
	numbers := strings.Split(numbersPart, " | ")
	if len(numbers) != 2 {
		return result, errors.New("invalid numbers format")
	}

	// Parse winning numbers
	winningNbrs, err := parseNumbers(numbers[0])
	if err != nil {
		return result, errors.New("failed to parse winning numbers: " + err.Error())
	}
	result.winningNbrs = winningNbrs

	// Parse all numbers
	cardNumbers, err := parseNumbers(numbers[1])
	if err != nil {
		return result, errors.New("failed to parse card numbers: " + err.Error())
	}
	result.numbers = cardNumbers

	return result, nil
}

func countNumbersInSlice(slice1, slice2 []int) int {
	count := 0
	// Create a map to store the presence of numbers in the first slice
	numbersInSlice1 := make(map[int]bool)

	for _, num := range slice1 {
		numbersInSlice1[num] = true
	}

	for _, num := range slice2 {
		if numbersInSlice1[num] {
			count++
		}
	}

	return count
}

func doubleOnEachMatch(matchNbr int) int {

	if matchNbr < 1 {
		return 0
	}

	result := 1
	for i := 1; i < matchNbr; i++ {
		result *= 2
	}

	return result
}

func processCardsPoints(cards []card) []card {

	var pocceedCards []card

	for _, card := range cards {
		card.winningTimes = countNumbersInSlice(card.winningNbrs, card.numbers)
		card.points = doubleOnEachMatch(card.winningTimes)
		pocceedCards = append(pocceedCards, card)
	}

	return pocceedCards
}

func getCards(lines []string) ([]card, error) {

	var cards []card

	for _, line := range lines {

		card, err := parseLine(line)
		if err != nil {
			return nil, errors.New("parse line fail with error: " + err.Error())
		}
		cards = append(cards, card)

	}

	return processCardsPoints(cards), nil
}

func sortCardsByNumber(cards []card) map[int][]card {

	cardsSorted := make(map[int][]card)

	for _, curentCard := range cards {

		cardsSorted[curentCard.cardNumber] = []card{curentCard}

	}
	return cardsSorted
}

func processPart2(sortedCards map[int][]card) int {

	cardNb := len(sortedCards)
	totalCards := sortedCards

	//get all cards
	for _, cardTbl := range totalCards {
		for _, card := range cardTbl {
			fmt.Println("I've card number", card.cardNumber, "with", card.winningTimes, "winning numbers")
			for i := 1; i <= card.winningTimes; i++ {
				cardNb++
				fmt.Println("I add card number", card.cardNumber+i)
				totalCards[card.cardNumber+i] = append(totalCards[card.cardNumber+i], sortedCards[card.cardNumber+i][0])
			}
		}
	}

	return cardNb
}

func main() {

	//read input file
	lines, err := readLines("./input_test.txt")
	if err != nil {
		panic(err.Error())
	}

	cards, err := getCards(lines)
	if err != nil {
		panic(err.Error())
	}

	pointsSum := 0
	for _, card := range cards {
		pointsSum += card.points
	}

	fmt.Println("Part1 - Sum of cards points: ", pointsSum)

	sortedCards := sortCardsByNumber(cards)

	fmt.Println("Part2 - Sum of cards: ", processPart2(sortedCards))

}

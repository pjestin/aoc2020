package day23

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pjestin/aoc2020/lib"
)

type listNode struct {
	value int
	next  *listNode
}

func parseCups(input string) (*listNode, []*listNode, error) {
	splitLine := strings.Split(input, "")
	cupMap := make([]*listNode, len(input))
	listStart := listNode{}
	currentCup := &listStart
	for _, numberString := range splitLine {
		cup, err := strconv.Atoi(numberString)
		if err != nil {
			return &listStart, cupMap, err
		}
		nextCup := listNode{value: cup}
		cupMap[cup-1] = &nextCup
		currentCup.next = &nextCup
		currentCup = &nextCup
	}
	currentCup.next = listStart.next
	return listStart.next, cupMap, nil
}

func mod(a int, b int) int {
	return ((a % b) + b) % b
}

func makeMove(currentCup *listNode, cupMap []*listNode) *listNode {
	n := len(cupMap)

	// Get cups to move
	cupsToMoveBegin := currentCup.next
	cupsToMoveEnd := cupsToMoveBegin.next.next

	// Detach from main list
	currentCup.next = cupsToMoveEnd.next
	cupsToMoveEnd.next = nil

	// Find destination label
	pickedCupLabels := []int{cupsToMoveBegin.value, cupsToMoveBegin.next.value, cupsToMoveEnd.value}
	destinationCupLabel := mod(currentCup.value-2, n) + 1
	for lib.ContainsInt(pickedCupLabels, destinationCupLabel) {
		destinationCupLabel = mod(destinationCupLabel-2, n) + 1
	}

	// Find destination cup
	destinationCup := cupMap[destinationCupLabel-1]

	// Insert cups to move
	afterDestinationCup := destinationCup.next
	destinationCup.next = cupsToMoveBegin
	cupsToMoveEnd.next = afterDestinationCup

	// Return next cup
	return currentCup.next
}

func getOrderAfterOne(cupMap []*listNode) string {
	n := len(cupMap)
	currentCup := cupMap[0].next
	stringNumbers := make([]string, n)
	for i := 0; i < n-1; i++ {
		stringNumbers[i] = fmt.Sprint(currentCup.value)
		currentCup = currentCup.next
	}
	return strings.Join(stringNumbers, "")
}

// GetCupOrderAfterMoves moves the cups a number of times and returns their final order
func GetCupOrderAfterMoves(input string, moves int) (string, error) {
	currentCup, cupMap, err := parseCups(input)
	if err != nil {
		return "", err
	}
	for round := 0; round < moves; round++ {
		currentCup = makeMove(currentCup, cupMap)
	}
	return getOrderAfterOne(cupMap), nil
}

func insertRemainingCups(firstCup *listNode, n int) []*listNode {
	maxValue := 0
	currentCup := firstCup
	cupMap := make([]*listNode, n)
	for currentCup.next != firstCup {
		if currentCup.value > maxValue {
			maxValue = currentCup.value
		}
		cupMap[currentCup.value-1] = currentCup
		currentCup = currentCup.next
	}
	cupMap[currentCup.value-1] = currentCup
	for cup := maxValue + 1; cup <= n; cup++ {
		nextCup := listNode{value: cup}
		currentCup.next = &nextCup
		currentCup = &nextCup
		cupMap[cup-1] = currentCup
	}
	currentCup.next = firstCup
	return cupMap
}

func findTwoCupsAfterOneProduct(cupMap []*listNode) uint64 {
	cupOne := cupMap[0]
	return uint64(cupOne.next.value) * uint64(cupOne.next.next.value)
}

// GetFirstTwoCupsAfterTenMillionMoves parse input cups, pads them with numbers until 1 million, makes 10 million moves, and returns the product of the two cups after 1
func GetFirstTwoCupsAfterTenMillionMoves(input string) (uint64, error) {
	currentCup, _, err := parseCups(input)
	if err != nil {
		return 0, err
	}
	cupMap := insertRemainingCups(currentCup, 1000000)
	for round := 0; round < 10000000; round++ {
		currentCup = makeMove(currentCup, cupMap)
	}
	return findTwoCupsAfterOneProduct(cupMap), nil
}

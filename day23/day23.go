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

// func displayCups(firstCup *listNode) string {
// 	currentCup := firstCup
// 	var cupValues []string
// 	for currentCup != nil && currentCup.next != firstCup {
// 		cupValues = append(cupValues, fmt.Sprint(currentCup.value))
// 		currentCup = currentCup.next
// 	}
// 	if currentCup == nil {
// 		cupValues = append(cupValues, "!")
// 	} else {
// 		cupValues = append(cupValues, fmt.Sprint(currentCup.value))
// 	}
// 	return strings.Join(cupValues, " ")
// }

func parseCups(input string) (*listNode, int, error) {
	splitLine := strings.Split(input, "")
	listStart := listNode{}
	currentCup := &listStart
	for _, numberString := range splitLine {
		cup, err := strconv.Atoi(numberString)
		if err != nil {
			return &listStart, 0, err
		}
		nextCup := listNode{value: cup}
		currentCup.next = &nextCup
		currentCup = &nextCup
	}
	currentCup.next = listStart.next
	return listStart.next, len(input), nil
}

func mod(a int, b int) int {
	return ((a % b) + b) % b
}

func makeMove(currentCup *listNode, n int) *listNode {
	// Get cups to move
	cupsToMoveBegin := currentCup.next
	cupsToMoveEnd := cupsToMoveBegin.next.next

	// Detach from main list
	currentCup.next = cupsToMoveEnd.next
	afterCupsToMove := cupsToMoveEnd.next
	cupsToMoveEnd.next = nil

	// Find destination label
	pickedCupLabels := []int{cupsToMoveBegin.value, cupsToMoveBegin.next.value, cupsToMoveEnd.value}
	destinationCupLabel := mod(currentCup.value-2, n) + 1
	for lib.ContainsInt(pickedCupLabels, destinationCupLabel) {
		destinationCupLabel = mod(destinationCupLabel-2, n) + 1
	}

	// Find destination cup
	destinationCup := afterCupsToMove
	for destinationCup.value != destinationCupLabel {
		destinationCup = destinationCup.next
	}

	// Insert cups to move
	afterDestinationCup := destinationCup.next
	destinationCup.next = cupsToMoveBegin
	cupsToMoveEnd.next = afterDestinationCup

	// Return next cup
	return currentCup.next
}

func getOrderAfterOne(currentCup *listNode, n int) string {
	for currentCup.value != 1 {
		currentCup = currentCup.next
	}
	currentCup = currentCup.next
	stringNumbers := make([]string, n)
	for i := 0; i < n-1; i++ {
		stringNumbers[i] = fmt.Sprint(currentCup.value)
		currentCup = currentCup.next
	}
	return strings.Join(stringNumbers, "")
}

// GetCupOrderAfterMoves moves the cups a number of times and returns their final order
func GetCupOrderAfterMoves(input string, moves int) (string, error) {
	currentCup, n, err := parseCups(input)
	if err != nil {
		return "", err
	}
	for round := 0; round < moves; round++ {
		currentCup = makeMove(currentCup, n)
	}
	return getOrderAfterOne(currentCup, n), nil
}

package day22

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type deck struct {
	playerID int
	cards    []int
}

func (d *deck) pop() int {
	card := d.cards[0]
	d.cards = d.cards[1:]
	return card
}

func (d *deck) push(card int) {
	d.cards = append(d.cards, card)
}

func parseDecks(lines []string) ([]deck, error) {
	var decks []deck
	currentDeck := deck{}
	for _, line := range lines {
		if len(line) == 0 {
			decks = append(decks, currentDeck)
			currentDeck = deck{}
		} else if strings.HasPrefix(line, "Player") {
			idString := strings.Split(strings.Split(line, " ")[1], ":")[0]
			id, err := strconv.Atoi(idString)
			if err != nil {
				return nil, err
			}
			currentDeck.playerID = id
		} else {
			card, err := strconv.Atoi(line)
			if err != nil {
				return nil, err
			}
			currentDeck.push(card)
		}
	}
	decks = append(decks, currentDeck)
	return decks, nil
}

func playRound(decks []deck) int {
	maxCard := 0
	maxCardDeckIndex := 0
	roundCards := make([]int, len(decks))
	for deckIndex := range decks {
		card := decks[deckIndex].pop()
		roundCards[deckIndex] = card
		if card > maxCard {
			maxCard = card
			maxCardDeckIndex = deckIndex
		}
	}
	sort.Slice(roundCards, func(i, j int) bool { return roundCards[i] > roundCards[j] })
	for _, roundCard := range roundCards {
		decks[maxCardDeckIndex].push(roundCard)
	}
	return maxCardDeckIndex
}

func gameFinished(decks []deck) bool {
	cardCount := 0
	for _, d := range decks {
		cardCount += len(d.cards)
	}
	for _, d := range decks {
		if len(d.cards) == cardCount {
			return true
		}
	}
	return false
}

func getDeckScore(d deck) int {
	score := 0
	for cardIndex, card := range d.cards {
		score += (len(d.cards) - cardIndex) * card
	}
	return score
}

// GetCombatWinningPlayerScore plays Combat game until one of the players has all cards and returns their score
func GetCombatWinningPlayerScore(lines []string) (int, error) {
	decks, err := parseDecks(lines)
	if err != nil {
		return 0, err
	}
	winningDeckIndex := 0
	for !gameFinished(decks) {
		winningDeckIndex = playRound(decks)
	}
	return getDeckScore(decks[winningDeckIndex]), nil
}

func playRecursiveRound(decks []deck) int {
	maxCard := 0
	maxCardDeckIndex := 0
	roundCards := make([]int, len(decks))
	canRecurse := true
	for deckIndex := range decks {
		card := decks[deckIndex].pop()
		roundCards[deckIndex] = card
		canRecurse = canRecurse && len(decks[deckIndex].cards) >= card
		if card > maxCard {
			maxCard = card
			maxCardDeckIndex = deckIndex
		}
	}
	winningDeckIndex := -1
	if canRecurse {
		deckCopies := make([]deck, len(decks))
		for deckIndex, d := range decks {
			cardCopies := make([]int, roundCards[deckIndex])
			for cardIndex := 0; cardIndex < roundCards[deckIndex]; cardIndex++ {
				cardCopies[cardIndex] = d.cards[cardIndex]
			}
			deckCopies[deckIndex] = deck{playerID: d.playerID, cards: cardCopies}
		}
		winningDeckIndex = playRecursiveGame(deckCopies)
	} else {
		winningDeckIndex = maxCardDeckIndex
	}
	decks[winningDeckIndex].push(roundCards[winningDeckIndex])
	for _, roundCard := range roundCards {
		if roundCard != roundCards[winningDeckIndex] {
			decks[winningDeckIndex].push(roundCard)
		}
	}
	return winningDeckIndex
}

func playRecursiveGame(decks []deck) int {
	winningDeckIndex := 0
	visitedConfigs := make(map[string]bool)
	for !gameFinished(decks) {
		winningDeckIndex = playRecursiveRound(decks)
		configHash := getDeckConfigHash(decks)
		_, visited := visitedConfigs[configHash]
		if visited {
			winningDeckIndex = getDeckIndexForPlayerID(decks, 1)
			break
		}
		visitedConfigs[configHash] = true
	}
	return winningDeckIndex
}

func getDeckConfigHash(decks []deck) string {
	deckHashes := make([]string, 0, len(decks))
	for _, d := range decks {
		deckCardStrings := make([]string, len(d.cards))
		for cardIndex, card := range d.cards {
			deckCardStrings[cardIndex] = fmt.Sprint(card)
		}
		deckHashes = append(deckHashes, fmt.Sprintf("%d: %s", d.playerID, strings.Join(deckCardStrings, ",")))
	}
	return strings.Join(deckHashes, "; ")
}

func getDeckIndexForPlayerID(decks []deck, playerID int) int {
	for deckIndex, deck := range decks {
		if deck.playerID == playerID {
			return deckIndex
		}
	}
	return -1
}

// GetRecursiveCombatWinningPlayerScore is similar to GetCombatWinningPlayerScore except it applies the rules of Recursive Combat
func GetRecursiveCombatWinningPlayerScore(lines []string) (int, error) {
	decks, err := parseDecks(lines)
	if err != nil {
		return 0, err
	}
	winningDeckIndex := 0
	for !gameFinished(decks) {
		winningDeckIndex = playRecursiveGame(decks)
	}
	return getDeckScore(decks[winningDeckIndex]), nil
}

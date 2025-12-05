package server

import (
	"fmt"
	"strconv"
	"strings"
)

// Rules Validation
type RulesViolation error

// showAction lets the active player play a set of cards from their hand to the active set.
func (g *Game) showAction(params string) RulesViolation {
	p := g.ActivePlayer
	// get subset of cards from hand
	firstIndex, length := parseParams(params)
	if firstIndex < 0 || firstIndex+length > len(p.Hand) {
		return RulesViolation(fmt.Errorf("index out of range"))
	}
	set := p.Hand[firstIndex : firstIndex+length]

	// validate set is valid
	err := validateSet(set)
	if err != nil {
		return err
	}

	// check set beats active set
	if setComparison(set, g.ActiveSet) > 0 {
		// remove set from player's hand
		newHand := make([]*Card, 0)
		for i, card := range p.Hand {
			if i < firstIndex || i >= firstIndex+length {
				newHand = append(newHand, card)
			}
		}
		p.Hand = newHand

		// gain points equal to number of cards in active set
		p.Score += len(set)

		// update active set
		g.ActiveSet = set
		g.ActiveSetPlayer = p

		return nil
	}

	return RulesViolation(fmt.Errorf("set does not beat active set"))
}

// returns 1 if set beats activeSet, -1 otherwise
func setComparison(set, activeSet []*Card) int {
	// always beat the empty active set
	if len(activeSet) == 0 {
		return 1
	}

	// always beat a smaller set
	if len(set) > len(activeSet) {
		return 1
	} else if len(set) < len(activeSet) {
		return -1
	}

	// matching beats consecutive
	isSetConsecutive := func(s []*Card) bool {
		return s[0].Value1 != s[1].Value1
	}
	if !isSetConsecutive(set) && isSetConsecutive(activeSet) {
		return 1
	}

	// lowest number is tie breaker
	minValue := func(s []*Card) int {
		min := s[0].Value1
		for _, card := range s {
			if card.Value1 < min {
				min = card.Value1
			}
		}
		return min
	}
	if minValue(set) > minValue(activeSet) {
		return 1
	}
	return -1
}

// scoutAction lets the active player take a card from the active set and place it in their hand.
// params is "takeIndex,putIndex", where takeIndex is the index of the card in the active set to take,
// and putIndex is the index in the player's hand to insert the taken card.
func (g *Game) scoutAction(params string) RulesViolation {
	p := g.ActivePlayer
	takeIndex, putIndex := parseParams(params)
	if takeIndex < 0 || takeIndex >= len(g.ActiveSet) {
		return RulesViolation(fmt.Errorf("index out of range"))
	}

	card := g.ActiveSet[takeIndex]

	// add to player's hand
	p.Hand = append(p.Hand[:putIndex], append([]*Card{card}, p.Hand[putIndex:]...)...)

	// remove card from active set
	newActiveSet := make([]*Card, 0)
	for i, c := range g.ActiveSet {
		if i != takeIndex {
			newActiveSet = append(newActiveSet, c)
		}
	}
	g.ActiveSet = newActiveSet

	// award point to active set player
	if g.ActiveSetPlayer != nil {
		g.ActiveSetPlayer.Score += 1
	}

	// TODO: check for game end here, too

	return nil
}

// parseParams splits params string into two integers. Returns -1,0 on error.
func parseParams(params string) (int, int) {
	parts := strings.Split(params, ",")
	if len(parts) < 2 {
		return -1, 0
	}

	firstIndex, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return -1, 0
	}
	length, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return -1, 0
	}

	return firstIndex, length
}

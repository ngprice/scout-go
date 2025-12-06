package server

import (
	"fmt"
	"strconv"
	"strings"
)

// Game Action
type GameAction func(params string) RulesViolation

// Rules Validation
type RulesViolation error

// showAction lets the active player play a set of cards from their hand to the active set.
func (g *Game) showAction(params string) RulesViolation {
	firstIndex, length := parseParams(params)
	p := g.ActivePlayer

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
	if setComparison(set, g.ActiveSet) {
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

// returns true if set beats compSet, false otherwise
func setComparison(set, compSet []*Card) bool {
	// always beat the empty set
	if len(compSet) == 0 {
		return true
	}

	// always beat a smaller set
	if len(set) > len(compSet) {
		return true
	} else if len(set) < len(compSet) {
		return false
	}

	// matching beats consecutive
	isSetConsecutive := func(s []*Card) bool {
		return s[0].Value1 != s[1].Value1
	}
	if !isSetConsecutive(set) {
		if isSetConsecutive(compSet) {
			return true
		}
	} else {
		if !isSetConsecutive(compSet) {
			return false
		}
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

	return minValue(set) > minValue(compSet)
}

// scoutAction lets the active player take a card from the active set and place it in their hand.
// params is "takeIndex,putIndex", where takeIndex is the index of the card in the active set to take,
// and putIndex is the index in the player's hand to insert the taken card.
// if takeIndex is greater than the length of the active set, it indicates reversing
// the values of the card at [takeIndex - len(activeSet)].
func (g *Game) scoutAction(params string) RulesViolation {
	takeIndex, putIndex := parseParams(params)
	p := g.ActivePlayer

	reverse := takeIndex >= len(g.ActiveSet)
	if reverse {
		takeIndex -= len(g.ActiveSet)
	}

	// can only scout from the 'ends' of the active set
	if takeIndex != 0 && takeIndex != len(g.ActiveSet)-1 {
		return RulesViolation(fmt.Errorf("can only scout from ends of active set"))
	}

	card := g.ActiveSet[takeIndex]

	if reverse {
		card.ReverseValues()
	}

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

	// advance the ConsecutiveScouts counter
	g.ConsecutiveScouts += 1

	return nil
}

// parseParams splits params string into two integers. Returns -1,0 on error.
func parseParams(params string) (int, int) {
	parts := strings.Split(params, ",")
	if len(parts) < 2 {
		return 1, 0
	}

	firstIndex, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 1, 0
	}
	length, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 1, 0
	}

	return firstIndex, length
}

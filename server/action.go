package server

import (
	"fmt"
)

// scoutAction lets the active player take a card from the active set and place it in their hand.
// params are (takeIndex, putIndex), where takeIndex is the index of the card in the active set to take,
// and putIndex is the index in the player's hand to insert the taken card.
func (g *Game) scoutAction(takeIndex, putIndex int) RulesViolation {
	p := g.ActivePlayer

	if !g.isValidScout(p, takeIndex, putIndex) {
		return RulesViolation(fmt.Errorf("not a valid scout action"))
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

	// advance the ConsecutiveScouts counter
	g.ConsecutiveScouts += 1

	return nil
}

// same as scoutAction, but reverses the card
func (g *Game) scoutActionReverse(takeIndex, putIndex int) RulesViolation {
	p := g.ActivePlayer

	if !g.isValidScout(p, takeIndex, putIndex) {
		return RulesViolation(fmt.Errorf("not a valid scout action"))
	}

	card := g.ActiveSet[takeIndex]

	card.ReverseValues()

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

func (g *Game) isValidScout(p *Player, takeIndex, putIndex int) bool {
	if len(g.ActiveSet) == 0 {
		return false
	}
	if putIndex > len(p.Hand) {
		return false
	}
	// can only scout from the 'ends' of the active set
	return !(takeIndex != 0 && takeIndex != len(g.ActiveSet)-1)
}

// showAction lets the active player play a set of cards from their hand to the active set.
// params are (firstIndex, length), where firstIndex is the index in the active player's hand
// to start the set, and length is the length of the set.
func (g *Game) showAction(firstIndex, length int) RulesViolation {
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

		// gain points equal to number of cards in active set you beat
		if g.ActiveSet != nil {
			p.Score += len(g.ActiveSet)
		}

		// update active set
		g.ActiveSet = set
		g.ActiveSetPlayer = p

		return nil
	}

	return RulesViolation(fmt.Errorf("set does not beat active set"))
}

func (g *Game) isValidShow(hand []*Card, firstIndex, length int) bool {
	if firstIndex < 0 || firstIndex+length > len(hand) {
		return false
	}
	set := hand[firstIndex : firstIndex+length]
	err := validateSet(set)
	if err != nil {
		return false
	}
	return setComparison(set, g.ActiveSet)
}

// scoutAndShowAction lets the active player perform a scout action, immediately followed by
// a show action. it can be used once per round.
// params are ({scoutParams}, {showParams}), where each arg is the params string passed to the
// action. for example, (0,1,1,3) is equivalent to scoutAction(0, 1), then showAction(1, 3)
func (g *Game) scoutAndShowAction(takeIndex, putIndex, firstIndex, length int) RulesViolation {
	err := g.scoutAction(takeIndex, putIndex)
	if err != nil {
		return err
	}
	err = g.showAction(firstIndex, length)
	if err != nil {
		return err
	}

	g.ActivePlayer.CanScoutAndShow = false
	return nil
}

func (g *Game) isValidScoutAndShow(p *Player, takeIndex, putIndex, startIndex, length int) bool {
	if g.isValidScout(p, takeIndex, putIndex) {
		// assemble the hand as it would be after scout
		hand := make([]*Card, len(p.Hand))
		copy(hand, p.Hand)
		cardCopy := &Card{
			Value1: g.ActiveSet[takeIndex].Value1,
			Value2: g.ActiveSet[takeIndex].Value2,
		}
		hand = append(hand[:putIndex], append([]*Card{cardCopy}, hand[putIndex:]...)...)
		return g.isValidShow(hand, startIndex, length)
	}
	return false
}

// same as scoutAndShow, but reverses the card
func (g *Game) scoutAndShowActionReverse(takeIndex, putIndex, firstIndex, length int) RulesViolation {
	err := g.scoutActionReverse(takeIndex, putIndex)
	if err != nil {
		return err
	}
	err = g.showAction(firstIndex, length)
	if err != nil {
		return err
	}

	g.ActivePlayer.CanScoutAndShow = false
	return nil
}

func (g *Game) isValidScoutAndShowReverse(p *Player, takeIndex, putIndex, startIndex, length int) bool {
	if g.isValidScout(p, takeIndex, putIndex) {
		// assemble the hand as it would be after scout
		hand := make([]*Card, len(p.Hand))
		copy(hand, p.Hand)
		cardCopy := &Card{
			Value1: g.ActiveSet[takeIndex].Value1,
			Value2: g.ActiveSet[takeIndex].Value2,
		}
		cardCopy.ReverseValues()
		hand = append(hand[:putIndex], append([]*Card{cardCopy}, hand[putIndex:]...)...)
		return g.isValidShow(hand, startIndex, length)
	}
	return false
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
		if len(s) < 2 {
			return false
		}
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

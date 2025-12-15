package server

import (
	"fmt"
)

// scoutAction lets the active player take a card from the active set and place it in their hand.
// params are (takeIndex, putIndex), where takeIndex is the index of the card in the active set to take,
// and putIndex is the index in the player's hand to insert the taken card.
func (g *Game) scoutAction(takeIndex, putIndex int) RulesViolation {
	p := g.ActivePlayer

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

// showAction lets the active player play a set of cards from their hand to the active set.
// params are (firstIndex, length), where firstIndex is the index in the active player's hand
// to start the set, and length is the length of the set.
func (g *Game) showAction(firstIndex, length int) RulesViolation {
	p := g.ActivePlayer

	if firstIndex < 0 || firstIndex+length > len(p.Hand) {
		return RulesViolation(fmt.Errorf("index out of range"))
	}
	set := p.Hand[firstIndex : firstIndex+length]

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

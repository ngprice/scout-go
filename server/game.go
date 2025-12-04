package server

import (
	"fmt"
	"strconv"
	"strings"
)

type Game struct {
	Players         []*Player
	ActiveSet       []*Card
	ActiveSetPlayer *Player
	Complete        bool
}

func NewGame(numPlayers int) *Game {
	// init players
	players := make([]*Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		player, err := NewPlayer("Player"+string(i+1), i)
		if err != nil {
			panic(err)
		}
		players[i] = player
	}

	// init deck
	deckSize := 20 // example for now
	deck := make([]*Card, deckSize)
	for i := 0; i < deckSize; i++ {
		card, err := NewCard(i+1, i+2)
		if err != nil {
			panic(err)
		}
		deck[i] = card
	}

	// deal cards to each player; players get same number of cards
	for i := 0; i < numPlayers; i++ {
		for j := 0; j < deckSize/numPlayers; j++ {
			players[i].Hand = append(players[i].Hand, deck[i*numPlayers+j])
		}
	}

	return &Game{
		Players: players,
	}
}

func (g *Game) PlayerAction(playerIndex int, action string, params string) RulesViolation {
	player := g.Players[playerIndex]

	switch action {
	case "scout":
		return g.scoutAction(player, params)
	case "show":
		return g.showAction(player, params)
	case "scoutandshow":
		// TODO: implement scoutandshow logic
	default:
		return nil
	}

	return nil
}

// Rules Validation
type RulesViolation error

// showAction lets a player to play a set of cards from their hand to the active set.
func (g *Game) showAction(p *Player, params string) RulesViolation {
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
		// remove set from hand and update active set
		newHand := make([]*Card, 0)
		for i, card := range p.Hand {
			if i < firstIndex || i >= firstIndex+length {
				newHand = append(newHand, card)
			}
		}
		p.Hand = newHand

		g.ActiveSet = set
		g.ActiveSetPlayer = p

		// gain points equal to number of cards in active set
		p.Score += len(set)

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
	}

	// matching beats consecutive
	isSetConsecutive := func(s []*Card) bool {
		return s[0] != s[1]
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

// scoutAction lets a player take a card from the active set and place it in their hand.
func (g *Game) scoutAction(p *Player, params string) RulesViolation {
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

package server

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/google/uuid"
)

type Game struct {
	Id                string
	NumPlayers        int
	Players           []*Player
	ActivePlayer      *Player
	ActiveSet         []*Card
	ActiveSetPlayer   *Player
	ConsecutiveScouts int
	Round             int
	Complete          bool
	mu                sync.RWMutex
}

func NewGame(numPlayers int) (*Game, RulesViolation) {
	// init players
	players := make([]*Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		player, err := NewPlayer("Player"+strconv.Itoa(i+1), i)
		if err != nil {
			panic(err)
		}
		players[i] = player
	}

	deck, _ := NewGameDeck(numPlayers)

	// deal cards to each player; players get same number of cards
	for i := 0; i < len(deck); i++ {
		players[i%numPlayers].Hand = append(players[i%numPlayers].Hand, deck[i])
	}

	return &Game{
		Id:           uuid.New().String(),
		NumPlayers:   numPlayers,
		Players:      players,
		ActivePlayer: players[0],
	}, nil
}

func (g *Game) PlayerAction(playerIndex int, action *ActionSpec) RulesViolation {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Complete {
		return RulesViolation(fmt.Errorf("game is complete"))
	}

	if playerIndex != g.ActivePlayer.Index {
		return RulesViolation(fmt.Errorf("not your turn"))
	}

	if !g.IsActionValid(playerIndex, action) {
		return RulesViolation(fmt.Errorf("action is invalid"))
	}

	var err RulesViolation

	switch action.Type {
	case ActionScout:
		err = g.scoutAction(action.ScoutTakeIndex, action.ScoutPutIndex)
	case ActionScoutReverse:
		err = g.scoutActionReverse(action.ScoutTakeIndex, action.ScoutPutIndex)
	case ActionShow:
		err = g.showAction(action.ShowFirstIndex, action.ShowLength)
	case ActionScoutAndShow:
		err = g.scoutAndShowAction(action.ScoutTakeIndex, action.ScoutPutIndex, action.ShowFirstIndex, action.ShowLength)
	case ActionScoutAndShowReverse:
		err = g.scoutAndShowActionReverse(action.ScoutTakeIndex, action.ScoutPutIndex, action.ShowFirstIndex, action.ShowLength)
	case ActionReverseHand:
		g.ActivePlayer.ReverseHand()
		return nil
	default:
		return RulesViolation(fmt.Errorf("unknown action"))
	}

	if err != nil {
		return err
	}

	// prevent reverse hand after the first round
	g.ActivePlayer.CanReverseHand = false

	// set the next active player
	g.ActivePlayer = g.Players[(g.ActivePlayer.Index+1)%len(g.Players)]

	if g.checkRoundCompletion(); g.Complete {
		g.calculateScores()
		if g.checkGameCompletion(); g.Complete {
			return nil // game over
		}
	}

	return nil
}

func (g *Game) IsActionValid(playerIndex int, action *ActionSpec) bool {
	p := g.Players[playerIndex]

	switch action.Type {
	case ActionScout:
		return g.isValidScout(p, action.ScoutTakeIndex, action.ScoutPutIndex)
	case ActionScoutReverse:
		return g.isValidScout(p, action.ScoutTakeIndex, action.ScoutPutIndex)
	case ActionShow:
		return g.isValidShow(p.Hand, action.ShowFirstIndex, action.ShowLength)
	case ActionScoutAndShow:
		return g.isValidScoutAndShow(p, action.ScoutTakeIndex, action.ScoutPutIndex, action.ShowFirstIndex, action.ShowLength)
	case ActionScoutAndShowReverse:
		return g.isValidScoutAndShowReverse(p, action.ScoutTakeIndex, action.ScoutPutIndex, action.ShowFirstIndex, action.ShowLength)
	case ActionReverseHand:
		return p.CanReverseHand
	default:
		return false
	}
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

func (g *Game) isValidScoutAndShow(p *Player, takeIndex, putIndex, startIndex, length int) bool {
	if !p.CanScoutAndShow {
		return false
	}
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

func (g *Game) isValidScoutAndShowReverse(p *Player, takeIndex, putIndex, startIndex, length int) bool {
	if !p.CanScoutAndShow {
		return false
	}
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

func (g *Game) checkRoundCompletion() {
	// all others players have scouted in succession
	if g.ConsecutiveScouts == len(g.Players)-1 {
		g.Complete = true
	}
	// active player has emptied their hand
	if len(g.ActivePlayer.Hand) == 0 {
		g.Complete = true
	}
}

func (g *Game) checkGameCompletion() {
	g.Round++
	// play a number of rounds equal to number of players
	if g.Round >= len(g.Players) {
		g.Complete = true
	} else { // reset for next round
		g.ActivePlayer = g.Players[g.Round]
		g.ActiveSet = []*Card{}
		g.ActiveSetPlayer = nil
		g.ConsecutiveScouts = 0
		g.Complete = false
		for _, p := range g.Players {
			p.CanReverseHand = true
			p.CanScoutAndShow = true
		}
	}
}

func (g *Game) calculateScores() {
	// lose a point for each card in hand, unless you played the set that ended the game
	for _, p := range g.Players {
		penalty := len(p.Hand)
		if g.ActiveSetPlayer != nil && p.Index == g.ActiveSetPlayer.Index {
			penalty = 0
		}
		p.Score -= penalty
	}
}

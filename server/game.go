package server

import (
	"fmt"
	"strconv"
)

type Game struct {
	Players           []*Player
	ActivePlayer      *Player
	ActiveSet         []*Card
	ActiveSetPlayer   *Player
	ConsecutiveScouts int
	Round             int
	Complete          bool
}

func NewGame(numPlayers int) *Game {
	// init players
	players := make([]*Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		player, err := NewPlayer("Player"+strconv.Itoa(i+1), i)
		if err != nil {
			panic(err)
		}
		players[i] = player
	}

	// init deck
	deckSize := 20 // example for now
	cards := make([]*Card, deckSize)
	for i := 0; i < deckSize; i++ {
		card, err := NewCard(i+1, i+2)
		if err != nil {
			panic(err)
		}
		cards[i] = card
	}

	deck := NewDeck(cards)
	deck.Shuffle()

	// deal cards to each player; players get same number of cards
	for i := 0; i < numPlayers; i++ {
		for j := 0; j < deckSize/numPlayers; j++ {
			players[i].Hand = append(players[i].Hand, deck[i*numPlayers+j])
		}
	}

	return &Game{
		Players:      players,
		ActivePlayer: players[0],
	}
}

func (g *Game) PlayerAction(playerIndex int, action string, params string) RulesViolation {
	if g.Complete {
		return RulesViolation(fmt.Errorf("game is complete"))
	}

	if playerIndex != g.ActivePlayer.Index {
		return RulesViolation(fmt.Errorf("not your turn"))
	}

	var err RulesViolation

	switch action {
	case "scout":
		err = g.scoutAction(params)
	case "show":
		err = g.showAction(params)
	case "scoutandshow":
		// TODO: implement scoutandshow logic
	default:
		return RulesViolation(fmt.Errorf("unknown action"))
	}

	if err != nil {
		return err
	}

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
	// play a number of rounds equal to number of players
	if g.Round >= len(g.Players) {
		g.Complete = true
	} else { // reset for next round
		g.Round++
		g.ActivePlayer = g.Players[g.Round]
		g.ActiveSet = []*Card{}
		g.ActiveSetPlayer = nil
		g.ConsecutiveScouts = 0
		g.Complete = false
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

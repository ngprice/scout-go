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

	// check for game completion
	if g.checkCompletion(); g.Complete {
		return nil // exit game loop
	}

	// set the next active player
	g.ActivePlayer = g.Players[(g.ActivePlayer.Index+1)%len(g.Players)]

	return nil
}

func (g *Game) checkCompletion() {
	if g.ConsecutiveScouts == len(g.Players)-1 {
		g.Complete = true
	}
	if len(g.ActivePlayer.Hand) == 0 {
		g.Complete = true
	}
}

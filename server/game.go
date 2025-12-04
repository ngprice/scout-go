package server

import "strconv"

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
	// TODO: actually deal the cards sequentially?
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

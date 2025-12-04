package server

type Game struct {
	Players   []*Player
	ActiveSet []*Card
}

func NewGame(numPlayers int) *Game {
	// init players
	players := make([]*Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = NewPlayer("Player"+string(i+1), i)
	}

	// init deck
	deckSize := 20 // example for now
	deck := make([]*Card, deckSize)
	for i := 0; i < deckSize; i++ {
		deck[i] = NewCard(i+1, i+2)
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

func (g *Game) ApplyAction(action string) RulesViolation {
	err := g.validateAction(action)
	if err != nil {
		return err
	}

	return nil
}

// Rules Validation
type RulesViolation error

func (g *Game) validateAction(action string) RulesViolation {
	return nil
}

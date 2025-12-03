package server

type Game struct {
	numPlayers int32
}

func NewGame(numPlayers int32) *Game {
	return &Game{
		// Initialize game state
		numPlayers: numPlayers,
	}
}

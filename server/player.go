package server

type Player struct {
	Name         string
	Index        int
	Score        int
	Hand         []*Card
	ScoutAndShow bool
}

func NewPlayer(name string, index int) *Player {
	return &Player{
		Name:  name,
		Index: index,
		Score: 0,
		Hand:  make([]*Card, 0),
	}
}

package server

type Player struct {
	Name         string
	Score        int32
	Hand         []*Card
	ScoutAndShow bool
}

func NewPlayer(name string) *Player {
	return &Player{
		Name:  name,
		Score: 0,
		Hand:  make([]*Card, 0),
	}
}

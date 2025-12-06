package server

type Player struct {
	Name              string
	Index             int
	Score             int
	Hand              []*Card
	FlipHandAvail     bool
	ScoutAndShowAvail bool
}

func NewPlayer(name string, index int) (*Player, error) {
	return &Player{
		Name:              name,
		Index:             index,
		Score:             0,
		Hand:              make([]*Card, 0),
		FlipHandAvail:     true,
		ScoutAndShowAvail: true,
	}, nil
}

func (p *Player) FlipHand() {
	for _, card := range p.Hand {
		card.ReverseValues()
	}
	p.FlipHandAvail = false
}

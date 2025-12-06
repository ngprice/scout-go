package server

type Player struct {
	Name            string
	Index           int
	Score           int
	Hand            []*Card
	CanReverseHand  bool
	CanScoutAndShow bool
}

func NewPlayer(name string, index int) (*Player, error) {
	return &Player{
		Name:            name,
		Index:           index,
		Score:           0,
		Hand:            make([]*Card, 0),
		CanReverseHand:  true,
		CanScoutAndShow: true,
	}, nil
}

func (p *Player) ReverseHand() {
	for _, card := range p.Hand {
		card.ReverseValues()
	}
	p.CanReverseHand = false
}

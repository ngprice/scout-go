package server

type Card struct {
	Value1 int
	Value2 int
}

func NewCard(value1, value2 int) *Card {
	return &Card{
		Value1: value1,
		Value2: value2,
	}
}

func (c *Card) FlipValues() {
	c.Value1, c.Value2 = c.Value2, c.Value1
}

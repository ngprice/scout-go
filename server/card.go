package server

import "fmt"

type Card struct {
	Value1 int
	Value2 int
}

func NewCard(value1, value2 int) (*Card, error) {
	if value1 == 0 || value2 == 0 {
		return nil, fmt.Errorf("card values cannot be zero")
	}
	return &Card{
		Value1: value1,
		Value2: value2,
	}, nil
}

func (c *Card) FlipValues() {
	c.Value1, c.Value2 = c.Value2, c.Value1
}

type Deck []*Card

func NewDeck(cards []*Card) Deck {
	return Deck(cards)
}

func (d *Deck) Shuffle() {
	// TODO: randomize the order of cards in the deck
}

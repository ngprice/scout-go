package server

import (
	"fmt"
	"math/rand"
)

type Deck []*Card

func NewDeck(cards []*Card) Deck {
	deck := Deck(cards)
	deck.Shuffle()
	return deck
}

// NewGameDeck returns the deck used in the tabletop game
func NewGameDeck(numPlayers int) (Deck, RulesViolation) {
	if numPlayers < 2 || numPlayers > 5 {
		return nil, RulesViolation(fmt.Errorf("invalid number of players"))
	}
	cards := make([]*Card, 36)
	ctr := 0
	for i := 1; i <= 9; i++ {
		for j := i + 1; j <= 9; j++ {
			cards[ctr], _ = NewCard(i, j)
			ctr++
		}
	}

	if numPlayers == 3 || numPlayers == 5 {
		ninetencard, _ := NewCard(9, 10)
		cards = append(cards, ninetencard)
	}

	if numPlayers == 2 || numPlayers == 4 || numPlayers == 5 {
		for i := 1; i <= 8; i++ {
			tencard, _ := NewCard(i, 10)
			cards = append(cards, tencard)
		}
	}

	return NewDeck(cards), nil
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(*d), func(i, j int) {
		(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
	})

	// randomize orientation
	for _, card := range *d {
		if rand.Intn(2) == 0 {
			card.ReverseValues()
		}
	}
}

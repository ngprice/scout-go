package server

import (
	"fmt"
	"math"
)

// validateSet checks if the provided set of cards is valid according to game rules.
func validateSet(set []*Card) RulesViolation {
	if len(set) == 0 {
		return RulesViolation(fmt.Errorf("set cannot be empty"))
	}
	if len(set) == 1 {
		return nil // single card is always valid
	}

	var isMatching bool = false

	// gut-check the set for obvious invalid cases
	if set[0].Value1 == set[1].Value1 {
		isMatching = true
	} else if math.Abs(float64(set[0].Value1-set[1].Value1)) == 1 {
		isMatching = false
	} else {
		return RulesViolation(fmt.Errorf("set is neither consecutive nor matching"))
	}

	if isMatching {
		// check all cards match
		firstValue := set[0].Value1
		for _, card := range set {
			if card.Value1 != firstValue {
				return RulesViolation(fmt.Errorf("set is not matching"))
			}
		}
	} else {
		// check all cards are consecutive
		var ascending bool = true
		val := set[0].Value1
		for _, card := range set[1:] {
			if math.Abs(float64(card.Value1-val)) != 1 {
				return RulesViolation(fmt.Errorf("set is neither consecutive nor matching"))
			}

			val = card.Value1

			if card.Value1 < val { // descending set
				ascending = false
			}

			if ascending && card.Value1 < val {
				return RulesViolation(fmt.Errorf("set is not consistently ascending or descending"))
			} else if !ascending && card.Value1 > val {
				return RulesViolation(fmt.Errorf("set is not consistently ascending or descending"))
			}
		}
	}

	return nil
}

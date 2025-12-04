package server

import "fmt"

// validateSet checks if the provided set of cards is valid according to game rules.
func validateSet(set []*Card) RulesViolation {
	if len(set) == 0 {
		return RulesViolation(fmt.Errorf("set cannot be empty"))
	}
	if len(set) == 1 {
		return nil
	}

	// check consecutive numbers
	var val int = 0
	for _, card := range set {
		if val == 0 {
			val = card.Value1
		} else if val > 0 { // increasing set
			if card.Value1 == val+1 {
				val = card.Value1
				continue
			}
		} else if val < 0 { // decreasing set
			if card.Value1 == val-1 {
				val = card.Value1 * -1 // mark as decreasing set
				continue
			}
		} else if card.Value1 == val { // matching set
			continue
		}

	}

	return nil
}

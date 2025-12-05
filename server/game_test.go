package server

import (
	"testing"
)

func TestGameInitialization(t *testing.T) {
	game := NewGame(2)
	if game == nil {
		t.Fatal("NewGame returned nil")
	}
	if len(game.Players) != 2 {
		t.Fatalf("expected 2 players, got %d", len(game.Players))
	}
}

func TestHandDealtEvenly(t *testing.T) {
	numPlayers := 4
	game := NewGame(numPlayers)
	if len(game.Players) != numPlayers {
		t.Fatalf("expected %d players, got %d", numPlayers, len(game.Players))
	}

	var expected int = -1
	for i, p := range game.Players {
		if p == nil {
			t.Fatalf("player %d is nil", i)
		}
		hl := len(p.Hand)
		if expected == -1 {
			expected = hl
		} else if hl != expected {
			t.Fatalf("uneven hand sizes: player %d has %d cards, expected %d", i, hl, expected)
		}
	}
}

func TestSetComparison(t *testing.T) {
	type testCase struct {
		name      string
		set       []*Card
		activeSet []*Card
		expected  bool
	}
	testCases := []testCase{
		{
			name:      "empty set",
			set:       []*Card{},
			activeSet: []*Card{{Value1: 3}},
			expected:  false,
		},
		{
			name:      "empty active set",
			set:       []*Card{{Value1: 3}},
			activeSet: []*Card{},
			expected:  true,
		},
		{
			name:      "lower straight loses",
			set:       []*Card{{Value1: 2}, {Value1: 3}},
			activeSet: []*Card{{Value1: 4}, {Value1: 5}},
			expected:  false,
		},
		{
			name:      "higher straight wins",
			set:       []*Card{{Value1: 5}, {Value1: 6}},
			activeSet: []*Card{{Value1: 3}, {Value1: 4}},
			expected:  true,
		},
		{
			name:      "bigger set wins",
			set:       []*Card{{Value1: 3}, {Value1: 4}, {Value1: 5}},
			activeSet: []*Card{{Value1: 2}, {Value1: 3}},
			expected:  true,
		},
		{
			name:      "smaller set loses",
			set:       []*Card{{Value1: 3}, {Value1: 4}},
			activeSet: []*Card{{Value1: 2}, {Value1: 3}, {Value1: 4}},
			expected:  false,
		},
		{
			name:      "matching beats consecutive",
			set:       []*Card{{Value1: 5}, {Value1: 5}, {Value1: 5}},
			activeSet: []*Card{{Value1: 5}, {Value1: 6}, {Value1: 7}},
			expected:  true,
		},
		{
			name:      "smaller matching beats consecutive",
			set:       []*Card{{Value1: 1}, {Value1: 1}, {Value1: 1}},
			activeSet: []*Card{{Value1: 5}, {Value1: 6}, {Value1: 7}},
			expected:  true,
		},
		{
			name:      "consecutive loses to matching",
			set:       []*Card{{Value1: 4}, {Value1: 5}, {Value1: 6}},
			activeSet: []*Card{{Value1: 7}, {Value1: 7}, {Value1: 7}},
			expected:  false,
		},
		{
			name:      "higher matching wins",
			set:       []*Card{{Value1: 8}, {Value1: 8}, {Value1: 8}},
			activeSet: []*Card{{Value1: 5}, {Value1: 5}, {Value1: 5}},
			expected:  true,
		},
		// the actual tests from the manual
		{
			name:      "stronger set ok",
			set:       []*Card{{Value1: 2}, {Value1: 3}, {Value1: 4}},
			activeSet: []*Card{{Value1: 8}, {Value1: 8}},
			expected:  true,
		},
		{
			name:      "stronger set no",
			set:       []*Card{{Value1: 1}},
			activeSet: []*Card{{Value1: 8}, {Value1: 8}},
			expected:  false,
		},
		{
			name:      "type of set ok",
			set:       []*Card{{Value1: 2}, {Value1: 2}},
			activeSet: []*Card{{Value1: 4}, {Value1: 5}},
			expected:  true,
		},
		{
			name:      "type of set no",
			set:       []*Card{{Value1: 4}, {Value1: 5}},
			activeSet: []*Card{{Value1: 2}, {Value1: 2}},
			expected:  false,
		},
		{
			name:      "numbers in set ok",
			set:       []*Card{{Value1: 5}, {Value1: 6}},
			activeSet: []*Card{{Value1: 4}, {Value1: 5}},
			expected:  true,
		},
		{
			name:      "numbers in set no 1",
			set:       []*Card{{Value1: 4}, {Value1: 5}},
			activeSet: []*Card{{Value1: 4}, {Value1: 5}},
			expected:  false,
		},
		{
			name:      "numbers in set no 2",
			set:       []*Card{{Value1: 4}, {Value1: 5}},
			activeSet: []*Card{{Value1: 4}, {Value1: 5}},
			expected:  false,
		},
	}
	for i, tc := range testCases {
		result := setComparison(tc.set, tc.activeSet)
		if result != tc.expected {
			t.Fatalf("test case %d: expected %v, got %v", i, tc.expected, result)
		}
	}
}

func TestValidateSet(t *testing.T) {
	type testCase struct {
		name    string
		set     []*Card
		wantErr bool // true if valid, false if invalid
	}
	testCases := []testCase{
		{
			name:    "empty set",
			set:     []*Card{},
			wantErr: true,
		},
		{
			name:    "single card set",
			set:     []*Card{{Value1: 3}},
			wantErr: false,
		},
		{
			name:    "two card matching set",
			set:     []*Card{{Value1: 5}, {Value1: 5}},
			wantErr: false,
		},
		{
			name:    "two card consecutive set",
			set:     []*Card{{Value1: 7}, {Value1: 8}},
			wantErr: false,
		},
		{
			name:    "valid ascending consecutive set",
			set:     []*Card{{Value1: 2}, {Value1: 3}, {Value1: 4}},
			wantErr: false,
		},
		{
			name:    "valid descending consecutive set",
			set:     []*Card{{Value1: 5}, {Value1: 4}, {Value1: 3}},
			wantErr: false,
		},
		{
			name:    "valid matching set",
			set:     []*Card{{Value1: 5}, {Value1: 5}, {Value1: 5}},
			wantErr: false,
		},
		{
			name:    "invalid non-consecutive set",
			set:     []*Card{{Value1: 2}, {Value1: 4}, {Value1: 5}},
			wantErr: true,
		},
		{
			name:    "invalid mixed set 1",
			set:     []*Card{{Value1: 3}, {Value1: 3}, {Value1: 4}},
			wantErr: true,
		},
		{
			name:    "invalid mixed set 2",
			set:     []*Card{{Value1: 3}, {Value1: 4}, {Value1: 4}},
			wantErr: true,
		},
		{
			name:    "broken set 1",
			set:     []*Card{{Value1: 4}, {Value1: 6}, {Value1: 8}},
			wantErr: true,
		},
		{
			name:    "broken set 2",
			set:     []*Card{{Value1: 4}, {Value1: 4}, {Value1: 6}, {Value1: 8}, {Value1: 8}},
			wantErr: true,
		},
	}
	for i, tc := range testCases {
		err := validateSet(tc.set)
		isErr := err != nil
		if isErr != tc.wantErr {
			t.Fatalf("test case %d (%s): expected valid=%v, got valid=%v", i, tc.name, tc.wantErr, isErr)
		}
	}
}

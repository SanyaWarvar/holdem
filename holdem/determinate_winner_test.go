package holdem

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompareCards(t *testing.T) {
	cases := []struct {
		TestCaseName string
		A            []Card
		B            []Card
		Expected     int
	}{
		{
			TestCaseName: "Empty",
			A:            []Card{},
			B:            []Card{},
			Expected:     0,
		},
		{
			TestCaseName: "Empty And Not Empty",
			A:            []Card{},
			B:            []Card{{Suit: "Spades", Value: 2}},
			Expected:     -1,
		},
		{
			TestCaseName: "Not Empty And Empty",
			A:            []Card{{Suit: "Diamonds", Value: 12}},
			B:            []Card{},
			Expected:     1,
		},
		{
			TestCaseName: "Equal Same Suits",
			A: []Card{
				{Suit: "Spades", Value: 14}, {Suit: "Spades", Value: 13},
				{Suit: "Spades", Value: 12}, {Suit: "Spades", Value: 11},
			},
			B: []Card{
				{Suit: "Spades", Value: 14}, {Suit: "Spades", Value: 13},
				{Suit: "Spades", Value: 12}, {Suit: "Spades", Value: 11},
			},
			Expected: 0,
		},
		{
			TestCaseName: "Equal Different Suits",
			A: []Card{
				{Suit: "Spades", Value: 14}, {Suit: "Spades", Value: 13},
				{Suit: "Spades", Value: 12}, {Suit: "Spades", Value: 11},
			},
			B: []Card{
				{Suit: "Hearts", Value: 14}, {Suit: "Hearts", Value: 13},
				{Suit: "Hearts", Value: 12}, {Suit: "Hearts", Value: 11},
			},
			Expected: 0,
		},
		{
			TestCaseName: "Equal Many Different Suits",
			A: []Card{
				{Suit: "Clubs", Value: 14}, {Suit: "Spades", Value: 13},
				{Suit: "Diamonds", Value: 12}, {Suit: "Hearts", Value: 11},
			},
			B: []Card{
				{Suit: "Diamonds", Value: 14}, {Suit: "Clubs", Value: 13},
				{Suit: "Hearts", Value: 12}, {Suit: "Spades", Value: 11},
			},
			Expected: 0,
		},
		{
			TestCaseName: "Greater first",
			A: []Card{
				{Suit: "Clubs", Value: 14}, {Suit: "Spades", Value: 13},
				{Suit: "Diamonds", Value: 12}, {Suit: "Hearts", Value: 11},
			},
			B: []Card{
				{Suit: "Diamonds", Value: 10}, {Suit: "Clubs", Value: 13},
				{Suit: "Hearts", Value: 12}, {Suit: "Spades", Value: 11},
			},
			Expected: 1,
		},
		{
			TestCaseName: "Greater last",
			A: []Card{
				{Suit: "Clubs", Value: 6}, {Suit: "Spades", Value: 2},
				{Suit: "Diamonds", Value: 2}, {Suit: "Hearts", Value: 14},
			},
			B: []Card{
				{Suit: "Diamonds", Value: 6}, {Suit: "Clubs", Value: 2},
				{Suit: "Hearts", Value: 2}, {Suit: "Spades", Value: 2},
			},
			Expected: 1,
		},
		{
			TestCaseName: "Less first",
			A: []Card{
				{Suit: "Clubs", Value: 6}, {Suit: "Spades", Value: 2},
				{Suit: "Diamonds", Value: 2}, {Suit: "Hearts", Value: 14},
			},
			B: []Card{
				{Suit: "Diamonds", Value: 10}, {Suit: "Clubs", Value: 2},
				{Suit: "Hearts", Value: 2}, {Suit: "Spades", Value: 2},
			},
			Expected: -1,
		},
		{
			TestCaseName: "Less last",
			A: []Card{
				{Suit: "Clubs", Value: 4}, {Suit: "Spades", Value: 5},
				{Suit: "Diamonds", Value: 6}, {Suit: "Hearts", Value: 2},
			},
			B: []Card{
				{Suit: "Diamonds", Value: 4}, {Suit: "Clubs", Value: 5},
				{Suit: "Hearts", Value: 6}, {Suit: "Spades", Value: 7},
			},
			Expected: -1,
		},
	}
	for _, tCase := range cases {
		t.Run(tCase.TestCaseName, func(t *testing.T) {
			res := compareCards(tCase.A, tCase.B)
			require.Equal(t, res, tCase.Expected)
		})
	}
}

func TestDeterminateWinnerGood(t *testing.T) {
	cases := []struct {
		TestCaseName   string
		CommunityCards []Card
		Players        map[string]IPlayer
		Expected       []string
		ExpectedErr    error
	}{
		{
			TestCaseName: "Different ranks",
			CommunityCards: []Card{
				{Suit: "Hearts", Value: 11}, {Suit: "Clubs", Value: 8}, {Suit: "Diamonds", Value: 7},
				{Suit: "Hearts", Value: 4}, {Suit: "Clubs", Value: 6},
			},
			Players: map[string]IPlayer{
				"first":  &Player{Hand: Hand{[2]Card{{Suit: "Spades", Value: 10}, {Suit: "Spades", Value: 12}}}},
				"second": &Player{Hand: Hand{[2]Card{{Suit: "Spades", Value: 14}, {Suit: "Spades", Value: 14}}}},
				"third":  &Player{Hand: Hand{[2]Card{{Suit: "Spades", Value: 11}, {Suit: "Spades", Value: 4}}}},
			},
			Expected: []string{"third"},
		},
		{
			TestCaseName: "Same rank, different kickers",
			CommunityCards: []Card{
				{Suit: "Hearts", Value: 10}, {Suit: "Clubs", Value: 9}, {Suit: "Diamonds", Value: 8},
				{Suit: "Hearts", Value: 7}, {Suit: "Clubs", Value: 6},
			},
			Players: map[string]IPlayer{
				"first":  &Player{Hand: Hand{[2]Card{{Suit: "Spades", Value: 11}, {Suit: "Spades", Value: 12}}}},
				"second": &Player{Hand: Hand{[2]Card{{Suit: "Spades", Value: 11}, {Suit: "Spades", Value: 10}}}},
			},
			Expected:    []string{"first"},
			ExpectedErr: nil,
		},
		{
			TestCaseName: "One player folded",
			CommunityCards: []Card{
				{Suit: "Hearts", Value: 10}, {Suit: "Clubs", Value: 9}, {Suit: "Diamonds", Value: 8},
				{Suit: "Hearts", Value: 7}, {Suit: "Clubs", Value: 6},
			},
			Players: map[string]IPlayer{
				"first":  &Player{Hand: Hand{[2]Card{{Suit: "Spades", Value: 11}, {Suit: "Spades", Value: 12}}}, IsFold: true},
				"second": &Player{Hand: Hand{[2]Card{{Suit: "Spades", Value: 10}, {Suit: "Spades", Value: 9}}}},
			},
			Expected:    []string{"second"},
			ExpectedErr: nil,
		},
		{
			TestCaseName: "Same cards",
			CommunityCards: []Card{
				{Suit: "Hearts", Value: 14}, {Suit: "Clubs", Value: 13}, {Suit: "Diamonds", Value: 12},
				{Suit: "Hearts", Value: 11}, {Suit: "Clubs", Value: 10},
			},
			Players: map[string]IPlayer{
				"first":  &Player{Hand: Hand{[2]Card{{Suit: "Spades", Value: 11}, {Suit: "Spades", Value: 12}}}},
				"second": &Player{Hand: Hand{[2]Card{{Suit: "Spades", Value: 10}, {Suit: "Spades", Value: 9}}}},
			},
			Expected:    []string{"first", "second"},
			ExpectedErr: nil,
		},
		{
			TestCaseName: "Same cards in hands",
			CommunityCards: []Card{
				{Suit: "Spades", Value: 14}, {Suit: "Spades", Value: 13}, {Suit: "Spades", Value: 12},
				{Suit: "Hearts", Value: 11}, {Suit: "Clubs", Value: 10},
			},
			Players: map[string]IPlayer{
				"first":  &Player{Hand: Hand{[2]Card{{Suit: "Spades", Value: 5}, {Suit: "Spades", Value: 4}}}},
				"second": &Player{Hand: Hand{[2]Card{{Suit: "Spades", Value: 10}, {Suit: "Spades", Value: 9}}}},
			},
			Expected:    []string{"first", "second"},
			ExpectedErr: nil,
		},
		{
			TestCaseName: "Not enough community cards",
			CommunityCards: []Card{
				{Suit: "Hearts", Value: 10}, {Suit: "Clubs", Value: 9}, {Suit: "Diamonds", Value: 8},
			},
			Players: map[string]IPlayer{
				"first":  &Player{Hand: Hand{[2]Card{{Suit: "Spades", Value: 11}, {Suit: "Spades", Value: 12}}}},
				"second": &Player{Hand: Hand{[2]Card{{Suit: "Spades", Value: 10}, {Suit: "Spades", Value: 9}}}},
			},
			Expected:    []string{},
			ExpectedErr: ErrNotEnoughCommunityCards,
		},
		{
			TestCaseName: "Empty players map",
			CommunityCards: []Card{
				{Suit: "Hearts", Value: 10}, {Suit: "Clubs", Value: 9}, {Suit: "Diamonds", Value: 8},
				{Suit: "Hearts", Value: 7}, {Suit: "Clubs", Value: 6},
			},
			Players:     map[string]IPlayer{},
			Expected:    []string{},
			ExpectedErr: ErrEmptyPlayersMap,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.TestCaseName, func(t *testing.T) {
			res, err := DeterminateWinner(tCase.CommunityCards, tCase.Players)
			require.ElementsMatch(t, res, tCase.Expected)
			require.Equal(t, err, tCase.ExpectedErr)
		})
	}
}

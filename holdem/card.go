package holdem

import "fmt"

var NameFromValue = map[int]string{
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "10",
	11: "Jack",
	12: "Queen",
	13: "King",
	14: "Ace",
}

type Card struct {
	Suit  string
	Value int
}

func (c *Card) String() string {
	return fmt.Sprintf("%s %s", string(c.Suit[0]), NameFromValue[c.Value])
}

func GetStandardDeck() []Card {
	standardDeck := []Card{
		Card{Suit: "Spades", Value: 2}, Card{Suit: "Hearts", Value: 2},
		Card{Suit: "Diamonds", Value: 2}, Card{Suit: "Clubs", Value: 2},
		Card{Suit: "Spades", Value: 3}, Card{Suit: "Hearts", Value: 3},
		Card{Suit: "Diamonds", Value: 3}, Card{Suit: "Clubs", Value: 3},
		Card{Suit: "Spades", Value: 4}, Card{Suit: "Hearts", Value: 4},
		Card{Suit: "Diamonds", Value: 4}, Card{Suit: "Clubs", Value: 4},
		Card{Suit: "Spades", Value: 5}, Card{Suit: "Hearts", Value: 5},
		Card{Suit: "Diamonds", Value: 5}, Card{Suit: "Clubs", Value: 5},
		Card{Suit: "Spades", Value: 6}, Card{Suit: "Hearts", Value: 6},
		Card{Suit: "Diamonds", Value: 6}, Card{Suit: "Clubs", Value: 6},
		Card{Suit: "Spades", Value: 7}, Card{Suit: "Hearts", Value: 7},
		Card{Suit: "Diamonds", Value: 7}, Card{Suit: "Clubs", Value: 7},
		Card{Suit: "Spades", Value: 8}, Card{Suit: "Hearts", Value: 8},
		Card{Suit: "Diamonds", Value: 8}, Card{Suit: "Clubs", Value: 8},
		Card{Suit: "Spades", Value: 9}, Card{Suit: "Hearts", Value: 9},
		Card{Suit: "Diamonds", Value: 9}, Card{Suit: "Clubs", Value: 9},
		Card{Suit: "Spades", Value: 10}, Card{Suit: "Hearts", Value: 10},
		Card{Suit: "Diamonds", Value: 10}, Card{Suit: "Clubs", Value: 10},
		Card{Suit: "Spades", Value: 11}, Card{Suit: "Hearts", Value: 11},
		Card{Suit: "Diamonds", Value: 11}, Card{Suit: "Clubs", Value: 11},
		Card{Suit: "Spades", Value: 12}, Card{Suit: "Hearts", Value: 12},
		Card{Suit: "Diamonds", Value: 12}, Card{Suit: "Clubs", Value: 12},
		Card{Suit: "Spades", Value: 13}, Card{Suit: "Hearts", Value: 13},
		Card{Suit: "Diamonds", Value: 13}, Card{Suit: "Clubs", Value: 13},
		Card{Suit: "Spades", Value: 14}, Card{Suit: "Hearts", Value: 14},
		Card{Suit: "Diamonds", Value: 14}, Card{Suit: "Clubs", Value: 14},
	}
	return standardDeck
}

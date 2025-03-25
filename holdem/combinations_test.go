package holdem

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvaluateHandGood(t *testing.T) {
	cases := []struct {
		TestCaseName        string
		PlayerHand          []Card
		CommunityCards      []Card
		ExpectedCombination Combination
	}{
		//High_Card_test
		{
			TestCaseName: "Higher_Card_Test",
			PlayerHand: []Card{
				Card{Suit: "Spades", Value: 14},
				Card{Suit: "Spades", Value: 3},
			},
			CommunityCards: []Card{
				Card{Suit: "Clubs", Value: 11},
				Card{Suit: "Clubs", Value: 10},
				Card{Suit: "Diamonds", Value: 13},
				Card{Suit: "Spades", Value: 4},
				Card{Suit: "Spades", Value: 6},
			},
			ExpectedCombination: Combination{
				Rank: HighCard,
				CompareCards: []Card{
					Card{Suit: "Spades", Value: 14},
					Card{Suit: "Diamonds", Value: 13},
					Card{Suit: "Clubs", Value: 11},
					Card{Suit: "Clubs", Value: 10},
					Card{Suit: "Spades", Value: 6},
				},
			},
		},
		//One_Pair_Test
		{
			TestCaseName: "One_Pair_Test",
			PlayerHand: []Card{
				Card{Suit: "Spades", Value: 14},
				Card{Suit: "Clubs", Value: 4},
			},
			CommunityCards: []Card{
				Card{Suit: "Spades", Value: 14},
				Card{Suit: "Clubs", Value: 3},
				Card{Suit: "Diamonds", Value: 13},
				Card{Suit: "Spades", Value: 8},
				Card{Suit: "Spades", Value: 6},
			},
			ExpectedCombination: Combination{
				Rank: OnePair,
				CompareCards: []Card{
					Card{Suit: "Spades", Value: 14},
					Card{Suit: "Spades", Value: 14},
					Card{Suit: "Diamonds", Value: 13},
					Card{Suit: "Spades", Value: 8},
					Card{Suit: "Spades", Value: 6},
				},
			},
		},
		//Several_Two_Pair_Test
		{
			TestCaseName: "Several_Two_Pair_Test",
			PlayerHand: []Card{
				Card{Suit: "Spades", Value: 4},
				Card{Suit: "Spades", Value: 4},
			},
			CommunityCards: []Card{
				Card{Suit: "Diamonds", Value: 14},
				Card{Suit: "Diamonds", Value: 13},
				Card{Suit: "Diamonds", Value: 14},
				Card{Suit: "Spades", Value: 6},
				Card{Suit: "Spades", Value: 6},
			},
			ExpectedCombination: Combination{
				Rank: TwoPairs,
				CompareCards: []Card{
					Card{Suit: "Diamonds", Value: 14},
					Card{Suit: "Diamonds", Value: 14},
					Card{Suit: "Spades", Value: 6},
					Card{Suit: "Spades", Value: 6},
					Card{Suit: "Diamonds", Value: 13},
				},
			},
		},
		//Two_Pair_Test
		{
			TestCaseName: "Two_Pair_Test",
			PlayerHand: []Card{
				Card{Suit: "Spades", Value: 6},
				Card{Suit: "Clubs", Value: 4},
			},
			CommunityCards: []Card{
				Card{Suit: "Spades", Value: 6},
				Card{Suit: "Clubs", Value: 4},
				Card{Suit: "Diamonds", Value: 13},
				Card{Suit: "Spades", Value: 8},
				Card{Suit: "Spades", Value: 5},
			},
			ExpectedCombination: Combination{
				Rank: TwoPairs,
				CompareCards: []Card{
					Card{Suit: "Spades", Value: 6},
					Card{Suit: "Spades", Value: 6},
					Card{Suit: "Clubs", Value: 4},
					Card{Suit: "Clubs", Value: 4},
					Card{Suit: "Diamonds", Value: 13},
				},
			},
		},
		//Three_Of_A_Kind_Test
		{
			TestCaseName: "Three_Of_A_Kind_Test",
			PlayerHand: []Card{
				Card{Suit: "Diamonds", Value: 2},
				Card{Suit: "Diamonds", Value: 2},
			},
			CommunityCards: []Card{
				Card{Suit: "Spades", Value: 6},
				Card{Suit: "Clubs", Value: 4},
				Card{Suit: "Diamonds", Value: 2},
				Card{Suit: "Spades", Value: 8},
				Card{Suit: "Spades", Value: 5},
			},
			ExpectedCombination: Combination{
				Rank: ThreeOfAKind,
				CompareCards: []Card{
					Card{Suit: "Diamonds", Value: 2},
					Card{Suit: "Diamonds", Value: 2},
					Card{Suit: "Diamonds", Value: 2},
					Card{Suit: "Spades", Value: 8},
					Card{Suit: "Spades", Value: 6},
				},
			},
		},

		//A-5_Straight_Test
		{
			TestCaseName: "A-5_Straight_Test",
			PlayerHand: []Card{
				Card{Suit: "Diamonds", Value: 14},
				Card{Suit: "Diamonds", Value: 2},
			},
			CommunityCards: []Card{
				Card{Suit: "Spades", Value: 3},
				Card{Suit: "Clubs", Value: 4},
				Card{Suit: "Diamonds", Value: 5},
				Card{Suit: "Spades", Value: 10},
				Card{Suit: "Spades", Value: 7},
			},
			ExpectedCombination: Combination{
				Rank: Straight,
				CompareCards: []Card{
					Card{Suit: "Diamonds", Value: 5},
				},
			},
		},
		//2-7_Straight_Test
		{
			TestCaseName: "2-7_Straight_Test",
			PlayerHand: []Card{
				Card{Suit: "Diamonds", Value: 14},
				Card{Suit: "Spades", Value: 3},
			},
			CommunityCards: []Card{
				Card{Suit: "Diamonds", Value: 2},
				Card{Suit: "Clubs", Value: 4},
				Card{Suit: "Diamonds", Value: 5},
				Card{Suit: "Spades", Value: 6},
				Card{Suit: "Spades", Value: 7},
			},
			ExpectedCombination: Combination{
				Rank: Straight,
				CompareCards: []Card{
					Card{Suit: "Spades", Value: 7},
				},
			},
		},
		//10-A_Straight_Test
		{
			TestCaseName: "10-A_Straight_Test",
			PlayerHand: []Card{
				Card{Suit: "Diamonds", Value: 14},
				Card{Suit: "Spades", Value: 13},
			},
			CommunityCards: []Card{
				Card{Suit: "Diamonds", Value: 12},
				Card{Suit: "Clubs", Value: 11},
				Card{Suit: "Diamonds", Value: 10},
				Card{Suit: "Spades", Value: 6},
				Card{Suit: "Spades", Value: 7},
			},
			ExpectedCombination: Combination{
				Rank: Straight,
				CompareCards: []Card{
					Card{Suit: "Diamonds", Value: 14},
				},
			},
		},
		//Flush_Test
		{
			TestCaseName: "Flush_Test",
			PlayerHand: []Card{
				Card{Suit: "Diamonds", Value: 6},
				Card{Suit: "Spades", Value: 13},
			},
			CommunityCards: []Card{
				Card{Suit: "Diamonds", Value: 12},
				Card{Suit: "Diamonds", Value: 8},
				Card{Suit: "Diamonds", Value: 10},
				Card{Suit: "Diamonds", Value: 6},
				Card{Suit: "Spades", Value: 7},
			},
			ExpectedCombination: Combination{
				Rank: Flush,
				CompareCards: []Card{
					Card{Suit: "Diamonds", Value: 12},
				},
			},
		},
		//Full_House_Test
		{
			TestCaseName: "Full_House_Test",
			PlayerHand: []Card{
				Card{Suit: "Diamonds", Value: 6},
				Card{Suit: "Diamonds", Value: 6},
			},
			CommunityCards: []Card{
				Card{Suit: "Diamonds", Value: 6},
				Card{Suit: "Spades", Value: 14},
				Card{Suit: "Clubs", Value: 10},
				Card{Suit: "Clubs", Value: 11},
				Card{Suit: "Spades", Value: 14},
			},
			ExpectedCombination: Combination{
				Rank: FullHouse,
				CompareCards: []Card{
					Card{Suit: "Diamonds", Value: 6},
					Card{Suit: "Diamonds", Value: 6},
					Card{Suit: "Diamonds", Value: 6},
					Card{Suit: "Spades", Value: 14},
					Card{Suit: "Spades", Value: 14},
				},
			},
		},
		//Few_Full_House_Test
		{
			TestCaseName: "Few_Full_House_Test",
			PlayerHand: []Card{
				Card{Suit: "Clubs", Value: 14},
				Card{Suit: "Hearts", Value: 13},
			},
			CommunityCards: []Card{
				Card{Suit: "Clubs", Value: 14},
				Card{Suit: "Clubs", Value: 14},
				Card{Suit: "Hearts", Value: 13},
				Card{Suit: "Hearts", Value: 13},
				Card{Suit: "Diamonds", Value: 7},
			},
			ExpectedCombination: Combination{
				Rank: FullHouse,
				CompareCards: []Card{
					Card{Suit: "Clubs", Value: 14},
					Card{Suit: "Clubs", Value: 14},
					Card{Suit: "Clubs", Value: 14},
					Card{Suit: "Hearts", Value: 13},
					Card{Suit: "Hearts", Value: 13},
				},
			},
		},
		//Four_Of_A_Kind_Test
		{
			TestCaseName: "Four_Of_A_Kind_Test",
			PlayerHand: []Card{
				Card{Suit: "Diamonds", Value: 6},
				Card{Suit: "Diamonds", Value: 6},
			},
			CommunityCards: []Card{
				Card{Suit: "Diamonds", Value: 6},
				Card{Suit: "Diamonds", Value: 6},
				Card{Suit: "Spades", Value: 14},
				Card{Suit: "Clubs", Value: 11},
				Card{Suit: "Spades", Value: 14},
			},
			ExpectedCombination: Combination{
				Rank: FourOfAKind,
				CompareCards: []Card{
					Card{Suit: "Diamonds", Value: 6},
					Card{Suit: "Diamonds", Value: 6},
					Card{Suit: "Diamonds", Value: 6},
					Card{Suit: "Diamonds", Value: 6},
					Card{Suit: "Spades", Value: 14},
				},
			},
		},
		//A-5_Straight_Flush_Test
		{
			TestCaseName: "A-5_Straight_Flush_Test",
			PlayerHand: []Card{
				Card{Suit: "Diamonds", Value: 14},
				Card{Suit: "Diamonds", Value: 2},
			},
			CommunityCards: []Card{
				Card{Suit: "Diamonds", Value: 3},
				Card{Suit: "Diamonds", Value: 4},
				Card{Suit: "Diamonds", Value: 5},
				Card{Suit: "Diamonds", Value: 10},
				Card{Suit: "Diamonds", Value: 7},
			},
			ExpectedCombination: Combination{
				Rank: StraightFlush,
				CompareCards: []Card{
					Card{Suit: "Diamonds", Value: 5},
				},
			},
		},
		//2-7_Straight_Flush_Test
		{
			TestCaseName: "2-7_Straight_Test",
			PlayerHand: []Card{
				Card{Suit: "Diamonds", Value: 14},
				Card{Suit: "Diamonds", Value: 3},
			},
			CommunityCards: []Card{
				Card{Suit: "Diamonds", Value: 2},
				Card{Suit: "Diamonds", Value: 4},
				Card{Suit: "Diamonds", Value: 5},
				Card{Suit: "Diamonds", Value: 6},
				Card{Suit: "Diamonds", Value: 7},
			},
			ExpectedCombination: Combination{
				Rank: StraightFlush,
				CompareCards: []Card{
					Card{Suit: "Diamonds", Value: 7},
				},
			},
		},
		//9-K_Straight_Flush_Test
		{
			TestCaseName: "9-K_Straight_Flush_Test",
			PlayerHand: []Card{
				Card{Suit: "Spades", Value: 9},
				Card{Suit: "Spades", Value: 13},
			},
			CommunityCards: []Card{
				Card{Suit: "Spades", Value: 12},
				Card{Suit: "Spades", Value: 11},
				Card{Suit: "Spades", Value: 10},
				Card{Suit: "Spades", Value: 6},
				Card{Suit: "Spades", Value: 7},
			},
			ExpectedCombination: Combination{
				Rank: StraightFlush,
				CompareCards: []Card{
					Card{Suit: "Spades", Value: 13},
				},
			},
		},
		//Royal_Flush_Test
		{
			TestCaseName: "Royal_Flush_Test",
			PlayerHand: []Card{
				Card{Suit: "Hearts", Value: 14},
				Card{Suit: "Hearts", Value: 13},
			},
			CommunityCards: []Card{
				Card{Suit: "Hearts", Value: 12},
				Card{Suit: "Hearts", Value: 11},
				Card{Suit: "Hearts", Value: 10},
				Card{Suit: "Diamonds", Value: 6},
				Card{Suit: "Diamonds", Value: 7},
			},
			ExpectedCombination: Combination{
				Rank: RoyalFlush,
				CompareCards: []Card{
					Card{Suit: "Hearts", Value: 14},
				},
			},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.TestCaseName, func(t *testing.T) {
			comb := EvaluateHand(tCase.PlayerHand, tCase.CommunityCards)
			require.Equal(t, comb, tCase.ExpectedCombination)
		})
	}
}

package holdem

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCardString(t *testing.T) {
	cases := []struct {
		TestCaseName   string
		Data           Card
		ExpectedString string
	}{
		{
			TestCaseName:   "S 2",
			Data:           Card{Suit: "Spades", Value: 2},
			ExpectedString: "S 2",
		},
		{
			TestCaseName:   "D 6",
			Data:           Card{Suit: "Diamonds", Value: 6},
			ExpectedString: "D 6",
		},
		{
			TestCaseName:   "C 10",
			Data:           Card{Suit: "Clubs", Value: 10},
			ExpectedString: "C 10",
		},
		{
			TestCaseName:   "H Jack",
			Data:           Card{Suit: "Hearts", Value: 11},
			ExpectedString: "H Jack",
		},
	}
	for _, tCase := range cases {
		t.Run(tCase.TestCaseName, func(t *testing.T) {
			str := tCase.Data.String()
			require.Equal(t, str, tCase.ExpectedString)
		})
	}
}

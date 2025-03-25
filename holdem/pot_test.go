package holdem

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPot(t *testing.T) {
	cases := []struct {
		TestCaseName string
		Data         map[string]IPlayer
		Expected     []Pot
	}{
		{
			TestCaseName: "Empty Players",
			Data:         map[string]IPlayer{},
			Expected:     []Pot{},
		},
		{
			TestCaseName: "One pot",
			Data: map[string]IPlayer{
				"1": &Player{
					LastBet: 500,
					IsFold:  false,
				},
				"2": &Player{
					LastBet: 500,
					IsFold:  false,
				},
				"3": &Player{
					LastBet: 500,
					IsFold:  false,
				},
			},
			Expected: []Pot{
				Pot{Amount: 1500, Applicants: []string{"1", "2", "3"}},
			},
		},
		{
			TestCaseName: "One pot with fold",
			Data: map[string]IPlayer{
				"1": &Player{
					LastBet: 500,
					IsFold:  false,
				},
				"2": &Player{
					LastBet: 500,
					IsFold:  false,
				},
				"3": &Player{
					LastBet: 500,
					IsFold:  true,
				},
			},
			Expected: []Pot{
				Pot{Amount: 1000, Applicants: []string{"1", "2"}},
			},
		},
		{
			TestCaseName: "two pots",
			Data: map[string]IPlayer{
				"1": &Player{
					LastBet: 500,
					IsFold:  false,
				},
				"2": &Player{
					LastBet: 400,
					IsFold:  false,
				},
			},
			Expected: []Pot{
				Pot{Amount: 800, Applicants: []string{"1", "2"}},
				Pot{Amount: 100, Applicants: []string{"1"}},
			},
		},
		{
			TestCaseName: "several pots",
			Data: map[string]IPlayer{
				"1": &Player{
					LastBet: 500,
					IsFold:  false,
				},
				"2": &Player{
					LastBet: 400,
					IsFold:  false,
				},
				"3": &Player{
					LastBet: 300,
					IsFold:  false,
				},
			},
			Expected: []Pot{
				Pot{Amount: 900, Applicants: []string{"1", "2", "3"}},
				Pot{Amount: 200, Applicants: []string{"1", "2"}},
				Pot{Amount: 100, Applicants: []string{"1"}},
			},
		},
		{
			TestCaseName: "several pots with fold",
			Data: map[string]IPlayer{
				"1": &Player{
					LastBet: 500,
					IsFold:  true,
				},
				"2": &Player{
					LastBet: 400,
					IsFold:  false,
				},
				"3": &Player{
					LastBet: 300,
					IsFold:  false,
				},
			},
			Expected: []Pot{
				Pot{Amount: 600, Applicants: []string{"2", "3"}},
				Pot{Amount: 100, Applicants: []string{"2"}},
			},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.TestCaseName,
			func(t *testing.T) {
				res := CreatePots(tCase.Data)
				for k, _ := range res {
					require.ElementsMatch(t, res[k].Applicants, tCase.Expected[k].Applicants)
					require.Equal(t, res[k].Amount, tCase.Expected[k].Amount)
				}

			},
		)
	}
}

func TestUnionPots(t *testing.T) {
	cases := []struct {
		TestCaseName string
		Data         []Pot
		Expected     []Pot
	}{
		{
			TestCaseName: "No pots",
			Data:         []Pot{},
			Expected:     []Pot{},
		},
		{
			TestCaseName: "1 -> 1 pot",
			Data: []Pot{
				{Amount: 100, Applicants: []string{"1", "2"}},
			},
			Expected: []Pot{
				{Amount: 100, Applicants: []string{"1", "2"}},
			},
		},
		{
			TestCaseName: "2 -> 1 pot",
			Data: []Pot{
				{Amount: 100, Applicants: []string{"1", "2"}},
				{Amount: 200, Applicants: []string{"1", "2"}},
			},
			Expected: []Pot{
				{Amount: 300, Applicants: []string{"1", "2"}},
			},
		},
		{
			TestCaseName: "5 -> 2 pots",
			Data: []Pot{
				{Amount: 150, Applicants: []string{"1", "2"}},
				{Amount: 800, Applicants: []string{"1", "2"}},
				{Amount: 500, Applicants: []string{"1", "2", "3"}},
			},
			Expected: []Pot{
				{Amount: 950, Applicants: []string{"1", "2"}},
				{Amount: 500, Applicants: []string{"1", "2", "3"}},
			},
		},
		{
			TestCaseName: "3 -> 3 pots",
			Data: []Pot{
				{Amount: 150, Applicants: []string{"1"}},
				{Amount: 800, Applicants: []string{"1", "2"}},
				{Amount: 500, Applicants: []string{"1", "2", "3"}},
			},
			Expected: []Pot{
				{Amount: 150, Applicants: []string{"1"}},
				{Amount: 800, Applicants: []string{"1", "2"}},
				{Amount: 500, Applicants: []string{"1", "2", "3"}},
			},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.TestCaseName,
			func(t *testing.T) {
				res := UnionPots(tCase.Data)
				require.ElementsMatch(t, res, tCase.Expected)
			},
		)
	}
}

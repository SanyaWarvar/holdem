package holdem

import (
	"errors"
)

var (
	ErrEmptyPlayersMap         = errors.New("empty players map")
	ErrNotEnoughCommunityCards = errors.New("len of community cards must be 5")
	ErrNotEnoughCardsInHand    = errors.New("len of player cards must be 2") //TODO add
)

func DeterminateWinner(communityCards []Card, players map[string]IPlayer) ([]string, error) {
	if len(players) == 0 {
		return []string{}, ErrEmptyPlayersMap
	}

	if len(communityCards) != 5 {
		return []string{}, ErrNotEnoughCommunityCards
	}

	bestPlayers := make([]string, len(players))
	var bestCombination Combination

	for id, player := range players {
		hand := player.GetHand()
		if player.GetFold() {
			continue
		}
		combination := EvaluateHand(hand.Cards[:], communityCards)
		if combination.Rank > bestCombination.Rank ||
			(combination.Rank == bestCombination.Rank && compareCards(combination.CompareCards, bestCombination.CompareCards) > 0) {
			bestPlayers = append(bestPlayers[:0], id)
			bestCombination = combination
		} else if combination.Rank == bestCombination.Rank && compareCards(combination.CompareCards, bestCombination.CompareCards) == 0 {
			bestPlayers = append(bestPlayers, id)
		}
	}

	return bestPlayers, nil
}

// Compare two combinations (kickers)
// return 1 if a > b
// return -1 if a < b
// return 0 if a == b
func compareCards(a, b []Card) int {
	if len(a) == 0 && len(b) != 0 {
		return -1
	}
	if len(a) != 0 && len(b) == 0 {
		return 1
	}
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i].Value > b[i].Value {
			return 1
		} else if a[i].Value < b[i].Value {
			return -1
		}
	}
	return 0
}

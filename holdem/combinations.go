package holdem

import (
	"sort"
)

type Combination struct {
	Rank         int
	CompareCards []Card
}

const (
	HighCard = iota + 1
	OnePair
	TwoPairs
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

// EvaluateHand.
// Функция для определения комбинации из двух карт игрока (параметр playerHand) и пяти карт на столке (параметр communityCards)
func EvaluateHand(playerHand []Card, communityCards []Card) Combination {
	allCards := append(playerHand, communityCards...)
	sort.Slice(allCards, func(i, j int) bool {
		return allCards[i].Value > allCards[j].Value
	})

	flushCards := checkFlush(allCards)
	if len(flushCards) >= 5 {
		straightFlushCards := checkStraight(flushCards)
		if len(straightFlushCards) >= 5 {
			if straightFlushCards[0].Value == 14 {
				return Combination{Rank: RoyalFlush, CompareCards: straightFlushCards[:1]} // Роял-флеш
			}
			return Combination{Rank: StraightFlush, CompareCards: straightFlushCards[:1]} // Стрит-флеш
		}
		return Combination{Rank: Flush, CompareCards: flushCards[:1]} // Флеш
	}

	straightCards := checkStraight(allCards)
	if len(straightCards) >= 5 {
		return Combination{Rank: Straight, CompareCards: straightCards[:1]} // Стрит
	}

	valueCounts := make(map[int]int)
	for _, card := range allCards {
		valueCounts[card.Value]++
	}

	var pairs, threes, fours []int
	for value, count := range valueCounts {
		switch count {
		case 2:
			pairs = append(pairs, value)
		case 3:
			threes = append(threes, value)
		case 4:
			fours = append(fours, value)
		}
	}

	sort.Slice(pairs, func(i, j int) bool { return pairs[i] > pairs[j] })
	sort.Slice(threes, func(i, j int) bool { return threes[i] > threes[j] })

	if len(fours) > 0 {
		fourCards := getCardsByValue(allCards, fours[0])
		kicker := getKickers(allCards, fourCards, 1)
		return Combination{Rank: FourOfAKind, CompareCards: append(fourCards, kicker...)} // Каре
	}

	if len(threes) >= 2 {
		threeCards := getCardsByValue(allCards, threes[0])
		kicker := getKickers(allCards, threeCards, 2)
		return Combination{Rank: FullHouse, CompareCards: append(threeCards, kicker...)} // Фулл-хаус
	}

	if len(threes) >= 1 && len(pairs) >= 1 {
		threeCards := getCardsByValue(allCards, threes[0])
		pairCards := getCardsByValue(allCards, pairs[0])
		return Combination{Rank: FullHouse, CompareCards: append(threeCards, pairCards...)} // Фулл-хаус
	}

	if len(threes) > 0 {
		threeCards := getCardsByValue(allCards, threes[0])
		kickers := getKickers(allCards, threeCards, 2)
		return Combination{Rank: ThreeOfAKind, CompareCards: append(threeCards, kickers...)} // Сет
	}

	if len(pairs) >= 2 {
		pair1Cards := getCardsByValue(allCards, pairs[0])
		pair2Cards := getCardsByValue(allCards, pairs[1])
		kicker := getKickers(allCards, append(pair1Cards, pair2Cards...), 1)
		return Combination{Rank: TwoPairs, CompareCards: append(append(pair1Cards, pair2Cards...), kicker...)} // Две пары
	}

	if len(pairs) > 0 {
		pairCards := getCardsByValue(allCards, pairs[0])
		kickers := getKickers(allCards, pairCards, 3)
		return Combination{Rank: OnePair, CompareCards: append(pairCards, kickers...)} // Пара
	}

	return Combination{Rank: HighCard, CompareCards: allCards[:5]} // Старшая карта
}

// Helper functions (unchanged)
func checkFlush(cards []Card) []Card {
	suitCounts := make(map[string][]Card)
	for _, card := range cards {
		suitCounts[card.Suit] = append(suitCounts[card.Suit], card)
	}
	for _, flushCards := range suitCounts {
		if len(flushCards) >= 5 {
			sort.Slice(flushCards, func(i, j int) bool {
				return flushCards[i].Value > flushCards[j].Value
			})
			return flushCards
		}
	}
	return nil
}

func checkStraight(cards []Card) []Card {
	uniqueValues := make(map[int]bool)
	var uniqueCards []Card
	for _, card := range cards {
		if !uniqueValues[card.Value] {
			uniqueValues[card.Value] = true
			uniqueCards = append(uniqueCards, card)
		}
	}
	if len(uniqueCards) < 5 {
		return nil
	}
	sort.Slice(uniqueCards, func(i, j int) bool {
		return uniqueCards[i].Value > uniqueCards[j].Value
	})

	for i := 0; i <= len(uniqueCards)-5; i++ {
		if uniqueCards[i].Value == uniqueCards[i+1].Value+1 &&
			uniqueCards[i+1].Value == uniqueCards[i+2].Value+1 &&
			uniqueCards[i+2].Value == uniqueCards[i+3].Value+1 &&
			uniqueCards[i+3].Value == uniqueCards[i+4].Value+1 {
			return uniqueCards[i : i+5]
		}
	}

	hasAce := false
	aceSuit := ""
	var lowStraightCards []Card
	for _, card := range uniqueCards {
		if card.Value == 14 {
			hasAce = true
			aceSuit = card.Suit
		}
		if card.Value == 5 || card.Value == 4 || card.Value == 3 || card.Value == 2 {
			lowStraightCards = append(lowStraightCards, card)
		}
	}
	if hasAce && len(lowStraightCards) >= 4 {
		sort.Slice(lowStraightCards, func(i, j int) bool {
			return lowStraightCards[i].Value > lowStraightCards[j].Value
		})
		return append(lowStraightCards, Card{Suit: aceSuit, Value: 14})
	}

	return nil
}

func getCardsByValue(cards []Card, value int) []Card {
	var result []Card
	for _, card := range cards {
		if card.Value == value {
			result = append(result, card)
		}
	}
	return result
}

func getKickers(cards []Card, exclude []Card, count int) []Card {
	var kickers []Card
	excludeMap := make(map[Card]bool)
	for _, card := range exclude {
		excludeMap[card] = true
	}
	for _, card := range cards {
		if !excludeMap[card] {
			kickers = append(kickers, card)
			if len(kickers) >= count {
				break
			}
		}
	}
	return kickers
}

package holdem

import "fmt"

type Pot struct {
	Amount     int
	Applicants []string
}

func CreatePots(players map[string]IPlayer) []Pot {
	pots := []Pot{}

	for {
		minBet := -1
		applicants := make([]string, 0, len(players))
		for k, v := range players { // находим минимальную ставку
			if v.GetFold() || v.GetLastBet() == 0 {
				continue
			}
			applicants = append(applicants, k)
			if minBet == -1 {
				minBet = v.GetLastBet()
			}
			minBet = min(v.GetLastBet(), minBet)
		}

		for _, k := range applicants {
			v := players[k]
			players[k].SetLastBet(v.GetLastBet() - minBet)
		}

		if minBet == -1 {
			break
		}
		pots = append(pots, Pot{Amount: len(applicants) * minBet, Applicants: applicants})
	}
	return pots
}

func UnionPots(pots []Pot) []Pot {
	merged := make(map[string]Pot)

	for _, pot := range pots {
		key := fmt.Sprintf("%v", pot.Applicants)

		if existing, ok := merged[key]; ok {
			existing.Amount += pot.Amount
			merged[key] = existing
		} else {
			merged[key] = pot
		}
	}

	result := make([]Pot, 0, len(merged))
	for _, pot := range merged {
		result = append(result, pot)
	}

	return result
}

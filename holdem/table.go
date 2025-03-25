package holdem

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"sync"
	"time"
)

var (
	ErrMaxPlayers       = errors.New("count of players reached max value")
	ErrGameStarted      = errors.New("this game already started")
	ErrGameNotStarted   = errors.New("this game not started")
	ErrNotEnoughCards   = errors.New("not enough card in deck")
	ErrNotYourTurn      = errors.New("not your turn t")
	ErrPlayerIsFold     = errors.New("player already fold his cards")
	ErrCantCheck        = errors.New("you cant check")
	ErrCantRaise        = errors.New("raise must exceed the current bet by at least two times")
	ErrNotEnoughMoney   = errors.New("not enough money for this  action")
	ErrUnexpectedAction = errors.New("unexpected action")
	ErrPlayerNotFound   = errors.New("player not found")
)

type IPokerTable interface {
	StartGame()
	AddObserver(o *IObserver)
	AddPlayer(player *IPlayer)
	RemovePlayer(player *IPlayer)
	MakeMove(playerId, action string, amount int)
}

type TableConfig struct {
	BlindIncreaseTime time.Duration
	LastBlindIncrease time.Time
	MaxPlayers        int
	MinPlayers        int
	EnterAfterStart   bool
	BankAmount        int
}

// TODO add timeout for 1 move and time bank
type TableMeta struct {
	SmallBlind     int
	Ante           int
	DealerIndex    int
	PlayerTurnInd  int
	CurrentBet     int
	CommunityCards []Card
	PlayersOrder   []string
	Players        map[string]IPlayer
	Query          map[string]IPlayer
	Pots           []Pot
	Deck           []Card
	CurrentRound   int
	GameStarted    bool
	Seed           int64
}

type PokerTable struct {
	observers []IObserver
	mu        sync.Mutex
	Config    *TableConfig
	Meta      *TableMeta
}

func NewTableConfig(BlindIncreaseTime time.Duration, maxPlayers, minPlayers, bankAmount int, enterAfteStart bool) *TableConfig {
	return &TableConfig{
		BlindIncreaseTime: BlindIncreaseTime,
		LastBlindIncrease: time.Now(),
		MaxPlayers:        maxPlayers,
		MinPlayers:        minPlayers,
		EnterAfterStart:   enterAfteStart,
		BankAmount:        bankAmount,
	}
}

func NewTableMeta(smallBlind int, ante int, seed int64) *TableMeta {
	return &TableMeta{
		SmallBlind:     smallBlind,
		Ante:           ante,
		DealerIndex:    0,
		PlayerTurnInd:  0,
		CurrentBet:     0,
		CommunityCards: []Card{},
		PlayersOrder:   make([]string, 0, 10),
		Players:        make(map[string]IPlayer),
		Query:          make(map[string]IPlayer),
		Pots:           []Pot{},
		Deck:           []Card{},
		CurrentRound:   -1,
		GameStarted:    false,
		Seed:           seed,
	}
}

func NewPokerTable(config *TableConfig, meta *TableMeta) *PokerTable {
	return &PokerTable{
		observers: []IObserver{},
		mu:        sync.Mutex{},
		Config:    config,
		Meta:      meta,
	}
}

func (m *TableMeta) refreshDeck() {
	m.Deck = GetStandardDeck()
	var r *rand.Rand
	if m.Seed != 0 {
		r = rand.New(rand.NewSource(m.Seed))
	}
	r.Shuffle(len(m.Deck), func(i, j int) {
		m.Deck[i], m.Deck[j] = m.Deck[j], m.Deck[i]
	})
}

func (m *TableMeta) addPlayerInGame(p IPlayer) {
	m.Players[p.GetId()] = p
}

//TODO remove player

func (m *TableMeta) addPlayerInQuery(p IPlayer) {
	m.Query[p.GetId()] = p
}

func (t *PokerTable) AddObserver(obs IObserver) {
	t.observers = append(t.observers, obs)
}

func (t *PokerTable) NotifyObservers(event string) {
	for _, obs := range t.observers {
		obs.Update(event)
	}
}

func (t *PokerTable) AddPlayer(p IPlayer) error {
	if t.Meta.GameStarted && !t.Config.EnterAfterStart {
		return ErrGameStarted
	}

	if t.Config.MaxPlayers <= len(t.Meta.Players)+len(t.Meta.Query)+1 {
		return ErrMaxPlayers
	}

	if t.Meta.GameStarted {
		t.Meta.addPlayerInQuery(p)
	} else {
		t.Meta.addPlayerInGame(p)
		t.Meta.PlayersOrder = append(t.Meta.PlayersOrder, p.GetId())
	}
	t.NotifyObservers(fmt.Sprintf("Player %s enter the game", p.GetId()))
	return nil
}

func (t *PokerTable) enterPlayersFromQuery() {
	for k, v := range t.Meta.Query {
		t.Meta.Players[k] = v
		t.Meta.PlayersOrder = append(t.Meta.PlayersOrder, k)
	}
}

func (t *PokerTable) StartGame() error {
	if t.Meta.GameStarted {
		return ErrGameStarted
	}
	//TODO change balance if bankamount != 0
	t.Meta.GameStarted = true
	t.Meta.CurrentRound = -1
	t.Meta.refreshDeck()
	t.NotifyObservers("Game started")
	t.NewRound()
	return nil
}

func (t *PokerTable) NewRound() error {
	if !t.Meta.GameStarted {
		return ErrGameNotStarted
	}
	t.createPots()
	t.Meta.CurrentRound += 1
	t.Meta.CurrentBet = 0
	t.NotifyObservers(fmt.Sprintf("New round started. Current round: %d", t.Meta.CurrentRound))

	refreshPlayers(t.Meta.Players, t.Meta.CurrentRound == 4)
	switch t.Meta.CurrentRound {
	case 0: //pre flop
		t.enterPlayersFromQuery()
		t.betAnte()
		for _, k := range t.Meta.PlayersOrder {
			cards, _ := t.drawCard(2)
			t.Meta.Players[k].SetHand(Hand{[2]Card{cards[0], cards[1]}})
			t.NotifyObservers(fmt.Sprintf("Player %s get cards: %v", t.Meta.Players[k].GetId(), cards))
		}
		t.choiceDealer()
		t.betBlinds()
	case 1: // flop
		t.Meta.CommunityCards, _ = t.drawCard(3)
		t.NotifyObservers(fmt.Sprintf("Community cards: %v", t.Meta.CommunityCards))
		t.Meta.PlayerTurnInd = (t.Meta.DealerIndex + 1) % len(t.Meta.PlayersOrder)

	case 2: // turn
		cards, _ := t.drawCard(1)
		t.Meta.CommunityCards = append(t.Meta.CommunityCards, cards...)
		t.NotifyObservers(fmt.Sprintf("Community cards: %v", t.Meta.CommunityCards))

	case 3: // river
		cards, _ := t.drawCard(1)
		t.Meta.CommunityCards = append(t.Meta.CommunityCards, cards...)
		t.NotifyObservers(fmt.Sprintf("Community cards: %v", t.Meta.CommunityCards))

	case 4: // determinate winner
		t.PayMoney()
		t.Meta.updateSeed()
		t.Meta.GameStarted = false
		t.Meta.CurrentRound = -1
		t.Meta.Pots = t.Meta.Pots[:0]
	}
	t.choiceFirstMovePlayer()

	return nil
}

func (m *TableMeta) updateSeed() {
	if m.Seed != 0 {
		r := rand.New(rand.NewSource(m.Seed))
		for {
			newSeed := r.Int63()
			if newSeed == 0 {
				continue
			}
			m.Seed = newSeed
			break
		}
	}
}

func (t *PokerTable) PayMoney() {
	for ind, pot := range t.Meta.Pots {
		applicants := make(map[string]IPlayer)
		for _, k := range pot.Applicants {
			p := t.Meta.Players[k]
			if p.GetFold() { // если игрок сбросил то он не претендует на банк
				continue
			}
			applicants[k] = p
		}
		winners, _ := DeterminateWinner(t.Meta.CommunityCards, applicants)
		winAmount := pot.Amount / len(winners)
		for _, winner := range winners {
			t.Meta.Players[winner].ChangeBalance(winAmount)
		}
		t.NotifyObservers(fmt.Sprintf("Winners of pot %.2d with %d amount: %v", ind+1, winAmount, winners))
		if winAmount*len(winners) == pot.Amount {
			continue
		}
		counter := pot.Amount - winAmount*len(winners)
		for i := 1; counter > 0; i++ {
			targetPlayer := t.Meta.PlayersOrder[(t.Meta.DealerIndex+i)%len(t.Meta.Players)]
			if t.Meta.Players[targetPlayer].GetFold() || !slices.Contains(winners, t.Meta.Players[targetPlayer].GetId()) {
				continue
			}
			t.Meta.Players[targetPlayer].ChangeBalance(1)
			counter--
		}
	}
}

func (t *PokerTable) createPots() error {
	if !t.Meta.GameStarted {
		return ErrGameNotStarted
	}

	pots := CreatePots(t.Meta.Players)
	t.Meta.Pots = append(t.Meta.Pots, pots...)

	return nil
}

func (t *PokerTable) RemovePlayer(playerId string) error {
	_, ok1 := t.Meta.Players[playerId]
	_, ok2 := t.Meta.Query[playerId]
	if !(ok1 || ok2) {
		return ErrPlayerNotFound
	}
	if ok2 {
		delete(t.Meta.Query, playerId)
		return nil
	}
	delete(t.Meta.Players, playerId)
	ind := slices.Index(t.Meta.PlayersOrder, playerId)
	t.Meta.PlayersOrder = append(t.Meta.PlayersOrder[:ind], t.Meta.PlayersOrder[ind+1:]...)
	return nil
}

func (t *PokerTable) betAnte() error {
	if !t.Meta.GameStarted {
		return ErrGameNotStarted
	}
	//TODO check if not 0 round
	toRemove := []string{}
	for k, v := range t.Meta.Players {
		if v.GetBalance() < t.Meta.Ante {
			v.GetFold()
			t.NotifyObservers(fmt.Sprintf("Player %s cant bet ante", k))
			toRemove = append(toRemove, k)
		}
	}
	for _, id := range toRemove {
		t.RemovePlayer(id)
	}

	t.Meta.Pots = append(t.Meta.Pots, Pot{Amount: t.Meta.Ante * len(t.Meta.Players), Applicants: t.Meta.PlayersOrder})
	t.NotifyObservers(fmt.Sprintf("Get ante: %d", t.Meta.Ante*len(t.Meta.Players)))
	return nil
}

func (t *PokerTable) betBlinds() error {
	if !t.Meta.GameStarted {
		return ErrGameNotStarted
	}
	var smallBlindPlayer, bigBlindPlayer string
	if len(t.Meta.PlayersOrder) > 2 {
		smallBlindPlayer = t.Meta.PlayersOrder[(t.Meta.DealerIndex+1)%len(t.Meta.PlayersOrder)]
		bigBlindPlayer = t.Meta.PlayersOrder[(t.Meta.DealerIndex+2)%len(t.Meta.PlayersOrder)]
	} else {
		smallBlindPlayer = t.Meta.PlayersOrder[t.Meta.DealerIndex] // дилер ставит малый блайнд в хендз апе
		bigBlindPlayer = t.Meta.PlayersOrder[(t.Meta.DealerIndex+1)%len(t.Meta.PlayersOrder)]
	}

	smallBlindPlayerBet := min(t.Meta.SmallBlind, t.Meta.Players[smallBlindPlayer].GetBalance())
	t.Meta.Players[smallBlindPlayer].ChangeBalance(-smallBlindPlayerBet)
	t.Meta.Players[smallBlindPlayer].SetLastBet(smallBlindPlayerBet)
	t.NotifyObservers(fmt.Sprintf("Player %s bet %d as small blind", smallBlindPlayer, smallBlindPlayerBet))

	bigBlindPlayerBet := min(t.Meta.SmallBlind*2, t.Meta.Players[bigBlindPlayer].GetBalance())
	t.Meta.Players[bigBlindPlayer].ChangeBalance(-bigBlindPlayerBet)
	t.Meta.Players[bigBlindPlayer].SetLastBet(bigBlindPlayerBet)
	t.NotifyObservers(fmt.Sprintf("Player %s bet %d as big blind", bigBlindPlayer, bigBlindPlayerBet))
	t.Meta.CurrentBet = max(bigBlindPlayerBet, smallBlindPlayerBet)
	return nil
}

func (t *PokerTable) getNextPlayer() {
	for i := 1; i < len(t.Meta.PlayersOrder); i++ {
		nextIndex := (t.Meta.PlayerTurnInd + i) % len(t.Meta.PlayersOrder)
		nextPlayer := t.Meta.PlayersOrder[nextIndex]
		if !t.Meta.Players[nextPlayer].GetFold() && !t.Meta.Players[nextPlayer].GetReadyStatus() {
			t.Meta.PlayerTurnInd = nextIndex
			t.NotifyObservers(fmt.Sprintf("Next move expect from %s player", nextPlayer))
			return
		}
	}
}

func (t *PokerTable) choiceDealer() error {
	if !t.Meta.GameStarted {
		return ErrGameNotStarted
	}
	t.Meta.DealerIndex = (t.Meta.DealerIndex + 1) % len(t.Meta.PlayersOrder)
	t.NotifyObservers(fmt.Sprintf("dealer is %s", t.Meta.PlayersOrder[t.Meta.DealerIndex]))
	return nil
}

func (t *PokerTable) drawCard(n int) ([]Card, error) {
	output := make([]Card, 0, n)
	if len(t.Meta.Deck) < n {
		return output, ErrNotEnoughCards
	}
	output = append(output, t.Meta.Deck[:n]...)
	t.Meta.Deck = t.Meta.Deck[n:]
	return output, nil
}

func (t *PokerTable) choiceFirstMovePlayer() error {
	if !t.Meta.GameStarted {
		return ErrGameNotStarted
	}
	if t.Meta.CurrentRound == 0 { //utg
		t.Meta.PlayerTurnInd = (t.Meta.DealerIndex + 3) % len(t.Meta.PlayersOrder)
	} else {
		t.Meta.PlayerTurnInd = (t.Meta.DealerIndex + 1) % len(t.Meta.PlayersOrder)
	}
	return nil
}

func (t *PokerTable) MakeMove(playerId, action string, amount int) error {
	if !t.Meta.GameStarted {
		return ErrGameNotStarted
	}

	if t.Meta.PlayersOrder[t.Meta.PlayerTurnInd] != playerId {
		return ErrNotYourTurn
	}

	if t.Meta.Players[playerId].GetFold() {
		return ErrPlayerIsFold
	}

	switch action {
	case "check":
		t.handleCheck(playerId)
	case "raise":
		t.handleRaise(playerId, amount)
	case "call":
		t.handleCall(playerId)
	case "fold":
		t.handleFold(playerId)
	default:
		return ErrUnexpectedAction
	}
	t.Meta.Players[playerId].SetStatus(true)
	t.getNextPlayer()
	if t.checkReady() {
		t.NewRound()
	} else {
		t.notifyNext()
	}
	return nil
}

func (t *PokerTable) notifyNext() error {
	if !t.Meta.GameStarted {
		return ErrGameNotStarted
	}
	pId := t.Meta.PlayersOrder[t.Meta.PlayerTurnInd]
	if t.Meta.CurrentBet != 0 {
		t.NotifyObservers(fmt.Sprintf("player %s can do call with %d", pId, t.Meta.CurrentBet))
	} else {
		t.NotifyObservers(fmt.Sprintf("player %s can do check", pId))
	}
	return nil
}

func (t *PokerTable) checkReady() bool {
	if !t.Meta.GameStarted {
		return false
	}
	for _, v := range t.Meta.Players {
		if (!v.GetFold() && !v.GetReadyStatus()) || v.GetBalance() == 0 {
			return false
		}
	}
	return true
}

func (t *PokerTable) handleCheck(playerId string) error {
	if !t.Meta.GameStarted {
		return ErrGameNotStarted
	}
	if t.Meta.Players[playerId].GetFold() {
		return ErrPlayerIsFold
	}

	if t.Meta.CurrentBet != 0 {
		return ErrCantCheck
	}
	t.Meta.Players[playerId].SetStatus(true)
	t.NotifyObservers(fmt.Sprintf("Player %s do check", playerId))
	return nil
}

func (t *PokerTable) handleFold(playerId string) error {
	if !t.Meta.GameStarted {
		return ErrGameNotStarted
	}

	t.Meta.Players[playerId].SetStatus(true)
	t.Meta.Players[playerId].SetFold(true)
	t.NotifyObservers(fmt.Sprintf("Player %s do fold", playerId))
	return nil
}

func (t *PokerTable) handleRaise(playerId string, amount int) error {
	if !t.Meta.GameStarted {
		return ErrGameNotStarted
	}
	if t.Meta.Players[playerId].GetFold() {
		return ErrPlayerIsFold
	}
	if !(amount > t.Meta.CurrentBet*2 && amount > t.Meta.Players[playerId].GetLastBet() && amount > 0) {
		return ErrCantRaise
	}
	delta := amount - t.Meta.Players[playerId].GetLastBet()
	if delta > t.Meta.Players[playerId].GetBalance() {
		return ErrNotEnoughMoney
	}
	t.resetPlayersStatus()
	t.Meta.Players[playerId].SetLastBet(amount)
	t.Meta.Players[playerId].ChangeBalance(-delta)
	t.Meta.Players[playerId].SetStatus(true)
	t.Meta.CurrentBet = amount

	t.NotifyObservers(fmt.Sprintf("Player %s do raise with %d amount", playerId, amount))
	return nil
}

func (t *PokerTable) resetPlayersStatus() error {
	if !t.Meta.GameStarted {
		return ErrGameNotStarted
	}

	for k, v := range t.Meta.Players {
		if !v.GetFold() {
			t.Meta.Players[k].SetStatus(false)
		}
	}
	return nil
}

func refreshPlayers(players map[string]IPlayer, fold bool) error {
	for _, v := range players {
		if fold {
			v.SetFold(false)
		}
		v.SetLastBet(0)
		v.SetStatus(false)
	}
	return nil
}

func (t *PokerTable) handleCall(playerId string) error {
	if !t.Meta.GameStarted {
		return ErrGameNotStarted
	}

	if t.Meta.CurrentBet == 0 {
		return t.handleCheck(playerId)
	}
	if t.Meta.Players[playerId].GetFold() {
		return ErrPlayerIsFold
	}
	needToBet := t.Meta.CurrentBet - t.Meta.Players[playerId].GetLastBet()
	possibleBet := min(needToBet, t.Meta.Players[playerId].GetBalance())

	t.Meta.Players[playerId].ChangeBalance(-possibleBet)
	t.Meta.Players[playerId].SetStatus(true)
	if t.Meta.Players[playerId].GetBalance() > 0 {
		t.Meta.Players[playerId].SetLastBet(t.Meta.CurrentBet)
	} else {
		t.Meta.Players[playerId].SetLastBet(possibleBet)
	}

	t.NotifyObservers(fmt.Sprintf("Player %s do call with %d amount", playerId, t.Meta.CurrentBet))
	return nil
}

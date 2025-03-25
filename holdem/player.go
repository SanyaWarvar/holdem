package holdem

import (
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
)

var (
	ErrNotEnoughBalance = errors.New("player dont have enough money")
)

type IPlayer interface {
	GetBalance() int
	ChangeBalance(delta int) error
	GetId() string
	GetFold() bool
	SetFold(status bool)
	GetReadyStatus() bool // статус показывает сделал ли игрок check || raise || call
	SetStatus(status bool)
	GetHand() Hand
	SetHand(h Hand)
	GetLastBet() int
	SetLastBet(bet int)
	fmt.Stringer
}

type Hand struct {
	Cards [2]Card
}

type Player struct {
	Balance int
	Id      uuid.UUID
	Status  bool
	LastBet int
	Hand    Hand
	IsFold  bool
}

func (p *Player) String() string {
	return fmt.Sprintf(
		"Player %s:\n balance = %d\n ready status = %v\n last bet = %d\n card in hands: %v\n fold his cards = %v",
		p.Id.String(), p.Balance, p.Status, p.LastBet, p.Hand, p.IsFold,
	)
}

type Number interface {
	int
}

func Abs[T Number](value T) T {
	return T(math.Abs(float64(value)))
}

func (p *Player) GetBalance() int {
	return p.Balance
}

func (p *Player) ChangeBalance(delta int) error {
	if delta < 0 && int(math.Abs(float64(delta))) > p.Balance {
		return ErrNotEnoughBalance
	}
	p.Balance += delta

	return nil
}

func (p *Player) GetId() string {
	return p.Id.String()
}

func (p *Player) GetReadyStatus() bool {
	return p.Status
}

func (p *Player) SetStatus(status bool) {
	p.Status = status
}

func (p *Player) GetLastBet() int {
	return p.LastBet
}

func (p *Player) SetLastBet(bet int) { // только положительные?
	p.LastBet = bet
}

func (p *Player) SetHand(h Hand) {
	p.Hand = h
}

func (p *Player) GetHand() Hand {
	return p.Hand
}

func (p *Player) GetFold() bool {
	return p.IsFold
}

func (p *Player) SetFold(status bool) {
	p.IsFold = status
}

package holdem

import (
	"fmt"
)

type IObserver interface {
	Update(event string)
}

type Logger struct{}

func (l Logger) Update(event string) {
	fmt.Println("Event:", event)
}

package holdem

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTableGame1Good(t *testing.T) {
	t.Run("game 1", func(t *testing.T) {
		meta := NewTableMeta(50, 0, 1488)
		config := NewTableConfig(time.Hour, 10, 2, -1, false)
		table := NewPokerTable(config, meta)
		p1 := &Player{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Balance: 1000} //bb
		p2 := &Player{Id: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Balance: 1000} //dealer
		p3 := &Player{Id: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Balance: 1000} //sb
		p1Id := p1.GetId()
		p2Id := p2.GetId()
		p3Id := p3.GetId()
		err := table.AddPlayer(p1)
		fmt.Println(err)
		err = table.AddPlayer(p2)
		fmt.Println(err)
		err = table.AddPlayer(p3)
		fmt.Println(err)

		table.StartGame()
		table.MakeMove(p2Id, "call", 100)
		table.MakeMove(p3Id, "call", 100)
		table.MakeMove(p1Id, "call", 100)

		table.MakeMove(p3Id, "raise", 200)
		table.MakeMove(p1Id, "fold", 0)
		table.MakeMove(p2Id, "call", 200)

		table.MakeMove(p3Id, "check", 0)
		table.MakeMove(p2Id, "check", 0)

		table.MakeMove(p3Id, "check", 0)
		table.MakeMove(p2Id, "check", 0)

		fmt.Println(p1.Balance, p2.Balance, p3.Balance)

		require.Equal(t, p3.Balance, 1400)
		require.Equal(t, p1.Balance, 900)
		require.Equal(t, p2.Balance, 700)
		require.Equal(t, table.Meta.GameStarted, false)
	})

}
func TestTableGame2Good(t *testing.T) {
	t.Run("game 2", func(t *testing.T) {
		meta := NewTableMeta(50, 0, 1488)
		config := NewTableConfig(time.Hour, 10, 2, -1, false)
		table := NewPokerTable(config, meta)
		p1 := &Player{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Balance: 1000} //bb
		p2 := &Player{Id: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Balance: 1000} //dealer
		p3 := &Player{Id: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Balance: 1000} //sb
		p1Id := p1.GetId()
		p2Id := p2.GetId()
		p3Id := p3.GetId()
		err := table.AddPlayer(p1)
		fmt.Println(err)
		err = table.AddPlayer(p2)
		fmt.Println(err)
		err = table.AddPlayer(p3)
		fmt.Println(err)

		table.StartGame()
		table.MakeMove(p2Id, "call", 0)
		table.MakeMove(p3Id, "call", 0)
		table.MakeMove(p1Id, "call", 0)
		table.MakeMove(p3Id, "raise", 200)
		table.MakeMove(p1Id, "fold", 0)
		table.MakeMove(p2Id, "call", 0)

		table.MakeMove(p3Id, "check", 0)
		table.MakeMove(p2Id, "check", 0)

		table.MakeMove(p3Id, "check", 0)
		table.MakeMove(p2Id, "check", 0)

		require.Equal(t, p3.Balance, 1400)
		require.Equal(t, p1.Balance, 900)
		require.Equal(t, p2.Balance, 700)
		require.Equal(t, table.Meta.GameStarted, false)

		table.StartGame()
		fmt.Println(table.MakeMove(p3Id, "fold", 0))
		fmt.Println(table.MakeMove(p1Id, "call", 0))
		fmt.Println(table.MakeMove(p2Id, "call", 0))

		fmt.Println(table.MakeMove(p1Id, "raise", 300))
		fmt.Println(table.MakeMove(p2Id, "call", 300))

		fmt.Println(table.MakeMove(p1Id, "check", 0))
		fmt.Println(table.MakeMove(p2Id, "check", 0))
		fmt.Println(table.MakeMove(p1Id, "check", 0))
		fmt.Println(table.MakeMove(p2Id, "check", 0))

		require.Equal(t, p1.Balance, 500)
		require.Equal(t, p2.Balance, 1100)
		require.Equal(t, p3.Balance, 1400)
		require.Equal(t, table.Meta.GameStarted, false)
	})
}

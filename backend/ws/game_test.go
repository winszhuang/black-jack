package ws

import (
	"black-jack/game"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGame_calculateFinalWinner(t *testing.T) {
	t.Run("無人或一人無法獲勝", func(t *testing.T) {
		newGame := initGame()
		winners, isExist := newGame.calculateFinalWinners()
		require.False(t, isExist)
		require.Equal(t, 0, len(winners))

		client01 := genClientAndAddCards("test01", "5S", "KS", "QS")
		newGame.OnJoin(client01)

		winners, isExist = newGame.calculateFinalWinners()
		require.False(t, isExist)
		require.Equal(t, 0, len(winners))
	})
	t.Run("沒人獲勝", func(t *testing.T) {
		newGame := initGame()

		client01 := genClientAndAddCards("test01", "5S", "KS", "QS")
		client02 := genClientAndAddCards("test02", "5H", "KC", "QC")

		newGame.OnJoin(client01)
		newGame.OnJoin(client02)
		winners, isExist := newGame.calculateFinalWinners()
		require.False(t, isExist)
		require.Equal(t, 0, len(winners))
	})
	t.Run("一人獲勝", func(t *testing.T) {
		newGame := initGame()

		client01 := genClientAndAddCards("test01", "4D", "KS")
		client02 := genClientAndAddCards("test02", "5H", "KC")

		newGame.OnJoin(client01)
		newGame.OnJoin(client02)
		winners, isExist := newGame.calculateFinalWinners()
		require.True(t, isExist)
		require.Equal(t, 1, len(winners))
		require.Equal(t, winners[0].GetID(), "test02")
	})
	t.Run("兩人同分", func(t *testing.T) {
		dealer := game.NewCardDealerMock(game.Deck{})
		newGame := NewGame(dealer)

		client01 := genClientAndAddCards("test01", "3D", "KS")
		client02 := genClientAndAddCards("test02", "6H", "7C")

		newGame.OnJoin(client01)
		newGame.OnJoin(client02)
		winners, isExist := newGame.calculateFinalWinners()
		require.True(t, isExist)
		require.Equal(t, 2, len(winners))
	})

}

func initGame() *Game {
	dealer := game.NewCardDealerMock(game.Deck{})
	return NewGame(dealer)
}

func genClientAndAddCards(id string, cards ...string) *ClientMock {
	client := NewClientMock(id)
	for _, name := range cards {
		client.AddCard(game.NewCardByName(name))
	}
	return client
}

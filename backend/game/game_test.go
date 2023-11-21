package game

import (
	"black-jack/utils"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestGame_AddPlayers(t *testing.T) {
	t.Run("init player cards", func(t *testing.T) {
		game := initGameAndPlayers(5)
		err := game.BuildPlayersCards()
		require.NoError(t, err)

		for _, player := range game.players {
			require.Len(t, player.Deck, 2)
		}
	})
	t.Run("hit player", func(t *testing.T) {
		game := initGameAndPlayers(5)
		err := game.BuildPlayersCards()
		require.NoError(t, err)

		currPlayerID := getFirstPlayer(game)
		err = game.DealCardToPlayer(currPlayerID)
		require.NoError(t, err)

		player, err := game.getPlayer(currPlayerID)
		require.NoError(t, err)
		require.Len(t, player.Deck, 3)
	})
	t.Run("hit players in same time (concurrency)", func(t *testing.T) {
		game := initGameAndPlayers(10)
		err := game.BuildPlayersCards()
		require.NoError(t, err)

		for id, _ := range game.players {
			go func(c string) {
				err = game.DealCardToPlayer(c)
				require.NoError(t, err)
			}(id)
		}
	})
}

func getFirstPlayer(game *Game) string {
	var currPlayerID string
	for playerID, _ := range game.players {
		currPlayerID = playerID
		break
	}
	return currPlayerID
}

func initGameAndPlayers(playerCount int) *Game {
	r := rand.New(rand.NewSource(123456))
	cardDealer := NewCardDealer(r)
	game := NewGame(cardDealer)

	for i := 0; i < playerCount; i++ {
		game.AddPlayer(utils.RandomPlayerName())
	}
	return game
}

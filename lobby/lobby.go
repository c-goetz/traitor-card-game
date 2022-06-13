package lobby

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/c-goetz/traitor-card-game/game"
)

/*
Manage the lobby, don't refer to http things here, just channels.
For each lobby start a go routine to handle events and notify players.
"Connection" to player needs the possibility to timeout (maybe sth. like 5 min?).
Or is just timing out the lobby enough?
One Player should be host. Hos should have special rights like removing players from lobby.
*/

type Lobby struct {
	game         *game.Game
	playerNames  []string
	playerTokens []string
	// TODO channels for notification
}

func (l *Lobby) Join(name string) error {
	if n := len(l.playerNames); n == 10 {
		return errors.New("max lobby size reached")
	}
	for _, n := range l.playerNames {
		if n == name {
			return fmt.Errorf("player with name %s already joined", n)
		}
	}
	l.playerNames = append(l.playerNames, name)
	token, err := rand.Int(rand.Reader, big.NewInt(1_000_000_000))
	if err != nil {
		return fmt.Errorf("generating token: %v", err)
	}
	l.playerTokens = append(l.playerTokens, token.String())
	// TODO channels
	// TODO update player views?
	return nil
}

func (l *Lobby) Start() error {
	n := len(l.playerNames)
	if n < 3 || n > 10 {
		return fmt.Errorf("can't start with %d players", n)
	}
	g, err := game.NewGame(n)
	if err != nil {
		return err
	}
	l.game = &g
	// TODO update player views?
	return nil
}

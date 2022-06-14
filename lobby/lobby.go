package lobby

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/c-goetz/traitor-card-game/game"
)

/*
Manage the lobby, don't refer to http things here, just channels.
For each lobby start a go routine to handle events and notify players.
"Connection" to player needs the possibility to timeout (maybe sth. like 5 min?).
Or is just timing out the lobby enough?
One Player should be host. Hos should have special rights like removing players from lobby.
*/

var lobbies struct {
	sync.RWMutex
	ls map[uint64]Lobby
}

type Player struct {
	name     string
	token    string
	lastSeen time.Time
	channel  *chan string
	uuid     game.Player
}

type Lobby struct {
	game    *game.Game
	uuid    uint64
	players []Player
}

func (player *Player) register(channel *chan string) {
	player.channel = channel
}

func (player *Player) unregisterChannel() {
	player.channel = nil
}

func (l *Lobby) NewPlayer(name string, token string, channel *chan string) Player {
	return Player{name, token, time.Now(), channel, game.Player(len(l.players))}
}

func createLobby(host string) (Player, Lobby, error) {
	lobbyUID, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt))
	if err != nil {
		return Player{}, Lobby{}, fmt.Errorf("generating Lobby UID: %v", err)
	}
	lobby := Lobby{nil,
		lobbyUID.Uint64(),
		[]Player{},
	}
	player, err := lobby.Join(host)
	if err != nil {
		return Player{}, Lobby{}, fmt.Errorf("could not create player with name: %v", host)
	}
	return player, lobby, nil
}

func (l *Lobby) Join(name string) (Player, error) {
	if n := len(l.players); n == 10 {
		return Player{}, errors.New("max lobby size reached")
	}
	for _, player := range l.players {
		if player.name == name {
			return Player{}, fmt.Errorf("player with name %s already joined", player.name)
		}
	}
	token, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt))
	if err != nil {
		return Player{}, fmt.Errorf("generating token: %v", err)
	}
	player := l.NewPlayer(name, token.String(), nil)
	l.players = append(l.players, player)

	// TODO update player views?
	return player, nil
}

func (l *Lobby) Start() error {
	lobbies.Lock()
	defer lobbies.Unlock()
	n := len(l.players)
	if n < 3 || n > 10 {
		return fmt.Errorf("can't start with %d players", n)
	}
	g, err := game.NewGame(n)
	if err != nil {
		return err
	}
	l.game = &g
	lobbies.ls[l.uuid] = *l
	// TODO update player views?
	return nil
}
func (l *Lobby) Close() error {
	lobbies.Lock()
	defer lobbies.Unlock()

	if _, ok := lobbies.ls[l.uuid]; ok {
		delete(lobbies.ls, l.uuid)
	} else {
		return fmt.Errorf("lobby not found UID: %v", l.uuid)
	}

	return nil
}

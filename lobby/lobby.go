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

var lobbies = struct {
	sync.RWMutex
	ls map[uint]*Lobby
}{
	ls: map[uint]*Lobby{},
}

type Player struct {
	Name     string
	token    string
	lastSeen time.Time
	channel  *chan Message
	position game.Player
}

type Lobby struct {
	sync.RWMutex
	game    *game.Game
	Uuid    uint
	players []Player
}

func Register(lobby, player uint, channel *chan Message) {
	lobbies.RLock()
	defer lobbies.RUnlock()
	l := lobbies.ls[lobby]
	l.Lock()
	defer l.Unlock()

	l.players[player].channel = channel
}

func UnregisterChannel(lobby, player uint) {
	lobbies.RLock()
	defer lobbies.RUnlock()
	l := lobbies.ls[lobby]
	l.Lock()
	defer l.Unlock()

	l.players[player].channel = nil
}

func SetName(lobby, player uint, name string) {
	lobbies.RLock()
	defer lobbies.RUnlock()
	l := lobbies.ls[lobby]
	l.Lock()
	defer l.Unlock()

	l.players[player].Name = name
}

func (l *Lobby) NewPlayer(name string, token string, channel *chan Message) Player {
	return Player{name, token, time.Now(), channel, game.Player(len(l.players))}
}

// CreateLobby Creates Lobby and Host
// First player is always Host
// Return Lobby number
func CreateLobby(host string) (uint, error) {
	lobbies.Lock()
	lobbyUID, err := rand.Int(rand.Reader, big.NewInt(math.MaxUint32))
	if err != nil {
		return 0, fmt.Errorf("generating Lobby UID: %v", err)
	}
	id := uint(lobbyUID.Uint64())
	lobby := Lobby{
		sync.RWMutex{},
		nil,
		id,
		[]Player{},
	}
	lobbies.ls[id] = &lobby
	lobbies.Unlock()
	_, err = Join(id, host)
	if err != nil {
		return 0, fmt.Errorf("could not create player with name: %v %s", host, err)
	}
	return id, nil
}

func Join(lobby uint, name string) (uint, error) {
	lobbies.RLock()
	defer lobbies.RUnlock()
	if _, ok := lobbies.ls[lobby]; !ok {
		return 0, fmt.Errorf("lobby not found %d", lobby)
	}
	l := lobbies.ls[lobby]
	l.Lock()
	defer l.Unlock()

	if n := len(l.players); n == 10 {
		return 0, errors.New("max lobby size reached")
	}
	for _, player := range l.players {
		if player.Name == name {
			return 0, fmt.Errorf("player with name %s already joined", player.Name)
		}
	}
	token, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return 0, fmt.Errorf("generating token: %v", err)
	}
	player := l.NewPlayer(name, token.String(), nil)
	l.players = append(l.players, player)

	// TODO update player views?
	return uint(player.position), nil
}

func Claim(lobby, player uint, claim game.Cards) {
	lobbies.RLock()
	defer lobbies.RUnlock()
	l := lobbies.ls[lobby]

	err := l.game.Claim(game.Player(player), claim)
	if err != nil {
		l.broadcast(&ClaimMessage{err: err})
	} else {
		l.broadcast(&ClaimMessage{Cards: *l.game.Claims[player]})
	}
}

func Play(lobby, from, to uint) {
	lobbies.RLock()
	defer lobbies.RUnlock()
	l := lobbies.ls[lobby]
	err := l.game.Play(game.Player(from), game.Player(to))
	if err != nil {
		l.broadcast(&RevealCardMessage{err: err})
	} else {
		l.broadcast(&RevealCardMessage{Cards: l.game.RevealedCards})
	}
}

func GetRole(lobby, from uint) {
	lobbies.RLock()
	defer lobbies.RUnlock()
	l := lobbies.ls[lobby]
	l.RLock()
	defer l.RUnlock()
	player := l.players[from]
	if len(l.game.Roles) < int(from) {
		err := fmt.Errorf("player position is not seated %d", player.position)
		*player.channel <- &RoleMessage{err: err}
	} else {
		*player.channel <- &RoleMessage{Role: l.game.Roles[player.position]}
	}
}
func GetHand(lobby, from uint) error {
	lobbies.RLock()
	defer lobbies.RUnlock()
	l := lobbies.ls[lobby]
	l.RLock()
	defer l.RUnlock()
	player := l.players[from]

	if len(l.game.Hands) < int(player.position) {
		return fmt.Errorf("player position is not seated %d", player.position)
	}
	*player.channel <- &HandMessage{Cards: l.game.Hands[player.position]}
	return nil
}

func GetGameState(lobby uint) {
	lobbies.RLock()
	defer lobbies.RUnlock()
	l := lobbies.ls[lobby]
	l.broadcast(&StateMessage{State: l.game.State()})
}

func (l *Lobby) broadcast(message Message) {
	for _, p := range l.players {
		if p.channel == nil {
			message.SetError(fmt.Errorf("player %d has no channel attached", p.position))
		}
		*p.channel <- message
	}
}

func Start(lobby uint) error {
	lobbies.Lock()
	defer lobbies.Unlock()
	l := lobbies.ls[lobby]
	l.Lock()
	defer l.Unlock()
	n := len(l.players)
	if n < 3 || n > 10 {
		return fmt.Errorf("can't start with %d players", n)
	}
	g, err := game.NewGame(n)
	if err != nil {
		return err
	}
	l.game = &g
	return nil
}
func Close(lobby uint) {
	lobbies.Lock()
	defer lobbies.Unlock()

	if _, ok := lobbies.ls[lobby]; ok {
		delete(lobbies.ls, lobby)
	}
}

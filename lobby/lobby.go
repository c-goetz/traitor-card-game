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
	ls map[uint64]Lobby
}{
	ls: map[uint64]Lobby{},
}

type Player struct {
	Name     string
	token    string
	lastSeen time.Time
	channel  *chan Message
	position game.Player
}

type Lobby struct {
	game    *game.Game
	Uuid    uint64
	players []Player
}

func (player *Player) Register(channel *chan Message) {
	player.channel = channel
}

func (player *Player) UnregisterChannel() {
	player.channel = nil
}

func (player *Player) SetName(name string) {
	player.Name = name
}

func (l *Lobby) NewPlayer(name string, token string, channel *chan Message) Player {
	return Player{name, token, time.Now(), channel, game.Player(len(l.players))}
}

func CreateLobby(host string) (*Player, *Lobby, error) {
	lobbyUID, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, nil, fmt.Errorf("generating Lobby UID: %v", err)
	}
	id := lobbyUID.Uint64()
	lobby := Lobby{
		nil,
		id,
		[]Player{},
	}
	player, err := lobby.Join(host)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create player with name: %v", host)
	}
	lobbies.Lock()
	defer lobbies.Unlock()
	lobbies.ls[id] = lobby
	return &player, &lobby, nil
}

func (l *Lobby) Join(name string) (Player, error) {
	if n := len(l.players); n == 10 {
		return Player{}, errors.New("max lobby size reached")
	}
	for _, player := range l.players {
		if player.Name == name {
			return Player{}, fmt.Errorf("player with name %s already joined", player.Name)
		}
	}
	token, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return Player{}, fmt.Errorf("generating token: %v", err)
	}
	player := l.NewPlayer(name, token.String(), nil)
	l.players = append(l.players, player)

	// TODO update player views?
	return player, nil
}

func (l *Lobby) Claim(player *Player, claim game.Cards) {
	err := l.game.Claim(player.position, claim)
	if err != nil {
		l.broadcast(&ClaimMessage{err: err})
	} else {
		l.broadcast(&ClaimMessage{Cards: *l.game.Claims[player.position]})
	}
}

func (l *Lobby) Play(from, to *Player) {
	err := l.game.Play(from.position, to.position)
	if err != nil {
		l.broadcast(&RevealCardMessage{err: err})
	} else {
		l.broadcast(&RevealCardMessage{Cards: l.game.RevealedCards})
	}
}

func (l *Lobby) GetRole(from *Player) {
	if len(l.game.Roles) < int(from.position) {
		err := fmt.Errorf("player position is not seated %d", from.position)
		*from.channel <- &RoleMessage{err: err}
	} else {
		*from.channel <- &RoleMessage{Role: l.game.Roles[from.position]}
	}
}
func (l *Lobby) GetHand(from *Player) error {
	if len(l.game.Hands) < int(from.position) {
		return fmt.Errorf("player position is not seated %d", from.position)
	}
	*from.channel <- &HandMessage{Cards: l.game.Hands[from.position]}
	return nil
}

func (l *Lobby) GetGameState() {
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
	lobbies.ls[l.Uuid] = *l
	// TODO update player views?
	return nil
}
func (l *Lobby) Close() error {
	lobbies.Lock()
	defer lobbies.Unlock()

	if _, ok := lobbies.ls[l.Uuid]; ok {
		delete(lobbies.ls, l.Uuid)
	} else {
		return fmt.Errorf("lobby not found UID: %v", l.Uuid)
	}

	return nil
}

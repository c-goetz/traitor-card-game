package lobby

import "testing"

func TestNewLobby(t *testing.T) {
	lobby, err := CreateLobby("test")
	defer Close(lobby)
	if err != nil {
		t.Fatalf("expected lobby to be created successfully %s", err)
	}
	l := getLobby(lobby)
	if len(l.players) != 1 {
		t.Fatalf("expected host to be joined")
	}

}

func TestChangePlayerName(t *testing.T) {
	lobby := CreateTestLobby()
	defer Close(lobby)
	SetName(lobby, 1, "Changed")
	l := getLobby(lobby)
	if l.players[1].Name != "Changed" {
		t.Fatalf("Could not change Players name")
	}
}

func TestNewPlayer(t *testing.T) {
	lobby, _ := CreateLobby("test")
	l := getLobby(lobby)
	l.RLock()
	defer l.RUnlock()
	if l.players[0].position != 0 {
		t.Fatalf("expected first player to have position 0 but it has %d", l.players[0].position)
	}
}

func TestStartLobby(t *testing.T) {
	lobby := CreateTestLobby()
	defer Close(lobby)
	l := getLobby(lobby)
	l.RLock()
	defer l.RUnlock()
	if len(l.players) != 4 {
		t.Fatalf("expected four players but it has %d", len(l.players))
	}

	if len(l.players) != len(l.game.Roles) {
		t.Fatalf("expected player count to match role count")
	}
}

func TestGameState(t *testing.T) {
	lobby := CreateTestLobby()
	channels := make([]chan Message, 4)
	setupChannels(lobby, channels)
	go GetGameState(lobby)
	for i, _ := range channels {
		message := <-channels[i]
		if message.GetKind() != "StateMessage" {
			t.Fatalf("expected state message to be broadcast")
		}
	}
}

func TestGetHand(t *testing.T) {
	lobby := CreateTestLobby()
	channels := make([]chan Message, 4)
	setupChannels(lobby, channels)
	for i, _ := range channels {
		go GetHand(lobby, uint(i))
		message := <-channels[i]
		if message.GetKind() != "HandMessage" {
			t.Fatalf("expected Hand message to be broadcast")
		}
	}
}

func setupChannels(lobby uint, channels []chan Message) {
	for i, _ := range channels {
		channels[i] = make(chan Message)
		Register(lobby, uint(i), &channels[i])
	}
}

func getLobby(lobby uint) *Lobby {
	lobbies.RLock()
	defer lobbies.RUnlock()
	return lobbies.ls[lobby]
}

func CreateTestLobby() uint {
	lobby, _ := CreateLobby("test")
	Join(lobby, "test2")
	Join(lobby, "test3")
	Join(lobby, "test4")
	Start(lobby)
	return lobby
}

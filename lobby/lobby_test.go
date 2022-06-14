package lobby

import "testing"

func TestNewLobby(t *testing.T) {
	_, lobby, _ := createLobby("test")
	if len(lobby.players) != 1 {
		t.Fatal("expected one player to be in lobby")
	}

}

func TestNewPlayer(t *testing.T) {
	player, _, _ := createLobby("test")
	if player.position != 0 {
		t.Fatalf("expected first player to have position 0 but it has %d", player.position)
	}
}

func TestStartLobby(t *testing.T) {
	lobby := createTestLobby()

	if len(lobby.players) != 4 {
		t.Fatalf("expected four players but it has %d", len(lobby.players))
	}

	if len(lobby.players) != len(lobby.game.Roles) {
		t.Fatalf("expected player count to match role count")
	}
}

func TestClaim(t *testing.T) {
	lobby := createTestLobby()
	channels := make([]chan Message, len(lobby.players))
	for i, _ := range lobby.players {
		channels[i] = make(chan Message)
		lobby.players[i].Register(&channels[i])
	}
	go lobby.GetGameState()
	for i, _ := range channels {
		message := <-channels[i]
		if message.GetKind() != "StateMessage" {
			t.Fatalf("expected state message to be broadcast")
		}
	}
}

func createTestLobby() Lobby {
	_, lobby, _ := createLobby("test")
	lobby.Join("test2")
	lobby.Join("test3")
	lobby.Join("test4")
	lobby.Start()
	return lobby
}

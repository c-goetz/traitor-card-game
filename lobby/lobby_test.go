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
	if player.uuid != 0 {
		t.Fatalf("expected first player to have uuid 0 but it has %d", player.uuid)
	}

}

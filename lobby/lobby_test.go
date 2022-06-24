package lobby

import "testing"

func TestNewLobby(t *testing.T) {
	_, lobby, _ := CreateLobby("test")
	defer lobby.Close()
	if len(lobby.players) != 1 {
		t.Fatal("expected one player to be in lobby")
	}

}

func TestChangePlayerName(t *testing.T) {
	lobby := CreateTestLobby()
	lobby.players[1].SetName("Changed")
	if lobby.players[1].Name != "Changed" {
		t.Fatalf("Could not change Players name")
	}
}

func TestNewPlayer(t *testing.T) {
	player, _, _ := CreateLobby("test")
	if player.position != 0 {
		t.Fatalf("expected first player to have position 0 but it has %d", player.position)
	}
}

func TestStartLobby(t *testing.T) {
	lobby := CreateTestLobby()
	defer lobby.Close()
	if len(lobby.players) != 4 {
		t.Fatalf("expected four players but it has %d", len(lobby.players))
	}

	if len(lobby.players) != len(lobby.game.Roles) {
		t.Fatalf("expected player count to match role count")
	}
}

func TestGameState(t *testing.T) {
	lobby := CreateTestLobby()
	defer lobby.Close()
	channels := make([]chan Message, len(lobby.players))
	setupChannels(lobby, channels)
	go lobby.GetGameState()
	for i, _ := range channels {
		message := <-channels[i]
		if message.GetKind() != "StateMessage" {
			t.Fatalf("expected state message to be broadcast")
		}
	}
}

func TestGetHand(t *testing.T) {
	lobby := CreateTestLobby()
	defer lobby.Close()
	channels := make([]chan Message, len(lobby.players))
	setupChannels(lobby, channels)
	for i, _ := range lobby.players {
		go lobby.GetHand(&lobby.players[i])
		message := <-channels[i]
		if message.GetKind() != "HandMessage" {
			t.Fatalf("expected Hand message to be broadcast")
		}
	}
}

func setupChannels(lobby Lobby, channels []chan Message) {
	for i, _ := range lobby.players {
		channels[i] = make(chan Message)
		lobby.players[i].Register(&channels[i])
	}
}

func CreateTestLobby() Lobby {
	_, lobby, _ := CreateLobby("test")
	lobby.Join("test2")
	lobby.Join("test3")
	lobby.Join("test4")
	lobby.Start()
	return *lobby
}

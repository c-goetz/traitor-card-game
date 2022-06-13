package lobby

import "testing"

func TestNewLobby(t *testing.T) {
	_, err := newLobby(2)
	if err == nil {
		t.Fatal("expected game to error with less than 3 players")
	}
	_, err = newLobby(11)
	if err == nil {
		t.Fatal("expected game to error with more than 10 players")
	}
}

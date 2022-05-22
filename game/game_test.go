package game

import "testing"

func TestRoleDraw(t *testing.T) {
	deck := roleDeck(4)
	good := uint8(0)
	bad := uint8(0)
	for deck.sum() > 0 {
		role := deck.draw()
		switch role {
		case RoleGood:
			good++
		case RoleBad:
			bad++
		}
	}
	deck = roleDeck(4)
	if deck.Good != good {
		t.Fatalf("expected to draw %d good cards, got: %d", deck.Good, good)
	}
	if deck.Bad != bad {
		t.Fatalf("expected to draw %d bad cards, got: %d", deck.Bad, bad)
	}
}

func TestCardDraw(t *testing.T) {
	deck := cardDeck(4)
	var draws Cards
	for deck.sum() > 0 {
		card := deck.draw()
		switch card {
		case CardNeutral:
			draws.Neutral++
		case CardGood:
			draws.Good++
		case CardBad:
			draws.Bad++
		}
	}
	deck = cardDeck(4)
	if deck != draws {
		t.Fatalf("expected: %v, got: %v", deck, draws)
	}
}

func TestNewGame(t *testing.T) {
	for i := 3; i <= 10; i++ {
		g, err := NewGame(i)
		if err != nil {
			t.Fatalf("expected game to be created without error, got: %v", err)
		}
		g.invariants(t)
	}
	_, err := NewGame(2)
	if err == nil {
		t.Fatal("expected game to error with less than 3 players")
	}
	_, err = NewGame(11)
	if err == nil {
		t.Fatal("expected game to error with more than 10 players")
	}
}

func (g *Game) invariants(t *testing.T) {
	t.Helper()
	g.ensureSizes(t)
	g.ensureCards(t)
	if t.Failed() {
		t.FailNow()
	}
}

func (g *Game) ensureSizes(t *testing.T) {
	t.Helper()
	if len(g.hands) != int(g.playerCount) {
		t.Errorf("expected as many hands: %d as players %d", len(g.hands), g.playerCount)
	}
	if len(g.roles) != int(g.playerCount) {
		t.Errorf("expected as many roles: %d as players: %d", len(g.roles), g.playerCount)
	}
	if len(g.claims) != int(g.playerCount) {
		t.Errorf("expected as many claims: %d as players: %d", len(g.claims), g.playerCount)
	}
}

func (g *Game) ensureCards(t *testing.T) {
	t.Helper()
	sum := g.revealedCards
	for i := Player(0); i < g.playerCount; i++ {
		hand := g.hands[i]
		sum.Neutral += hand.Neutral
		sum.Good += hand.Good
		sum.Bad += hand.Bad
	}
	deck := cardDeck(g.playerCount)
	if sum != deck {
		t.Errorf("expected all cards: %v to be somewhere, got: %v", deck, sum)
	}
}

package game

import (
	"math/rand"
	"testing"
)

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
		t.Fatalf("expected: %+v, got: %+v", deck, draws)
	}
}

func TestNewGame(t *testing.T) {
	for i := 3; i <= 10; i++ {
		g, err := NewGame(i)
		if err != nil {
			t.Fatalf("expected game to be created without error, got: %v", err)
		}
		(&testGame{g, t}).invariants()
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

func TestHappyPathGoodWin(t *testing.T) {
	rand.Seed(42)
	game, err := NewGame(4)
	if err != nil {
		t.Fatalf("could not create game: %v", err)
	}
	g := testGame{game, t}
	g.claimTruth()
	// manipulate hands to get some known state
	// 4 player deck: 12, 6, 2
	g.hands[0] = Cards{0, 5, 0}
	g.hands[1] = Cards{5, 0, 0}
	g.hands[2] = Cards{2, 1, 2}
	g.hands[3] = Cards{5, 0, 0}
	g.tPlay(0, 1)
	g.tPlay(1, 0)
	g.tPlay(0, 1)
	g.tPlay(1, 0)
	g.claimTruth()
	// deck: 10, 4, 2
	g.hands[0] = Cards{0, 4, 0}
	g.hands[1] = Cards{4, 0, 0}
	g.hands[2] = Cards{4, 0, 0}
	g.hands[3] = Cards{2, 0, 2}
	g.tPlay(0, 1)
	g.tPlay(1, 0)
	g.tPlay(0, 1)
	g.tPlay(1, 0)
	g.claimTruth()
	// deck: 8, 2, 2
	g.hands[0] = Cards{1, 2, 0}
	g.hands[1] = Cards{3, 0, 0}
	g.hands[2] = Cards{3, 0, 0}
	g.hands[3] = Cards{1, 0, 2}
	g.tPlay(0, 1)
	g.tPlay(1, 0)
	g.tPlay(0, 1)
	g.tPlay(1, 0)
	g.claimTruth()
	// deck: 5, 1, 2
	g.hands[0] = Cards{1, 1, 0}
	g.hands[1] = Cards{2, 0, 0}
	g.hands[2] = Cards{2, 0, 0}
	g.hands[3] = Cards{0, 0, 2}
	g.tPlay(0, 1)
	g.tPlay(1, 0)
	if g.state() != StateWinGood {
		t.Fatal("expected good to win")
	}
}

func TestHappyPathBadWinBadCards(t *testing.T) {
	rand.Seed(42)
	game, err := NewGame(4)
	if err != nil {
		t.Fatalf("could not create game: %v", err)
	}
	g := testGame{game, t}
	g.claimTruth()
	// 4 player deck: 12, 6, 2
	g.hands[0] = Cards{0, 5, 0}
	g.hands[1] = Cards{5, 0, 0}
	g.hands[2] = Cards{2, 1, 2}
	g.hands[3] = Cards{5, 0, 0}
	g.tPlay(0, 2)
	g.tPlay(2, 0)
	g.tPlay(0, 2)
	g.tPlay(2, 0)
	g.claimTruth()
	// deck: 12, 3, 1
	g.hands[0] = Cards{1, 3, 0}
	g.hands[1] = Cards{4, 0, 0}
	g.hands[2] = Cards{3, 0, 1}
	g.hands[3] = Cards{4, 0, 0}
	g.tPlay(0, 2)
	g.tPlay(2, 0)
	g.tPlay(0, 2)
	g.tPlay(2, 0)
	g.claimTruth()
	// deck: 9, 2, 1
	g.hands[0] = Cards{1, 2, 0}
	g.hands[1] = Cards{3, 0, 0}
	g.hands[2] = Cards{2, 0, 1}
	g.hands[3] = Cards{3, 0, 0}
	g.tPlay(0, 2)
	g.tPlay(2, 0)
	g.tPlay(0, 2)
	if g.state() != StateWinBad || g.revealedCards.Bad != 2 {
		t.Fatal("expected bad to win by discovering 2 bad cards")
	}
}

type testGame struct {
	Game
	*testing.T
}

func (g *testGame) tPlay(from, to Player) {
	g.Helper()
	g.invariants()
	err := g.Play(from, to)
	if err != nil {
		g.Fatal(err)
	}
	g.invariants()
}

func (g *testGame) tClaim(p Player, c Cards) {
	g.Helper()
	g.invariants()
	err := g.Claim(p, c)
	if err != nil {
		g.Fatal(err)
	}
	g.invariants()
}

func (g *testGame) claimTruth() {
	g.Helper()
	if g.state() != StateClaiming {
		g.Fatal("expected to be claiming")
	}
	for p := Player(0); p < g.playerCount; p++ {
		g.tClaim(p, g.hands[p])
	}
	if g.state() != StatePlaying {
		g.Fatal("expected to be playing")
	}
}

func (g *testGame) invariants() {
	g.Helper()
	g.ensureSizes()
	g.ensureCards()
	g.ensureHandSize()
	if g.Failed() {
		g.FailNow()
	}
}

func (g *testGame) ensureSizes() {
	g.Helper()
	if len(g.hands) != int(g.playerCount) {
		g.Errorf("expected as many hands: %d as players %d", len(g.hands), g.playerCount)
	}
	if len(g.roles) != int(g.playerCount) {
		g.Errorf("expected as many roles: %d as players: %d", len(g.roles), g.playerCount)
	}
	if len(g.claims) != int(g.playerCount) {
		g.Errorf("expected as many claims: %d as players: %d", len(g.claims), g.playerCount)
	}
}

func (g *testGame) ensureCards() {
	g.Helper()
	sum := g.revealedCards
	for i := Player(0); i < g.playerCount; i++ {
		hand := g.hands[i]
		sum.Neutral += hand.Neutral
		sum.Good += hand.Good
		sum.Bad += hand.Bad
	}
	deck := cardDeck(g.playerCount)
	if sum != deck {
		g.Errorf("expected all cards: %v to be somewhere, got: %v", deck, sum)
	}
}

func (g *testGame) ensureHandSize() {
	g.Helper()
	if g.cardsPlayedInRound() != 0 {
		return
	}
	for p := Player(1); p < g.playerCount; p++ {
		if g.hands[p-1].sum() != g.hands[p].sum() {
			g.Errorf(
				"expected player: %d with cards %+v and player: %d with cards %+v to have same amount of cards",
				p-1,
				g.hands[p-1],
				p,
				g.hands[p],
			)
		}
	}
}

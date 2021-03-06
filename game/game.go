package game

import (
	"fmt"
	"math/rand"
)

type Card uint8

const (
	CardNeutral Card = iota
	CardGood
	CardBad
)

type Cards struct {
	Neutral, Good, Bad uint8
}

func (c *Cards) sum() uint8 {
	return c.Neutral + c.Good + c.Bad
}

func (c *Cards) draw() Card {
	num := uint8(rand.Int31n(int32(c.sum())))
	switch {
	case num < c.Neutral:
		c.Neutral--
		return CardNeutral
	case num < c.Neutral+c.Good:
		c.Good--
		return CardGood
	default:
		c.Bad--
		return CardBad
	}
}

type Role uint8

const (
	RoleGood Role = iota
	RoleBad
)

type Roles struct {
	Good, Bad uint8
}

func roleDeck(players Player) Roles {
	switch players {
	case 3:
		return Roles{2, 2}
	case 4:
		return Roles{3, 2}
	case 5:
		return Roles{3, 2}
	case 6:
		return Roles{4, 2}
	case 7:
		return Roles{5, 3}
	case 8:
		return Roles{6, 3}
	case 9:
		return Roles{6, 3}
	case 10:
		return Roles{7, 4}
	}
	panic("unreachable")
}

func (r *Roles) sum() uint8 {
	return r.Good + r.Bad
}

func (r *Roles) draw() Role {
	num := rand.Int31n(int32(r.sum()))
	if uint8(num) < r.Good {
		r.Good--
		return RoleGood
	}
	r.Bad--
	return RoleBad
}

func cardDeck(players Player) Cards {
	switch players {
	case 3:
		return Cards{8, 5, 2}
	case 4:
		return Cards{12, 6, 2}
	case 5:
		return Cards{16, 7, 2}
	case 6:
		return Cards{20, 8, 2}
	case 7:
		return Cards{26, 7, 2}
	case 8:
		return Cards{30, 8, 2}
	case 9:
		return Cards{34, 9, 2}
	case 10:
		return Cards{37, 10, 3}
	}
	panic("unreachable")
}

type State uint8

const (
	StateClaiming State = iota
	StatePlaying
	StateWinGood
	StateWinBad
)

type Player uint8

type Game struct {
	playerCount   Player
	currentPlayer Player
	RevealedCards Cards
	Hands         []Cards
	Roles         []Role
	Claims        []*Cards
}

func NewGame(players int) (Game, error) {
	var g Game
	if players < 3 || 10 < players {
		return g, fmt.Errorf("invalid player count: %d, must be 3-10", players)
	}
	g.playerCount = Player(players)
	g.Claims = make([]*Cards, players)
	g.Hands = make([]Cards, players)
	roles := roleDeck(g.playerCount)
	for i := Player(0); i < g.playerCount; i++ {
		g.Roles = append(g.Roles, roles.draw())
	}
	g.deal()
	return g, nil
}

func (g *Game) Claim(player Player, claim Cards) error {
	if s := g.State(); s != StateClaiming {
		return fmt.Errorf("palyer: %d tried to claim in State %v", player, s)
	}
	g.Claims[player] = &claim
	return nil
}

func (g *Game) Play(from, to Player) error {
	if s := g.State(); s != StatePlaying {
		return fmt.Errorf("player: %d tried to play in State %v", from, s)
	}
	if g.currentPlayer != from {
		return fmt.Errorf("player: %d tried to play, but currentPlayer is: %d", from, g.currentPlayer)
	}
	g.currentPlayer = to
	card := g.Hands[to].draw()
	switch card {
	case CardNeutral:
		g.RevealedCards.Neutral++
	case CardGood:
		g.RevealedCards.Good++
	case CardBad:
		g.RevealedCards.Bad++
	}
	if g.cardsPlayedInRound() != 0 {
		return nil
	}
	// round ended
	for p := Player(0); p < g.playerCount; p++ {
		g.Claims[p] = nil
	}
	g.deal()
	return nil
}

func (g *Game) deal() {
	deck := cardDeck(g.playerCount)
	deck.Neutral -= g.RevealedCards.Neutral
	deck.Good -= g.RevealedCards.Good
	deck.Bad -= g.RevealedCards.Bad
	toDraw := 5 - g.round()
	for p := Player(0); p < g.playerCount; p++ {
		g.Hands[p] = Cards{}
		cards := &g.Hands[p]
		for i := uint8(0); i < toDraw; i++ {
			card := deck.draw()
			switch card {
			case CardNeutral:
				cards.Neutral++
			case CardGood:
				cards.Good++
			case CardBad:
				cards.Bad++
			}
		}
	}
}

func (g *Game) round() uint8 {
	return g.RevealedCards.sum() / uint8(g.playerCount)
}

func (g *Game) cardsPlayedInRound() uint8 {
	return g.RevealedCards.sum() % uint8(g.playerCount)
}

func (g *Game) State() State {
	deck := cardDeck(g.playerCount)
	if g.RevealedCards.Bad == deck.Bad {
		return StateWinBad
	}
	if g.RevealedCards.Good == deck.Good {
		return StateWinGood
	}
	if g.round() == 4 {
		return StateWinBad
	}
	for _, c := range g.Claims {
		if c == nil {
			return StateClaiming
		}
	}
	return StatePlaying
}

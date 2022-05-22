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
	revealedCards Cards
	hands         []Cards
	roles         []Role
	claims        []*Cards
}

func NewGame(players int) (Game, error) {
	var g Game
	if players < 3 || 10 < players {
		return g, fmt.Errorf("invalid player count: %d, must be 3-10", players)
	}
	g.playerCount = Player(players)
	g.claims = make([]*Cards, players)
	deck := cardDeck(g.playerCount)
	roles := roleDeck(g.playerCount)
	for i := Player(0); i < g.playerCount; i++ {
		g.roles = append(g.roles, roles.draw())
		g.hands = append(g.hands, Cards{})
		for j := 0; j < 5; j++ {
			card := deck.draw()
			cards := &g.hands[i]
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
	return g, nil
}

func (g *Game) state() State {
	for _, c := range g.claims {
		if c == nil {
			return StateClaiming
		}
	}
	// TODO continue here
	return StatePlaying
}

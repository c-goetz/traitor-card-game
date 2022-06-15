package lobby

import "github.com/c-goetz/traitor-card-game/game"

type Message interface {
	GetKind() string
	GetError() error
	SetError(err error)
}

type ClaimMessage struct {
	err   error
	Cards game.Cards
}
type RevealCardMessage struct {
	err   error
	Cards game.Cards
}

type RoleMessage struct {
	err  error
	Role game.Role
}

type StateMessage struct {
	err   error
	State game.State
}

type HandMessage struct {
	err   error
	Cards game.Cards
}

func (m *HandMessage) GetKind() string {
	return "HandMessage"
}

func (m *StateMessage) GetKind() string {
	return "StateMessage"
}

func (m *RoleMessage) GetKind() string {
	return "RoleMessage"
}

func (m *RevealCardMessage) GetKind() string {
	return "RevealCardMessage"
}

func (m *ClaimMessage) GetKind() string {
	return "ClaimMessage"
}

func (m *HandMessage) GetError() error {
	return m.err
}

func (m *StateMessage) GetError() error {
	return m.err
}

func (m *RoleMessage) GetError() error {
	return m.err
}

func (m *RevealCardMessage) GetError() error {
	return m.err
}

func (m *ClaimMessage) GetError() error {
	return m.err
}

func (m *HandMessage) SetError(err error) {
	m.err = err
}

func (m *StateMessage) SetError(err error) {
	m.err = err
}

func (m *RoleMessage) SetError(err error) {
	m.err = err
}

func (m *RevealCardMessage) SetError(err error) {
	m.err = err
}

func (m *ClaimMessage) SetError(err error) {
	m.err = err
}

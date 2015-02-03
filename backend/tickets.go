package backend

import (
	"errors"

	"gatewayd/utils"
)

// This was implemented, and later I found out it's no use. Happens :/

// WARNING!!! OBSOLETE!!!

// Ticket is used to let user to connect to the session
type Ticket struct {
	Session *Session
}

// TicketManager is used to handle tickets
type TicketManager struct {
	tickets map[string]Ticket // stores tickets, not exported
}

// NewTicketManager initializes new ticket manager
func NewTicketManager() *TicketManager {
	return &TicketManager{map[string]Ticket{}}
}

// AddTicket adds ticket
func (t *TicketManager) AddTicket(ticket string, session *Session) error {
	if _, ok := t.tickets[ticket]; ok {
		return errors.New("This ticket is already out!")
	}
	t.tickets[ticket] = Ticket{session}
	return nil
}

// RemoveTicket removes ticket
func (t *TicketManager) RemoveTicket(ticket string) {
	delete(t.tickets, ticket)
}

// Claim tries to claim ticket
// It removes ticket if successfully claimed
func (t *TicketManager) Claim(ticket string) (*Session, error) {
	tckt, ok := t.tickets[ticket]
	if !ok {
		return nil, errors.New("Invalid ticket")
	}
	delete(t.tickets, ticket)
	return tckt.Session, nil
}

// Generate generates ticket, adds it to the system and returns
// it's key, which then should be passed to user
func (t *TicketManager) Generate(session *Session) (ticket string) {
	for {
		ticket = utils.RandStr(32)
		if t.AddTicket(ticket, session) == nil {
			return
		}
	}
}

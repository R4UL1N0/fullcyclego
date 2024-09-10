package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTicketPriceZero = errors.New("price must be greater than zero")
	ErrTicketTypeInvalid = errors.New("ticket Type must be full or half")
)

type TicketType string

const (
	TicketTypeHalf TicketType = "half"
	TicketTypeFull TicketType = "full"
)

type Ticket struct {
	ID         string
	EventID    string
	Spot       *Spot
	TicketType TicketType
	Price      float64
}

func NewTicket(event *Event, spot *Spot, ticketType TicketType) (*Ticket, error) {
	if !IsValidTicketType(ticketType) {
		return nil, ErrTicketTypeInvalid
	}
	
	ticket := &Ticket{
		ID: uuid.NewString(),
		EventID: event.ID,
		Spot: spot,
		TicketType: ticketType,
		Price: event.Price,
	}

	ticket.CalculatePrice()
	if err := ticket.Validate(); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (t *Ticket) Validate() error {
	if t.Price <= 0 {
		return ErrTicketPriceZero
	}

	if IsValidTicketType(t.TicketType) {
		return ErrTicketTypeInvalid
	}

	return nil
}

func (t *Ticket) CalculatePrice() {
	if t.TicketType == TicketTypeHalf {
		t.Price /= 2
	}
}

func IsValidTicketType(ticketType TicketType) bool {
	return ticketType == TicketTypeFull || ticketType == TicketTypeHalf
}
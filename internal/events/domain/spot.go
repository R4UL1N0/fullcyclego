package domain

import (
	"errors"

	"github.com/google/uuid" 
)

var (
	ErrInvalidSpotNumber = errors.New("Invalid Spot number")
	ErrSpotNotFound = errors.New("Spot not found")
	ErrSpotAlreadyReserved = errors.New("Spot is already reserved")

	ErrSpotNameRequired = errors.New("Spot name is required")
	ErrSpotNameLength = errors.New("Spot name must be 2 characters long")
	ErrSpotNameStartWithLetter = errors.New("Spot name must start with a letter")
	ErrSpotNameEndWithNumber = errors.New("Spot name must end with a number")
)

type SpotStatus string

const (
	SpotStatusAvailable SpotStatus = "available"
	SpotStatusReserved  SpotStatus = "reserved"
	SpotStatusSold      SpotStatus = "sold"
)

type Spot struct {
	ID       string
	EventID  string
	Name     string
	Status   SpotStatus
	TicketID string
}

func NewSpot(event *Event, spotName string) (*Spot, error) {
	s := &Spot{ // & = Local na memória
		ID: uuid.New().String(),
		EventID: event.ID,
		Name: spotName,
		Status: SpotStatusAvailable,
	}

	err := s.Validate(); if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Spot) ReserveSpot(ticketID string) error {
	if s.Status == SpotStatusSold {
		return ErrSpotAlreadyReserved
	}

	s.Status = SpotStatusSold
	s.TicketID = ticketID

	return nil
}

func (s *Spot) Validate() error {
	if len(s.Name) == 0 {
		return ErrSpotNameRequired
	}

	if len(s.Name) < 2 || len(s.Name) > 2 {
		return ErrSpotNameLength
	}

	if s.Name[0] < 'A' || s.Name[0] > 'Z' {
		return ErrSpotNameStartWithLetter
	}

	if s.Name[1] < '0' || s.Name[1] > '9' {
		return ErrSpotNameEndWithNumber
	}

	return nil
}
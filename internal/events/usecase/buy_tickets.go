package usecase

import (
	"study/go/internal/events/domain"
	"study/go/internal/events/infra/service"
)

type BuyTicketsInputDTO struct {
	EventID    string   `json:"event_id"`
	Spots      []string `json:"spots"`
	TicketType string   `json:"ticket_type"`
	CardHash   string   `json:"card_hash"`
	Email      string   `json:"email"`
}

type BuyTicketsOutputDTO struct {
	Tickets []TicketDTO `json:"tickets"`
}

type BuyTicketsUseCase struct {
	repo domain.EventRepository
	partnerFactory service.PartnerFactory
}

func NewBuyTicketsUseCase(repo domain.EventRepository, partnerFactory service.PartnerFactory) *BuyTicketsUseCase {
	return &BuyTicketsUseCase{repo: repo, partnerFactory: partnerFactory}
}

func (uc *BuyTicketsUseCase) Execute(input BuyTicketsInputDTO) (*BuyTicketsOutputDTO, error) {
	event, err := uc.repo.FindEventByID(input.EventID)

	if err != nil {
		return nil, err
	}

	req := &service.ReservationRequest{
		EventID: input.EventID,
		Spots: input.Spots,
		TicketType: input.TicketType,
		CardHash: input.CardHash,
		Email: input.Email,
	}

	partnerService, err := uc.partnerFactory.CreatePartner(event.PartnerID)

	if err != nil {
		return nil, err
	}

	reservationRes, err := partnerService.MakeReservation(req)

	tickets := make([]domain.Ticket, len(reservationRes))
	for i, reservation := range reservationRes {
		spot, err := uc.repo.FindSpotByName(event.ID, reservation.Spot)

		if err != nil {
			return nil, err
		}

		ticket, err := domain.NewTicket(event, spot, domain.TicketType(req.TicketType))

		if err != nil {
			return nil, err
		}

		err = uc.repo.CreateTicket(ticket)

		if err != nil {
			return nil, err
		}

		spot.ReserveSpot(ticket.ID)
		err = uc.repo.ReserveSpot(spot.ID, ticket.ID)

		if err != nil {
			return nil, err
		}

		tickets[i] = *ticket
	}

	ticketsDTO := make([]TicketDTO, len(tickets))
	for i, ticket := range tickets {
		ticketDTO := TicketDTO{
			ID: ticket.ID,
			SpotID: ticket.Spot.ID,
			TicketType: string(ticket.TicketType),
			Price: ticket.Price,
		}

		ticketsDTO[i] = ticketDTO 
	}

	return &BuyTicketsOutputDTO{Tickets: ticketsDTO}, nil
}

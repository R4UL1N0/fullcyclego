package service

type PartnerFactory interface {
	CreatePartner(partnerID int) (*Partner, error)
}

type Partner struct {

}

func (p *Partner) MakeReservation(req *ReservationRequest) ([]ReservationResponse, error) {
	
}

type ReservationResponse struct {
	Spot      string `json:"spots"`
}

type ReservationRequest struct {
	EventID    string   `json:"event_id"`
	Spots      []string `json:"spots"`
	TicketType string   `json:"ticket_type"`
	CardHash   string   `json:"card_hash"`
	Email      string   `json:"email"`
}

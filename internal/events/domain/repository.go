package domain

type EventRepository interface {
	ListEvents() ([]Event, error)
	FindEventByID(eventId string) (*Event, error)
	FindSpotsByEventID(eventId string) ([]Spot, error)
	FindSpotByName(eventId, spotName string) (*Spot, error)
	CreateEvent(event *Event) error
	CreateSpot(spot *Spot) error
	CreateTicket(ticket *Ticket) error
	ReserveSpot(spotID, ticketID string) error 
}
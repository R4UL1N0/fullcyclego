package repository

import (
	"database/sql"
	"errors"
	"time"

	"study/go/internal/events/domain"

	_ "github.com/go-sql-driver/mysql"

	_ "study/go/internal/events/domain"
)

type mysqlEventRepository struct {
	db *sql.DB
}

func NewMysqlEventRepository(db *sql.DB) (*mysqlEventRepository, error) {
	return &mysqlEventRepository{db: db}, nil
}

func (r *mysqlEventRepository) CreateEvent(e *domain.Event) error {
	query := `
		INSERT INTO events (id, name, location, rating, date, image_url, capacity, price, partner_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query, e.ID,
		e.Name, e.Location,
		e.Rating, e.Date,
		e.ImageURL, e.Capacity,
		e.Price, e.PartnerID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlEventRepository) CreateSpot(spot *domain.Spot) error {
	query := `
		INSERT INTO spots (id, event_id, name, status, ticket_id) 
		VALUES (?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, spot.ID, spot.EventID, spot.Name, spot.Status, spot.TicketID) // first value returns register quantities

	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlEventRepository) CreateTicket(t *domain.Ticket) error {
	query := `
		INSERT INTO tickets (id, event_id, spot_id, ticket_type, price)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query, t.ID, t.EventID,
		t.Spot.ID, t.TicketType,
		t.Price,
	)

	return err
}

func (r *mysqlEventRepository) ListEvents() ([]domain.Event, error) {
	query := `
		SELECT e.id, e.name, e.location, e.rating, e.date, e.image_url, e.capacity, e.price, e.partner_id,
		s.id, s.event_id, s.name, s.status, s.ticket_id,
		t.id, t.event_id, t.spot_id, t.ticket_type, t.price
		FROM events e
		LEFT JOIN spots s ON s.event_id = e.id
		LEFT JOIN tickets t ON t.spot_id = s.id
	`

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	eventMap := make(map[string]*domain.Event)
	spotMap := make(map[string]*domain.Spot)

	for rows.Next() {
		var eventID, eventName, eventLocation, eventOrganization,
		eventRating, eventDate, eventImageURL, 
		eventPartnerID, spotID, spotEventID, 
		spotName, spotStatus, spotTicketID,
		ticketID, ticketEventID, ticketSpotID, ticketType sql.NullString
		var eventCapacity int
		var eventPrice, ticketPrice sql.NullFloat64
		var partnerID sql.NullInt32

		err := rows.Scan(
			&eventID, &eventName, &eventLocation, &eventCapacity, &eventRating, &eventDate, &eventImageURL,
			&eventPartnerID, &eventPrice, &spotID, &spotEventID, &spotName, &spotStatus, &spotTicketID, 
			&ticketID, &ticketEventID, &ticketSpotID, &ticketType, &ticketPrice,
		)

		if err != nil {
			return nil, err
		}

		if !eventID.Valid {
			continue
		}

		event, exists := eventMap[eventID.String]

		if !exists {
			eventDateParsed, err := time.Parse("2006-01-02 15:04:05", eventDate.String)
			if err != nil {
				return nil, err
			}

			event = &domain.Event{
				ID: eventID.String,
				Name: eventName.String,
				Location: eventLocation.String,
				Organization: eventOrganization.String,
				Date: eventDateParsed,
				Rating: domain.Rating(eventRating.String),
				Capacity: eventCapacity,
				ImageURL: eventImageURL.String,
				Price: eventPrice.Float64,
				PartnerID: int(partnerID.Int32),
				Spots: []domain.Spot{},
				Tickets: []domain.Ticket{},
			}

			eventMap[eventID.String] = event
		}

		if spotID.Valid {
			spot, spotExists := spotMap[spotID.String]

			if !spotExists {
				spot = &domain.Spot{
					ID: spotID.String,
					EventID: spotEventID.String,
					Name: spotName.String,
					Status: domain.SpotStatus(spotStatus.String),
					TicketID: spotTicketID.String,
				}
				event.Spots = append(event.Spots, *spot)
				spotMap[spotID.String] = spot
			}

			if ticketID.Valid {
				ticket := &domain.Ticket{
					ID: ticketID.String,
					EventID: ticketEventID.String,
					Spot: spot,
					TicketType: domain.TicketType(ticketType.String),
					Price: ticketPrice.Float64,
				}

				event.Tickets = append(event.Tickets, *ticket)
			}
		}
	}

	events := make([]domain.Event, 0, len(eventMap))
	for _, event := range eventMap {
		events = append(events, *event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *mysqlEventRepository) ReserveSpot(spotID, ticketID string) error {
	query := `
		UPDATE spots
		SET status = ?, ticket_id = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, domain.SpotStatusSold, ticketID, spotID)
	return err
}

func (r *mysqlEventRepository) FindEventByID(eventID string) (*domain.Event, error) {
	query := `
		SELECT id, name, location, organization, rating, 
		date, image_url, capacity, price, partner_id 
		FROM events
		WHERE id = ?
	`

	rows, err := r.db.Query(query, eventID)

	if err != nil {
		return nil, err
	}

	defer rows.Close() // avoids data (?) leakage

	var event *domain.Event

	err = rows.Scan(
		&event.ID, &event.Name, &event.Location, &event.Organization, &event.Rating,
		&event.Date, &event.ImageURL, &event.Capacity, &event.Price, &event.PartnerID)

	if err != nil {
		return nil, err
	}

	return event, nil
}

func (r *mysqlEventRepository) FindSpotsByEventID(eventID string) ([]domain.Spot, error) {
	query := `
		SELECT id, event_id, name, status, ticket_id
		FROM spots
		WHERE event_id = ? 
	`

	rows, err := r.db.Query(query, eventID)

	if err != nil {
		return nil, err
	}

	var spots []domain.Spot

	for rows.Next() {
		var spot domain.Spot

		err = rows.Scan(&spot.ID, &spot.Name, &spot.Status, &spot.TicketID); if err != nil {
			return nil, err
		}

		spots = append(spots, spot)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return spots, nil
}

func (r *mysqlEventRepository) FindSpotByName(eventID, name string) (*domain.Spot, error) {
	query := `
		SELECT s.id, s.event_id, s.name, s.status, s.ticket_id,
		t.id, t.event_id, t.spot_id, t.ticket_type, t.price
		FROM spots s
		LEFT JOIN tickets t ON s.id = t.spot_id
		WHERE s.event_id = ? AND s.name = ?
	`

	row := r.db.QueryRow(query, eventID, name)

	var spot domain.Spot
	var ticket domain.Ticket
	var ticketID, ticketEventID, ticketSpotID, ticketType sql.NullString
	var ticketPrice sql.NullFloat64

	err := row.Scan(
		&spot.ID, &spot.EventID, &spot.Name, &spot.Status, &spot.TicketID,
		&ticketID, &ticketEventID, &ticketSpotID, &ticketType, &ticketPrice,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSpotNotFound
		}
	}

	if ticketID.Valid {
		ticket.ID = ticketID.String
		ticket.EventID = ticketEventID.String
		ticket.Spot = &spot
		ticket.TicketType = domain.TicketType(ticketType.String)
		ticket.Price = ticketPrice.Float64
		spot.TicketID = ticket.ID
	}

	return &spot, nil

}

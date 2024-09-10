package events

import (
	"database/sql"
	"net/http"
	httpHandler "study/go/internal/events/infra/http"
	"study/go/internal/events/infra/repository"
	"study/go/internal/events/usecase"
)

func main() {
	db, err := sql.Open(
		"mysql", 
		"test_user:test_pass@tcp(localhost:3306)/test_db",
	)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	eventRepo, err := repository.NewMysqlEventRepository(db)

	if err != nil {
		panic(err)
	}

	partnerBaseURLs := map[int]string{
		1: "http://localhost:9080/api1",
		2: "http://localhost:9080/api2",
	}

	partnerFactory := service.NewPartnerFactory(partnerBaseURLs)

	listEventsUseCase := usecase.NewListEventsUsecase(eventRepo)
	getEventUseCase := usecase.NewGetEventUseCase(eventRepo)
	listSpotsUseCase := usecase.NewListSpotsUseCase(eventRepo)
	buyTicketsUseCase := usecase.NewBuyTicketsUseCase(eventRepo)

	eventsHandler := httpHandler.NewEventsHandler(listEventsUseCase, listSpotsUseCase, getEventUseCase, buyTicketsUseCase)

	r := http.NewServeMux()
	r.HandleFunc("GET /events", eventsHandler.ListEvents)
	r.HandleFunc("GET /events/{eventID}", eventsHandler.GetEvent)
	r.HandleFunc("GET /events/{eventID}/spots", eventsHandler.ListSpots)
	r.HandleFunc("POST /checkout", eventsHandler.BuyTickets)

	http.ListenAndServe(":8080", r)
}
package model

import (
	"log"
	"scala-disaster-adviser/event-service/database"
	"time"
)

type EventSearch struct {
	Start     int64
	End       int64
	Radius    int64
	Longitude float64
	Latitude  float64
}

type Event struct {
	id        int
	Summary   string    `json:"description"`
	Uuid      string    `json:"name"`
	Location  string    `json:"location"`
	Latitude  float64   `json:"lat"`
	Longitude float64   `json:"lon"`
	StartDate time.Time `json:"date"`
	endDate   time.Time
}

type EventResponse struct {
	Events []Event `json:"events"`
}

func (e EventSearch) FindEvents() EventResponse {
	var empty = make([]Event, 0)
	var response = EventResponse{empty}
	var events []Event

	db, err := database.Connect()
	if err != nil {
		log.Printf("Database error %s", err)
		return response
	}

	rows, err := db.Query("SELECT e.id, e.summary, u.uuid, e.location, e.latitude, e.longitude, e.start_date, e.end_date "+
		" FROM events e "+
		" INNER JOIN users u ON u.id = e.user_id "+
		" WHERE e.start_date >= $1 AND e.end_date <= $2 "+
		" AND ST_DWithin(e.geom, ST_MakePoint($3,$4)::geography, $5)", time.Unix(e.Start, 0), time.Unix(e.End, 0), e.Longitude, e.Latitude, e.Radius)
	if err != nil {
		log.Printf("Wrong request to the db, %v", err)
		return response
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.id, &event.Summary, &event.Uuid, &event.Location, &event.Latitude, &event.Longitude, &event.StartDate, &event.endDate)
		if err == nil {
			events = append(events, event)
		} else {
			log.Printf("Error %s", err)
		}
	}

	if len(events) == 0 {
		return response
	}

	return EventResponse{events}
}

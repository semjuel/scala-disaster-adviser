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
	Id        int       `json:"id"`
	Summary   string    `json:"summary"`
	Location  string    `json:"location"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	StartDate time.Time `json:"start"`
	EndDate   time.Time `json:"end"`
}

func (e EventSearch) FindEvents() []Event {
	var response = make([]Event, 0)
	var events []Event

	db, err := database.Connect()
	if err != nil {
		log.Printf("Database error %s", err)
		return response
	}

	rows, err := db.Query("SELECT id, summary, location, latitude, longitude, start_date, end_date"+
		" FROM events WHERE start_date >= $1 AND end_date <= $2 "+
		" AND ST_DWithin(geom, ST_MakePoint($3,$4)::geography, $5)", time.Unix(e.Start, 0), time.Unix(e.End, 0), e.Longitude, e.Latitude, e.Radius)
	if err != nil {
		log.Printf("Wrong request to the db, %v", err)
		return response
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.Summary, &event.Location, &event.Latitude, &event.Longitude, &event.StartDate, &event.EndDate)
		if err == nil {
			events = append(events, event)
		} else {
			log.Printf("Error %s", err)
		}
	}

	if len(events) == 0 {
		return response
	}

	return events
}

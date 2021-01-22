package model

import (
	"log"
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
	UserId    int64
	EventId   string
	Summary   string  `json:"description"`
	Uuid      string  `json:"name"`
	Location  string  `json:"location"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Date      int64   `json:"date"`
	StartDate time.Time
	EndDate   time.Time
}

type EventResponse struct {
	Events []Event `json:"events"`
}

func (e EventSearch) FindEvents() EventResponse {
	var empty = make([]Event, 0)
	var response = EventResponse{empty}
	var events []Event

	err := DB.Connect()
	if err != nil {
		log.Printf("database error %s", err)
		return response
	}

	rows, err := DB.Instance.Query("SELECT e.id, e.summary, u.uuid, e.location, e.latitude, e.longitude, e.start_date, e.end_date "+
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
		err := rows.Scan(&event.id, &event.Summary, &event.Uuid, &event.Location, &event.Latitude, &event.Longitude, &event.StartDate, &event.EndDate)
		if err == nil {
			event.Date = event.StartDate.UnixNano() / int64(time.Millisecond)

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

func SaveEvent(event Event) error {
	insertQuery := "INSERT INTO events (user_id, event_id, summary, location, latitude, longitude, start_date, end_date, created, updated) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	_, err := DB.Instance.Exec(insertQuery, event.UserId, event.EventId, event.Summary, event.Location, event.Latitude,
		event.Longitude, event.StartDate.Format(time.RFC3339), event.EndDate.Format(time.RFC3339),
		event.StartDate.Format(time.RFC3339), event.StartDate.Format(time.RFC3339))
	if err != nil {
		return err
	}

	_, err = DB.Instance.Exec("UPDATE events SET geom = ST_SetSRID(ST_MakePoint(longitude, latitude), 4326) WHERE geom IS NULL;")

	return err
}

func DeleteUserEvents(userId int64) {
	err := DB.Connect()
	if err != nil {
		log.Printf("database error %s", err)
		return
	}

	_, err = DB.Instance.Query("DELETE FROM events WHERE user_id = $1", userId)
	if err != nil {
		log.Printf("delete events error %s", err)
		return
	}
}

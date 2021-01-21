package task

import (
	"encoding/json"
	"log"
	"net/http"
	"scala-disaster-adviser/event-service/broker"
	"scala-disaster-adviser/event-service/model"
	"time"
)

type CalendarResponse struct {
	Items []CalendarEvent `json:"items"`
}

type CalendarEvent struct {
	Id       string   `json:"id"`
	Summary  string   `json:"summary"`
	Location string   `json:"location"`
	Start    DateTime `json:"start"`
	End      DateTime `json:"end"`
}

type DateTime struct {
	DateTime string `json:"dateTime"`
}

func FetchEvents() {
	log.Printf("FetchEvents...")

	users, err := model.UserAll()
	if err != nil {
		log.Printf("error %s", err)
		return
	}

	for _, user := range users {
		events, err := makeRequest(user.Token)
		if err != nil {
			log.Println(err)
			log.Printf("error %s", err)
			continue
		}
		log.Println(events)

		model.DeleteUserEvents(user.Id)
		for _, e := range events {
			log.Printf("Start processing %s", e.Summary)
			if e.Location == "" {
				log.Printf("empty location")
				return
			}
			lat, lon := getCoordinates(e.Location)

			start, _ := time.Parse(time.RFC3339, e.Start.DateTime)
			end, _ := time.Parse(time.RFC3339, e.End.DateTime)

			event := model.Event{
				UserId:    user.Id,
				EventId:   e.Id,
				Summary:   e.Summary,
				Location:  e.Location,
				Latitude:  lat,
				Longitude: lon,
				StartDate: start,
				EndDate:   end,
			}
			err = model.SaveEvent(event)
			if err != nil {
				log.Printf("error %s", err)
				continue
			}

			msg := broker.Message{
				Name:        user.Uuid,
				Date:        start.Unix(),
				Lat:         lat,
				Lon:         lon,
				Description: e.Summary,
			}

			broker.SendEvent(msg)
		}
	}
}

func makeRequest(token string) ([]CalendarEvent, error) {
	var res []CalendarEvent
	url := "https://www.googleapis.com/calendar/v3/calendars/primary/events?singleEvents=true"

	var bearer = "Bearer " + token
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return res, err
	}
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}

	decoder := json.NewDecoder(resp.Body)
	var data CalendarResponse
	err = decoder.Decode(&data)
	if err != nil {
		return res, err
	}

	return data.Items, nil
}

func getCoordinates(location string) (float64, float64) {
	return 49.841952, 24.0315921
}

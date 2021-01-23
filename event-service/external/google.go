package external

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"scala-disaster-adviser/event-service/model"
	"time"
)

func GoogleCalendarEvents(token string) []model.CalendarEvent {
	var res []model.CalendarEvent

	today := time.Now().Format(time.RFC3339)
	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/primary/events?singleEvents=true&timeMin=%s&maxResults=1000", url.QueryEscape(today))

	var bearer = "Bearer " + token
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return res
	}
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return res
	}

	decoder := json.NewDecoder(resp.Body)
	var data model.CalendarResponse
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		return res
	}

	return data.Items
}

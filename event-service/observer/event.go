package observer

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"log"
	"scala-disaster-adviser/event-service/broker"
	"scala-disaster-adviser/event-service/external"
	"scala-disaster-adviser/event-service/model"
	"time"
)

func EventObserver(eventCh chan rxgo.Item) {
	observable := rxgo.FromChannel(eventCh)

	observable.
		Filter(func(item interface{}) bool {
			event := item.(model.CalendarEvent)

			return event.Location != ""
		}).
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			e := item.(model.CalendarEvent)
			log.Printf("start processing %s", e.Summary)

			lon, lat := external.Coordinates(e.Location)

			start, _ := time.Parse(time.RFC3339, e.Start.DateTime)
			if e.Start.DateTime == "" {
				start, _ = time.Parse("2006-01-02", e.Start.Date)
			}
			end, _ := time.Parse(time.RFC3339, e.End.DateTime)
			if e.End.DateTime == "" {
				end, _ = time.Parse("2006-01-02", e.End.Date)
			}

			event := model.Event{
				UserId:    e.User.Id,
				EventId:   e.Id,
				Summary:   e.Summary,
				Location:  e.Location,
				Latitude:  lat,
				Longitude: lon,
				StartDate: start,
				EndDate:   end,
			}

			err := model.SaveEvent(event)
			if err != nil {
				log.Printf("error %s", err)
			}

			msg := broker.Message{
				Name:        e.User.Uuid,
				Date:        start.Unix(),
				Lat:         lat,
				Lon:         lon,
				Description: e.Summary,
			}

			broker.SendEvent(msg)

			return e, nil
		},
			rxgo.WithPool(8),
			rxgo.Serialize(func(item interface{}) int {
				event := item.(model.CalendarEvent)
				return int(event.Uuid)
			}),
			rxgo.WithBufferedChannel(1))

	observable.Observe()
}

package producer

import (
	"github.com/reactivex/rxgo/v2"
	"scala-disaster-adviser/event-service/external"
	"scala-disaster-adviser/event-service/model"
	"time"
)

func GoogleCalendarEventsProducer(user model.User, ch chan rxgo.Item) {
	events := external.GoogleCalendarEvents(user.Token)

	for _, event := range events {
		event.User = user
		event.Uuid = time.Now().UnixNano()

		ch <- rxgo.Of(event)
	}
}

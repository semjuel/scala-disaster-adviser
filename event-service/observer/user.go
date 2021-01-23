package observer

import (
	"context"
	"github.com/reactivex/rxgo/v2"
	"scala-disaster-adviser/event-service/model"
	"scala-disaster-adviser/event-service/producer"
)

func UserObserver(eventCh chan rxgo.Item) {
	ch := make(chan rxgo.Item)
	go producer.UserProducer(ch)
	observable := rxgo.FromChannel(ch)

	observable.
		Filter(func(item interface{}) bool {
			user := item.(model.User)

			return user.Token != ""
		}).
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			user := item.(model.User)
			model.DeleteUserEvents(user.Id)
			producer.GoogleCalendarEventsProducer(user, eventCh)

			return user, nil
		},
			rxgo.WithPool(8),
			rxgo.Serialize(func(item interface{}) int {
				user := item.(model.User)
				return int(user.Id)
			}),
			rxgo.WithBufferedChannel(1))

	observable.Observe()
}

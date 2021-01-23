package task

import (
	"github.com/reactivex/rxgo/v2"
	"log"
	"scala-disaster-adviser/event-service/observer"
)

func FetchEvents(eventCh chan rxgo.Item) {
	log.Printf("start fetching events ...")

	observer.UserObserver(eventCh)
	observer.EventObserver(eventCh)
}

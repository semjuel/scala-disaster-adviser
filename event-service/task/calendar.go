package task

import (
	"scala-disaster-adviser/event-service/broker"
)

func FetchEvents() {
	broker.SendEvent()
}

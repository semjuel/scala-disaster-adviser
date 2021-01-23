package producer

import (
	"github.com/reactivex/rxgo/v2"
	"log"
	"scala-disaster-adviser/event-service/model"
)

func UserProducer(ch chan rxgo.Item) {
	users, err := model.UserAll()
	if err != nil {
		log.Printf("error %s", err)
		return
	}

	for _, user := range users {
		ch <- rxgo.Of(user)
	}
}

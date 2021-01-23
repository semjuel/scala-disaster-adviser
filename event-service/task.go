package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"github.com/reactivex/rxgo/v2"
	"log"
	"scala-disaster-adviser/event-service/model"
	"scala-disaster-adviser/event-service/task"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = model.DB.Connect()
	if err != nil {
		log.Printf("database error %s", err)
		return
	}
	defer model.DB.Instance.Close()

	// Drop  all data in the database.
	model.DropAll()

	eventCh := make(chan rxgo.Item)

	gocron.Every(1).Minute().Do(task.FetchEvents, eventCh)
	<-gocron.Start()
}

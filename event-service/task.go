package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
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

	gocron.Every(10).Seconds().Do(task.FetchEvents)
	<-gocron.Start()
}

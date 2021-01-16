package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"log"
	"scala-disaster-adviser/event-service/task"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	gocron.Every(10).Seconds().Do(task.FetchEvents)
	<-gocron.Start()
}

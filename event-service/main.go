package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"scala-disaster-adviser/event-service/handler"
	"scala-disaster-adviser/event-service/util"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	path := util.GetProjectDir()
	static := fmt.Sprintf("%s/static", path)

	http.Handle("/", http.FileServer(http.Dir(static)))
	http.HandleFunc("/users", handler.UsersHandler)
	http.HandleFunc("/events", handler.EventsHandler)
	log.Fatal(http.ListenAndServe(":8002", nil))
}

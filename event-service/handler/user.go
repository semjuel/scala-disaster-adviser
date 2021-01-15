package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"scala-disaster-adviser/event-service/model"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	if r.Method == "OPTIONS" {
		return
	}
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var u model.User
		err := decoder.Decode(&u)
		if err != nil {
			log.Printf("Decode error %s \n", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err = model.UserSave(u)
		if err != nil {
			log.Printf("Error %s", err)
		}

		w.WriteHeader(http.StatusOK)
		break

	default:
		log.Printf("Request method %s not allowed \n", r.Method)
		http.Error(w, "Request method not allowed", http.StatusMethodNotAllowed)
		break
	}
}

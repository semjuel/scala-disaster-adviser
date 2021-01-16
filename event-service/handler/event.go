package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"scala-disaster-adviser/event-service/model"
	"strconv"
	"time"
)

type request struct {
	Date     date     `json:"date"`
	Location location `json:"location"`
}

type date struct {
	Timestamp string `json:"timestamp"`
	Range     int64  `json:"range"`
}

type location struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Range float64 `json:"range"`
}

func EventsHandler(w http.ResponseWriter, r *http.Request) {
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
		var req request
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("Decode error %s \n", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		eventSearch := model.EventSearch{
			Start:     getStartDate(req),
			End:       getEndDate(req),
			Radius:    getRadius(req),
			Longitude: getLongitude(req),
			Latitude:  getLatitude(req),
		}

		events := eventSearch.FindEvents()

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(events)
		if err != nil {
			log.Printf("Encoder error %v \n", err)
		}
		break

	default:
		log.Printf("Request method %s not allowed \n", r.Method)
		http.Error(w, "Request method not allowed", http.StatusMethodNotAllowed)
		break
	}
}

func getRadius(r request) int64 {
	return int64(r.Location.Range * 1000)
}

func getStartDate(r request) int64 {
	now := time.Now()
	sec := now.Unix()

	value, err := strconv.ParseInt(r.Date.Timestamp, 10, 32)
	if err != nil {
		return sec - r.Date.Range
	}

	return value - r.Date.Range
}

func getEndDate(r request) int64 {
	now := time.Now()
	sec := now.Unix()

	value, err := strconv.ParseInt(r.Date.Timestamp, 10, 32)
	if err != nil {
		return sec + r.Date.Range
	}

	return value + r.Date.Range
}

func getLongitude(r request) float64 {
	return r.Location.Lon
}

func getLatitude(r request) float64 {
	return r.Location.Lat
}

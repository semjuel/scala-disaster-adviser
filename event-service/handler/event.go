package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"scala-disaster-adviser/event-service/model"
	"strconv"
	"time"
)

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
		// @TODO remove this
		//kw := &kafka.Writer{
		//	Addr:  kafka.TCP("localhost:9092"),
		//	Topic: "user-events",
		//}
		//
		//err := kw.WriteMessages(context.Background(),
		//	kafka.Message{
		//		Key:   []byte("Key-A"),
		//		Value: []byte("Hello World!"),
		//	},
		//)
		//if err != nil {
		//	log.Printf("Kafka connection error %s", err)
		//	w.WriteHeader(http.StatusOK)
		//	return
		//}
		//
		//kw.Close()

		log.Printf("Success")
		w.WriteHeader(http.StatusOK)
		break

	case http.MethodGet:
		eventSearch := model.EventSearch{
			Start:     getStartDate(r),
			End:       getEndDate(r),
			Radius:    getRadius(r),
			Longitude: getLongitude(r),
			Latitude:  getLatitude(r),
		}

		events := eventSearch.FindEvents()

		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(events)
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

func getRadius(r *http.Request) int64 {
	// Default radius = 5km.
	return getIntParamFromQuery(r, "radius", 5000)
}

func getStartDate(r *http.Request) int64 {
	now := time.Now()
	sec := now.Unix()
	return getIntParamFromQuery(r, "start", sec)
}

func getEndDate(r *http.Request) int64 {
	now := time.Now()
	sec := now.Unix() + 3600
	return getIntParamFromQuery(r, "end", sec)
}

func getLongitude(r *http.Request) float64 {
	return getFloatParamFromQuery(r, "longitude", 0)
}

func getLatitude(r *http.Request) float64 {
	return getFloatParamFromQuery(r, "latitude", 0)
}

func getIntParamFromQuery(r *http.Request, name string, defaultValue int64) int64 {
	params, ok := r.URL.Query()[name]
	if !ok || len(params) < 1 {
		return defaultValue
	}

	value, err := strconv.ParseInt(params[0], 10, 32)
	if err != nil {
		return defaultValue
	}

	return value
}

func getFloatParamFromQuery(r *http.Request, name string, defaultValue float64) float64 {
	params, ok := r.URL.Query()[name]
	if !ok || len(params) < 1 {
		return defaultValue
	}

	value, err := strconv.ParseFloat(params[0], 32)
	if err != nil {
		return defaultValue
	}

	return value
}

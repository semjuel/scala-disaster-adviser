package external

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type MapBoxResponse struct {
	Type     string           `json:"FeatureCollection"`
	Features []MapBoxFeatures `json:"features"`
}

type MapBoxFeatures struct {
	Id     string    `json:"id"`
	Center []float64 `json:"center"`
}

func Coordinates(location string) (float64, float64) {
	url := fmt.Sprintf("https://api.mapbox.com/geocoding/v5/mapbox.places/%s.json?access_token=%s",
		url.QueryEscape(location),
		"pk.eyJ1IjoiZGlzYXN0ZXItYWR2aXNlciIsImEiOiJja2s3MGRmMTMwN3lnMnZvMnpvczJ2YXlwIn0.7LrCNyVfyG3GmStj2Pl6NA")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return 0, 0
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return 0, 0
	}

	decoder := json.NewDecoder(resp.Body)
	var data MapBoxResponse
	err = decoder.Decode(&data)
	if err != nil {
		log.Println(err)
		return 0, 0
	}

	features := data.Features
	if len(features) > 0 {
		feature := features[0]
		if len(feature.Center) > 0 {
			return feature.Center[0], feature.Center[1]
		}
	}

	return 0, 0
}

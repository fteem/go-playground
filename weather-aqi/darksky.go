package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	darkskyAPIKey = "YOUR_API_KEY_HERE"
	darkskyAPIURL = "https://api.darksky.net/forecast/%s/%s,%s"
)

type darkskyResponse struct {
	Currently struct {
		Temperature         float64 `json:"temperature"`
		ApparentTemperature float64 `json:"apparentTemperature"`
	} `json:"currently"`
}

type darkskyClient struct {
	latitude  string
	longitude string
}

func newDarkskyClient(longitude string, latitude string) *darkskyClient {
	return &darkskyClient{
		longitude: longitude,
		latitude:  latitude,
	}
}

func (d *darkskyClient) apiURL() string {
	return fmt.Sprintf(darkskyAPIURL, darkskyAPIKey, d.latitude, d.longitude)
}

func (d *darkskyClient) temperature() string {
	var dsResponse darkskyResponse

	response, err := http.Get(d.apiURL())
	if err != nil {
		log.Fatalf("Error fetching temperature: %s", err)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
	}
	defer response.Body.Close()

	err = json.Unmarshal(bodyBytes, &dsResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling response: %s", err)
	}

	return fmt.Sprintf("%f", dsResponse.Currently.Temperature)
}

func (d *darkskyClient) work(out chan<- string) {
	out <- d.temperature()
}

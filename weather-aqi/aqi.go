package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	aqiURL = "https://api.openaq.org/v1/latest?city=%s&parameter=pm10"
)

type aqiResponse struct {
	Results []struct {
		Measurements []struct {
			Value float64 `json:"value"`
		} `json:"measurements"`
	}
}

type aqiClient struct {
	city string
}

func newAqiClient(city string) *aqiClient {
	return &aqiClient{city}
}

func (a *aqiClient) apiURL() string {
	return fmt.Sprintf(aqiURL, a.city)
}

func (a *aqiClient) index() string {
	var aqi aqiResponse

	response, err := http.Get(a.apiURL())
	if err != nil {
		log.Fatalf("Error fetching AQI: %s", err)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
	}
	defer response.Body.Close()

	err = json.Unmarshal(bodyBytes, &aqi)
	if err != nil {
		log.Fatalf("Error unmarshalling response: %s", err)
	}

	var index float64
	if len(aqi.Results) > 0 {
		index = aqi.Results[0].Measurements[0].Value
	}
	return fmt.Sprintf("%f", index)
}

func (a *aqiClient) work(out chan<- string) {
	out <- a.index()
}

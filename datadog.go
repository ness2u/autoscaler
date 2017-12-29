package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// todo: add serialization hints
type Response struct {
	To_date      uint
	Resp_version int
	Query        string
	Message      string
	// group by

	Status    string
	Res_type  string
	From_date uint
	Series    []Metric
}

type Metric struct {
	Metric      string
	QueryIndex  int
	DisplayName string
	// unit
	PointList [][]float64
}

type DDConfig struct {
	BaseUri string
	ApiKey  string
	AppKey  string
}

func GetDataDogMetrics(query string, seconds int64) float64 {

	response := queryDatadog(query, seconds)
	avg := calcAvg(response)
	fmt.Printf("metric: %s; %f\n", response.Series[0].Metric, avg)
	return avg
}

func calcAvg(resp Response) float64 {

	if len(resp.Series) != 1 {
		fmt.Printf("series len: %d\n", len(resp.Series))
		panic("i can only handle single-series queries at the moment")
	}

	metric := resp.Series[0]

	total := 0.0
	for _, points := range metric.PointList {
		total += points[1]
	}
	avg := total / float64(len(metric.PointList))
	return avg
}

func loadDDConfigFromEnv() DDConfig {
	baseUri := requiredEnv("DD_URI")
	apiKey := requiredEnv("DD_API_KEY")
	appKey := requiredEnv("DD_APP_KEY")

	return DDConfig{baseUri, apiKey, appKey}
}

func queryDatadog(query string, seconds int64) Response {

	config := loadDDConfigFromEnv()
	to := int64(time.Now().Unix())
	from := to - seconds

	//	fmt.Printf("from: %d, to: %d\n", from, to)

	uri := fmt.Sprintf("%s?api_key=%s&application_key=%s&from=%d&to=%d&query=%s",
		config.BaseUri, config.ApiKey, config.AppKey, from, to, url.QueryEscape(query))

	//	fmt.Printf("uri: %s\n", uri)

	response := new(Response)
	if err := getJson(uri, response); err != nil {
		panic(err.Error())
	}
	return *response
}

func getJson(uri string, target interface{}) error {
	req, _ := http.NewRequest("GET", uri, nil)

	var client = &http.Client{Timeout: 10 * time.Second}
	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

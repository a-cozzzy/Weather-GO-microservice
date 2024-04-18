package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey`
}

type weatherData struct {
    Name string `json:"name"`
    Main struct {
        Temp float64 `json:"temp"`
        FeelsLike float64 `json:"feels_like"`
        MaxTemp float64 `json:"temp_max"`
        MinTemp float64 `json:"temp_min"`
        Humidity int `json:"humidity"`
    } `json:"main"`
    Wind struct {
        Speed float64 `json:"speed"`
    } `json:"wind"`
    Rain struct {
        Rain float64 `json:"1h"` // Adjust the field name based on your API response
    } `json:"rain"`
	Weather []struct {
        Description string `json:"description"`
        Icon        string `json:"icon"`
    } `json:"weather"`
}


func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}

func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    http.HandleFunc("/", serveIndex)
    http.HandleFunc("/weather/",
        func(w http.ResponseWriter, r *http.Request) {
            city := strings.SplitN(r.URL.Path, "/", 3)[2]
            data, err := query(city)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            w.Header().Set("Content-Type", "application/json; charset=utf-8")
            json.NewEncoder(w).Encode(data)
        })

    http.ListenAndServe(":8080", nil)
}

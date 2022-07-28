package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	BaseWeatherURL = "https://wttr.in/%s?format=j1"
)

type TemperatureData struct {
	TempCelsius    string `json:"temp_C"`
	TempFahrenheit string `json:"temp_F"`
}

type WeatherData struct {
	CurrentCondition []TemperatureData `json:"current_condition"`
}

func main() {
	location := flag.String("location", "05401", "Where would you like to track temperature? (e.g. 05401, 'Moscow', 'SFO')")

	weatherURL := fmt.Sprintf(BaseWeatherURL, *location)
	weatherData, err := loadWeatherData(weatherURL)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%v\n", weatherData)
}

func loadWeatherData(url string) (WeatherData, error) {
	var weatherData WeatherData

	response, err := http.Get(url)
	if err != nil {
		return weatherData, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return weatherData, err
	}

	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return weatherData, err
	}

	return weatherData, nil
}

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
	location := flag.String("loc", "05401", "The location where would you like to track temperature (e.g. 05401, Moscow, SFO)")
	verbose := flag.Bool("v", false, "Verbose mode enables additional logging")
	flag.Parse()

	weatherURL := fmt.Sprintf(BaseWeatherURL, *location)
	weatherData, err := loadWeatherData(weatherURL, *verbose)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%v\n", weatherData)
}

func loadWeatherData(url string, verbose bool) (WeatherData, error) {
	var weatherData WeatherData

	if verbose {
		fmt.Printf("Fetching weather data from %s\n", url)
	}

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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	useCelsius := flag.Bool("c", false, "Use celsius temperature measurements, otherwise defaults to Fahrenheit")
	tempThreshold := flag.Int("temp", 75, "The temperature threshold to watch for, in fahrenheit. Use the '-c' flag to use celsius instead.")
	flag.Parse()

	passedThreshold, currentTemp, err := checkTempThreshold(*location, *tempThreshold, *useCelsius, *verbose)
	if err != nil {
		log.Fatalf("Something went wrong trying to fetch the weather: %v", err)
	}

	if passedThreshold {
		fmt.Printf("It's %d˚ out! That's hotter than your %d˚ threshold.\n", currentTemp, *tempThreshold)
	} else {
		fmt.Printf("It's %d˚ degrees out, which is below your %d˚ degree threshold.\n", currentTemp, *tempThreshold)
	}
}

// Checks if the current temperature at a specific location has surpassed a
// specified temperature threshold. If it has, returns `true` as the first argument,
// otherwise it returns `false`. The current temperature is returned as the second
// argument. An error is returned if it fails to fetch weather data.
func checkTempThreshold(location string, tempThreshold int, useCelsius, verbose bool) (passedThreshold bool, currentTemp int, err error) {
	weatherURL := fmt.Sprintf(BaseWeatherURL, location)
	weatherData, err := loadWeatherData(weatherURL, verbose)
	if err != nil {
		return false, currentTemp, err
	}

	if useCelsius {
		currentTemp, err = strconv.Atoi(weatherData.CurrentCondition[0].TempCelsius)
	} else {
		currentTemp, err = strconv.Atoi(weatherData.CurrentCondition[0].TempFahrenheit)
	}

	if err != nil {
		return false, currentTemp, nil
	}

	return currentTemp > tempThreshold, currentTemp, nil
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

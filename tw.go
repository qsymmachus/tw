package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
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

type FlagOptions struct {
	Location      string
	TempThreshold int
	Interval      int
	UseCelsius    bool
	Verbose       bool
}

func main() {
	opts := parseFlags()

	fmt.Printf("Watching for the temperature to exceed %d˚...\n", opts.TempThreshold)

	err := checkTempThreshold(opts.Location, opts.TempThreshold, opts.UseCelsius, opts.Verbose)
	if err != nil {
		log.Fatalf("Something went wrong trying to fetch the weather: %v", err)
	}

	tick := time.Tick(time.Duration(opts.Interval) * time.Minute)
	for {
		select {
		case <-tick:
			err := checkTempThreshold(opts.Location, opts.TempThreshold, opts.UseCelsius, opts.Verbose)
			if err != nil {
				log.Fatalf("Something went wrong trying to fetch the weather: %v", err)
			}
		}
	}
}

// Parse option flags from the command line.
func parseFlags() FlagOptions {
	location := flag.String("loc", "05401", "The location where would you like to track temperature (e.g. 05401, Moscow, SFO)")
	verbose := flag.Bool("v", false, "Verbose mode enables additional logging")
	useCelsius := flag.Bool("c", false, "Use celsius temperature measurements, otherwise defaults to Fahrenheit")
	tempThreshold := flag.Int("temp", 75, "The temperature threshold to watch for, in fahrenheit. Use the '-c' flag to use celsius instead.")
	minutes := flag.Int("m", 1, "How frequently to check the temperature, in minutes")
	flag.Parse()

	return FlagOptions{
		Location:      *location,
		TempThreshold: *tempThreshold,
		Interval:      *minutes,
		Verbose:       *verbose,
		UseCelsius:    *useCelsius,
	}
}

// Checks if the current temperature at a specific location has surpassed a
// specified temperature threshold, and prints a message with the result.
// An error is returned if it fails to fetch weather data.
func checkTempThreshold(location string, tempThreshold int, useCelsius, verbose bool) error {
	weatherURL := fmt.Sprintf(BaseWeatherURL, location)
	weatherData, err := loadWeatherData(weatherURL, verbose)
	if err != nil {
		return err
	}

	var currentTemp int
	if useCelsius {
		currentTemp, err = strconv.Atoi(weatherData.CurrentCondition[0].TempCelsius)
	} else {
		currentTemp, err = strconv.Atoi(weatherData.CurrentCondition[0].TempFahrenheit)
	}

	if err != nil {
		return err
	}

	if currentTemp > tempThreshold {
		fmt.Printf("It's %d˚ out! That's hotter than your %d˚ threshold.\n", currentTemp, tempThreshold)
		fmt.Print("\a") // Ding the bell
	} else if currentTemp < tempThreshold {
		fmt.Printf("It's %d˚ degrees out, which is below your %d˚ degree threshold.\n", currentTemp, tempThreshold)
	} else {
		fmt.Printf("It's %d˚ degrees out, which equal to your %d˚ degree threshold.\n", currentTemp, tempThreshold)
	}

	return nil
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

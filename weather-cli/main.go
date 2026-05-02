package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type GeoResponse struct {
	Results []Location `json:"results"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
}

type WeatherResponse struct {
	CurrentWeather CurrentCondition `json:"current_weather"`
}

type CurrentCondition struct {
	Time        string  `json:"time"`
	Temperature float64 `json:"temperature"`
	Windspeed   float64 `json:"windspeed"`
	Weathercode int     `json:"weathercode"`
}

func getWeatherDescription(code int) string {
	switch {
	case code == 0:
		return "Clear sky"
	case code == 1 || code == 2 || code == 3:
		return "Partly cloudy"
	case code == 45 || code == 48:
		return "Fog"
	case (code >= 51 && code <= 67) || (code >= 80 && code <= 82):
		return "Rain / Drizzle"
	case (code >= 71 && code <= 77) || code == 85 || code == 86:
		return "Snow"
	case code >= 95 && code <= 99:
		return "Thunderstorm"
	default:
		return "Unknown weather"
	}
}

func main() {
	cityPtr := flag.String("city", "", "The name of city to get weather for")
	flag.Parse()

	if *cityPtr == "" {
		fmt.Println("Error: provide a city")
		fmt.Println("Usage: weather --city \"New Delhi\"")
		os.Exit(1)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	escapedCity := url.QueryEscape(*cityPtr)
	geoURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1", escapedCity)

	respGeo, err := client.Get(geoURL)
	if err != nil {
		fmt.Println("Error fetching coordinates:", err)
		return
	}

	defer respGeo.Body.Close()

	var geoData GeoResponse
	err = json.NewDecoder(respGeo.Body).Decode(&geoData)

	if len(geoData.Results) == 0 {
		fmt.Println("enter correct city")
		return
	}

	lat := geoData.Results[0].Latitude
	lon := geoData.Results[0].Longitude

	weatherURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true&timezone=auto", lat, lon)

	resp, err := client.Get(weatherURL)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer resp.Body.Close()

	var weatherData WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherData)

	if err != nil {
		fmt.Print(err)
	}
	apiTimeLayout := "2006-01-02T15:04"
	parsedTime, err := time.Parse(apiTimeLayout, weatherData.CurrentWeather.Time)
	if err != nil {
		fmt.Println(err)
		return
	}
	displayTime := parsedTime.Format("Jan 02, 2006 at 03:04 PM")
	condition := getWeatherDescription(weatherData.CurrentWeather.Weathercode)

	fmt.Printf(
		"--- Weather for %s ---\nTime: %s\nTemperature: %.1f°C\nWindspeed: %.1f km/h\nDescription: %s\n",
		geoData.Results[0].Name,
		displayTime,
		weatherData.CurrentWeather.Temperature,
		weatherData.CurrentWeather.Windspeed,
		condition,
)	
}

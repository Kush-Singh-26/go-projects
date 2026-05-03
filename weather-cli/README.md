# Weather CLI

A simple, fast command-line interface (CLI) application written in Go that fetches the current weather for any city. Powered by the free [Open-Meteo API](https://open-meteo.com/).

## Features

Prints a clean summary of the current weather, including:
- **Place** (City Name)
- **Local Time** (Adjusted to the city's timezone)
- **Temperature** (°C)
- **Windspeed** (km/h)
- **Description** (e.g., Clear sky, Rain, Snow)

## Usage 

You can run the program directly using Go:
```bash
go run main.go --city "New Delhi"
```

**Example Output:**
```text
--- Weather for New Delhi ---
Time: May 02, 2026 at 06:00 PM
Temperature: 35.2°C
Windspeed: 12.7 km/h
Description: Clear sky
```

### Build it yourself
To compile the app into a standalone executable that you can use anywhere:
```bash
go build -o bin/weather
./weather --city "Tokyo"
```

---

## Under the Hood

1. **CLI Flags:** Uses Go's standard `flag` package to parse the `--city` input.
2. **Geocoding:** Calls the [Geocoding API](https://open-meteo.com/en/docs/geocoding-api) to convert the city string into latitude and longitude coordinates safely.
3. **Weather Data:** Passes those coordinates into the Weather API to get the current forecast.
4. **Resiliency:** Implements a custom `http.Client` with a 5-second timeout to prevent the terminal from hanging forever if the API server goes down:

```go
client := &http.Client{
    Timeout: 5 * time.Second,
}
```
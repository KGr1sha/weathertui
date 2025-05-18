package main

import (
	"io"
	"os"
	"bytes"
	"net/http"
	"encoding/json"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type WeatherResponse struct {
	CurrentCondition []CurrentCondition `json:"current_condition"`
	Weather          []WeatherForecast  `json:"weather"`
}

type CurrentCondition struct {
	FeelsLikeC        string       `json:"FeelsLikeC"`
	FeelsLikeF        string       `json:"FeelsLikeF"`
	Cloudcover        string       `json:"cloudcover"`
	Humidity          string       `json:"humidity"`
	LocalObsDateTime  string       `json:"localObsDateTime"`
	ObservationTime   string       `json:"observation_time"`
	PrecipInches      string       `json:"precipInches"`
	PrecipMM          string       `json:"precipMM"`
	Pressure          string       `json:"pressure"`
	PressureInches    string       `json:"pressureInches"`
	TempC             string       `json:"temp_C"`
	TempF             string       `json:"temp_F"`
	UvIndex           string       `json:"uvIndex"`
	Visibility        string       `json:"visibility"`
	VisibilityMiles   string       `json:"visibilityMiles"`
	WeatherCode       string       `json:"weatherCode"`
	WeatherDesc       []WeatherDesc `json:"weatherDesc"`
	WeatherIconUrl    []WeatherIcon `json:"weatherIconUrl"`
	Winddir16Point    string       `json:"winddir16Point"`
	WinddirDegree     string       `json:"winddirDegree"`
	WindspeedKmph     string       `json:"windspeedKmph"`
	WindspeedMiles    string       `json:"windspeedMiles"`
}

type WeatherForecast struct {
	Date             string       `json:"date"`
	MaxtempC         string       `json:"maxtempC"`
	MaxtempF         string       `json:"maxtempF"`
	MintempC         string       `json:"mintempC"`
	MintempF         string       `json:"mintempF"`
	SunHour          string       `json:"sunHour"`
	TotalSnowCM      string       `json:"totalSnow_cm"`
	UvIndex          string       `json:"uvIndex"`
	Hourly           []HourlyForecast `json:"hourly"`
}

type HourlyForecast struct {
	Time              string       `json:"time"`
	TempC             string       `json:"tempC"`
	TempF             string       `json:"tempF"`
	FeelsLikeC        string       `json:"FeelsLikeC"`
	FeelsLikeF        string       `json:"FeelsLikeF"`
	DewPointC         string       `json:"DewPointC"`
	DewPointF         string       `json:"DewPointF"`
	HeatIndexC        string       `json:"HeatIndexC"`
	HeatIndexF        string       `json:"HeatIndexF"`
	WindChillC        string       `json:"WindChillC"`
	WindChillF        string       `json:"WindChillF"`
	WindGustKmph      string       `json:"WindGustKmph"`
	WindGustMiles     string       `json:"WindGustMiles"`
	Cloudcover        string       `json:"cloudcover"`
	Humidity          string       `json:"humidity"`
	PrecipInches      string       `json:"precipInches"`
	PrecipMM          string       `json:"precipMM"`
	Pressure          string       `json:"pressure"`
	PressureInches    string       `json:"pressureInches"`
	Visibility        string       `json:"visibility"`
	VisibilityMiles   string       `json:"visibilityMiles"`
	WeatherCode       string       `json:"weatherCode"`
	WeatherDesc       []WeatherDesc `json:"weatherDesc"`
	WeatherIconUrl    []WeatherIcon `json:"weatherIconUrl"`
	Winddir16Point    string       `json:"winddir16Point"`
	WinddirDegree     string       `json:"winddirDegree"`
	WindspeedKmph     string       `json:"windspeedKmph"`
	WindspeedMiles    string       `json:"windspeedMiles"`
	DiffRad           string       `json:"diffRad"`
	ShortRad          string       `json:"shortRad"`
	ChanceofFog       string       `json:"chanceoffog"`
	ChanceofFrost     string       `json:"chanceoffrost"`
	ChanceofHighTemp  string       `json:"chanceofhightemp"`
	ChanceofOvercast  string       `json:"chanceofovercast"`
	ChanceofRain      string       `json:"chanceofrain"`
	ChanceofRemDry    string       `json:"chanceofremdry"`
	ChanceofSnow      string       `json:"chanceofsnow"`
	ChanceofSunshine  string       `json:"chanceofsunshine"`
	ChanceofThunder   string       `json:"chanceofthunder"`
	ChanceofWindy     string       `json:"chanceofwindy"`
}

type WeatherDesc struct {
	Value string `json:"value"`
}

type WeatherIcon struct {
	Value string `json:"value"`
}


type errMsg struct {
	err error
}

func (e errMsg) Error() string { return e.err.Error() }

type weatherMsg struct {
	Weather WeatherResponse
}

func getWeather(city string) tea.Cmd {
	return func() tea.Msg {
		city = strings.ReplaceAll(city, " ", "+")
		resp, err := http.Get("https://wttr.in/" + city + "?format=j1")
		if err != nil {
			return errMsg{err}
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errMsg{err}
		}
		var weather WeatherResponse
		err = json.Unmarshal(body, &weather)
		if err != nil {
			return errMsg{err}
		}

		return weatherMsg{Weather: weather}
	}
}

type filterMsg struct {
	Filtered []string
}

func filterCmd(cities []string, filterStr string) tea.Cmd {
	return func() tea.Msg {
		return filterMsg{filter(cities, filterStr)}
	}
}

type City struct {
	Name string `json:"name"`
}

type citiesLoadMsg struct {
	Cities []string
}

func loadCities() tea.Msg {
	data, err := os.ReadFile("cities.txt")
	if err != nil {
		return errMsg{err}
	}

	lines := bytes.Split(data, []byte{'\n'})
	cities := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(line) > 0 {
			cities = append(cities, strings.TrimSpace(string(line)))
		}
	}
	return citiesLoadMsg{cities}
}


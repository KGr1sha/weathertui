package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cities []string
	filtered []string
	cursor int
	selected string
	weather WeatherResponse
	err error
	hasWeather bool
}

var initialModel = model{
	cities: []string{
		"Moscow",
		"Saint's Petersburg",
		"New York",
		"Cape Town",
		"Paris",
		"Batumi",
	},
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.cities)-1 {
				m.cursor++
			}
		case " ", "enter":
			m.selected = m.cities[m.cursor]
			return m, getWeather(m.selected)
		case "c", "s":
			m.selected = ""
			m.cursor = 0
			m.hasWeather = false
		}
	case weatherMsg:
		m.weather = msg.Weather
		m.hasWeather = true
	case errMsg:
		m.err = msg
	}
	return m, nil
}

func parseWeather(w WeatherResponse) string {
	s := ""
	current := w.CurrentCondition[0]
	s += fmt.Sprintf("Temp: %s°C\n", current.TempC)
	s += fmt.Sprintf("Wind: %s(km/h)\n", current.WindspeedKmph)
	return s
}

func (m model) View() string {
	s := "Weather TUI!\n\n"
	
	if m.err != nil {
		s += m.err.Error()
		return s
	}

	if m.selected != "" {
		s += "Selected " + m.selected + "\n"
		weather := "Waiting for weather...\n"
		if m.hasWeather {
			weather = parseWeather(m.weather)
		}
		s += weather
	} else {
		for i, c := range m.cities {
			cursor := "  "
			if m.cursor == i {
				cursor = "->"
			}
			s += fmt.Sprintf("%s%d. %s\n",cursor, i+1, c)
		}
	}
	help := "\npress q to quit\npress c or s to select a city\n"
	return s + help
}

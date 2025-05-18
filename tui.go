package main

import (
	"fmt" 

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))    // Purple
	filterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))    // Cyan
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))    // Yellow
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true) // Green
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))      // Grey
    temperatureStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("202")) // Warm orange-red
    windStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81")) // Cool blue
    uvIndexStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("129")) // Purple tone for UV
)

type model struct {
	cities []string
	filter string
	filtering bool
	filtered []string
	cursor int
	selected string
	weather WeatherResponse
	err error
	showingWeather bool
}

var initialModel = model{}

func (m model) Init() tea.Cmd {
	return loadCities
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showingWeather {
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "c", "s":
				m.selected = ""
				m.cursor = 0
				m.showingWeather = false
		}
		} else {
		if m.filtering {
			switch msg.String() {
			case "esc":
				m.filtering = false
				return m, nil
			case "enter":
				m.selected = m.filtered[m.cursor]
				return m, getWeather(m.selected)
			case "backspace":
				if len(m.filter) > 0 {
					m.filter = m.filter[:len(m.filter)-1]
				}
				return m, filterCmd(m.cities, m.filter)
			default:
			m.filter += msg.String()
			m.cursor = 0
			return m, filterCmd(m.cities, m.filter)
			}
		}
		switch msg.String() {
		case "/":
			m.filtering = true
			m.filter = ""
			m.filtered = make([]string, 0)
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
			if m.filtered != nil {
				m.selected = m.filtered[m.cursor]
			} else {
				m.selected = m.cities[m.cursor]
			}
			return m, getWeather(m.selected)
		}
	}
	case citiesLoadMsg:
		m.cities = msg.Cities
	case weatherMsg:
		m.weather = msg.Weather
		m.showingWeather = true
	case filterMsg:
		m.filtered = msg.Filtered
	case errMsg:
		m.err = msg
	}
	return m, nil
}

func (m model) View() string {
	banner := `
					  \   /
					   .-.
	Weather TUI 	― (   ) ―
					   '-'
					  /   \   
	`
	s := headerStyle.Render(banner) + "\n\n"

	if m.err != nil {
		s += fmt.Sprintf("Runtime error: %v\n", m.err)
		return s
	}

	if m.selected != "" {
		// Weather
		s += selectedStyle.Render(m.selected) + "\n"
		weather := "Asking Zeus for weather report...\n"
		if m.showingWeather {
			weather = parseWeather(m.weather)
		}
		s += weather
	} else {
		if m.cities == nil {
			s += "Loading..."
			return s
		}
		// Filter
		filterLine := filterStyle.Render("/ to filter: ") + filterStyle.Render(m.filter)
		if m.filtering {
			filterLine += filterStyle.Render("|")
		}
		s += filterLine + "\n\n"
		s += fmt.Sprintf("%d\n", len(m.filtered))

		// List
		cities := m.cities
		if len(m.filtered) != 0 || (m.filtering && len(m.filter) > 0) {
			cities = m.filtered
		}
		if len(cities) == 0 {
			s += "No such location...\n"
		}

		first := max(m.cursor - 5, 0)
		last := min(max(m.cursor + 5, first + 10), len(cities))
		for i := first; i < last; i++ {
			c := cities[i]
			line := fmt.Sprintf("%d. %s", i+1, c)
			if m.cursor == i {
				line = cursorStyle.Render("-> ") + line
			} else {
				line = "   " + line
			}
			s += line + "\n"
		}
	}

	// Help
	help := ""
	if m.selected == "" {
		if m.filtering {
			help = "\n" + helpStyle.Render("esc|enter • stop filtering")
		} else {
			help = "\n" + helpStyle.Render("q: quit • /: filter • ↑(k)/↓(j)|: navigate • enter|space: select")
		}
	} else {
		help = "\n" + helpStyle.Render("q: quit • s|c: select other city")
	}
	s += help

	return lipgloss.NewStyle().Padding(1, 2).Render(s)
}

func parseWeather(w WeatherResponse) string {
	current := w.CurrentCondition[0]
	wcode := current.WeatherCode
	key := WeatherCodes[wcode]

	icon := WeatherSymbols[key]
	desc := current.WeatherDesc[0].Value + " " + icon + "\n"
	asciiArt := ""
	lines := WeatherAsciiSymbols[key]
	for _, line := range lines {
		asciiArt += line + "\n"
	}

	temp := "Temp "
	temp += temperatureStyle.Render(current.TempC) + "°C"
	temp += " (Feels like " + temperatureStyle.Render(current.FeelsLikeC) + "°C)"
	temp += "\n"
	wind := "Wind " + windStyle.Render(current.WindspeedKmph) + "(km/h)\n"
	uv := "UV index " + uvIndexStyle.Render(current.UvIndex) + "\n"

	s := ""
	s += desc
	s += asciiArt
	s += temp + "\n"
	s += wind
	s += uv
	return s
}

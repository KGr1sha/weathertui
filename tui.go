package main

import (
	"fmt" 
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))    // Purple
	filterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))    // Cyan
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))    // Yellow
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true) // Green
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))      // Grey
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
	hasWeather bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.filtering {
			switch msg.String() {
			case "esc", "enter":
				m.filtering = false
				return m, nil
			case "backspace":
				if len(m.filter) > 0 {
					m.filter = m.filter[:len(m.filter)-1]
				}
				return m, nil
			default:
			m.filter += msg.String()
			m.filtered = filter(m.cities, m.filter) // TODO: move it to command, so it runs as goroutine
			m.cursor = 0
			return m, nil
			}
		}
		switch msg.String() {
		case "/":
			m.filtering = true
			m.filter = ""
			clear(m.filtered)
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

func (m model) View() string {
	banner := `
					  \   /
					   .-.
	Weather TUI 	― (   ) ―
					   '-'
					  /   \   
	`
	s := headerStyle.Render(banner) + "\n\n"

	// Filter/selected block
	if m.selected != "" {
		s += selectedStyle.Render(m.selected) + "\n"
		weather := "Asking Zeus for weather report...\n"
		if m.hasWeather {
			weather = parseWeather(m.weather)
		}
		s += weather
	} else {
		filterLine := filterStyle.Render("/ to filter: ") + filterStyle.Render(m.filter)
		if m.filtering {
			filterLine += filterStyle.Render("|")
		}
		s += filterLine + "\n\n"

		// List
		cities := m.cities
		if len(m.filtered) != 0 || (m.filtering && len(m.filter) > 0) {
			cities = m.filtered
		}
		if len(cities) == 0 {
			s += "No such location...\n"
		}

		for i, c := range cities {
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
	s := ""
	current := w.CurrentCondition[0]
	s += fmt.Sprintf("Temp: %s°C\n", current.TempC)
	s += fmt.Sprintf("Wind: %s(km/h)\n", current.WindspeedKmph)
	return s
}

func filter(cities []string, filter string) []string {
	type item struct {
		word string
		distance int // levenshtein
	}
	filter = strings.ToLower(filter)
	withCommonRunes := withCommonRunes(cities, filter)
	withDistances := make([]item, len(withCommonRunes))
	for i, w := range withCommonRunes {
		withDistances[i] = item{word: w, distance: levenshteinDistance(w, filter)}
	}

	slices.SortFunc(withDistances, func(a, b item) int {
		if a.distance < b.distance {
			return -1
		} else if a.distance > b.distance {
			return 1
		}
		return 0
	})

	ans := make([]string, len(withDistances))
	for i, item := range withDistances {
		ans[i] = item.word
	}
	return ans
}

func levenshteinDistance(word1, word2 string) int {
	cache := make(map[[2]string]int)
	var lev func(string, string) int
	lev = func(a, b string) int {
		if val, ok := cache[[2]string{a,b}]; ok {
			return val
		}
		if len(b) == 0 {
			ans := len(a)
			cache[[2]string{a,b}] = ans
			return ans
		}
		if len(a) == 0 {
			ans := len(b)
			cache[[2]string{a,b}] = ans
			return ans
		}
		if a[0] == b[0] {
			ans := lev(a[1:], b[1:])
			cache[[2]string{a,b}] = ans
			return ans
		}
		ans := 1 + min(
			lev(a[1:], b),
			lev(a, b[1:]),
			lev(a[1:], b[1:]),
		)
		cache[[2]string{a,b}] = ans
		return ans
	}
	return lev(word1, word2)
}

func withCommonRunes(words []string, word string) []string {
	var ans []string
	word = strings.ToLower(word)
	searchingFor := countRunes(word)
	for _, w := range words {
		low := strings.ToLower(w)
		runes := countRunes(low)
		for r := range searchingFor {
			if _, ok := runes[r]; ok {
				ans = append(ans, w)
				break;
			}
		}
	}
	return ans
}

func countRunes(word string) map[rune]int {
	ans := make(map[rune]int)
	for _, r := range word {
		ans[r] += 1
	}
	return ans
}

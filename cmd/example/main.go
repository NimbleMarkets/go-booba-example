// A simple BubbleTea program that runs both natively and in the browser via WASM.
package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	booba "github.com/NimbleMarkets/go-booba"
)

func main() {
	if err := booba.Run(initialModel()); err != nil {
		log.Fatal(err)
	}
}

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	itemStyle     = lipgloss.NewStyle().PaddingLeft(2)
	selectedStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	tickStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
)

type tickMsg struct{}

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	ticks    int
	quitting bool
	width    int
	height   int
}

func initialModel() model {
	return model{
		choices: []string{
			"Plant carrots 🥕",
			"Go to the market 🛒",
			"Read a book 📚",
			"See friends 👋",
		},
		selected: make(map[int]struct{}),
		ticks:    30,
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", "space":
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}

	case tickMsg:
		if m.ticks <= 0 {
			m.quitting = true
			return m, tea.Quit
		}
		m.ticks--
		return m, tick()
	}

	return m, nil
}

func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("\n  Bye! 👋\n\n")
	}

	s := "\n"
	s += titleStyle.Render("What to do today?") + "\n\n"

	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = "▸ "
		}

		checked := "○"
		style := itemStyle
		if _, ok := m.selected[i]; ok {
			checked = "●"
			style = selectedStyle
		}

		s += style.Render(fmt.Sprintf("%s%s %s", cursor, checked, choice)) + "\n"
	}

	s += "\n"
	s += subtleStyle.Render("  j/k, ↑/↓: move • enter/space: select • q: quit") + "\n"
	s += subtleStyle.Render("  Time left: ") + tickStyle.Render(strconv.Itoa(m.ticks)+"s") + "\n"

	return tea.NewView(s)
}

package model

import (
	"fmt"
	"strconv"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	TitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	ItemStyle     = lipgloss.NewStyle().PaddingLeft(2)
	SelectedStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	SubtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	TickStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
)

type tickMsg struct{}

type Model struct {
	Choices  []string
	Cursor   int
	Selected map[int]struct{}
	Ticks    int
	Quitting bool
	Width    int
	Height   int
}

func InitialModel() Model {
	return Model{
		Choices: []string{
			"Plant carrots 🥕",
			"Go to the market 🛒",
			"Read a book 📚",
			"See friends 👋",
		},
		Selected: make(map[int]struct{}),
		Ticks:    30,
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		case "enter", "space":
			if _, ok := m.Selected[m.Cursor]; ok {
				delete(m.Selected, m.Cursor)
			} else {
				m.Selected[m.Cursor] = struct{}{}
			}
		}

	case tickMsg:
		if m.Ticks <= 0 {
			m.Quitting = true
			return m, tea.Quit
		}
		m.Ticks--
		return m, tick()
	}

	return m, nil
}

func (m Model) View() tea.View {
	if m.Quitting {
		return tea.NewView("\n  Bye! 👋\n\n")
	}

	s := "\n"
	s += TitleStyle.Render("What to do today?") + "\n\n"

	for i, choice := range m.Choices {
		cursor := "  "
		if m.Cursor == i {
			cursor = "▸ "
		}

		checked := "○"
		style := ItemStyle
		if _, ok := m.Selected[i]; ok {
			checked = "●"
			style = SelectedStyle
		}

		s += style.Render(fmt.Sprintf("%s%s %s", cursor, checked, choice)) + "\n"
	}

	s += "\n"
	s += SubtleStyle.Render("  j/k, ↑/↓: move • enter/space: select • q: quit") + "\n\n"
	s += SubtleStyle.Render("  Time left: ") + TickStyle.Render(strconv.Itoa(m.Ticks)+"s") + "\n"
	s += SubtleStyle.Render(fmt.Sprintf("  Size: %d x %d", m.Width, m.Height)) + "\n"

	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

package mcd

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"os"
	"strings"
)

func initModel(initialSearchText string) model {
	searchInput := textinput.New()

	searchInput.Placeholder = "Search"
	searchInput.Focus()
	searchInput.CharLimit = 156
	searchInput.Width = 32
	searchInput.Prompt = ": "

	if initialSearchText != "" {
		searchInput.SetValue(initialSearchText)
	}

	candidates := getCandidates()

	return model{
		searchInput: searchInput,
		searchText:  "",
		candidates:  candidates,

		state: &State{
			selected:           false,
			filteredCandidates: getFilteredCandidates(&candidates, initialSearchText),
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

var nameStyle = lipgloss.NewStyle().Bold(true)
var pathStyle = lipgloss.NewStyle().Faint(true)

func (m model) candidatesView() string {
	buf := bytes.NewBufferString("")

	l := len(m.candidates)
	current := m.state.cursor % l

	for i, fc := range *m.state.filteredCandidates {
		name := nameStyle.Render(fc.candidate.name)

		bookmarkPath := pathStyle.Render(
			strings.Replace(fc.candidate.path, os.Getenv("HOME"), "~", 1),
		)

		row := fmt.Sprintf("%s %s\n", name, bookmarkPath)

		if i == current {
			buf.WriteString("> " + row)
		} else {
			buf.WriteString("  " + row)
		}

	}

	return buf.String()
}

func (m model) View() string {
	if m.state.exited {
		return ""
	} else {
		return fmt.Sprintf(
			"%s\n%s",
			m.searchInput.View(),
			m.candidatesView(),
		)
	}
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			m.state.exited = true

			if message.Type == tea.KeyEnter {
				m.state.selected = true
			}

			return m, tea.Quit

		case tea.KeyDown:
			m.state.cursor++

			if m.state.cursor >= len(m.candidates) {
				m.state.cursor = 0
			}

			return m, nil

		case tea.KeyUp:
			if m.state.cursor > 0 {
				m.state.cursor--
			} else {
				m.state.cursor = len(*m.state.filteredCandidates) - 1
			}

			return m, nil
		}
	}

	m.searchInput, cmd = m.searchInput.Update(message)
	searchText := m.searchInput.Value()

	if m.searchText != searchText {
		m.state.filteredCandidates = getFilteredCandidates(&m.candidates, searchText)
		m.searchText = searchText
		m.state.cursor = 0
	}

	return m, cmd
}

func Run(initialSearchText string) {
	m := initModel(initialSearchText)

	if len(*m.state.filteredCandidates) == 1 {
		os.Stdout.WriteString((*m.state.filteredCandidates)[0].candidate.path)
		return
	}

	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	_, err := p.Run()

	if err != nil {
		log.Fatal(err)
	}

	if m.state.selected {
		if len(*m.state.filteredCandidates) > 0 {
			targetPath := (*m.state.filteredCandidates)[m.state.cursor].candidate.path
			os.Stdout.WriteString(targetPath)
		}
	}
}

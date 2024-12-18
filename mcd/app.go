package mcd

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
			maxHeight:          4,
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

	filteredCandidates := *m.state.filteredCandidates
	filteredCandidateCount := len(filteredCandidates)

	current := m.state.cursor

	availableLines := m.state.maxHeight - 3
	availableCandidates := len(m.candidates)

	candidateRenderCount := int(math.Min(float64(availableLines), float64(availableCandidates)))

	renderTillIndex := int(math.Min(float64(filteredCandidateCount), float64(current+candidateRenderCount/2)))
	renderFromIndex := int(math.Max(0, float64(renderTillIndex-candidateRenderCount)))
	renderTillIndex = renderFromIndex + candidateRenderCount

	for i := range candidateRenderCount {
		if renderFromIndex+i >= filteredCandidateCount {
			buf.WriteString("\n")
			continue
		}

		fc := filteredCandidates[renderFromIndex+i]
		name := nameStyle.Render(fc.candidate.name)

		directoryPath := pathStyle.Render(
			strings.Replace(fc.candidate.path, os.Getenv("HOME"), "~", 1),
		)

		row := fmt.Sprintf("%s %s", name, directoryPath)

		if renderFromIndex+i == current {
			buf.WriteString("> " + row + "\n")
		} else {
			buf.WriteString("  " + row + "\n")
		}

	}

	hiddenResultCountBelow := filteredCandidateCount - renderTillIndex

	if hiddenResultCountBelow <= 0 {
		buf.WriteString("\n")
	} else {
		buf.WriteString(pathStyle.Render(fmt.Sprintf("  -- %d more below (%d total results)", hiddenResultCountBelow, filteredCandidateCount)) + "\n")
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

		case tea.KeyDown, tea.KeyCtrlP:
			m.state.cursor++

			if m.state.cursor >= len(*m.state.filteredCandidates) {
				m.state.cursor = 0
			}

			return m, nil

		case tea.KeyUp, tea.KeyCtrlN:
			if m.state.cursor > 0 {
				m.state.cursor--
			} else {
				m.state.cursor = len(*m.state.filteredCandidates) - 1
			}

			return m, nil
		}
	case tea.WindowSizeMsg:
		envHeight := os.Getenv("MONOCD_MAX_HEIGHT")

		if envHeight == "" {
			m.state.maxHeight = message.Height
		} else {
			envHeightInt, err := strconv.Atoi(envHeight)

			if err != nil {
				m.state.maxHeight = message.Height
			} else {
				m.state.maxHeight = envHeightInt
			}
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

package mcd

import "github.com/charmbracelet/bubbles/textinput"

type Candidate struct {
	name string
	path string
}

type FilteredCandidate struct {
	candidate Candidate
	rank      int
}

type State struct {
	selected  bool
	cursor    int
	exited    bool
	maxHeight int

	filteredCandidates *[]FilteredCandidate
}

type model struct {
	searchInput textinput.Model
	searchText  string
	candidates  []Candidate // directories of interest

	state *State
}

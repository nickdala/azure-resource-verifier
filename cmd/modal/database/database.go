package database

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	REDIS = iota
	POSTGRESQL
	POSTGRESQL_HA
)

type databaseChoice struct {
	id          int
	description string
}

type model struct {
	choices  []databaseChoice // database choices
	cursor   int              // which database item our cursor is pointing at
	selected map[int]struct{} // which database items are selected
}

func ShowDatabaseModalAndGetChoices() ([]int, error) {
	databaseChoices := []int{}

	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		return databaseChoices, fmt.Errorf("error running the modal: %v", err)
	}

	// Assert the final tea.Model to our local model and print the choice.
	if m, ok := m.(model); !ok {
		return databaseChoices, fmt.Errorf("error asserting the database model")
	} else {
		// return selected choices
		for k := range m.selected {
			databaseChoices = append(databaseChoices, m.choices[k].id)
		}
		return databaseChoices, nil
	}
}

func initialModel() model {
	return model{
		choices: []databaseChoice{
			{REDIS, "Azure Cache for Redis"},
			{POSTGRESQL, "Azure PostgreSQL Flexible Server"},
			{POSTGRESQL_HA, "Azure PostgreSQL Flexible Server with HA"},
		},

		// A map which indicates which choices are selected. The keys
		// refer to the indexes of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "c":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}

				// If the user selects the PostgreSQL HA option, also deselect the PostgreSQL option
				if m.cursor == POSTGRESQL_HA {
					delete(m.selected, POSTGRESQL)
				} else if m.cursor == POSTGRESQL { // If the user selects the PostgreSQL option, also deselect the PostgreSQL HA option
					delete(m.selected, POSTGRESQL_HA)
				}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "What databases are you deploying?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.description)
	}

	// The footer
	s += "\nPress c to confirm.\n"

	// Send the UI for rendering
	return s
}

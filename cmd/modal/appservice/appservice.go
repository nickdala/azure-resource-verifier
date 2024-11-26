package appservice

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	APP_SERVICE_NONE = iota
	APP_SERVICE_LINUX_CODE
	APP_SERVICE_LINUX_CONTAINER
	APP_SERVICE_WINDOWS_CODE
	APP_SERVICE_WINDOWS_CONTAINER
)

type appServiceChoice struct {
	id          int
	description string
}

type model struct {
	choices  []appServiceChoice // app service choices
	cursor   int                // which app service item our cursor is pointing at
	selected *appServiceChoice  // which app service item is selected
}

func ShowAppServiceModalAndGetChoices() (int, error) {

	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		return -1, fmt.Errorf("error running the modal: %v", err)
	}

	// Assert the final tea.Model to our local model and print the choice.
	appServiceModel, ok := m.(model)
	if !ok {
		return -1, fmt.Errorf("error asserting the app service model")
	}

	// return selected choice
	if appServiceModel.selected == nil {
		return APP_SERVICE_NONE, nil
	}

	return appServiceModel.selected.id, nil

}

func initialModel() model {
	return model{
		choices: []appServiceChoice{
			{APP_SERVICE_LINUX_CODE, "Azure App Service - Linux Code"},
			{APP_SERVICE_LINUX_CONTAINER, "Azure App Service - Linux Container"},
			{APP_SERVICE_WINDOWS_CODE, "Azure App Service - Windows Code"},
			{APP_SERVICE_WINDOWS_CONTAINER, "Azure App Service - Windows Container"},
		},

		selected: nil,
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
			selected := m.choices[m.cursor]

			// If the item is already selected, unselect it
			if m.selected != nil && m.selected.id == selected.id {
				m.selected = nil
			} else { // Otherwise select it
				m.selected = &selected
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "What type of App Service are you deploying?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if m.selected != nil && m.selected.id == choice.id {
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

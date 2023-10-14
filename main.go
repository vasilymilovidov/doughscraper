package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
)

type Path struct {
	question string
	answer   string
	input    Input
}

type model struct {
	styles         *Styles
	paths          []Path
	MainMenu       int
	RenameMenu     bool
	CreateJsonMenu bool
	Quitting       bool
	width          int
	height         int
	index          int
	done           string
}

func main() {
	questions := []Path{
		newLocalFolder("Local folder path"),
		newRemoteFolder("Web folder path"),
	}
	styles := DefaultStyles()
	initialModel := model{styles, questions, 0, false, false, false, 0, 0, 0, ""}
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

func newPath(q string) Path {
	return Path{question: q}
}

func newLocalFolder(q string) Path {
	path := newPath(q)
	model := NewLocalFolderPathField()
	path.input = model
	return path
}

func newRemoteFolder(q string) Path {
	path := newPath(q)
	model := NewWebFolderPathField()
	path.input = model
	return path
}

func (m model) Init() tea.Cmd {
	return nil
}

// Next input
func (m *model) Next() {
	if m.index < len(m.paths)-1 {
		m.index++
	} else {
		m.index = 0
	}
}

// Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	if !m.RenameMenu && !m.CreateJsonMenu {
		return updateChoices(msg, m)
	} else if m.RenameMenu {
		return updateRename(msg, m)
	} else if m.CreateJsonMenu {
		return updateCreateJson(msg, m)
	}

	return m, nil
}

// View logic for the main menu
func (m model) View() string {
	var s string
	if m.Quitting {
		return "Quitting..."
	}
	if !m.RenameMenu && !m.CreateJsonMenu {
		s = choicesView(m)
	} else if m.RenameMenu {
		s = renameView(m)
	} else if m.CreateJsonMenu {
		s = createJsonView(m)
	}
	return indent.String("\n"+s+"\n\n", 2)
}

// View for the main menu
func choicesView(m model) string {
	c := m.MainMenu
	tpl := "Doughscraper\n\n"
	tpl += "%s\n\n"
	tpl += subtle(m.done)
	tpl += "\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s",
		checkbox("Rename pitched files", c == 0),
		checkbox("Generate JSON", c == 1),
	)

	return fmt.Sprintf(tpl, choices)
}

// Update for the main menu
func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.MainMenu++
			if m.MainMenu > 1 {
				m.MainMenu = 0
			}
		case "k", "up":
			m.MainMenu--
			if m.MainMenu < 0 {
				m.MainMenu = 1
			}
		case "enter":
			switch m.MainMenu {
			case 0:
				m.RenameMenu = true
			case 1:
				m.CreateJsonMenu = true
			}
		}
	}
	return m, nil
}

// View for the "Create JSON" menu
func createJsonView(m model) string {
	current := m.paths[m.index]
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			current.question,
			m.styles.InputField.Render(current.input.View()),
			lipgloss.JoinHorizontal(lipgloss.Bottom, subtle("Provide paths to the roots of your local and remote sample folders")),
		),
	)
}

// View for the "Rename files" menu
func renameView(m model) string {
	current := m.paths[0]
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			current.question,
			m.styles.InputField.Render(current.input.View()),
			lipgloss.JoinHorizontal(lipgloss.Bottom, subtle("Provide the path to the folder with files to rename.")),
		),
	)
}

// Update for the "Create JSON" menu
func updateCreateJson(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	current := &m.paths[m.index]
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			current.answer = current.input.Value()
			if m.index == len(m.paths)-1 {
				err := GenerateJson(m.paths[0].answer, m.paths[1].answer)
				if err != nil {
					fmt.Printf("Error: %s", err.Error())
					m.done = fmt.Sprintf("Error: %s", err.Error())
				}
				m.done = fmt.Sprintf("Generated JSON in %s", m.paths[0].answer)
				m.CreateJsonMenu = false
			}
			m.Next()
			return m, current.input.Blur
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

// Update for the "Rename files" menu
func updateRename(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	current := &m.paths[0]
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			current.answer = current.input.Value()
			err := RenameFiles(current.answer)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
				m.done = fmt.Sprintf("Error: %s", err.Error())
			} else {
				m.done = fmt.Sprintf("Renamed files in %s", current.answer)
			}
			m.RenameMenu = false
			return m, current.input.Blur
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

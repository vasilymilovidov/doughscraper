package main

import (
	"errors"
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
)

type pitchDetectorFinishedMsg struct{err error}

func PitchDetectAndRenameFiles(path string) tea.Cmd {
	// pitchdetector, err := exec.LookPath("pitchdetector")
	c := exec.Command("pitchdetector", path)
	if errors.Is(c.Err, exec.ErrDot) {
		c.Err = nil
	}
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return pitchDetectorFinishedMsg{err}
})
}

type PathInput struct {
	question string
	answer   string
	input    Input
}

type model struct {
	styles               *Styles
	paths                []PathInput
	MainMenu             int
	PitchDetectAndRename bool
	RenameMenu           bool
	CreateJsonMenu       bool
	Quitting             bool
	width                int
	height               int
	index                int
	doneMessage          string
}

func main() {
	paths := []PathInput{
		newLocalFolder("Local folder path"),
		newLocalRoot("Local root path"),
		newRemoteFolder("Remote root folder path"),
		newLocalFolder("Local folder path"),
	}
	styles := DefaultStyles()
	initialModel := model{
		styles:               styles,
		paths:                paths,
		MainMenu:             0,
		PitchDetectAndRename: false,
		RenameMenu:           false,
		CreateJsonMenu:       false,
		Quitting:             false,
		width:                0,
		height:               0,
		index:                0,
		doneMessage:          "",
	}
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

func newPath(q string, placeholder string) PathInput {
	path := PathInput{question: q, input: NewPath(placeholder)}
	return path
}

func newLocalRoot(q string) PathInput {
	placeholder := "example: /Users/username/samples"
	return newPath(q, placeholder)
}

func newLocalFolder(q string) PathInput {
	placeholder := "example: /Users/username/samples/piano"
	return newPath(q, placeholder)
}

func newRemoteFolder(q string) PathInput {
	placeholder := "example: https://raw.githubusercontent.com/username/samples/main/"
	return newPath(q, placeholder)
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
	if !m.RenameMenu && !m.CreateJsonMenu && !m.PitchDetectAndRename {
		return updateChoices(msg, m)
	} else if m.RenameMenu {
		return updateRename(msg, m)
	} else if m.CreateJsonMenu {
		return updateCreateJson(msg, m)
	} else if m.PitchDetectAndRename {
		return updatePitchDetectAndRename(msg, m)
	}

	return m, nil
}

// View logic for the main menu
func (m model) View() string {
	var s string
	if m.Quitting {
		return "Quitting..."
	}
	if !m.RenameMenu && !m.CreateJsonMenu && !m.PitchDetectAndRename {
		s = choicesView(m)
	} else if m.RenameMenu {
		s = renameView(m)
	} else if m.CreateJsonMenu {
		s = createJsonView(m)
	} else if m.PitchDetectAndRename {
		s = PitchDetectAndRename(m)
	}

	return indent.String("\n"+s+"\n\n", 2)
}

// View for the main menu
func choicesView(m model) string {
	c := m.MainMenu
	tpl := "Doughscraper\n\n"
	tpl += "%s\n\n"
	tpl += subtle(m.doneMessage)
	tpl += "\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n",
		checkbox("Rename pitched files", c == 0),
		checkbox("Generate JSON", c == 1),
		checkbox("Detect pitch and rename", c == 2),
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
			if m.MainMenu > 2 {
				m.MainMenu = 0
			}
		case "k", "up":
			m.MainMenu--
			if m.MainMenu < 0 {
				m.MainMenu = 2
			}
		case "enter":
			switch m.MainMenu {
			case 0:
				m.RenameMenu = true
			case 1:
				m.CreateJsonMenu = true
				m.index = 1
			case 2:
				m.PitchDetectAndRename = true
			}
		}
	}
	return m, nil
}

// View for the "Rename pitched files" menu
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
			lipgloss.JoinHorizontal(lipgloss.Bottom, subtle("Path to the folder with files to rename")),
		),
	)
}

// Update for the "Rename pitched files" menu
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
				m.doneMessage = fmt.Sprintf("Error: %s", err.Error())
			} else {
				m.doneMessage = fmt.Sprintf("Renamed files in %s", current.answer)
			}
			m.RenameMenu = false
			return m, current.input.Blur()
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

// View for the "Detect pitch and rename" menu
func PitchDetectAndRename(m model) string {
	current := m.paths[3]
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			current.question,
			m.styles.InputField.Render(current.input.View()),
			lipgloss.JoinHorizontal(lipgloss.Bottom, subtle("Path to the folder with files to pitch detect and rename")),
		),
	)
}

// Update for the "Detect pitch and rename" menu
func updatePitchDetectAndRename(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	current := &m.paths[3]
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
			current.input.SetValue("Working...")
			return m, PitchDetectAndRenameFiles(current.answer)
		}
	case pitchDetectorFinishedMsg:
		if msg.err != nil {
			fmt.Printf("Error: %v", msg.err)
			m.doneMessage = fmt.Sprintf("Error: %v", msg.err)
			m.PitchDetectAndRename = false
			} else {
				m.doneMessage = fmt.Sprintf("Detected and renamed files in %s", current.answer)
				m.PitchDetectAndRename = false
			}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

// View for the "Generate JSON" menu
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

// Update for the "Generate JSON" menu
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
			if m.index == len(m.paths)-2 {
				err := GenerateJson(m.paths[1].answer, m.paths[2].answer)
				if err != nil {
					fmt.Printf("Error: %s", err.Error())
					m.doneMessage = fmt.Sprintf("Error: %s", err.Error())
				}
				m.doneMessage = fmt.Sprintf("Generated JSON in %s", m.paths[1].answer)
				m.CreateJsonMenu = false
			}
				m.Next()
			
			return m, current.input.Blur()
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

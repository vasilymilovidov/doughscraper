package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Blink() tea.Msg
	Blur() tea.Cmd
	Focus() tea.Cmd
	SetValue(string)
	Value() string
	Update(tea.Msg) (Input, tea.Cmd)
	View() string
}

type Path struct {
	textinput textinput.Model
	placeholder string
}

func NewPath(placeholder string) *Path {
	a := Path{}
	model := textinput.New()
	model.Placeholder = placeholder
	model.Focus()

	a.textinput = model
	return &a
}

func (a *Path) Blink() tea.Msg {
	return textinput.Blink()
}

func (a *Path) Focus() tea.Cmd {
	return a.textinput.Focus()
}

func (a *Path) SetValue(s string) {
	a.textinput.SetValue(s)
}

func (a *Path) Value() string {
	return a.textinput.Value()
}

func (a *Path) Blur() tea.Cmd {
    a.textinput.Blur()
	return nil
}

func (a *Path) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	a.textinput, cmd = a.textinput.Update(msg)
	return a, cmd
}

func (a *Path) View() string {
	return a.textinput.View()
}

type LocalFolderPath struct {
	*Path
}

func NewLocalFolderPathField() *LocalFolderPath {
	placeholder := "example: /Users/username/samples/piano"
	return &LocalFolderPath{NewPath(placeholder)}
}

type WebFolderPath struct {
	*Path
}

func NewWebFolderPathField() *WebFolderPath {
	placeholder := "example: https://raw.githubusercontent.com/username/samples/main/"
	return &WebFolderPath{NewPath(placeholder)}
}

type LocalRootPath struct {
	*Path
}

func NewLocalRootPathField() *LocalRootPath {
	placeholder := "example: /Users/username/samples"
	return &LocalRootPath{NewPath(placeholder)}
}

package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Blink() tea.Msg
	Blur() tea.Msg
	Focus() tea.Cmd
	SetValue(string)
	Value() string
	Update(tea.Msg) (Input, tea.Cmd)
	View() string
}

type LocalFolderPath struct {
	textinput textinput.Model
}

func NewLocalFolderPathField() *LocalFolderPath {
	a := LocalFolderPath{}

	model := textinput.New()
	model.Placeholder = "example: /Users/username/samples"
	model.Focus()

	a.textinput = model
	return &a
}

func (a *LocalFolderPath) Blink() tea.Msg {
	return textinput.Blink()
}

func (a *LocalFolderPath) Init() tea.Cmd {
	return nil
}

func (a *LocalFolderPath) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	a.textinput, cmd = a.textinput.Update(msg)
	return a, cmd
}

func (a *LocalFolderPath) View() string {
	return a.textinput.View()
}

func (a *LocalFolderPath) Focus() tea.Cmd {
	return a.textinput.Focus()
}

func (a *LocalFolderPath) SetValue(s string) {
	a.textinput.SetValue(s)
}

func (a *LocalFolderPath) Blur() tea.Msg {
	return a.textinput.Blur
}

func (a *LocalFolderPath) Value() string {
	return a.textinput.Value()
}

type WebFolderPath struct {
	textinput textinput.Model
}

func NewWebFolderPathField() *WebFolderPath {
	a := WebFolderPath{}

	model := textinput.New()
	model.Placeholder = "example: https://raw.githubusercontent.com/username/samples/main/"
	model.Focus()

	a.textinput = model
	return &a
}

func (a *WebFolderPath) Blink() tea.Msg {
	return textinput.Blink()
}

func (a *WebFolderPath) Init() tea.Cmd {
	return nil
}

func (a *WebFolderPath) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	a.textinput, cmd = a.textinput.Update(msg)
	return a, cmd
}

func (a *WebFolderPath) View() string {
	return a.textinput.View()
}

func (a *WebFolderPath) Focus() tea.Cmd {
	return a.textinput.Focus()
}

func (a *WebFolderPath) SetValue(s string) {
	a.textinput.SetValue(s)
}

func (a *WebFolderPath) Blur() tea.Msg {
	return a.textinput.Blur
}

func (a *WebFolderPath) Value() string {
	return a.textinput.Value()
}

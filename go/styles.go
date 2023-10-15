package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

// DefaultStyles Default styles for the input fields
func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = "#cae797"
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

// General stuff for styling the view
var (
	term   = termenv.EnvColorProfile()
	subtle = makeFgStyle("241")
	dot    = colorFg(" • ", "136")
)

func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg(" > "+label, "#cae797")
	}
	return fmt.Sprintf("   %s", label)
}

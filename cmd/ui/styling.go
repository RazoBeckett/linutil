package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	DocStyle     = lipgloss.NewStyle().Margin(1, 2)
	HelpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 2)
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A3BE8C")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#A3BE8C")).
			Padding(1).
			Align(lipgloss.Center)
	ErrStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BF616A")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#BF616A")).
			Padding(1).
			Align(lipgloss.Center)
	StdErrStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#BF616A"))
	StdErrStyleTitle = lipgloss.NewStyle().Bold(true).Inherit(StdErrStyle)
)

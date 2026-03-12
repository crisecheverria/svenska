package ui

import "github.com/charmbracelet/lipgloss"

var (
	colorBlue   = lipgloss.Color("#006AA7")
	colorYellow = lipgloss.Color("#FECC00")
	colorGreen  = lipgloss.Color("#98C379")
	colorRed    = lipgloss.Color("#E06C75")
	colorPurple = lipgloss.Color("#C678DD")
	colorCyan   = lipgloss.Color("#4FC1FF")
	colorDim    = lipgloss.Color("#666666")

	titleStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	swedishStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	englishStyle = lipgloss.NewStyle().
			Foreground(colorGreen)

	correctStyle = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	wrongStyle = lipgloss.NewStyle().
			Foreground(colorRed).
			Bold(true)

	progressStyle = lipgloss.NewStyle().
			Foreground(colorBlue)

	promptStyle = lipgloss.NewStyle().
			Foreground(colorPurple).
			Bold(true)

	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true).
			PaddingLeft(2)

	timerStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	timerLowStyle = lipgloss.NewStyle().
			Foreground(colorRed).
			Bold(true)

	hardcoreStyle = lipgloss.NewStyle().
			Foreground(colorPurple).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	flagBlue = lipgloss.NewStyle().
			Foreground(colorBlue)

	flagYellow = lipgloss.NewStyle().
			Foreground(colorYellow)
)

var flag = []string{
	flagBlue.Render("████████") + flagYellow.Render("██") + flagBlue.Render("████████████"),
	flagBlue.Render("████████") + flagYellow.Render("██") + flagBlue.Render("████████████"),
	flagYellow.Render("██████████████████████"),
	flagBlue.Render("████████") + flagYellow.Render("██") + flagBlue.Render("████████████"),
	flagBlue.Render("████████") + flagYellow.Render("██") + flagBlue.Render("████████████"),
}

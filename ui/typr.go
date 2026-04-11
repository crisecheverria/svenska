package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/crisecheverria/svenska/data"
)

type typrState struct {
	story     data.Story
	runes     []rune
	input     []rune
	pos       int
	started   bool
	startTime time.Time
	finished  bool
	endTime   time.Time
}

func (t *typrState) accuracy() float64 {
	if len(t.input) == 0 {
		return 0
	}
	correct := 0
	for i, r := range t.input {
		if i < len(t.runes) && r == t.runes[i] {
			correct++
		}
	}
	return float64(correct) / float64(len(t.input)) * 100
}

func (t *typrState) elapsed() time.Duration {
	if !t.started {
		return 0
	}
	if t.finished {
		return t.endTime.Sub(t.startTime)
	}
	return time.Since(t.startTime)
}

func (t *typrState) wpm() float64 {
	secs := t.elapsed().Seconds()
	if secs == 0 {
		return 0
	}
	words := float64(t.pos) / 5.0
	return words / (secs / 60.0)
}

// --- Story Selection ---

func (m Model) updateSelectStory(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "esc", "b":
			m.screen = screenMenu
			m.cursor = 5
			return m, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(data.Stories)-1 {
				m.cursor++
			}
		case "enter":
			story := data.Stories[m.cursor]
			m.typr = typrState{
				story: story,
				runes: []rune(story.Sv),
				input: make([]rune, 0, len([]rune(story.Sv))),
			}
			m.screen = screenTypr
		}
	}
	return m, nil
}

func (m Model) viewSelectStory() string {
	var b strings.Builder
	b.WriteString("\n" + titleStyle.Render("  Story Typer") + "  " + dimStyle.Render("Type over Swedish stories") + "\n\n")

	for i, s := range data.Stories {
		cursor := "  "
		style := menuItemStyle
		if i == m.cursor {
			cursor = promptStyle.Render("▸ ")
			style = selectedStyle
		}
		words := len(strings.Fields(s.Sv))
		b.WriteString(cursor + style.Render(s.Title) + "  " +
			dimStyle.Render(fmt.Sprintf("%d words", words)) + "\n")
	}

	b.WriteString("\n" + dimStyle.Render("  ↑↓/jk navigate  •  enter select  •  b back") + "\n")
	return b.String()
}

// --- Typr (active typing) ---

func (m Model) updateTypr(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.typr.finished {
		return m.updateTyprResults(msg)
	}

	kmsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	case kmsg.String() == "esc":
		m.screen = screenMenu
		m.cursor = 5
		return m, nil

	case kmsg.String() == "ctrl+r":
		story := m.typr.story
		m.typr = typrState{
			story: story,
			runes: []rune(story.Sv),
			input: make([]rune, 0, len([]rune(story.Sv))),
		}
		return m, nil

	case kmsg.Type == tea.KeyBackspace:
		if m.typr.pos > 0 {
			m.typr.pos--
			m.typr.input = m.typr.input[:m.typr.pos]
		}
		return m, nil

	default:
		var typed rune
		var hasChar bool

		switch kmsg.Type {
		case tea.KeyRunes:
			if len(kmsg.Runes) > 0 {
				typed = kmsg.Runes[0]
				hasChar = true
			}
		case tea.KeySpace:
			typed = ' '
			hasChar = true
		}

		if hasChar && m.typr.pos < len(m.typr.runes) {
			var cmd tea.Cmd
			if !m.typr.started {
				m.typr.started = true
				m.typr.startTime = time.Now()
				cmd = tickCmd()
			}

			m.typr.input = append(m.typr.input, typed)
			m.typr.pos++

			if m.typr.pos >= len(m.typr.runes) {
				m.typr.finished = true
				m.typr.endTime = time.Now()
			}

			return m, cmd
		}
	}

	return m, nil
}

func (m Model) viewTypr() string {
	if m.typr.finished {
		return m.viewTyprResults()
	}

	var b strings.Builder

	// Title
	b.WriteString("\n" + titleStyle.Render("  Story Typer") + "  " +
		swedishStyle.Render(m.typr.story.Title) + "\n\n")

	// Render Swedish text with per-character coloring and word wrap
	maxWidth := m.width - 6
	if maxWidth < 40 {
		maxWidth = 40
	}

	runes := m.typr.runes
	b.WriteString("  ")
	col := 0

	for i, r := range runes {
		// Word wrap: if at a space, check if the next word would overflow
		if r == ' ' && col > 0 && i+1 < len(runes) {
			nextLen := 0
			for j := i + 1; j < len(runes) && runes[j] != ' '; j++ {
				nextLen++
			}
			if col+1+nextLen > maxWidth {
				writeTyprChar(&b, m.typr, i, r)
				b.WriteString("\n  ")
				col = 0
				continue
			}
		}

		writeTyprChar(&b, m.typr, i, r)
		col++
	}
	b.WriteString("\n")

	// Separator + English translation
	b.WriteString("\n  " + dimStyle.Render("── Translation ──") + "\n")
	for _, line := range wrapText(m.typr.story.En, maxWidth) {
		b.WriteString("  " + dimStyle.Render(line) + "\n")
	}

	// Timer + live stats
	b.WriteString("\n")
	if m.typr.started {
		elapsed := int(m.typr.elapsed().Seconds())
		b.WriteString(fmt.Sprintf("  "+dimStyle.Render("⏱ %ds")+"    "+
			dimStyle.Render("%.0f WPM")+"    "+
			dimStyle.Render("%.0f%% accuracy"),
			elapsed, m.typr.wpm(), m.typr.accuracy()))
	} else {
		b.WriteString("  " + dimStyle.Render("Start typing to begin..."))
	}
	b.WriteString("\n")

	// Controls
	b.WriteString("\n" + dimStyle.Render("  esc quit  •  ctrl+r restart") + "\n")

	return b.String()
}

func writeTyprChar(b *strings.Builder, t typrState, i int, r rune) {
	ch := string(r)
	if i < t.pos {
		if t.input[i] == r {
			b.WriteString(ch) // default terminal color for correct
		} else {
			b.WriteString(typrErrorStyle.Render(ch))
		}
	} else if i == t.pos && !t.finished {
		b.WriteString(typrCursorStyle.Render(ch))
	} else {
		b.WriteString(dimStyle.Render(ch))
	}
}

// --- Typr Results ---

func (m Model) updateTyprResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "esc", "m":
			m.screen = screenMenu
			m.cursor = 5
		case "r":
			story := m.typr.story
			m.typr = typrState{
				story: story,
				runes: []rune(story.Sv),
				input: make([]rune, 0, len([]rune(story.Sv))),
			}
			m.screen = screenTypr
		case "enter", "b":
			m.screen = screenSelectStory
			m.cursor = 0
		}
	}
	return m, nil
}

func (m Model) viewTyprResults() string {
	var b strings.Builder

	elapsed := m.typr.elapsed()
	wpm := m.typr.wpm()
	acc := m.typr.accuracy()

	b.WriteString("\n" + titleStyle.Render("  Story Complete!") + "  " +
		swedishStyle.Render(m.typr.story.Title) + "\n\n")

	// Stats
	b.WriteString(fmt.Sprintf("  Time:      %s\n", progressStyle.Render(fmt.Sprintf("%ds", int(elapsed.Seconds())))))
	b.WriteString(fmt.Sprintf("  WPM:       %s\n", swedishStyle.Render(fmt.Sprintf("%.0f", wpm))))

	accStyle := correctStyle
	if acc < 90 {
		accStyle = wrongStyle
	}
	b.WriteString(fmt.Sprintf("  Accuracy:  %s\n", accStyle.Render(fmt.Sprintf("%.1f%%", acc))))

	// Error count
	errors := 0
	for i, r := range m.typr.input {
		if i < len(m.typr.runes) && r != m.typr.runes[i] {
			errors++
		}
	}
	if errors > 0 {
		b.WriteString(fmt.Sprintf("  Errors:    %s\n", wrongStyle.Render(fmt.Sprintf("%d", errors))))
	} else {
		b.WriteString(fmt.Sprintf("  Errors:    %s\n", correctStyle.Render("0 — perfect!")))
	}

	b.WriteString("\n")

	// Show the full story for review
	maxWidth := m.width - 6
	if maxWidth < 40 {
		maxWidth = 40
	}

	b.WriteString("  " + dimStyle.Render("── Swedish ──") + "\n")
	for _, line := range wrapText(m.typr.story.Sv, maxWidth) {
		b.WriteString("  " + swedishStyle.Render(line) + "\n")
	}
	b.WriteString("\n  " + dimStyle.Render("── English ──") + "\n")
	for _, line := range wrapText(m.typr.story.En, maxWidth) {
		b.WriteString("  " + englishStyle.Render(line) + "\n")
	}

	b.WriteString("\n" + dimStyle.Render("  r replay  •  enter stories  •  m menu") + "\n")
	return b.String()
}

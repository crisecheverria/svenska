package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/crisecheverria/svenska/ai"
	"github.com/crisecheverria/svenska/data"
	"github.com/crisecheverria/svenska/game"
	"github.com/crisecheverria/svenska/updater"
)

type screen int

const (
	screenMenu screen = iota
	screenSelectCategory
	screenSelectLevel
	screenSelectDirection
	screenPlaying
	screenFeedback
	screenResults
	screenStats
	screenHelp
)

// aiHelpMsg is sent when the AI help request completes
type aiHelpMsg struct {
	text string
	err  error
}

// updateAvailableMsg is sent when the background update check completes
type updateAvailableMsg struct {
	info updater.UpdateInfo
}

type Model struct {
	screen          screen
	cursor          int
	mode            game.Mode
	direction       game.Direction
	round           *game.Round
	stats           *game.Stats
	input           textinput.Model
	lastRight       bool
	width           int
	height          int
	helpText        string
	helpLoad        bool
	newAchievements []game.Achievement
	version         string
	updateAvailable string // new version string, empty if up to date
}

func NewModel(version string) Model {
	ti := textinput.New()
	ti.Placeholder = "Type your answer..."
	ti.CharLimit = 200
	ti.Width = 50

	return Model{
		screen:  screenMenu,
		stats:   game.LoadStats(),
		input:   ti,
		version: version,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.checkForUpdate())
}

func (m Model) checkForUpdate() tea.Cmd {
	version := m.version
	return func() tea.Msg {
		info := updater.CheckForUpdate(version)
		return updateAvailableMsg{info: info}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case aiHelpMsg:
		m.helpLoad = false
		if msg.err != nil {
			m.helpText = "Error: " + msg.err.Error()
		} else {
			m.helpText = msg.text
		}
		return m, nil
	case updateAvailableMsg:
		if msg.info.Available {
			m.updateAvailable = msg.info.NewVersion
		}
		return m, nil
	}

	switch m.screen {
	case screenMenu:
		return m.updateMenu(msg)
	case screenSelectCategory:
		return m.updateSelectCategory(msg)
	case screenSelectLevel:
		return m.updateSelectLevel(msg)
	case screenSelectDirection:
		return m.updateSelectDirection(msg)
	case screenPlaying:
		return m.updatePlaying(msg)
	case screenFeedback:
		return m.updateFeedback(msg)
	case screenResults:
		return m.updateResults(msg)
	case screenStats:
		return m.updateStats(msg)
	case screenHelp:
		return m.updateHelp(msg)
	}
	return m, nil
}

func (m Model) View() string {
	switch m.screen {
	case screenMenu:
		return m.viewMenu()
	case screenSelectCategory:
		return m.viewSelectCategory()
	case screenSelectLevel:
		return m.viewSelectLevel()
	case screenSelectDirection:
		return m.viewSelectDirection()
	case screenPlaying:
		return m.viewPlaying()
	case screenFeedback:
		return m.viewFeedback()
	case screenResults:
		return m.viewResults()
	case screenStats:
		return m.viewStats()
	case screenHelp:
		return m.viewHelp()
	}
	return ""
}

// --- Menu ---

func (m Model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 4 {
				m.cursor++
			}
		case "enter":
			switch m.cursor {
			case 0:
				m.mode = game.ModeVocabulary
				m.screen = screenSelectCategory
				m.cursor = 0
			case 1:
				m.mode = game.ModeTyping
				m.screen = screenSelectCategory
				m.cursor = 0
			case 2:
				m.mode = game.ModeTranslate
				m.screen = screenSelectLevel
				m.cursor = 0
			case 3:
				m.screen = screenStats
			case 4:
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m Model) viewMenu() string {
	var b strings.Builder

	b.WriteString("\n")
	for _, line := range flag {
		b.WriteString("  " + line + "\n")
	}
	b.WriteString("\n")

	b.WriteString(titleStyle.Render("  Svenska") + "  " + subtitleStyle.Render("Learn Swedish in the Terminal") + "\n\n")

	items := []struct{ title, desc string }{
		{"Vocabulary", "Translate words (SV↔EN)"},
		{"Typing", "Type Swedish words exactly"},
		{"Translate", "Translate full sentences"},
		{"Statistics", "View your progress"},
		{"Quit", "Exit the program"},
	}

	for i, item := range items {
		cursor := "  "
		style := menuItemStyle
		if i == m.cursor {
			cursor = promptStyle.Render("▸ ")
			style = selectedStyle
		}
		b.WriteString(cursor + style.Render(item.title) + "  " + dimStyle.Render(item.desc) + "\n")
	}

	b.WriteString("\n")
	lvl, lvlName := m.stats.Level()
	nextXP := m.stats.NextLevelXP()
	streakStr := ""
	if m.stats.Streak > 0 {
		streakStr = fmt.Sprintf("  |  Streak: %d day", m.stats.Streak)
		if m.stats.Streak > 1 {
			streakStr += "s"
		}
	}
	b.WriteString(dimStyle.Render(fmt.Sprintf("  Lv.%d %s  |  XP: %d/%d  |  Accuracy: %.0f%%%s",
		lvl, lvlName, m.stats.XP, nextXP, m.stats.Accuracy(), streakStr)) + "\n")

	if m.updateAvailable != "" {
		b.WriteString("\n" + swedishStyle.Render(fmt.Sprintf("  New version v%s available!", m.updateAvailable)) +
			dimStyle.Render("  Run: ") + progressStyle.Render("svenska update") + "\n")
	}

	b.WriteString("\n" + dimStyle.Render("  ↑↓/jk navigate  •  enter select  •  q quit") + "\n")

	return b.String()
}

// --- Category selection ---

func (m Model) categoryItems() []struct{ key, name string } {
	items := []struct{ key, name string }{
		{"all", "All Categories (mixed)"},
	}
	for _, cat := range data.Categories {
		items = append(items, struct{ key, name string }{cat.Key, fmt.Sprintf("%s (%d words)", cat.Name, len(cat.Words))})
	}
	return items
}

func (m Model) updateSelectCategory(msg tea.Msg) (tea.Model, tea.Cmd) {
	items := m.categoryItems()
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "esc", "b":
			m.screen = screenMenu
			m.cursor = 0
			return m, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(items)-1 {
				m.cursor++
			}
		case "enter":
			cat := items[m.cursor].key
			if m.mode == game.ModeTyping {
				m.round = game.NewTypingRound(cat)
				m.screen = screenPlaying
				m.input.SetValue("")
				m.input.Focus()
				return m, textinput.Blink
			}
			m.round = nil
			m.screen = screenSelectDirection
			// store category temporarily in a round stub
			m.round = &game.Round{Category: cat}
			m.cursor = 0
		}
	}
	return m, nil
}

func (m Model) viewSelectCategory() string {
	var b strings.Builder
	items := m.categoryItems()

	b.WriteString("\n" + titleStyle.Render("  Select Category") + "  " + dimStyle.Render(m.mode.String()) + "\n\n")

	// Calculate visible window for scrolling
	maxVisible := m.height - 6
	if maxVisible < 5 {
		maxVisible = 5
	}
	start := 0
	if m.cursor >= maxVisible {
		start = m.cursor - maxVisible + 1
	}
	end := start + maxVisible
	if end > len(items) {
		end = len(items)
	}

	if start > 0 {
		b.WriteString(dimStyle.Render("  ↑ more...") + "\n")
	}

	for i := start; i < end; i++ {
		cursor := "  "
		style := menuItemStyle
		if i == m.cursor {
			cursor = promptStyle.Render("▸ ")
			style = selectedStyle
		}
		b.WriteString(cursor + style.Render(items[i].name) + "\n")
	}

	if end < len(items) {
		b.WriteString(dimStyle.Render("  ↓ more...") + "\n")
	}

	b.WriteString("\n" + dimStyle.Render("  ↑↓/jk navigate  •  enter select  •  b back") + "\n")
	return b.String()
}

// --- Level selection ---

func (m Model) levelItems() []struct{ key, name string } {
	items := []struct{ key, name string }{
		{"all", "All Levels (mixed)"},
	}
	for _, lvl := range data.Levels {
		items = append(items, struct{ key, name string }{lvl.Key, fmt.Sprintf("%s (%d sentences)", lvl.Name, len(lvl.Sentences))})
	}
	return items
}

func (m Model) updateSelectLevel(msg tea.Msg) (tea.Model, tea.Cmd) {
	items := m.levelItems()
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "esc", "b":
			m.screen = screenMenu
			m.cursor = 0
			return m, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(items)-1 {
				m.cursor++
			}
		case "enter":
			lvl := items[m.cursor].key
			m.round = &game.Round{Category: lvl}
			m.screen = screenSelectDirection
			m.cursor = 0
		}
	}
	return m, nil
}

func (m Model) viewSelectLevel() string {
	var b strings.Builder
	items := m.levelItems()

	b.WriteString("\n" + titleStyle.Render("  Select Level") + "  " + dimStyle.Render("Translate") + "\n\n")

	for i, item := range items {
		cursor := "  "
		style := menuItemStyle
		if i == m.cursor {
			cursor = promptStyle.Render("▸ ")
			style = selectedStyle
		}
		b.WriteString(cursor + style.Render(item.name) + "\n")
	}

	b.WriteString("\n" + dimStyle.Render("  ↑↓/jk navigate  •  enter select  •  b back") + "\n")
	return b.String()
}

// --- Direction selection ---

func (m Model) updateSelectDirection(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "esc", "b":
			if m.mode == game.ModeVocabulary {
				m.screen = screenSelectCategory
			} else {
				m.screen = screenSelectLevel
			}
			m.cursor = 0
			return m, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 1 {
				m.cursor++
			}
		case "enter":
			cat := m.round.Category
			if m.cursor == 0 {
				m.direction = game.SvToEn
			} else {
				m.direction = game.EnToSv
			}
			if m.mode == game.ModeVocabulary {
				m.round = game.NewVocabularyRound(cat, m.direction)
			} else {
				m.round = game.NewTranslateRound(cat, m.direction)
			}
			m.screen = screenPlaying
			m.input.SetValue("")
			m.input.Focus()
			return m, textinput.Blink
		}
	}
	return m, nil
}

func (m Model) viewSelectDirection() string {
	var b strings.Builder
	b.WriteString("\n" + titleStyle.Render("  Translation Direction") + "\n\n")

	dirs := []string{"Svenska → English", "English → Svenska"}
	for i, d := range dirs {
		cursor := "  "
		style := menuItemStyle
		if i == m.cursor {
			cursor = promptStyle.Render("▸ ")
			style = selectedStyle
		}
		b.WriteString(cursor + style.Render(d) + "\n")
	}

	b.WriteString("\n" + dimStyle.Render("  ↑↓/jk navigate  •  enter select  •  b back") + "\n")
	return b.String()
}

// --- Playing ---

func (m Model) updatePlaying(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "esc":
			m.screen = screenMenu
			m.cursor = 0
			return m, nil
		case "enter":
			answer := m.input.Value()
			if answer == "" {
				return m, nil
			}
			if strings.TrimSpace(answer) == "?" {
				m.input.SetValue("")
				m.helpText = ""
				m.helpLoad = true
				m.screen = screenHelp
				ch := m.round.CurrentChallenge()
				return m, func() tea.Msg {
					result := ai.GetHelp(ch.Sv, ch.En, m.mode.String())
					return aiHelpMsg{text: result.Text, err: result.Err}
				}
			}
			m.lastRight = m.round.Submit(answer)
			m.input.SetValue("")
			m.screen = screenFeedback
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) viewPlaying() string {
	var b strings.Builder
	ch := m.round.CurrentChallenge()

	progress := fmt.Sprintf("%d/%d", m.round.Current+1, m.round.Total())
	score := fmt.Sprintf("✓ %d  ✗ %d", m.round.Correct, m.round.Wrong)

	b.WriteString("\n" + titleStyle.Render("  "+m.mode.String()) + "  " +
		progressStyle.Render(progress) + "  " + dimStyle.Render(score) + "\n\n")

	switch m.mode {
	case game.ModeVocabulary:
		if m.direction == game.SvToEn {
			b.WriteString("  " + swedishStyle.Render(ch.Prompt) + "\n")
			b.WriteString("  " + dimStyle.Render("Translate to English") + "\n\n")
		} else {
			b.WriteString("  " + englishStyle.Render(ch.Prompt) + "\n")
			b.WriteString("  " + dimStyle.Render("Translate to Swedish") + "\n\n")
		}
	case game.ModeTyping:
		b.WriteString("  " + swedishStyle.Render(ch.Prompt) + "  " + dimStyle.Render("("+ch.En+")") + "\n")
		b.WriteString("  " + dimStyle.Render("Type the Swedish word exactly") + "\n\n")
	case game.ModeTranslate:
		if m.direction == game.SvToEn {
			b.WriteString("  " + swedishStyle.Render(ch.Prompt) + "\n")
			b.WriteString("  " + dimStyle.Render("Translate to English") + "\n\n")
		} else {
			b.WriteString("  " + englishStyle.Render(ch.Prompt) + "\n")
			b.WriteString("  " + dimStyle.Render("Translate to Swedish") + "\n\n")
		}
	}

	b.WriteString("  " + promptStyle.Render("→ ") + m.input.View() + "\n")
	b.WriteString("\n" + dimStyle.Render("  enter submit  •  ? ai help  •  esc quit") + "\n")

	return b.String()
}

// --- Feedback ---

func (m Model) updateFeedback(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "enter", " ":
			if m.round.Done() {
				m.stats.RecordRound(m.round)
				m.newAchievements = m.stats.CheckAchievements(m.round)
				m.stats.Save()
				m.screen = screenResults
				m.cursor = 0
			} else {
				m.screen = screenPlaying
				m.input.Focus()
				return m, textinput.Blink
			}
		}
	}
	return m, nil
}

func (m Model) viewFeedback() string {
	var b strings.Builder
	last := m.round.Answers[len(m.round.Answers)-1]

	b.WriteString("\n")
	if m.lastRight {
		b.WriteString("  " + correctStyle.Render("✓ Correct!") + "\n\n")
	} else {
		b.WriteString("  " + wrongStyle.Render("✗ Wrong!") + "\n\n")
		b.WriteString("  " + dimStyle.Render("Your answer: ") + wrongStyle.Render(last.Given) + "\n")
	}

	b.WriteString("  " + dimStyle.Render("SV: ") + swedishStyle.Render(last.Sv) + "\n")
	b.WriteString("  " + dimStyle.Render("EN: ") + englishStyle.Render(last.En) + "\n")

	b.WriteString("\n" + dimStyle.Render("  Press enter to continue") + "\n")
	return b.String()
}

// --- Results ---

func (m Model) updateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "m", "esc":
			m.screen = screenMenu
			m.cursor = 0
		case "r":
			// Replay same settings
			cat := m.round.Category
			switch m.round.Mode {
			case game.ModeVocabulary:
				m.round = game.NewVocabularyRound(cat, m.round.Direction)
			case game.ModeTyping:
				m.round = game.NewTypingRound(cat)
			case game.ModeTranslate:
				m.round = game.NewTranslateRound(cat, m.round.Direction)
			}
			m.screen = screenPlaying
			m.input.SetValue("")
			m.input.Focus()
			return m, textinput.Blink
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) viewResults() string {
	var b strings.Builder

	perfect := m.round.Correct == m.round.Total()
	pct := float64(m.round.Correct) / float64(m.round.Total()) * 100

	if perfect {
		b.WriteString("\n")
		b.WriteString(correctStyle.Render("  ★ ★ ★  PERFEKT!  ★ ★ ★") + "\n")
		b.WriteString(swedishStyle.Render("  Fantastiskt! Du fick alla rätt!") + "\n\n")
	} else {
		b.WriteString("\n" + titleStyle.Render("  Results") + "\n\n")
	}

	b.WriteString(fmt.Sprintf("  Score: %s / %d (%.0f%%)\n",
		correctStyle.Render(fmt.Sprintf("%d", m.round.Correct)),
		m.round.Total(), pct))

	// XP earned
	xp := m.stats.XPForRound(m.round)
	xpLine := fmt.Sprintf("  +%d XP", xp)
	if perfect {
		xpLine += " (includes perfect bonus!)"
	}
	b.WriteString(progressStyle.Render(xpLine) + "\n\n")

	for i, a := range m.round.Answers {
		mark := correctStyle.Render("✓")
		if !a.Correct {
			mark = wrongStyle.Render("✗")
		}
		b.WriteString(fmt.Sprintf("  %s %2d. %s → %s",
			mark, i+1,
			swedishStyle.Render(a.Sv),
			englishStyle.Render(a.En)))
		if !a.Correct {
			b.WriteString("  " + dimStyle.Render("(you: "+a.Given+")"))
		}
		b.WriteString("\n")
	}

	// Show newly unlocked achievements
	if len(m.newAchievements) > 0 {
		b.WriteString("\n" + swedishStyle.Render("  ── Achievement Unlocked! ──") + "\n\n")
		for _, a := range m.newAchievements {
			b.WriteString(fmt.Sprintf("  %s  %s  %s\n",
				swedishStyle.Render(a.Icon),
				correctStyle.Render(a.Name),
				dimStyle.Render(a.Desc)))
		}
	}

	b.WriteString("\n" + dimStyle.Render("  m menu  •  r replay  •  q quit") + "\n")
	return b.String()
}

// --- Stats ---

func (m Model) updateStats(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "esc", "b", "q", "enter":
			m.screen = screenMenu
			m.cursor = 0
		}
	}
	return m, nil
}

func (m Model) viewStats() string {
	var b strings.Builder
	s := m.stats

	lvl, lvlName := s.Level()
	nextXP := s.NextLevelXP()

	b.WriteString("\n" + titleStyle.Render("  Statistics") + "\n\n")

	// Level & XP
	b.WriteString(fmt.Sprintf("  Level:               %s\n", swedishStyle.Render(fmt.Sprintf("Lv.%d %s", lvl, lvlName))))
	b.WriteString(fmt.Sprintf("  XP:                  %s\n", progressStyle.Render(fmt.Sprintf("%d / %d", s.XP, nextXP))))
	xpBar := renderProgressBar(s.XP, nextXP, 20)
	b.WriteString(fmt.Sprintf("                       %s\n", xpBar))

	// Streak
	streakStr := fmt.Sprintf("%d day", s.Streak)
	if s.Streak != 1 {
		streakStr += "s"
	}
	b.WriteString(fmt.Sprintf("  Daily streak:        %s\n", swedishStyle.Render(streakStr)))

	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  Sessions completed:  %s\n", correctStyle.Render(fmt.Sprintf("%d", s.Sessions))))
	b.WriteString(fmt.Sprintf("  Perfect rounds:      %s\n", correctStyle.Render(fmt.Sprintf("%d", s.PerfectRounds))))
	b.WriteString(fmt.Sprintf("  Total questions:     %d\n", s.TotalPlayed))
	b.WriteString(fmt.Sprintf("  Correct answers:     %s\n", correctStyle.Render(fmt.Sprintf("%d", s.TotalCorrect))))
	b.WriteString(fmt.Sprintf("  Wrong answers:       %s\n", wrongStyle.Render(fmt.Sprintf("%d", s.TotalWrong))))
	b.WriteString(fmt.Sprintf("  Accuracy:            %.1f%%\n", s.Accuracy()))
	b.WriteString(fmt.Sprintf("  Unique words learned: %d\n", len(s.WordsLearned)))

	// Achievements
	b.WriteString("\n" + titleStyle.Render("  Achievements") + "\n\n")
	for _, a := range game.AllAchievements {
		if s.HasAchievement(a.Key) {
			b.WriteString(fmt.Sprintf("  %s  %s  %s  %s\n",
				swedishStyle.Render(a.Icon),
				correctStyle.Render(a.Name),
				dimStyle.Render(a.Desc),
				dimStyle.Render("("+s.Achievements[a.Key]+")")))
		} else {
			b.WriteString(fmt.Sprintf("  %s  %s  %s\n",
				dimStyle.Render("?"),
				dimStyle.Render(a.Name),
				dimStyle.Render(a.Desc)))
		}
	}

	// Category mastery (top practiced)
	if len(s.CategoryStats) > 0 {
		b.WriteString("\n" + titleStyle.Render("  Category Mastery") + "\n\n")
		for _, cat := range data.Categories {
			mastery := s.CategoryMastery(cat.Key)
			if mastery > 0 {
				bar := renderProgressBar(int(mastery), 100, 15)
				b.WriteString(fmt.Sprintf("  %-30s %s %.0f%%\n", cat.Name, bar, mastery))
			}
		}
	}

	b.WriteString("\n" + dimStyle.Render("  Press any key to go back") + "\n")
	return b.String()
}

func renderProgressBar(current, max, width int) string {
	if max <= 0 {
		max = 1
	}
	ratio := float64(current) / float64(max)
	if ratio > 1 {
		ratio = 1
	}
	filled := int(ratio * float64(width))
	empty := width - filled
	bar := progressStyle.Render(strings.Repeat("█", filled)) + dimStyle.Render(strings.Repeat("░", empty))
	return bar
}

// --- AI Help ---

func (m Model) updateHelp(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if !m.helpLoad {
			switch msg.String() {
			case "enter", "esc", "q", "b":
				m.screen = screenPlaying
				m.input.Focus()
				return m, textinput.Blink
			}
		}
	}
	return m, nil
}

func (m Model) viewHelp() string {
	var b strings.Builder
	ch := m.round.CurrentChallenge()

	b.WriteString("\n" + titleStyle.Render("  AI Help") + "\n\n")
	b.WriteString("  " + dimStyle.Render("SV: ") + swedishStyle.Render(ch.Sv) + "\n")
	b.WriteString("  " + dimStyle.Render("EN: ") + englishStyle.Render(ch.En) + "\n\n")

	if m.helpLoad {
		b.WriteString("  " + dimStyle.Render("Loading...") + "\n")
	} else {
		for _, line := range strings.Split(m.helpText, "\n") {
			b.WriteString("  " + line + "\n")
		}
	}

	if !m.helpLoad {
		b.WriteString("\n" + dimStyle.Render("  Press enter to go back and answer") + "\n")
	}

	return b.String()
}

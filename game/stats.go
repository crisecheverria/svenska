package game

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type CategoryStats struct {
	Correct int `json:"correct"`
	Total   int `json:"total"`
}

type Stats struct {
	TotalCorrect  int                      `json:"total_correct"`
	TotalWrong    int                      `json:"total_wrong"`
	TotalPlayed   int                      `json:"total_played"`
	Sessions      int                      `json:"sessions"`
	WordsLearned  map[string]int           `json:"words_learned"`
	XP            int                      `json:"xp"`
	Streak        int                      `json:"streak"`
	LastPlayed    string                   `json:"last_played"`
	PerfectRounds int                      `json:"perfect_rounds"`
	CategoryStats map[string]CategoryStats `json:"category_stats"`
	Achievements  map[string]string        `json:"achievements"` // key -> date unlocked
}

type Achievement struct {
	Key  string
	Name string
	Desc string
	Icon string
}

var AllAchievements = []Achievement{
	{"first_steps", "Första stegen", "Complete your first round", ">>"},
	{"perfect_round", "Perfekt!", "Get 10/10 in a round", "★★"},
	{"polyglot", "Polyglott", "Practice all 27 categories", "◆◆"},
	{"on_fire", "Eldsjäl", "Reach a 7-day streak", "▲▲"},
	{"century", "Hundra", "Answer 100 questions correctly", "●●"},
}

// XP rewards
const (
	XPPerCorrect  = 10
	XPStreakBonus  = 5  // extra per correct when on a 3+ answer streak
	XPPerfectRound = 50 // bonus for 10/10
)

// Level thresholds and names
var levels = []struct {
	Name string
	XP   int
}{
	{"Nybörjare", 0},
	{"Elev", 100},
	{"Studerande", 300},
	{"Praktikant", 600},
	{"Kunnig", 1000},
	{"Avancerad", 1800},
	{"Expert", 3000},
	{"Mästare", 5000},
}

func statsPath() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".local", "share", "svenska")
	os.MkdirAll(dir, 0o755)
	return filepath.Join(dir, "stats.json")
}

func LoadStats() *Stats {
	s := &Stats{
		WordsLearned:  make(map[string]int),
		CategoryStats: make(map[string]CategoryStats),
		Achievements:  make(map[string]string),
	}
	data, err := os.ReadFile(statsPath())
	if err != nil {
		return s
	}
	json.Unmarshal(data, s)
	if s.WordsLearned == nil {
		s.WordsLearned = make(map[string]int)
	}
	if s.CategoryStats == nil {
		s.CategoryStats = make(map[string]CategoryStats)
	}
	if s.Achievements == nil {
		s.Achievements = make(map[string]string)
	}
	return s
}

func (s *Stats) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(statsPath(), data, 0o644)
}

func (s *Stats) RecordRound(r *Round) {
	s.TotalCorrect += r.Correct
	s.TotalWrong += r.Wrong
	s.TotalPlayed += r.Total()
	s.Sessions++

	// XP: base points per correct answer
	roundXP := r.Correct * XPPerCorrect

	// XP: streak bonus — count consecutive correct answers
	streak := 0
	for _, a := range r.Answers {
		if a.Correct {
			streak++
			if streak >= 3 {
				roundXP += XPStreakBonus
			}
		} else {
			streak = 0
		}
	}

	// Perfect round bonus
	if r.Correct == r.Total() {
		roundXP += XPPerfectRound
		s.PerfectRounds++
	}

	s.XP += roundXP

	// Words learned
	for _, a := range r.Answers {
		if a.Correct {
			s.WordsLearned[a.Sv]++
		}
	}

	// Category mastery
	cat := s.CategoryStats[r.Category]
	cat.Correct += r.Correct
	cat.Total += r.Total()
	s.CategoryStats[r.Category] = cat

	// Daily streak
	today := time.Now().Format("2006-01-02")
	if s.LastPlayed == "" {
		s.Streak = 1
	} else if s.LastPlayed == today {
		// Already played today, streak unchanged
	} else {
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		if s.LastPlayed == yesterday {
			s.Streak++
		} else {
			s.Streak = 1
		}
	}
	s.LastPlayed = today
}

func (s *Stats) Accuracy() float64 {
	total := s.TotalCorrect + s.TotalWrong
	if total == 0 {
		return 0
	}
	return float64(s.TotalCorrect) / float64(total) * 100
}

func (s *Stats) Level() (int, string) {
	lvl := 0
	for i, l := range levels {
		if s.XP >= l.XP {
			lvl = i
		}
	}
	return lvl + 1, levels[lvl].Name
}

func (s *Stats) NextLevelXP() int {
	for _, l := range levels {
		if s.XP < l.XP {
			return l.XP
		}
	}
	return levels[len(levels)-1].XP
}

func (s *Stats) XPForRound(r *Round) int {
	xp := r.Correct * XPPerCorrect
	streak := 0
	for _, a := range r.Answers {
		if a.Correct {
			streak++
			if streak >= 3 {
				xp += XPStreakBonus
			}
		} else {
			streak = 0
		}
	}
	if r.Correct == r.Total() {
		xp += XPPerfectRound
	}
	return xp
}

func (s *Stats) CategoryMastery(catKey string) float64 {
	cat, ok := s.CategoryStats[catKey]
	if !ok || cat.Total == 0 {
		return 0
	}
	return float64(cat.Correct) / float64(cat.Total) * 100
}

// CheckAchievements evaluates all achievements and returns any newly unlocked ones.
func (s *Stats) CheckAchievements(r *Round) []Achievement {
	today := time.Now().Format("2006-01-02")
	var newlyUnlocked []Achievement

	unlock := func(key string) {
		if _, already := s.Achievements[key]; !already {
			s.Achievements[key] = today
			for _, a := range AllAchievements {
				if a.Key == key {
					newlyUnlocked = append(newlyUnlocked, a)
					break
				}
			}
		}
	}

	// First Steps: completed at least 1 session
	if s.Sessions >= 1 {
		unlock("first_steps")
	}

	// Perfect Round: got 10/10
	if r != nil && r.Correct == r.Total() {
		unlock("perfect_round")
	}

	// Polyglot: practiced all 27 categories (excluding "all")
	if len(s.CategoryStats) >= 27 {
		count := 0
		for k := range s.CategoryStats {
			if k != "all" {
				count++
			}
		}
		if count >= 27 {
			unlock("polyglot")
		}
	}

	// On Fire: 7-day streak
	if s.Streak >= 7 {
		unlock("on_fire")
	}

	// Century: 100 correct answers
	if s.TotalCorrect >= 100 {
		unlock("century")
	}

	return newlyUnlocked
}

func (s *Stats) HasAchievement(key string) bool {
	_, ok := s.Achievements[key]
	return ok
}

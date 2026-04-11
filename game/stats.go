package game

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/crisecheverria/svenska/data"
)

type CategoryStats struct {
	Correct int `json:"correct"`
	Total   int `json:"total"`
}

// WordReview holds SM-2 spaced repetition data for a single word/sentence.
type WordReview struct {
	EF          float64 `json:"ef"`          // easiness factor (≥1.3, default 2.5)
	Interval    int     `json:"interval"`    // days until next review
	Repetitions int     `json:"repetitions"` // consecutive correct answers
	NextReview  string  `json:"next_review"` // YYYY-MM-DD
}

type Stats struct {
	TotalCorrect   int                      `json:"total_correct"`
	TotalWrong     int                      `json:"total_wrong"`
	TotalPlayed    int                      `json:"total_played"`
	Sessions       int                      `json:"sessions"`
	WordsLearned   map[string]int           `json:"words_learned"`
	XP             int                      `json:"xp"`
	Streak         int                      `json:"streak"`
	LastPlayed     string                   `json:"last_played"`
	PerfectRounds  int                      `json:"perfect_rounds"`
	CategoryStats  map[string]CategoryStats `json:"category_stats"`
	Achievements   map[string]string        `json:"achievements"` // key -> date unlocked
	BestSpeedScore int                      `json:"best_speed_score"`
	HardcoreRounds int                      `json:"hardcore_rounds"`
	Reviews        map[string]*WordReview   `json:"reviews"`
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
	{"speed_demon", "Snabbis", "Complete a speed round", "▸▸"},
	{"hardcore_hero", "Hårding", "Complete a hardcore round", "■■"},
	{"speed_master", "Blixt", "Get 20+ in a speed round", "◇◇"},
}

// RoadmapLevel describes what each level unlocks/recommends.
type RoadmapLevel struct {
	Level int
	Name  string
	XP    int
	Desc  string
}

var RoadmapLevels = []RoadmapLevel{
	{1, "Nybörjare", 0, "Greetings, Numbers, Pronouns"},
	{2, "Elev", 100, "Colors, Family, Food & Drink"},
	{3, "Studerande", 300, "Verbs, Adjectives, Beginner sentences"},
	{4, "Praktikant", 600, "Professions, Shopping, Elementary sentences"},
	{5, "Kunnig", 1000, "Intermediate sentences, Speed Rounds"},
	{6, "Avancerad", 1800, "Challenge yourself with Hardcore mode"},
	{7, "Expert", 3000, "Master Advanced sentences"},
	{8, "Mästare", 5000, "Du har bemästrat svenskan!"},
}

// XP rewards
const (
	XPPerCorrect   = 10
	XPStreakBonus   = 5  // extra per correct when on a 3+ answer streak
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
		Reviews:       make(map[string]*WordReview),
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
	if s.Reviews == nil {
		s.Reviews = make(map[string]*WordReview)
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

// updateReview applies the SM-2 algorithm for a single word.
// quality: 4 for correct, 1 for incorrect.
func (s *Stats) updateReview(sv string, correct bool) {
	r, ok := s.Reviews[sv]
	if !ok {
		r = &WordReview{EF: 2.5}
		s.Reviews[sv] = r
	}

	quality := 1.0
	if correct {
		quality = 4.0
	}

	// Update easiness factor
	r.EF = r.EF + (0.1 - (5-quality)*(0.08+(5-quality)*0.02))
	if r.EF < 1.3 {
		r.EF = 1.3
	}

	if correct {
		switch r.Repetitions {
		case 0:
			r.Interval = 1
		case 1:
			r.Interval = 6
		default:
			r.Interval = int(math.Round(float64(r.Interval) * r.EF))
		}
		r.Repetitions++
	} else {
		r.Repetitions = 0
		r.Interval = 1
	}

	r.NextReview = time.Now().AddDate(0, 0, r.Interval).Format("2006-01-02")
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
	if r.Correct == r.Total() && !r.Timed {
		roundXP += XPPerfectRound
		s.PerfectRounds++
	}

	// Hardcore: double XP
	if r.Hardcore {
		roundXP *= 2
		s.HardcoreRounds++
	}

	// Speed round: track best score
	if r.Timed && r.Correct > s.BestSpeedScore {
		s.BestSpeedScore = r.Correct
	}

	s.XP += roundXP

	// Words learned + spaced repetition update
	for _, a := range r.Answers {
		if a.Correct {
			s.WordsLearned[a.Sv]++
		}
		s.updateReview(a.Sv, a.Correct)
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
	if r.Correct == r.Total() && !r.Timed {
		xp += XPPerfectRound
	}
	if r.Hardcore {
		xp *= 2
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

// PrioritizeWords sorts words by SM-2 review schedule.
// Due/overdue words come first, then new words, then words not yet due.
// Within each group, order is randomized.
func (s *Stats) PrioritizeWords(words []data.Word) []data.Word {
	today := time.Now().Truncate(24 * time.Hour)

	type bucket struct {
		word     data.Word
		priority int // 2=due/overdue, 1=new, 0=not yet due
		overdue  int // days overdue (for sub-sorting)
	}

	buckets := make([]bucket, len(words))
	for i, w := range words {
		r, ok := s.Reviews[w.Sv]
		if !ok {
			buckets[i] = bucket{word: w, priority: 1}
			continue
		}
		next, err := time.Parse("2006-01-02", r.NextReview)
		if err != nil {
			buckets[i] = bucket{word: w, priority: 1}
			continue
		}
		if !next.After(today) {
			days := int(today.Sub(next).Hours() / 24)
			buckets[i] = bucket{word: w, priority: 2, overdue: days}
		} else {
			buckets[i] = bucket{word: w, priority: 0}
		}
	}

	// Shuffle within each priority group for variety
	sort.SliceStable(buckets, func(i, j int) bool {
		if buckets[i].priority != buckets[j].priority {
			return buckets[i].priority > buckets[j].priority
		}
		if buckets[i].priority == 2 {
			return buckets[i].overdue > buckets[j].overdue
		}
		return false
	})

	result := make([]data.Word, len(buckets))
	for i, b := range buckets {
		result[i] = b.word
	}
	return result
}

// PrioritizeSentences sorts sentences by SM-2 review schedule.
func (s *Stats) PrioritizeSentences(sentences []data.Sentence) []data.Sentence {
	today := time.Now().Truncate(24 * time.Hour)

	type bucket struct {
		sentence data.Sentence
		priority int
		overdue  int
	}

	buckets := make([]bucket, len(sentences))
	for i, sent := range sentences {
		r, ok := s.Reviews[sent.Sv]
		if !ok {
			buckets[i] = bucket{sentence: sent, priority: 1}
			continue
		}
		next, err := time.Parse("2006-01-02", r.NextReview)
		if err != nil {
			buckets[i] = bucket{sentence: sent, priority: 1}
			continue
		}
		if !next.After(today) {
			days := int(today.Sub(next).Hours() / 24)
			buckets[i] = bucket{sentence: sent, priority: 2, overdue: days}
		} else {
			buckets[i] = bucket{sentence: sent, priority: 0}
		}
	}

	sort.SliceStable(buckets, func(i, j int) bool {
		if buckets[i].priority != buckets[j].priority {
			return buckets[i].priority > buckets[j].priority
		}
		if buckets[i].priority == 2 {
			return buckets[i].overdue > buckets[j].overdue
		}
		return false
	})

	result := make([]data.Sentence, len(buckets))
	for i, b := range buckets {
		result[i] = b.sentence
	}
	return result
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

	// Speed Demon: complete a speed round
	if r != nil && r.Timed {
		unlock("speed_demon")
	}

	// Hardcore Hero: complete a hardcore round
	if r != nil && r.Hardcore && !r.Timed {
		unlock("hardcore_hero")
	}

	// Speed Master: 20+ correct in a speed round
	if r != nil && r.Timed && r.Correct >= 20 {
		unlock("speed_master")
	}

	return newlyUnlocked
}

func (s *Stats) HasAchievement(key string) bool {
	_, ok := s.Achievements[key]
	return ok
}

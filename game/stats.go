package game

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Stats struct {
	TotalCorrect int            `json:"total_correct"`
	TotalWrong   int            `json:"total_wrong"`
	TotalPlayed  int            `json:"total_played"`
	Sessions     int            `json:"sessions"`
	WordsLearned map[string]int `json:"words_learned"`
}

func statsPath() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".local", "share", "svenska")
	os.MkdirAll(dir, 0o755)
	return filepath.Join(dir, "stats.json")
}

func LoadStats() *Stats {
	s := &Stats{WordsLearned: make(map[string]int)}
	data, err := os.ReadFile(statsPath())
	if err != nil {
		return s
	}
	json.Unmarshal(data, s)
	if s.WordsLearned == nil {
		s.WordsLearned = make(map[string]int)
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

	for _, a := range r.Answers {
		if a.Correct {
			s.WordsLearned[a.Sv]++
		}
	}
}

func (s *Stats) Accuracy() float64 {
	total := s.TotalCorrect + s.TotalWrong
	if total == 0 {
		return 0
	}
	return float64(s.TotalCorrect) / float64(total) * 100
}

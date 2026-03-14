package game

import (
	"math/rand"
	"strings"

	"github.com/crisecheverria/svenska/data"
)

type Mode int

const (
	ModeVocabulary Mode = iota
	ModeTyping
	ModeTranslate
)

func (m Mode) String() string {
	switch m {
	case ModeVocabulary:
		return "Vocabulary"
	case ModeTyping:
		return "Typing"
	case ModeTranslate:
		return "Translate"
	}
	return ""
}

type Direction int

const (
	SvToEn Direction = iota
	EnToSv
)

type Answer struct {
	Sv      string
	En      string
	Given   string
	Correct bool
}

type Challenge struct {
	Prompt   string // what to show the user
	Expected string // the correct answer
	Sv       string
	En       string
}

const ChallengesPerRound = 10
const SpeedChallengesPool = 100

type Round struct {
	Mode       Mode
	Direction  Direction
	Category   string // category key or level key
	Challenges []Challenge
	Answers    []Answer
	Current    int
	Correct  int
	Wrong    int
	Hardcore bool
	Timed    bool
}

func NewVocabularyRound(categoryKey string, dir Direction) *Round {
	var words []data.Word
	if categoryKey == "all" {
		words = data.AllWords()
	} else {
		for _, cat := range data.Categories {
			if cat.Key == categoryKey {
				words = cat.Words
				break
			}
		}
	}

	shuffled := make([]data.Word, len(words))
	copy(shuffled, words)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	count := ChallengesPerRound
	if count > len(shuffled) {
		count = len(shuffled)
	}

	challenges := make([]Challenge, count)
	for i := 0; i < count; i++ {
		w := shuffled[i]
		if dir == SvToEn {
			challenges[i] = Challenge{Prompt: w.Sv, Expected: w.En, Sv: w.Sv, En: w.En}
		} else {
			challenges[i] = Challenge{Prompt: w.En, Expected: w.Sv, Sv: w.Sv, En: w.En}
		}
	}

	return &Round{
		Mode:       ModeVocabulary,
		Direction:  dir,
		Category:   categoryKey,
		Challenges: challenges,
		Answers:    make([]Answer, 0, count),
	}
}

func NewTypingRound(categoryKey string) *Round {
	var words []data.Word
	if categoryKey == "all" {
		words = data.AllWords()
	} else {
		for _, cat := range data.Categories {
			if cat.Key == categoryKey {
				words = cat.Words
				break
			}
		}
	}

	shuffled := make([]data.Word, len(words))
	copy(shuffled, words)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	count := ChallengesPerRound
	if count > len(shuffled) {
		count = len(shuffled)
	}

	challenges := make([]Challenge, count)
	for i := 0; i < count; i++ {
		w := shuffled[i]
		challenges[i] = Challenge{Prompt: w.Sv, Expected: w.Sv, Sv: w.Sv, En: w.En}
	}

	return &Round{
		Mode:       ModeTyping,
		Category:   categoryKey,
		Challenges: challenges,
		Answers:    make([]Answer, 0, count),
	}
}

func NewTypingSentenceRound(levelKey string) *Round {
	var sentences []data.Sentence
	if levelKey == "all" {
		sentences = data.AllSentences()
	} else {
		for _, lvl := range data.Levels {
			if lvl.Key == levelKey {
				sentences = lvl.Sentences
				break
			}
		}
	}

	shuffled := make([]data.Sentence, len(sentences))
	copy(shuffled, sentences)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	count := ChallengesPerRound
	if count > len(shuffled) {
		count = len(shuffled)
	}

	challenges := make([]Challenge, count)
	for i := 0; i < count; i++ {
		s := shuffled[i]
		challenges[i] = Challenge{Prompt: s.Sv, Expected: s.Sv, Sv: s.Sv, En: s.En}
	}

	return &Round{
		Mode:       ModeTyping,
		Category:   levelKey,
		Challenges: challenges,
		Answers:    make([]Answer, 0, count),
	}
}

func NewTranslateRound(levelKey string, dir Direction) *Round {
	var sentences []data.Sentence
	if levelKey == "all" {
		sentences = data.AllSentences()
	} else {
		for _, lvl := range data.Levels {
			if lvl.Key == levelKey {
				sentences = lvl.Sentences
				break
			}
		}
	}

	shuffled := make([]data.Sentence, len(sentences))
	copy(shuffled, sentences)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	count := ChallengesPerRound
	if count > len(shuffled) {
		count = len(shuffled)
	}

	challenges := make([]Challenge, count)
	for i := 0; i < count; i++ {
		s := shuffled[i]
		if dir == SvToEn {
			challenges[i] = Challenge{Prompt: s.Sv, Expected: s.En, Sv: s.Sv, En: s.En}
		} else {
			challenges[i] = Challenge{Prompt: s.En, Expected: s.Sv, Sv: s.Sv, En: s.En}
		}
	}

	return &Round{
		Mode:       ModeTranslate,
		Direction:  dir,
		Category:   levelKey,
		Challenges: challenges,
		Answers:    make([]Answer, 0, count),
	}
}

func NewSpeedRound(dir Direction, hardcore bool) *Round {
	words := data.AllWords()
	shuffled := make([]data.Word, len(words))
	copy(shuffled, words)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	count := SpeedChallengesPool
	if count > len(shuffled) {
		count = len(shuffled)
	}

	challenges := make([]Challenge, count)
	for i := 0; i < count; i++ {
		w := shuffled[i]
		if dir == SvToEn {
			challenges[i] = Challenge{Prompt: w.Sv, Expected: w.En, Sv: w.Sv, En: w.En}
		} else {
			challenges[i] = Challenge{Prompt: w.En, Expected: w.Sv, Sv: w.Sv, En: w.En}
		}
	}

	return &Round{
		Mode:       ModeVocabulary,
		Direction:  dir,
		Category:   "all",
		Challenges: challenges,
		Answers:    make([]Answer, 0, count),
		Timed:      true,
		Hardcore:   hardcore,
	}
}

func (r *Round) CurrentChallenge() Challenge {
	return r.Challenges[r.Current]
}

func (r *Round) Submit(answer string) bool {
	ch := r.Challenges[r.Current]
	correct := checkAnswer(r.Mode, answer, ch.Expected, r.Hardcore)

	r.Answers = append(r.Answers, Answer{
		Sv:      ch.Sv,
		En:      ch.En,
		Given:   answer,
		Correct: correct,
	})

	if correct {
		r.Correct++
	} else {
		r.Wrong++
	}
	r.Current++
	return correct
}

func (r *Round) Done() bool {
	return r.Current >= len(r.Challenges)
}

// Finish ends a timed round early, truncating challenges to what was answered.
func (r *Round) Finish() {
	r.Challenges = r.Challenges[:len(r.Answers)]
	r.Current = len(r.Answers)
}

func (r *Round) Total() int {
	return len(r.Challenges)
}

func checkAnswer(mode Mode, given, expected string, hardcore bool) bool {
	switch mode {
	case ModeTyping:
		return given == expected
	case ModeVocabulary:
		g := normalize(given)
		if hardcore {
			return g == normalize(expected)
		}
		for _, alt := range expandAlternatives(expected) {
			if g == alt {
				return true
			}
		}
		return false
	case ModeTranslate:
		g := normalize(given)
		if hardcore {
			return g == normalize(expected)
		}
		alts := expandAlternatives(expected)
		for _, alt := range alts {
			if g == alt {
				return true
			}
		}
		// Fuzzy matching for sentences: allow minor differences like
		// missing apostrophes, small typos, or dropped articles.
		// Threshold scales with sentence length (min 2, ~10% of length).
		for _, alt := range alts {
			threshold := max(2, len([]rune(alt))/10)
			if levenshtein(g, alt) <= threshold {
				return true
			}
		}
		return false
	}
	return false
}

// levenshtein returns the edit distance between two strings.
func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		prev[j] = j
	}

	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			curr[j] = min(curr[j-1]+1, min(prev[j]+1, prev[j-1]+cost))
		}
		prev, curr = curr, prev
	}
	return prev[lb]
}

func normalize(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.TrimRight(s, ".!?")
	return strings.TrimSpace(s)
}

// expandAlternatives generates all accepted forms of an expected answer.
// For "yes (contradicting negative)" it accepts: the full string, "yes",
// and "contradicting negative".
// For "to go / walk" it accepts: "to go / walk", "to go", "walk".
// For "coffee break / fika" it accepts all slash-separated parts too.
func expandAlternatives(expected string) []string {
	seen := make(map[string]bool)
	add := func(s string) {
		s = normalize(s)
		if s != "" {
			seen[s] = true
		}
	}

	// Full string is always accepted
	add(expected)

	// Strip parenthetical suffix: "yes (contradicting negative)" -> "yes"
	if idx := strings.Index(expected, "("); idx > 0 {
		add(expected[:idx])
		// Also accept the content inside parens
		inner := expected[idx+1:]
		inner = strings.TrimRight(inner, ")")
		add(inner)
	}

	// Strip parens first for slash splitting
	stripped := expected
	if idx := strings.Index(stripped, "("); idx > 0 {
		stripped = strings.TrimSpace(stripped[:idx])
	}

	// Split on " / " for alternatives: "to go / walk" -> ["to go", "walk"]
	for _, src := range []string{expected, stripped} {
		parts := strings.Split(src, " / ")
		for _, p := range parts {
			add(p)
			if idx := strings.Index(p, "("); idx > 0 {
				add(p[:idx])
			}
		}
	}

	// Split on "/" (no spaces) for compact alternatives: "his/her own" -> ["his", "her own"]
	for _, src := range []string{expected, stripped} {
		parts := strings.Split(src, "/")
		if len(parts) > 1 {
			for _, p := range parts {
				add(p)
			}
		}
	}

	result := make([]string, 0, len(seen))
	for k := range seen {
		result = append(result, k)
	}
	return result
}

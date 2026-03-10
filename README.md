# Svenska 🇸🇪

A terminal app for practicing Swedish (A1-A2 level), built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

![Swedish flag in terminal]

## Features

- **Vocabulary** — Translate words between Swedish and English
- **Typing** — Type Swedish words exactly as shown
- **Translate** — Translate full sentences (SV↔EN)
- **AI Help** — Type `?` during any challenge to get an AI-powered explanation
- **Gamification** — XP, levels, daily streaks, and achievements
- **Statistics** — Track your progress, category mastery, and achievements across sessions
- **Flexible matching** — Accepts partial answers, strips parentheticals, splits alternatives on `/`
- 27 word categories and 4 sentence difficulty levels 

## Install

### With Go

```bash
go install github.com/crisecheverria/svenska@latest
```

### From releases (no Go required)

Download the binary for your platform from [Releases](https://github.com/crisecheverria/svenska/releases) and place it somewhere in your `PATH`.

## Usage

```bash
svenska
```

Navigate with `↑↓` or `j/k`, select with `enter`, go back with `b` or `esc`.

During a challenge, type `?` and press enter to get AI help with the current word or sentence.

## AI Help

AI help uses **free models** via [OpenRouter](https://openrouter.ai/) — no API key or account needed. Just type `?` during any challenge.

Models used (with automatic fallback):
- Primary: `arcee-ai/trinity-large-preview:free`
- Fallback: `nvidia/nemotron-3-nano-30b-a3b:free`

For higher rate limits, you can optionally set an OpenRouter API key (Recommended):

```bash
export OPENROUTER_API_KEY="sk-or-..."
```

Add it to your shell config (`~/.zshrc` or `~/.bashrc`) to persist across sessions.

## Data

- Greetings, pronouns, numbers, time, countries, professions, family, food, shopping, clothing, colors, home, furniture, body, health, weather, animals, verbs, adjectives, adverbs, prepositions, education, transport, leisure, everyday words, and conjunctions
- 4 sentence levels: Beginner (A1), Elementary (A1-A2), Intermediate (A2), Advanced (A2+)

## Gamification

- **XP & Levels** — Earn 10 XP per correct answer, bonus for streaks and perfect rounds
- **Daily Streak** — Play every day to build your streak
- **Achievements** — Unlock badges like Första stegen, Perfekt!, Polyglott, Eldsjäl, and Hundra
- **Category Mastery** — Track accuracy per category

## Stats

Your progress is saved to `~/.local/share/svenska/stats.json` and persists across sessions.

# Svenska 🇸🇪

A terminal app for practicing Swedish (A1-A2 level), built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

![Swedish flag in terminal]

## Features

- **Vocabulary** — Translate words between Swedish and English
- **Typing** — Type Swedish words exactly as shown
- **Translate** — Translate full sentences (SV↔EN)
- **AI Help** — Type `?` during any challenge to get an AI-powered explanation
- **Statistics** — Track your progress across sessions
- 27 word categories and 4 sentence difficulty levels from Rivstart A1+A2

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

## AI Help Setup

The AI help feature requires an API key from either Anthropic (Claude) or OpenAI. Set one of these environment variables:

### Option 1: Anthropic (Claude)

Get your API key from [console.anthropic.com](https://console.anthropic.com/)

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
```

### Option 2: OpenAI

Get your API key from [platform.openai.com](https://platform.openai.com/)

```bash
export OPENAI_API_KEY="sk-..."
```

### Make it permanent

Add the export line to your shell config so it persists across sessions:

```bash
# For zsh (~/.zshrc)
echo 'export ANTHROPIC_API_KEY="sk-ant-..."' >> ~/.zshrc

# For bash (~/.bashrc)
echo 'export ANTHROPIC_API_KEY="sk-ant-..."' >> ~/.bashrc
```

> **Note:** AI help is optional. The app works fully without it — you just won't be able to use the `?` help feature.

## Data

All vocabulary and sentences are sourced from **Rivstart A1+A2**, covering:

- Greetings, pronouns, numbers, time, countries, professions, family, food, shopping, clothing, colors, home, furniture, body, health, weather, animals, verbs, adjectives, adverbs, prepositions, education, transport, leisure, everyday words, and conjunctions
- 4 sentence levels: Beginner (A1), Elementary (A1-A2), Intermediate (A2), Advanced (A2+)

## Stats

Your progress is saved to `~/.local/share/svenska/stats.json` and persists across sessions.

package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type HelpResult struct {
	Text string
	Err  error
}

func GetHelp(sv, en, context string) HelpResult {
	prompt := fmt.Sprintf(
		"You are a Swedish language tutor for A1-A2 beginners.\n\n"+
			"The student needs help with:\n"+
			"  Swedish: %s\n"+
			"  English: %s\n"+
			"  Context: %s\n\n"+
			"Please provide:\n"+
			"1. The English translation\n"+
			"2. A word-by-word breakdown\n"+
			"3. Brief grammar notes if relevant\n\n"+
			"Keep it brief (max 8 lines). Use simple language.",
		sv, en, context,
	)

	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		return callAnthropic(key, prompt)
	}
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		return callOpenAI(key, prompt)
	}
	return HelpResult{Err: fmt.Errorf("set ANTHROPIC_API_KEY or OPENAI_API_KEY to enable AI help")}
}

func callAnthropic(key, prompt string) HelpResult {
	body := map[string]any{
		"model":      "claude-sonnet-4-20250514",
		"max_tokens": 400,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", key)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return HelpResult{Err: fmt.Errorf("API request failed: %w", err)}
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return HelpResult{Err: fmt.Errorf("API error %d: %s", resp.StatusCode, string(raw))}
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	json.Unmarshal(raw, &result)
	if len(result.Content) > 0 {
		return HelpResult{Text: result.Content[0].Text}
	}
	return HelpResult{Err: fmt.Errorf("empty response")}
}

func callOpenAI(key, prompt string) HelpResult {
	body := map[string]any{
		"model":      "gpt-4o-mini",
		"max_tokens": 400,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return HelpResult{Err: fmt.Errorf("API request failed: %w", err)}
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return HelpResult{Err: fmt.Errorf("API error %d: %s", resp.StatusCode, string(raw))}
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	json.Unmarshal(raw, &result)
	if len(result.Choices) > 0 {
		return HelpResult{Text: result.Choices[0].Message.Content}
	}
	return HelpResult{Err: fmt.Errorf("empty response")}
}

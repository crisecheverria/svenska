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

	// OpenRouter free models — try primary, fall back to secondary
	models := []string{
		"arcee-ai/trinity-large-preview:free",
		"nvidia/nemotron-3-nano-30b-a3b:free",
	}
	var lastErr error
	for _, model := range models {
		res := callOpenRouter(prompt, model)
		if res.Err == nil {
			return res
		}
		lastErr = res.Err
	}
	return HelpResult{Err: lastErr}
}

func callOpenRouter(prompt, model string) HelpResult {
	body := map[string]any{
		"model":      model,
		"max_tokens": 400,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/crisecheverria/svenska")
	req.Header.Set("X-Title", "Svenska CLI")

	// Optional: use API key if provided for higher rate limits
	if key := os.Getenv("OPENROUTER_API_KEY"); key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}

	client := &http.Client{Timeout: 20 * time.Second}
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

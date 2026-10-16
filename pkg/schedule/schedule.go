// Package schedule provides functionality to fetch and parse official schedules
// for the idol group "Boku ga Mitakatta Aozora" (僕が見たかった青空).
package schedule

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Event represents an official event or schedule entry.
type Event struct {
	Date     string `json:"date"`
	Week     string `json:"week"`
	Time     string `json:"time,omitempty"`
	Category string `json:"category"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	IsLive   bool   `json:"is_live"`
}

// Client handles fetching schedules from the official site or a mock server.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Client with standard defaults.
func NewClient() *Client {
	return &Client{
		BaseURL: "https://bokuao.com",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

var (
	reDate = regexp.MustCompile(`<strong>([0-9\.]+)</strong>\s*<span class="week">([A-Z]+)</span>`)
	reCat  = regexp.MustCompile(`<p class="category[^"]*">([^<]+)</p>`)
	reTime = regexp.MustCompile(`<p class="time">([^<]+(?:<span>.*?</span>)?)</p>`)
	reTit  = regexp.MustCompile(`<p class="tit">([^<]+)</p>`)
	reHref = regexp.MustCompile(`<a href="([^"]+)"`)
	reTag  = regexp.MustCompile(`<[^>]+>`)
)

// Non-live keywords used to filter out contact and media appearances.
var nonLiveKeywords = []string{
	"お話し会", "おはなし会", "サイン会", "握手会", "SASUKE", "舞台挨拶", "写真会", "チェキ", "誕生日",
}

// CheckIsLive determines whether an event represents a live musical/stage performance.
func CheckIsLive(title, category string) bool {
	upperCat := strings.ToUpper(category)
	if !strings.Contains(upperCat, "LIVE") && !strings.Contains(upperCat, "EVENT") {
		return false
	}

	for _, kw := range nonLiveKeywords {
		if strings.Contains(title, kw) {
			// If it includes live/performance explicitly alongside contact (e.g. "フリーライブ&握手会"), treat as live.
			if strings.Contains(title, "フリーライブ") ||
				strings.Contains(title, "ミニライブ") ||
				strings.Contains(title, "単独") ||
				strings.Contains(title, "公演") ||
				strings.Contains(title, "フェス") {
				return true
			}
			return false
		}
	}
	return true
}

// Fetch retrieves the schedule for the specified year and month.
func (c *Client) Fetch(ctx context.Context, year, month int) ([]Event, error) {
	url := fmt.Sprintf("%s/schedule/list/%d/%d/", strings.TrimSuffix(c.BaseURL, "/"), year, month)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "bokuao-cli/1.0 (+https://github.com/AobaIwaki123/bokuao)")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return ParseHTML(string(body), c.BaseURL)
}

// ParseHTML extracts Events from raw HTML schedule list content.
func ParseHTML(rawHTML, baseURL string) ([]Event, error) {
	blocks := strings.Split(rawHTML, "<li data-delighter>")
	if len(blocks) <= 1 {
		return nil, nil
	}

	var events []Event
	for _, b := range blocks[1:] {
		var ev Event
		if m := reDate.FindStringSubmatch(b); len(m) > 2 {
			ev.Date = m[1]
			ev.Week = m[2]
		}
		if m := reCat.FindStringSubmatch(b); len(m) > 1 {
			ev.Category = html.UnescapeString(strings.TrimSpace(m[1]))
		}
		if m := reTime.FindStringSubmatch(b); len(m) > 1 {
			cleanTime := reTag.ReplaceAllString(m[1], "")
			ev.Time = html.UnescapeString(strings.TrimSpace(cleanTime))
		}
		if m := reTit.FindStringSubmatch(b); len(m) > 1 {
			ev.Title = html.UnescapeString(strings.TrimSpace(m[1]))
		}
		if m := reHref.FindStringSubmatch(b); len(m) > 1 {
			path := m[1]
			if strings.HasPrefix(path, "http") {
				ev.URL = path
			} else {
				ev.URL = strings.TrimSuffix(baseURL, "/") + path
			}
		}

		if ev.Title != "" {
			ev.IsLive = CheckIsLive(ev.Title, ev.Category)
			events = append(events, ev)
		}
	}

	return events, nil
}

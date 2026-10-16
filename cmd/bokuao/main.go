// Package main provides the command-line interface for the bokuao tool.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/AobaIwaki123/bokuao/pkg/schedule"
)

func main() {
	now := time.Now()
	yearFlag := flag.Int("year", now.Year(), "Year to fetch (default current year)")
	monthFlag := flag.Int("month", int(now.Month()), "Month to fetch (default current month)")
	liveOnly := flag.Bool("live-only", false, "Filter to live performance and concert events only")
	jsonOutput := flag.Bool("json", false, "Output results in JSON format")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := schedule.NewClient()
	events, err := client.Fetch(ctx, *yearFlag, *monthFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var filtered []schedule.Event
	for _, ev := range events {
		if *liveOnly && !ev.IsLive {
			continue
		}
		filtered = append(filtered, ev)
	}

	if *jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(filtered); err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
			os.Exit(1)
		}
		return
	}

	modeStr := "全予定"
	if *liveOnly {
		modeStr = "ライブ・公演のみ"
	}

	fmt.Printf("=== 僕が見たかった青空 スケジュール (%d年%d月 / %s: %d件) ===\n\n", *yearFlag, *monthFlag, modeStr, len(filtered))
	if len(filtered) == 0 {
		fmt.Println("該当するイベントはありません。")
		return
	}

	for _, ev := range filtered {
		badge := ""
		if ev.IsLive {
			badge = " [LIVE]"
		}
		fmt.Printf("[%s (%s)] %s%s\n", ev.Date, ev.Week, ev.Title, badge)
		if ev.Time != "" && ev.Time != "00:00〜23:59" {
			fmt.Printf("  時間: %s\n", ev.Time)
		}
		fmt.Printf("  URL:  %s\n\n", ev.URL)
	}
}

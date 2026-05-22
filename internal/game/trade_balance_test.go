package game

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestTrade_RejectsWhenFoodAtOrBelowThreshold(t *testing.T) {
	s := newTestState()
	s.Food = MinFoodToTrade
	s.Credits = 0
	s.Morale = 50

	s.Trade()

	if s.Credits != 0 {
		t.Fatalf("Credits = %d, want unchanged 0", s.Credits)
	}
	if s.Food != MinFoodToTrade {
		t.Fatalf("Food = %d, want unchanged %d", s.Food, MinFoodToTrade)
	}
	if len(s.Log) == 0 || !containsAll(s.Log[0], "Trade refused", "food") {
		t.Fatalf("log = %q", s.Log[0])
	}
}

func TestTrade_AllowsWhenFoodAboveThreshold(t *testing.T) {
	s := newTestState()
	s.Food = MinFoodToTrade + 1
	s.Credits = 0
	s.Morale = 50

	s.Trade()

	if s.Credits != TradeCreditsGain {
		t.Fatalf("Credits = %d, want %d", s.Credits, TradeCreditsGain)
	}
	wantFood := MinFoodToTrade + 1 - TradeFoodCost
	if s.Food != wantFood {
		t.Fatalf("Food = %d, want %d", s.Food, wantFood)
	}
}

func TestTrade_RejectedTradeRecordsReasonInSessionLog(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "trade.jsonl")
	logger, err := OpenSessionLog(path)
	if err != nil {
		t.Fatalf("OpenSessionLog: %v", err)
	}

	s := newTestState()
	s.SessionLog = logger
	s.LogSessionStart()
	s.Food = 20
	s.Trade()
	_ = logger.Close()

	lines := readLogLines(t, path)
	var trade LogEntry
	for _, line := range lines {
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("parse: %v", err)
		}
		if entry.Type == "trade" {
			trade = entry
		}
	}
	if trade.Type != "trade" {
		t.Fatal("missing trade log entry")
	}
	if trade.Detail["ok"] != false || trade.Detail["reason"] != "low_food" {
		t.Fatalf("trade detail: %+v", trade.Detail)
	}
}

func TestReplaySession_TradeRejectedMatchesLog(t *testing.T) {
	content := testContentWithEvents()
	start := Snapshot{Day: 1, Power: 65, Food: 20, Morale: 70, Credits: 0, Population: 8, PopulationCap: 10, BeaconParts: 0, MaxBeacon: 5}

	entries := []LogEntry{
		{Type: "session_start", Snapshot: &start, Detail: map[string]any{"seed": "7"}},
		{
			Type:   "trade",
			Day:    1,
			Detail: map[string]any{"ok": false, "reason": "low_food"},
			Before: &start,
			After:  &start,
		},
	}

	final, err := ReplaySession(content, entries)
	if err != nil {
		t.Fatalf("ReplaySession: %v", err)
	}
	if final.Food != 20 || final.Credits != 0 {
		t.Fatalf("final stats changed: food=%d credits=%d", final.Food, final.Credits)
	}
}

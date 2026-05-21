package game

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenSessionLog_WritesSessionStartAndActionsAsJSONL(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.jsonl")

	logger, err := OpenSessionLog(path)
	if err != nil {
		t.Fatalf("OpenSessionLog: %v", err)
	}
	t.Cleanup(func() { _ = logger.Close() })

	start := Snapshot{Day: 1, Power: 65, Credits: 180}
	if err := logger.Record("session_start", 1, nil, start, start); err != nil {
		t.Fatalf("Record session_start: %v", err)
	}
	before := start
	after := Snapshot{Day: 1, Power: 83, Credits: 110, Buildings: map[string]int{"solar_array": 1}}
	if err := logger.Record("build", 1, map[string]any{"building_id": "solar_array", "ok": true}, before, after); err != nil {
		t.Fatalf("Record build: %v", err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d log lines, want 2:\n%s", len(lines), raw)
	}

	var first LogEntry
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatalf("parse line 1: %v", err)
	}
	if first.Type != "session_start" || first.SessionID == "" || first.Snapshot.Power != 65 {
		t.Fatalf("session_start entry: %+v", first)
	}

	var second LogEntry
	if err := json.Unmarshal([]byte(lines[1]), &second); err != nil {
		t.Fatalf("parse line 2: %v", err)
	}
	if second.Type != "build" || second.Before.Credits != 180 || second.After.Credits != 110 {
		t.Fatalf("build entry: %+v", second)
	}
	if second.Detail["building_id"] != "solar_array" {
		t.Fatalf("build detail: %+v", second.Detail)
	}
}

func TestState_BuildRecordsToSessionLog(t *testing.T) {
	dir := t.TempDir()
	logger, err := OpenSessionLog(filepath.Join(dir, "play.jsonl"))
	if err != nil {
		t.Fatalf("OpenSessionLog: %v", err)
	}

	s := newTestState()
	s.SessionLog = logger
	s.LogSessionStart()

	s.Build("solar_array")
	_ = logger.Close()

	lines := readLogLines(t, filepath.Join(dir, "play.jsonl"))
	if len(lines) < 2 {
		t.Fatalf("expected session_start + build, got %d lines", len(lines))
	}
	var build LogEntry
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &build); err != nil {
		t.Fatalf("parse build: %v", err)
	}
	if build.Type != "build" || build.After.Buildings["solar_array"] != 1 {
		t.Fatalf("build entry: %+v", build)
	}
}

func TestWorkOnBeacon_LogsGameEnd(t *testing.T) {
	dir := t.TempDir()
	logger, err := OpenSessionLog(filepath.Join(dir, "win.jsonl"))
	if err != nil {
		t.Fatalf("OpenSessionLog: %v", err)
	}

	s := newTestState()
	s.SessionLog = logger
	s.LogSessionStart()
	s.BeaconParts = 4
	s.Power = 50
	s.Credits = 100

	s.WorkOnBeacon()
	_ = logger.Close()

	lines := readLogLines(t, filepath.Join(dir, "win.jsonl"))
	var gameEnd LogEntry
	for _, line := range lines {
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("parse: %v", err)
		}
		if entry.Type == "game_end" {
			gameEnd = entry
		}
	}
	if gameEnd.Type != "game_end" || gameEnd.After == nil || !gameEnd.After.Won {
		t.Fatalf("game_end entry: %+v\nlines:\n%s", gameEnd, strings.Join(lines, "\n"))
	}
}

func readLogLines(t *testing.T, path string) []string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	return strings.Split(strings.TrimSpace(string(raw)), "\n")
}

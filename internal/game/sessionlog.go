package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Snapshot is a machine-readable colony state for session analysis.
type Snapshot struct {
	Day           int            `json:"day"`
	Power         int            `json:"power"`
	Food          int            `json:"food"`
	Morale        int            `json:"morale"`
	Credits       int            `json:"credits"`
	Population    int            `json:"population"`
	PopulationCap int            `json:"population_cap"`
	BeaconParts   int            `json:"beacon_parts"`
	MaxBeacon     int            `json:"max_beacon_parts"`
	GameOver      bool           `json:"game_over"`
	Won           bool           `json:"won"`
	Buildings     map[string]int  `json:"buildings,omitempty"`
	Damaged       map[string]bool `json:"damaged,omitempty"`
}

// LogEntry is one JSONL record for a play session.
type LogEntry struct {
	TS        string         `json:"ts"`
	SessionID string         `json:"session_id"`
	Type      string         `json:"type"`
	Day       int            `json:"day"`
	Log       string         `json:"log_path,omitempty"`
	Snapshot  *Snapshot      `json:"snapshot,omitempty"`
	Before    *Snapshot      `json:"before,omitempty"`
	After     *Snapshot      `json:"after,omitempty"`
	Detail    map[string]any `json:"detail,omitempty"`
}

// Recorder is the session-logging seam. State holds a Recorder; SessionLogger implements it.
// Pass nil for a no-op (tests, headless runs).
type Recorder interface {
	Record(typ string, day int, detail map[string]any, before, after Snapshot) error
	Close() error
}

// SessionLogger appends structured session events to a JSONL file.
type SessionLogger struct {
	Path      string
	sessionID string
	file      *os.File
	enc       *json.Encoder
}

// OpenSessionLog creates (or truncates) a JSONL session log at path.
func OpenSessionLog(path string) (*SessionLogger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("mkdir log dir: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open session log: %w", err)
	}
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	return &SessionLogger{
		Path:      path,
		sessionID: id,
		file:      f,
		enc:       json.NewEncoder(f),
	}, nil
}

// DefaultSessionLogPath returns a new log file under the user cache dir.
func DefaultSessionLogPath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	base := filepath.Join(dir, "outpost-404", "sessions")
	name := fmt.Sprintf("session-%s.jsonl", time.Now().Format("20060102-150405"))
	return filepath.Join(base, name), nil
}

// Record appends one JSONL event. Pass the same Snapshot for before/after when unchanged.
func (l *SessionLogger) Record(typ string, day int, detail map[string]any, before, after Snapshot) error {
	entry := LogEntry{
		TS:        time.Now().UTC().Format(time.RFC3339Nano),
		SessionID: l.sessionID,
		Type:      typ,
		Day:       day,
		Log:       l.Path,
		Detail:    detail,
	}
	if typ == "session_start" {
		s := before
		entry.Snapshot = &s
	} else {
		b, a := before, after
		entry.Before = &b
		entry.After = &a
	}
	if err := l.enc.Encode(entry); err != nil {
		return fmt.Errorf("encode log entry: %w", err)
	}
	return l.file.Sync()
}

// Close flushes and closes the log file.
func (l *SessionLogger) Close() error {
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

func (s *State) snapshot() Snapshot {
	buildings := make(map[string]int, len(s.Buildings))
	damaged := make(map[string]bool)
	for id, b := range s.Buildings {
		buildings[id] = b.Level
		if b.Damaged {
			damaged[id] = true
		}
	}
	return Snapshot{
		Day:           s.Day,
		Power:         s.Power,
		Food:          s.Food,
		Morale:        s.Morale,
		Credits:       s.Credits,
		Population:    s.Population,
		PopulationCap: s.PopulationCap,
		BeaconParts:   s.BeaconParts,
		MaxBeacon:     s.MaxBeaconParts,
		GameOver:      s.GameOver,
		Won:           s.Won,
		Buildings:     buildings,
		Damaged:       damaged,
	}
}

// LogSessionStart writes the opening snapshot for analysis tools.
func (s *State) LogSessionStart() {
	if s.SessionLog == nil {
		return
	}
	snap := s.snapshot()
	_ = s.SessionLog.Record("session_start", s.Day, map[string]any{
		"seed":       strconv.FormatInt(s.Seed, 10),
		"scenario":   s.ScenarioID,
		"difficulty": s.DifficultyID,
	}, snap, snap)
}

func (s *State) recordAction(typ string, detail map[string]any, before, after Snapshot) {
	if s.SessionLog == nil {
		return
	}
	_ = s.SessionLog.Record(typ, before.Day, detail, before, after)
}

// EndSession closes the session log file.
func (s *State) EndSession() {
	if s.SessionLog != nil {
		_ = s.SessionLog.Close()
		s.SessionLog = nil
	}
}

// AttachSessionLog opens path and binds it to state. Empty path uses DefaultSessionLogPath.
func AttachSessionLog(s *State, path string) (*SessionLogger, error) {
	if path == "" {
		var err error
		path, err = DefaultSessionLogPath()
		if err != nil {
			return nil, err
		}
	}
	logger, err := OpenSessionLog(path)
	if err != nil {
		return nil, err
	}
	s.SessionLog = logger
	s.LogSessionStart()
	return logger, nil
}

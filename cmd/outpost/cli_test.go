package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ericdahl/outpost-404/internal/game"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("repo root not found")
		}
		dir = parent
	}
}

func buildOutpost(t *testing.T, root string) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "outpost")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/outpost")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build outpost: %v\n%s", err, out)
	}
	return bin
}

type cliResult struct {
	stdout string
	stderr string
	err    error
}

func runCLI(t *testing.T, root, bin string, args ...string) cliResult {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = root
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return cliResult{stdout: stdout.String(), stderr: stderr.String(), err: err}
}

// recordSessionLog writes a minimal JSONL session using the same content source as the CLI (./data when present).
func recordSessionLog(t *testing.T, root string, seed int64) string {
	t.Helper()
	content, err := game.LoadContent(filepath.Join(root, "data"))
	if err != nil {
		t.Fatalf("LoadContent: %v", err)
	}
	path := filepath.Join(t.TempDir(), "session.jsonl")
	logger, err := game.OpenSessionLog(path)
	if err != nil {
		t.Fatalf("OpenSessionLog: %v", err)
	}
	s := game.NewStateWithSeed(content, seed)
	s.SessionLog = logger
	s.LogSessionStart()
	s.Build("solar_array")
	s.NextDay()
	if err := logger.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if _, err := game.ReplaySession(content, mustLoadLog(t, path)); err != nil {
		t.Fatalf("sanity replay before CLI: %v", err)
	}
	return path
}

func mustLoadLog(t *testing.T, path string) []game.LogEntry {
	t.Helper()
	entries, err := game.LoadSessionLog(path)
	if err != nil {
		t.Fatalf("LoadSessionLog: %v", err)
	}
	return entries
}

func TestSimulateCLI(t *testing.T) {
	root := repoRoot(t)
	res := runCLI(t, root, buildOutpost(t, root), "-simulate", "scripts/conservative.json", "-seed", "1")
	if res.err != nil {
		t.Fatalf("simulate exit: %v\nstderr: %s", res.err, res.stderr)
	}
	for _, want := range []string{"seed=", "day=", "won="} {
		if !strings.Contains(res.stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, res.stdout)
		}
	}
}

func TestSimulateCLI_badScript(t *testing.T) {
	root := repoRoot(t)
	res := runCLI(t, root, buildOutpost(t, root), "-simulate", "scripts/no-such.json", "-seed", "1")
	if res.err == nil {
		t.Fatal("expected non-zero exit for missing script")
	}
}

func TestSimulateCLI_seedsSweep(t *testing.T) {
	root := repoRoot(t)
	res := runCLI(t, root, buildOutpost(t, root),
		"-simulate", "scripts/conservative.json",
		"-seeds", "1,7",
	)
	if res.err != nil {
		t.Fatalf("sweep exit: %v\nstderr: %s\nstdout: %s", res.err, res.stderr, res.stdout)
	}
	if !strings.Contains(res.stderr, "sweep: ") {
		t.Fatalf("stderr missing sweep summary:\n%s", res.stderr)
	}
	for _, want := range []string{"seed=1", "seed=7"} {
		if !strings.Contains(res.stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, res.stdout)
		}
	}
}

func TestReplayCLI_recordedSession(t *testing.T) {
	root := repoRoot(t)
	logPath := recordSessionLog(t, root, 4242)
	res := runCLI(t, root, buildOutpost(t, root), "-replay", logPath)
	if res.err != nil {
		t.Fatalf("replay exit: %v\nstderr: %s", res.err, res.stderr)
	}
	for _, want := range []string{"replay ok:", "day=", "won="} {
		if !strings.Contains(res.stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, res.stdout)
		}
	}
}

func TestReplayCLI_missingFile(t *testing.T) {
	root := repoRoot(t)
	res := runCLI(t, root, buildOutpost(t, root), "-replay", filepath.Join(t.TempDir(), "missing.jsonl"))
	if res.err == nil {
		t.Fatal("expected non-zero exit for missing replay file")
	}
	if !strings.Contains(res.stderr, "replay failed:") {
		t.Fatalf("stderr = %q", res.stderr)
	}
}

func TestReplayCLI_invalidJSONL(t *testing.T) {
	root := repoRoot(t)
	path := filepath.Join(t.TempDir(), "bad.jsonl")
	if err := os.WriteFile(path, []byte("not json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	res := runCLI(t, root, buildOutpost(t, root), "-replay", path)
	if res.err == nil {
		t.Fatal("expected non-zero exit for invalid JSONL")
	}
	if !strings.Contains(res.stderr, "replay failed:") {
		t.Fatalf("stderr = %q", res.stderr)
	}
}

func TestVersionCLI(t *testing.T) {
	root := repoRoot(t)
	bin := buildOutpost(t, root)
	wantVersion := "test-version-54"
	binTagged := filepath.Join(t.TempDir(), "outpost-versioned")
	build := exec.Command("go", "build", "-o", binTagged, "-ldflags", "-X main.version="+wantVersion, "./cmd/outpost")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build versioned: %v\n%s", err, out)
	}

	res := runCLI(t, root, binTagged, "-version")
	if res.err != nil {
		t.Fatalf("version exit: %v\nstderr: %s", res.err, res.stderr)
	}
	if strings.TrimSpace(res.stdout) != wantVersion {
		t.Fatalf("stdout = %q, want %q", res.stdout, wantVersion)
	}

	// default dev build still prints something
	resDev := runCLI(t, root, bin, "-version")
	if resDev.err != nil || strings.TrimSpace(resDev.stdout) == "" {
		t.Fatalf("dev version: err=%v stdout=%q", resDev.err, resDev.stdout)
	}
}

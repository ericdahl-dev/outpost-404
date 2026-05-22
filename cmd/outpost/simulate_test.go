package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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

func TestSimulateCLI(t *testing.T) {
	root := repoRoot(t)
	bin := buildOutpost(t, root)

	cmd := exec.Command(bin, "-simulate", "scripts/conservative.json", "-seed", "1")
	cmd.Dir = root
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("simulate exit: %v\nstderr: %s", err, stderr.Bytes())
	}

	out := stdout.String()
	for _, want := range []string{"seed=", "day=", "won="} {
		if !strings.Contains(out, want) {
			t.Fatalf("stdout missing %q:\n%s", want, out)
		}
	}
}

func TestSimulateCLI_badScript(t *testing.T) {
	root := repoRoot(t)
	bin := buildOutpost(t, root)

	cmd := exec.Command(bin, "-simulate", "scripts/no-such.json", "-seed", "1")
	cmd.Dir = root
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit for missing script")
	}
}

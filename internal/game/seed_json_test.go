package game

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

func TestSeedFromDetail_stringRoundTrip(t *testing.T) {
	want := int64(1779403310247544000)
	got, err := seedFromDetail(map[string]any{"seed": strconv.FormatInt(want, 10)})
	if err != nil {
		t.Fatalf("seedFromDetail: %v", err)
	}
	if got != want {
		t.Fatalf("seed = %d, want %d", got, want)
	}
}

func TestSeedFromDetail_jsonNumber(t *testing.T) {
	got, err := seedFromDetail(map[string]any{"seed": json.Number("4242")})
	if err != nil {
		t.Fatalf("seedFromDetail: %v", err)
	}
	if got != 4242 {
		t.Fatalf("seed = %d, want 4242", got)
	}
}

func TestSeedFromDetail_missingSeed(t *testing.T) {
	_, err := seedFromDetail(nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "missing seed") {
		t.Fatalf("error = %q", err)
	}
}

func TestSeedFromDetail_invalidJsonNumber(t *testing.T) {
	_, err := seedFromDetail(map[string]any{"seed": json.Number("not-a-number")})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid seed") {
		t.Fatalf("error = %q", err)
	}
}

func TestSeedFromDetail_float64LosesPrecision(t *testing.T) {
	want := int64(1779403310247544000)
	raw, err := json.Marshal(map[string]any{"seed": want})
	if err != nil {
		t.Fatal(err)
	}
	var detail map[string]any
	if err := json.Unmarshal(raw, &detail); err != nil {
		t.Fatal(err)
	}
	got, err := seedFromDetail(detail)
	if err != nil {
		t.Fatalf("seedFromDetail: %v", err)
	}
	if got == want {
		t.Fatal("expected float64 JSON seed to differ from original int64")
	}
}

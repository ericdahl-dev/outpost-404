package game

import (
	"encoding/json"
	"strconv"
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

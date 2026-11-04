package sdk

import (
	"testing"
	"time"

	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

func TestWholeSecondDuration(t *testing.T) {
	if got, ok := WholeSecondDuration(nil); !ok || got != 0 {
		t.Fatalf("nil = %v, %v", got, ok)
	}
	if got, ok := WholeSecondDuration(durationpb.New(0)); !ok || got != 0 {
		t.Fatalf("zero = %v, %v", got, ok)
	}
	if got, ok := WholeSecondDuration(durationpb.New(5 * time.Second)); !ok || got != 5*time.Second {
		t.Fatalf("whole = %v, %v", got, ok)
	}
	if _, ok := WholeSecondDuration(durationpb.New(1500 * time.Millisecond)); ok {
		t.Fatal("fractional accepted")
	}
}

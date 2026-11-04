package sdk

import (
	"time"

	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

// WholeSecondDuration maps a command timeout for provider exec APIs.
// Nil and explicit zero become 0 (no command limit once the provider applies
// that value). Non-zero values must be whole seconds; ok is false otherwise.
func WholeSecondDuration(value *durationpb.Duration) (time.Duration, bool) {
	timeout := OptionalDuration(value)
	if timeout == nil {
		return 0, true
	}
	if *timeout < 0 || *timeout%time.Second != 0 {
		return 0, false
	}
	return *timeout, true
}

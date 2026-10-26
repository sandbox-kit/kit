package sdk

import (
	"time"

	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

// OptionalDuration converts a protobuf duration to a time.Duration pointer.
func OptionalDuration(value *durationpb.Duration) *time.Duration {
	if value == nil {
		return nil
	}
	converted := value.AsDuration()
	return &converted
}

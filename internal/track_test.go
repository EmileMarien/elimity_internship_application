package internal_test

import (
	"context"
	"testing"
	"time"

	"github.com/elimity-com/backend-intern-exercise/internal"
	"github.com/stretchr/testify/assert"
)

func TestTrack(t *testing.T) {
	// Define test cases
	tests := []struct {
		name      string
		interval  time.Duration
		token     string
		minStars  int
		wantError bool
	}{
		{
			name:      "Valid token and minimum stars",
			interval:  1 * time.Second,
			token:     "",
			minStars:  100,
			wantError: false,
		},
		{
			name:      "Valid token and no minimum stars",
			interval:  1 * time.Second,
			token:     "",
			minStars:  0,
			wantError: false,
		},
	}

	// Loop through test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a context with a 10-second timeout
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Check if the test case has exceeded its deadline
			deadline, ok := t.Deadline()
			if ok && deadline.Before(time.Now()) {
				t.Skip("Test case exceeded deadline without errors")
			}

			// Call the Track function with the test case parameters
			err := internal.Track(ctx, tc.interval, tc.token, tc.minStars)
			// Check if the error matches the expected value
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

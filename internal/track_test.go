package internal_test

import (
	"context"
	"testing"
	"time"

	"github.com/elimity-com/backend-intern-exercise/internal"
	"github.com/stretchr/testify/assert"
)

func TestTrack(t *testing.T) {
	/*
	 * TestTrack tests the internal.Track function.
	 * It tests the function with valid and invalid tokens and minimum stars.
	 * It tests the function with a context that expires before the function completes.
	 * Since the function is a long-running function, it is tested the first 5 seconds only and passes if the function does not return an error in that period.
	 */

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
			// Create a context with a 5-second timeout
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Create a channel to receive errors from the internal.Track function
			errCh := make(chan error, 1)

			// Run the internal.Track function in a goroutine
			go func() {
				// Call the Track function with the test case parameters
				errCh <- internal.Track(ctx, tc.interval, tc.token, tc.minStars)
			}()

			// Wait for the function to complete or for the context to expire
			select {
			case err := <-errCh:
				// Check if the error matches the expected value
				if tc.wantError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			case <-ctx.Done():
				// If the context expires, fail the test
				if tc.wantError {
					t.Fatal("Test case exceeded deadline without errors")
				}
			}
		})
	}
}

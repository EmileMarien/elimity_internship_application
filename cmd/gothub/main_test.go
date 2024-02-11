package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestReadTokenFromFile(t *testing.T) {
	// Create a temporary token file
	tempTokenFile, err := ioutil.TempFile("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempTokenFile.Name())

	// Write token to the temporary file
	token := "mytoken"
	_, err = tempTokenFile.WriteString(token)
	if err != nil {
		t.Fatalf("Failed to write token to temp file: %v", err)
	}
	tempTokenFile.Close()

	// Test reading the token from the file
	readToken, err := readTokenFromFile(tempTokenFile.Name())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if readToken != token {
		t.Errorf("Expected token %s, got %s", token, readToken)
	}
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedResult struct {
			interval   time.Duration
			tokenFile  string
			minStars   int
			expectErr  bool
			errMessage string
		}
	}{
		{
			name: "ValidArgs",
			args: []string{"-interval=15s", "-token-file=path/to/token", "-min-stars=5"},
			expectedResult: struct {
				interval   time.Duration
				tokenFile  string
				minStars   int
				expectErr  bool
				errMessage string
			}{
				interval:  15 * time.Second,
				tokenFile: "path/to/token",
				minStars:  5,
				expectErr: false,
			},
		},
		{
			name: "InvalidInterval",
			args: []string{"-interval=0", "-token-file=path/to/token", "-min-stars=5"},
			expectedResult: struct {
				interval   time.Duration
				tokenFile  string
				minStars   int
				expectErr  bool
				errMessage string
			}{
				expectErr:  true,
				errMessage: "invalid interval: must be greater than zero",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interval, tokenFile, minStars, err := parseArgs(tt.args)
			if tt.expectedResult.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.expectedResult.errMessage {
					t.Errorf("unexpected error message, got: %s, want: %s", err.Error(), tt.expectedResult.errMessage)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if interval != tt.expectedResult.interval {
					t.Errorf("unexpected interval, got: %v, want: %v", interval, tt.expectedResult.interval)
				}
				if tokenFile != tt.expectedResult.tokenFile {
					t.Errorf("unexpected token file, got: %s, want: %s", tokenFile, tt.expectedResult.tokenFile)
				}
				if minStars != tt.expectedResult.minStars {
					t.Errorf("unexpected min stars, got: %d, want: %d", minStars, tt.expectedResult.minStars)
				}
			}
		})
	}
}

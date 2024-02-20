package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "ValidArgs",
			args:    []string{"main.go", "track", "-interval=30s", "-min-stars=5"},
			wantErr: false,
		},
		{
			name:    "InvalidInterval",
			args:    []string{"main.go", "track", "-interval=0s", "-min-stars=5"},
			wantErr: true,
		},
		{
			name:    "InvalidMinStars",
			args:    []string{"track", "-interval=30s", "-min-stars=-1"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, _, _, err := parseArgs(tt.args)
			fmt.Printf("%v", err)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestReadTokenFromFile(t *testing.T) {
	// Create a temporary file
	file, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write the token to the file
	token := "my-github-token"
	if _, err := file.WriteString(token); err != nil {
		t.Fatalf("failed to write token to file: %v", err)
	}

	// Close the file
	if err := file.Close(); err != nil {
		t.Fatalf("failed to close file: %v", err)
	}

	// Read the token from the file
	readToken, err := readTokenFromFile(file.Name())
	if err != nil {
		t.Fatalf("readTokenFromFile() failed: %v", err)
	}

	// Check if the read token matches the expected token
	if readToken != token {
		t.Errorf("read token does not match expected token: read=%s, expected=%s", readToken, token)
	}
}

// Define a mock implementation of the internal.Track function
func mockTrack(ctx context.Context, interval time.Duration, token string, minStars int) error {
	fmt.Println("Mock track function called")
	return nil
}

/*
// Testcode only works with implementation of run(args string[]) function which differs from the source code

func TestRun(t *testing.T) {
	// Save original os.Args and defer resetting it
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Define test cases
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "HelpCommand",
			args:    []string{"", "", "help"},
			wantErr: false,
		},
		{
			name:    "TrackCommand",
			args:    []string{"", "", "track", "-interval=30s", "-min-stars=5"},
			wantErr: false,
		},
		{
			name:    "InvalidCommand",
			args:    []string{"", "", "invalid"},
			wantErr: false,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Create a pipe to capture stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("failed to create pipe: %v", err)
			}
			defer r.Close()
			defer w.Close()

			// Replace os.Stdout with the write end of the pipe
			oldStdout := os.Stdout
			defer func() { os.Stdout = oldStdout }()
			os.Stdout = w

			// Run the function
			err = run(tt.args)
			fmt.Printf("%v", err)
			// Close the write end of the pipe
			w.Close()

			// Read the captured output from the read end of the pipe
			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			if err != nil {
				t.Fatalf("failed to read captured output: %v", err)
			}

			// Check if the output is as expected
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
*/

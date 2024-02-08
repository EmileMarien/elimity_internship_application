package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/elimity-com/backend-intern-exercise/internal"
)

var (
	args          = os.Args
	tokenFilePath string
	useToken      bool
	minStars      int
)

var name = makeName()

func parseArgs() (time.Duration, bool, string, int, error) {
	// Create a new custom FlagSet
	set := flag.NewFlagSet("", flag.ContinueOnError)

	// Define flag variables
	var interval time.Duration
	var useToken bool
	var tokenFilePath string
	var minStars int

	// Add flags to the custom FlagSet
	set.DurationVar(&interval, "interval", 10*time.Second, "repository update interval")
	set.BoolVar(&useToken, "use-token", false, "set to true to use GitHub personal access token")
	set.StringVar(&tokenFilePath, "token", "", "GitHub personal access token")
	set.IntVar(&minStars, "min-stars", 0, "minimum stars")

	// Parse the flags
	if err := set.Parse(os.Args[1:]); err != nil {
		return 0, false, "", 0, fmt.Errorf("failed parsing flags: %v", err)
	}

	// Check if interval is a valid duration
	if interval <= 0 {
		return 0, false, "", 0, errors.New("invalid interval: must be greater than zero")
	}

	// Return parsed values
	return interval, useToken, tokenFilePath, minStars, nil
}

func log(message string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", name, message)
	os.Exit(1)
}

func main() {
	// Set up context and other variables as needed
	ctx := context.Background()

	// Parse command-line arguments
	interval, useToken, tokenFilePath, minStars, err := parseArgs()
	if err != nil {
		log(fmt.Sprintf("Error parsing arguments: %v", err))
		return
	}
	// Read token from file if needed
	var token string
	if useToken {
		token, err = readTokenFromFile(tokenFilePath)
		if err != nil {
			log(fmt.Sprintf("Error reading token from file: %v", err))
			return
		}
	} else {
		token = ""
	}

	// Run the tracking function in a separate goroutine
	go func() {
		if err := run(ctx, interval, useToken, token, minStars); err != nil {
			log(fmt.Sprintf("Error: %v", err))
		}
	}()

	// Print table output
	printTableOutput(ctx)
}

func readTokenFromFile(filePath string) (string, error) {
	token, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	// Trim any leading/trailing white spaces and newlines
	return strings.TrimSpace(string(token)), nil
}

func makeName() string {
	path := args[0]
	return filepath.Base(path)
}

func printTableOutput(ctx context.Context) {
	// Print table headers
	fmt.Println("Owner\t| Name\t| Updated at (UTC)\t| Star count")
	// Loop until context is cancelled
	for {
		select {
		case <-ctx.Done():
			// Exit loop if context is cancelled
			return
		default:
			// Fetch repository details and print in table format
			// You can implement this part based on your requirement and the Track function output
			// For demonstration, I'm just sleeping here
			time.Sleep(5 * time.Second)
		}
	}
}

func run(ctx context.Context, interval time.Duration, useToken bool, token string, minStars int) error {
	if nbArgs := len(args); nbArgs < 2 {
		return usageError{message: "missing command"}
	}
	switch args[1] {
	case "help":
		const usage = `
Simple CLI for tracking public GitHub repositories.

Usage:
  %[1]s help
  %[1]s track [-interval=<interval>] [-use-token=<yes|no>]

Commands:
  help  Show usage information
  track Track public GitHub repositories

Options:
  -interval=<interval> Repository update interval, greater than zero [default: 10s]
  -use-token=<yes|no>   Set to 'yes' to use GitHub personal access token, 'no' to run without token [default: no]
`
		fmt.Fprintf(os.Stdout, usage, name)
		return nil

	case "track":
		if useToken {
			if token == "" {
				return errors.New("GitHub personal access token is required for tracking")
			}
		}

		if err := internal.Track(ctx, interval, useToken, tokenFilePath, minStars); err != nil {
			return fmt.Errorf("failed tracking: %v", err)
		}
		return nil

	default:
		return usageError{message: "got invalid command"}
	}
}

type usageError struct {
	message string
}

func (e usageError) Error() string {
	return e.message
}

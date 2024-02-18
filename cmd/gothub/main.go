package main

///Users/emile/Documents/github_token_EmileMarien.txt
import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/elimity-com/backend-intern-exercise/internal"
)

var (
	args = os.Args
)

var name = makeName()

func parseArgs(args []string) (time.Duration, string, int, error) {
	/*
	 * parseArgs parses the command-line arguments and returns the interval, token file path, and minimum stars.
	 * It returns an error if the arguments are invalid.
	 * If the token file path is not provided, it returns an empty string for the token file path.
	 * If the interval is not provided or less than or equal to zero, it returns an error.
	 * If the minimum stars is not provided or below zero, it returns an error.
	 */
	// Create a new flag set
	set := flag.NewFlagSet("", flag.ContinueOnError)
	var interval time.Duration
	var minStars int
	var tokenFilePath string

	// Set the flag set output to discard the default output
	set.DurationVar(&interval, "interval", 10*time.Second, "")
	set.IntVar(&minStars, "min-stars", 0, "minimum stars")
	set.StringVar(&tokenFilePath, "tokenFile", "", "GitHub personal access token")
	set.SetOutput(io.Discard)

	// Parse the arguments
	args = args[2:]
	if err := set.Parse(args); err != nil {
		return 0, "", 0, errors.New("got invalid flags")
	}
	// Check if the interval is valid
	if interval <= 0 {
		return 0, "", 0, errors.New("got invalid interval")
	}
	// Check if the minimum stars is valid
	if minStars < 0 {
		return 0, "", 0, errors.New("got invalid min-stars")
	}
	return interval, tokenFilePath, minStars, nil
}

func log(message string) {
	/*
	 * log prints the message to the standard error.
	 */
	fmt.Fprintf(os.Stderr, "%s: %s\n", name, message)
}

func main() {
	// Run the command
	if err := run(args); err != nil {
		message := err.Error()
		log(message)
		if _, ok := err.(usageError); ok {
			message := fmt.Sprintf("run '%s help' for usage information", name)
			log(message)
		}
	}
}

func readTokenFromFile(filePath string) (string, error) {
	/*
	 * readTokenFromFile reads the GitHub personal access token from the provided file path.
	 * It returns the token as a string and an error if the file cannot be read.
	 */
	token, err := os.ReadFile(filePath)
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

func run(args []string) error {
	fmt.Print(args)
	if nbArgs := len(args); nbArgs < 2 {
		return usageError{message: "missing command"}
	}
	switch args[1] {
	case "help":
		const usage = `
Simple CLI for tracking public GitHub repositories.

Usage:
  %[1]s help
  %[1]s track [-interval=<interval>] [-min-stars=<min-stars>] [-token-file=<path/to/token>]

Commands:
  help  Show usage information
  track Track public GitHub repositories

Options:
  -interval=<interval> Repository update interval, greater than zero [default: 10s]
  -min-stars=<stars>   Minimum number of stars required for a repository to be tracked [default: 0]
  -token-file=<path>   File containing the GitHub personal access token [default: ""]`

		fmt.Fprintf(os.Stdout, usage, name)
		return nil

	case "track":
		// Set up context and other variables as needed
		ctx := context.Background()
		// Parse command-line arguments
		interval, tokenFilePath, minStars, err := parseArgs(args)

		// Check if there was an error parsing the arguments
		if err != nil {
			message := fmt.Sprintf("failed parsing argument: %v", err)
			return usageError{message: message}
		}

		// Read token from file if needed
		token := ""
		//Only read token from file if token file path is provided
		if tokenFilePath != "" {
			token, err = readTokenFromFile(tokenFilePath)
			if err != nil {
				message := fmt.Sprintf("Error reading token from file: %v", err)
				return usageError{message: message}
			}
		}

		// Track repositories
		if err := internal.Track(ctx, interval, token, minStars); err != nil {
			return fmt.Errorf("failed tracking: %v", err)
		}
		return nil

	default:
		return usageError{message: "got invalid command"}
	}
}

type usageError struct {
	/*
	 * usageError is an error type for usage errors.
	 */
	message string
}

func (e usageError) Error() string {
	/*
	 * Error returns the error message.
	 */
	return e.message
}

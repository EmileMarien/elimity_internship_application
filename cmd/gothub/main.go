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
)

var name = makeName()

func log(message string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", name, message)
}
func init() {
	flag.StringVar(&tokenFilePath, "token-file", "", "path to the file containing the GitHub personal access token")
	flag.BoolVar(&useToken, "use-token", false, "set to true to use GitHub personal access token")
	flag.Parse()
}

func main() {
	// Read token from file if provided and useToken flag is set
	var token string
	if useToken {
		if tokenFilePath != "" {
			t, err := readTokenFromFile(tokenFilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading token file: %v\n", err)
				os.Exit(1)
			}
			token = t
		} else {
			fmt.Fprintln(os.Stderr, "Error: GitHub personal access token file path is required when -use-token is set to true")
			os.Exit(1)
		}
	}

	// Create a cancellation context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancellation of context when main exits

	// Run the tracking function in a separate goroutine
	go func() {
		interval, err := parseInterval()
		if err != nil {
			log(fmt.Sprintf("failed parsing interval: %v", err))
			return
		}
		if err := run(ctx, interval, token); err != nil {
			log(fmt.Sprintf("failed running: %v", err))
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

func parseInterval() (time.Duration, error) {
	set := flag.NewFlagSet("", flag.ContinueOnError)
	var interval time.Duration
	set.DurationVar(&interval, "interval", 10*time.Second, "")
	set.SetOutput(ioutil.Discard)
	args := args[2:]
	if err := set.Parse(args); err != nil {
		return 0, errors.New("got invalid flags")
	}
	if interval <= 0 {
		return 0, errors.New("got invalid interval")
	}
	return interval, nil
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

func run(ctx context.Context, interval time.Duration, token string) error {
	if nbArgs := len(args); nbArgs < 2 {
		return usageError{message: "missing command"}
	}
	switch args[1] {
	case "help":
		const usage = `
Simple CLI for tracking public GitHub repositories.

Usage:
  %[1]s help
  %[1]s track [-interval=<interval>]

Commands:
  help  Show usage information
  track Track public GitHub repositories

Options:
  -interval=<interval> Repository update interval, greater than zero [default: 10s]
`
		fmt.Fprintf(os.Stdout, usage, name)
		return nil

	case "track":
		interval, err := parseInterval()
		if err != nil {
			message := fmt.Sprintf("failed parsing interval: %v", err)
			return usageError{message: message}
		}
		if err := internal.Track(ctx, interval, token); err != nil {
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

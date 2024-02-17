package main

///Users/emile/Documents/github_token_EmileMarien.rtf
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
	args = os.Args
)

var name = makeName()

func parseArgs() (time.Duration, string, int, error) {
	set := flag.NewFlagSet("", flag.ContinueOnError)
	var interval time.Duration
	var minStars int
	var tokenFilePath string

	set.DurationVar(&interval, "interval", 10*time.Second, "")
	set.IntVar(&minStars, "min-stars", 0, "minimum stars")
	set.StringVar(&tokenFilePath, "tokenFile", "", "GitHub personal access token")
	set.SetOutput(ioutil.Discard)
	args := args[2:]
	if err := set.Parse(args); err != nil {
		return 0, "", 0, errors.New("got invalid flags")
	}
	if interval <= 0 {
		return 0, "", 0, errors.New("got invalid interval")
	}
	return interval, tokenFilePath, minStars, nil
}
func log(message string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", name, message)
	os.Exit(1)
}

func main() {
	// Set up context and other variables as needed
	ctx := context.Background()

	// Parse command-line arguments
	interval, tokenFilePath, minStars, err := parseArgs()
	fmt.Fprintf(os.Stdout, "interval: %s\n", interval)
	fmt.Fprintf(os.Stdout, "tokenFilePath: %s\n", tokenFilePath)
	fmt.Fprintf(os.Stdout, "minStars: %d\n", minStars)

	if err != nil {
		log(fmt.Sprintf("Error parsing arguments: %v", err))
		return
	}
	// Read token from file if needed

	var token string
	if tokenFilePath != "" { //Only read token from file if token file path is provided
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
		if err := run(ctx, interval, token, minStars); err != nil {
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

func run(ctx context.Context, interval time.Duration, tokenFilePath string, minStars int) error {
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

		if err := internal.Track(ctx, interval, tokenFilePath, minStars); err != nil {
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

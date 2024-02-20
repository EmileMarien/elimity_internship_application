package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

func Track(ctx context.Context, interval time.Duration, token string, minStars int) error {
	/* Track periodically fetches the most recently updated public GitHub repositories
	 * that meet the minimum star requirement and prints their details to the standard output.
	 * It uses the provided interval to wait between each fetch.
	 * If the token is not empty, it uses the token to authenticate with GitHub.
	 * If the token is empty, it uses an unauthenticated client to fetch the repositories.
	 * It returns an error if the GitHub API request fails.
	 *
	 * - ctx: the context that the function should use to handle cancellations and timeouts
	 * - interval: the duration to wait between each fetch
	 * - token: the GitHub personal access token to authenticate with GitHub
	 * - minStars: the minimum number of stars required for a repository to be tracked
	 * - error: an error if the GitHub API request fails
	 */

	// Print table headers
	fmt.Println("Owner\t| Name\t| Updated at (UTC)\t| Star count")

	// Map to keep track of printed repositories to prevent printing duplicates
	printedRepositories := make(map[string]bool)

	// Loop until context is cancelled
	for ; ; <-time.Tick(interval) {
		// Create a GitHub client without authentication
		client := github.NewClient(nil)
		if len(token) > 0 { // if token is not empty, overwrite client with authenticated client

			// Create a GitHub client with authentication using the provided token
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(ctx, ts)
			client = github.NewClient(tc)
		}
		// Set the list options to limit the number of results
		listOptions := github.ListOptions{PerPage: 3}
		searchOptions := &github.SearchOptions{ListOptions: listOptions, Sort: "updated"}

		// Construct the search query with the minimum number of stars
		query := "is:public"
		if minStars > 0 {
			query = fmt.Sprintf("is:public stars:>=%d", minStars)
		}

		// Search for repositories using the provided query and search options
		result, _, err := client.Search.Repositories(ctx, query, searchOptions)
		if err != nil {
			return err
		}

		// Iterate over the repositories and print their details
		for _, repository := range result.Repositories {

			// Get owner details or organization details if the owner is an organization
			owner := repository.GetOwner().GetLogin()
			if repository.GetOwner().GetType() == "Organization" {
				owner = repository.GetOwner().GetLogin() + " (org)"
			}

			// Get name, last updated time, and star count of the repository
			name := repository.GetName()
			updatedAt := repository.GetUpdatedAt().UTC().Format("2006-01-02T15:04:05")
			starCount := repository.GetStargazersCount()

			// Check if the repository has already been printed
			if printedRepositories[name] {
				continue
			}
			// Print repository details
			fmt.Printf("%s\t| %s\t| %s\t| %d\n", owner, name, updatedAt, starCount)
			// Mark the repository as printed
			printedRepositories[name] = true
		}
	}
}

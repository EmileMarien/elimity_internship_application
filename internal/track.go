package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

// Track tracks public GitHub repositories, continuously updating according to the given interval.
//
// The given interval must be greater than zero.
func Track(ctx context.Context, interval time.Duration, token string) error {
	// Create a GitHub client with authentication using the provided token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	for {
		select {
		case <-ctx.Done():
			// If the context is canceled, stop the tracking loop.
			return ctx.Err()
		case <-time.After(interval):
			listOptions := github.ListOptions{PerPage: 3}
			searchOptions := &github.SearchOptions{ListOptions: listOptions, Sort: "updated"}
			result, _, err := client.Search.Repositories(ctx, "is:public", searchOptions)
			if err != nil {
				return err
			}
			// Print table headers
			fmt.Println("Owner\t| Name\t| Updated at (UTC)\t| Star count")
			// Print repository details
			for _, repository := range result.Repositories {
				owner := repository.GetOwner().GetLogin()
				if repository.GetOwner().GetType() == "Organization" {
					owner = repository.GetOwner().GetLogin()
				}
				name := repository.GetName()
				updatedAt := repository.GetUpdatedAt().UTC().Format("2006-01-02T15:04:05")
				starCount := repository.GetStargazersCount()
				fmt.Printf("%s\t| %s\t| %s\t| %d\n", owner, name, updatedAt, starCount)
			}
		}
	}
}

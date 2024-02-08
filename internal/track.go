package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v33/github"
)

// Track tracks public GitHub repositories, continuously updating according to the given interval.
//
// The given interval must be greater than zero.
/*
func Track(ctx context.Context, interval time.Duration, useToken bool, token string) error {
    var client *github.Client
    if useToken {
        // Create a GitHub client with authentication using the provided token
        ts := oauth2.StaticTokenSource(
            &oauth2.Token{AccessToken: token},
        )
        tc := oauth2.NewClient(ctx, ts)
        client = github.NewClient(tc)
    } else {
        // Create a GitHub client without authentication
        client = github.NewClient(nil)
    }

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

*/

func Track(ctx context.Context, interval time.Duration, useToken bool, token string) error {
	for ; ; <-time.Tick(interval) {
		// Create a GitHub client without authentication
		client := github.NewClient(nil)
		/*
			useToken=false
			if useToken {
				// Create a GitHub client with authentication using the provided token
				ts := oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: token},
				)
				tc := oauth2.NewClient(ctx, ts)
				client = github.NewClient(tc)
			}
		*/
		con := context.Background()
		listOptions := github.ListOptions{PerPage: 3}
		searchOptions := &github.SearchOptions{ListOptions: listOptions, Sort: "updated"}
		result, _, err := client.Search.Repositories(con, "is:public", searchOptions)
		if err != nil {
			return err
		}
		for _, repository := range result.Repositories {
			// Print repository details
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

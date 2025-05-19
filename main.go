package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-tfe"
)

var (
	teamNames = map[string]string{}
)

func main() {
	ctx := context.Background()

	// Parse command line flags
	delete := flag.Bool("delete", false, "Deletes the team tokens that fit the provided criteria for deletion. Defaults to false.")
	team := flag.String("team", "", "The team name to delete tokens for. If not provided, tokens from all teams will be considered for deletion.")
	deleteExpired := flag.Bool("expired", true, "Marks expired tokens for deletion, regardless of created_at or last_used_at.")
	lastUsedAt := flag.Int("last-used-days-ago", 30, "Duration of time in days for how long ago a resource should have been "+
		"last used before deleting.")
	createdAt := flag.Int("created-at-days-ago", 0, "Duration of time in days for how long ago a resource should have been "+
		"created before deleting.")
	flag.Parse()

	// Initialize TFE client
	config := &tfe.Config{
		RetryServerErrors: true,
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		fmt.Println("error initializing TFE client:", err)
		os.Exit(1)
	}

	orgName := os.Getenv("TFE_ORGANIZATION")
	if orgName == "" {
		fmt.Println("TFE_ORGANIZATION environment variable is not set.")
		os.Exit(1)
	}

	// List team tokens
	opts := &tfe.TeamTokenListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	}
	if *team != "" {
		opts.Query = *team
	}
	var tokens []*tfe.TeamToken
	for {
		resp, err := client.TeamTokens.List(ctx, orgName, opts)
		if err != nil {
			fmt.Printf("error listing team tokens for organization %s: %v\n", orgName, err)
			os.Exit(1)
		}
		tokens = append(tokens, resp.Items...)

		if resp.CurrentPage >= resp.TotalPages {
			break
		}

		opts.PageNumber = resp.NextPage
	}

	// Determine which tokens to delete
	var createdAtDuration time.Duration
	if *createdAt > 0 {
		createdAtDuration = time.Duration(*createdAt*24) * time.Hour
	}

	var lastUsedAtDuration time.Duration
	if *lastUsedAt > 0 {
		lastUsedAtDuration = time.Duration(*lastUsedAt*24) * time.Hour
	}

	now := time.Now()
	toDelete := make([]*tfe.TeamToken, 0)
	for _, token := range tokens {
		team := getTeamName(ctx, client, token.Team.ID)
		var identifier string
		if token.Description != nil && *token.Description != "" {
			identifier = *token.Description
		} else {
			identifier = token.ID
		}

		if *deleteExpired && !token.ExpiredAt.IsZero() && token.ExpiredAt.Before(now) {
			fmt.Printf("Marking token for deletion because expired: '%s' in team '%s' expired_at=%s \n", identifier, team, token.ExpiredAt.String())
			toDelete = append(toDelete, token)
			continue
		}

		if lastUsedAtDuration > 0 && time.Since(token.LastUsedAt) > lastUsedAtDuration {
			fmt.Printf("Marking token for deletion because last used too long ago: '%s' in team '%s' last_used_at=%s \n", identifier, team, token.LastUsedAt.String())
			toDelete = append(toDelete, token)
			continue
		}

		if createdAtDuration > 0 && time.Since(token.CreatedAt) > createdAtDuration {
			fmt.Printf("Marking token for deletion because created too long ago: '%s' in team '%s' created_at=%s \n", identifier, team, token.CreatedAt.String())
			toDelete = append(toDelete, token)
			continue
		}
	}
	fmt.Printf("\n%d tokens marked for deletion.\n", len(toDelete))

	// Exit if delete flag is not set or no tokens to delete
	if !*delete {
		fmt.Println("Use the --delete flag to delete the tokens that fit the specified criteria.")
		os.Exit(0)
	}
	if len(toDelete) == 0 {
		os.Exit(0)
	}

	// Confirm deletion
	r := bufio.NewReader(os.Stdin)
	fmt.Println("Are you sure you want to delete these team tokens? (y/n): ")
	resp, err := r.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	if strings.ToLower(strings.TrimSpace(resp)) != "y" && strings.ToLower(strings.TrimSpace(resp)) != "yes" {
		fmt.Println("Aborting deletion.")
		os.Exit(0)
	}

	// Delete tokens
	for _, token := range toDelete {
		if token.Description != nil && *token.Description != "" {
			fmt.Printf("Deleting token: %s (%s)\n", token.ID, *token.Description)
		} else {
			fmt.Printf("Deleting token: %s\n", token.ID)
		}

		err := client.TeamTokens.DeleteByID(ctx, token.ID)
		if err != nil {
			fmt.Printf("Error deleting token %s: %v\n", token.ID, err)
			// Continue trying to delete other tokens in case intermittent issue
			continue
		}
	}
	fmt.Println("Team tokens deleted.")
}

func getTeamName(ctx context.Context, client *tfe.Client, teamID string) string {
	if name, ok := teamNames[teamID]; ok {
		return name
	}

	team, err := client.Teams.Read(ctx, teamID)
	if err != nil {
		// The team name is just for improved logging. If we have issues readin
		// the team, just use the team ID.
		teamNames[teamID] = teamID
		return teamID
	}
	teamNames[teamID] = team.Name

	return team.Name
}

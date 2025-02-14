package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v69/github"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

var gh *github.Client

func InitializeGhClient() {
	if gh == nil {
		ctx := context.Background()
		token := os.Getenv("GITHUB_TOKEN") // Use a GitHub token for authentication if needed

		if token != "" {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			gh = github.NewClient(oauth2.NewClient(ctx, ts))
		} else {
			gh = github.NewClient(nil)
		}
	}
}

type Release struct {
	ID          int64     `json:"id"`
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	ZipballURL  string    `json:"zipball_url"`
	URL         string    `json:"url"`
	HTMLURL     string    `json:"html_url"`
	PublishedAt time.Time `json:"published_at"`
}

func fetchReleases(owner, repo string) ([]Release, error) {
	ctx := context.Background()
	var allReleases []Release
	opt := &github.ListOptions{PerPage: 100}

	for {
		releases, resp, err := gh.Repositories.ListReleases(ctx, owner, repo, opt)
		if err != nil {
			return nil, fmt.Errorf("error fetching releases: %w", err)
		}

		for _, r := range releases {
			if r.GetDraft() || r.GetPrerelease() {
				continue // Ignore draft and prerelease
			}

			release := Release{
				ID:          r.GetID(),
				TagName:     r.GetTagName(),
				Name:        r.GetName(),
				ZipballURL:  r.GetZipballURL(),
				URL:         r.GetURL(),
				HTMLURL:     r.GetHTMLURL(),
				PublishedAt: r.GetPublishedAt().Time,
			}
			allReleases = append(allReleases, release)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	slices.SortStableFunc(allReleases, func(a, b Release) int {
		if a.PublishedAt.After(b.PublishedAt) {
			return -1
		} else if a.PublishedAt.Before(b.PublishedAt) {
			return 1
		}
		return 0
	})

	return allReleases, nil
}

func saveReleasesToFile(releases []Release, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(releases); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func filterReleasesAfter(releases []Release, minDate time.Time) []Release {
	var filtered []Release
	for _, release := range releases {
		if release.PublishedAt.After(minDate) {
			filtered = append(filtered, release)
		}
	}
	return filtered
}

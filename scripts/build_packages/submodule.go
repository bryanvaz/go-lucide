package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	filepathPkg "path/filepath"
	"strings"
)

// getGitTags extracts all tags from the submodule
func getGitTags(path string) ([]string, error) {
	cmd := exec.Command("git", "tag")
	cmd.Dir = path
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Split by newlines and clean up the tags
	var tags []string
	for _, tag := range strings.Split(out.String(), "\n") {
		trimmedTag := strings.TrimSpace(tag)
		if trimmedTag != "" {
			tags = append(tags, trimmedTag)
		}
	}

	return tags, nil
}

func findMissingTags(tagList []string, releases []Release) []Release {
	tagSet := make(map[string]bool)
	missingTags := []Release{}
	latestMatchedTag := ""

	// Add all tags from the first list into a set
	for _, tag := range tagList {
		tagSet[strings.TrimPrefix(tag, "v")] = true
	}

	// Check for tags in the second list that are missing in the first list
	for _, rel := range releases {
		if !tagSet[strings.TrimPrefix(rel.TagName, "v")] { // If the tag is not in the first list, add it to missingTags
			missingTags = append(missingTags, rel)
		} else if latestMatchedTag == "" {
			latestMatchedTag = rel.TagName
			break
		}
	}

	fmt.Printf("Latest matched tag: %s\n", latestMatchedTag)

	return missingTags
}

func cloneRepo(url, path string) error {
	gitDir := filepathPkg.Join(path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", url, path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func fetchTags(path string) error {
	cmd := exec.Command("git", "-C", path, "fetch", "--tags")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func checkoutTag(tag, path string) {
	cmd := exec.Command("git", "-C", path, "checkout", "tags/"+tag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}

}

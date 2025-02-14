package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"os"
	filepathPkg "path/filepath"
	"time"
)

const (
	LUCIDE_OWNER         = "lucide-icons"
	LUCIDE_REPO          = "lucide"
	LUCIDE_DIR           = "./lucide"
	MIN_LUCIDE_VERSION   = "2024-01-01"
	LUCIDE_GIT_URL       = "https://github.com/lucide-icons/lucide.git"
	LUCIDE_GIT_DIR       = "./tmp/lucide"
	TEMPL_SUBMODULE_PATH = "./packages/go-templ-lucide-icons"
	TEMPL_UTILS_PATH     = "./src/templ"
)

func main() {
	InitializeGhClient()
	versionsAfterTime, err := time.Parse("2006-01-02", MIN_LUCIDE_VERSION)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	// sync lucide icon versions
	fmt.Println("Syncing lucide icon releases...")
	lucideReleases, err := fetchReleases(LUCIDE_OWNER, LUCIDE_REPO)
	if err != nil {
		fmt.Println("Error fetching releases:", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d releases for %s/%s\n", len(lucideReleases), LUCIDE_OWNER, LUCIDE_REPO)
	fmt.Printf("Latest release: %s (%s)\n", lucideReleases[0].TagName, lucideReleases[0].PublishedAt)

	releasesJsonPath := filepathPkg.Join(LUCIDE_DIR, "releases.json")
	saveReleasesToFile(lucideReleases, releasesJsonPath)
	fmt.Println("Releases saved to", releasesJsonPath)

	releases := filterReleasesAfter(lucideReleases, versionsAfterTime)
	fmt.Printf("Found %d releases after %s\n", len(releases), MIN_LUCIDE_VERSION)

	// sync templ package library
	fmt.Println("Updating submodule and fetching tags...")
	err = updateSubmodule(TEMPL_SUBMODULE_PATH)
	if err != nil {
		fmt.Println("Error updating submodule:", err)
		return
	}
	// get list of tags on templ package
	fmt.Println("Extracting all tags from submodule...")
	tags, err := getGitTags(TEMPL_SUBMODULE_PATH)
	if err != nil {
		fmt.Println("Error fetching tags:", err)
		return
	}

	// build list of missing tags starting from latest tag
	missingReleases := findMissingTags(tags, releases)
	fmt.Printf("Found %d missing tags\n", len(missingReleases))

	// ask if want to sync the next version tag or all
	var syncOption string = "Next Version"
	var pushToGitHub bool
	var publishToProxy bool

	syncPrompt := &survey.Select{
		Message: "Do you want to sync the next version tag or all?",
		Options: []string{"Next Version", "All Versions"},
	}
	survey.AskOne(syncPrompt, &syncOption)

	pushPrompt := &survey.Confirm{
		Message: "Do you want to push changes to GitHub? (default: no)",
		Default: false,
	}
	survey.AskOne(pushPrompt, &pushToGitHub)

	if pushToGitHub {
		publishPrompt := &survey.Confirm{
			Message: "Do you want to publish to proxy.golang.org? (default: yes)",
			Default: true,
		}
		survey.AskOne(publishPrompt, &publishToProxy)
	}

	fmt.Println("Responses:")
	fmt.Printf("Sync option: %s\n", syncOption)
	fmt.Printf("Push to GitHub: %t\n", pushToGitHub)
	fmt.Printf("Publish to proxy.golang.org: %t\n", publishToProxy)

	// sync the repo into the tmp directory
	fmt.Println("Syncing lucide icon repo...")
	os.MkdirAll("./tmp/lucide", os.ModePerm)
	err = cloneRepo(LUCIDE_GIT_URL, LUCIDE_GIT_DIR)
	if err != nil {
		fmt.Println("Error cloning lucide repo:", err)
		os.Exit(1)
	}
	err = fetchTags(LUCIDE_GIT_DIR)
	if err != nil {
		fmt.Println("Error fetching tags:", err)
		os.Exit(1)
	}

	currRel := missingReleases[len(missingReleases)-1]
	fmt.Println("working on lucide release:", currRel.TagName)
	// switch to tag to sync
	checkoutTag(currRel.TagName, LUCIDE_GIT_DIR)
	// read icons from icons directory
	svgIcons, err := injestIcons(LUCIDE_GIT_DIR)
	if err != nil {
		fmt.Println("Error reading icons:", err)
		os.Exit(1)
	}

	// Delete the templ output folder if it exists
	submoduleTemplPath := filepathPkg.Join(TEMPL_SUBMODULE_PATH, "icons")
	if _, err := os.Stat(submoduleTemplPath); err == nil {
		err := os.RemoveAll(submoduleTemplPath)
		if err != nil {
			fmt.Println("Error deleting folder:", err)
			os.Exit(1)
		}
	}
	if _, err := os.Stat(filepathPkg.Join(submoduleTemplPath, "VERSION")); err == nil {
		err := os.Remove(filepathPkg.Join(submoduleTemplPath, "VERSION"))
		if err != nil {
			fmt.Println("Error deleting file:", err)
			os.Exit(1)
		}
	}

	// generate templ file
	fmt.Printf("Generating %d templ and go files...\n", len(svgIcons))
	templFiles := make(map[string]string)
	goFiles := make(map[string]string)
	for _, icon := range svgIcons {
		templFunc, err := generateTemplFunc(icon)
		if err != nil {
			fmt.Println("Error generating templ func:", err)
			os.Exit(1)
		}
		templFile, err := generateTemplFile(templFunc)
		if err != nil {
			fmt.Println("Error generating templ file:", err)
			os.Exit(1)
		}
		templFiles[icon.Basename()] = templFile
		goFile, err := generateGoFromTempl(templFile)
		if err != nil {
			fmt.Println("Error generating go file:", err)
			os.Exit(1)
		}
		goFiles[icon.Basename()] = goFile
	}

	fmt.Printf("Writing %d templ and go files...\n", len(templFiles))
	if err := os.MkdirAll(submoduleTemplPath, os.ModePerm); err != nil {
		fmt.Println("Error creating folder:", err)
		os.Exit(1)
	}
	for basename, templFile := range templFiles {
		outputTemplPath := filepathPkg.Join(submoduleTemplPath, basename+".templ")
		outputGoPath := filepathPkg.Join(submoduleTemplPath, basename+"_templ.go")
		if err := os.WriteFile(outputTemplPath, []byte(templFile), 0644); err != nil {
			fmt.Println("Error writing to output file:", outputTemplPath, err)
			os.Exit(1)
		}
		if err := os.WriteFile(outputGoPath, []byte(goFiles[basename]), 0644); err != nil {
			fmt.Println("Error writing to output file:", outputGoPath, err)
			os.Exit(1)
		}
	}

	// Write VERSION file
	if err := os.WriteFile(filepathPkg.Join(TEMPL_SUBMODULE_PATH, "VERSION"), []byte(currRel.TagName), 0644); err != nil {
		fmt.Println("Error writing to VERSION file:", err)
		os.Exit(1)
	}

	// Copy utils.go file
	srcUtilsPath := filepathPkg.Join(TEMPL_UTILS_PATH, "utils.go")
	dstUtilsPath := filepathPkg.Join(TEMPL_SUBMODULE_PATH, "icons", "utils.go")
	if err := copyFile(srcUtilsPath, dstUtilsPath); err != nil {
		fmt.Println("Error copying utils.go:", err)
	}
	fmt.Println("Copied utils.go")
	srcDefaultAttributesPath := filepathPkg.Join(TEMPL_UTILS_PATH, "default_attributes.go")
	destDefaultAttributesPath := filepathPkg.Join(TEMPL_SUBMODULE_PATH, "icons", "default_attributes.go")
	if err := copyFile(srcDefaultAttributesPath, destDefaultAttributesPath); err != nil {
		fmt.Println("Error copying default_attributes.go:", err)
	}
	fmt.Println("Copied default_attributes.go")

	rollupFile, err := createRollupFile(svgIcons)
	if err != nil {
		fmt.Println("Error creating rollup file:", err)
		os.Exit(1)
	}
	rollupFilePath := filepathPkg.Join(TEMPL_SUBMODULE_PATH, "icons.go")
	if err := os.WriteFile(rollupFilePath, []byte(rollupFile), 0644); err != nil {
		fmt.Println("Error writing to rollup file:", err)
		os.Exit(1)
	}
	fmt.Println("Rollup file saved to", rollupFilePath)

	fmt.Println("Done writing files")

	// commit changes
	// add git tag

	// publish to git hub if applicable
	// publish to proxy.goland.org if applicable
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

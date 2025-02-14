package main

import (
	"fmt"
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
	LUCIDE_GIT_DIR       = "./dist/lucide"
	TEMPL_GIT_URL        = "git@github.com:bryanvaz/go-templ-lucide-icons.git"
	TEMPL_SUBMODULE_PATH = "./dist/go-templ-lucide-icons"
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
	releasesJsonPath := filepathPkg.Join(LUCIDE_DIR, "releases.json")
	saveReleasesToFile(lucideReleases, releasesJsonPath)
	releases := filterReleasesAfter(lucideReleases, versionsAfterTime)
	fmt.Printf("  Latest Lucide release: %s (%s)\n", lucideReleases[0].TagName, lucideReleases[0].PublishedAt)

	// sync the repo into the dist directory
	fmt.Println("Syncing templ-lucide icon repo...")
	os.MkdirAll("./dist", os.ModePerm)
	err = cloneRepo(TEMPL_GIT_URL, TEMPL_SUBMODULE_PATH)
	if err != nil {
		fmt.Println("Error cloning templ icon repo:", err)
		os.Exit(1)
	}
	err = fetchTags(TEMPL_SUBMODULE_PATH)
	if err != nil {
		fmt.Println("Error fetching tags:", err)
		os.Exit(1)
	}
	tags, err := getGitTags(TEMPL_SUBMODULE_PATH)
	if err != nil {
		fmt.Println("Error fetching tags:", err)
		return
	}
	if len(tags) == 0 {
		fmt.Println("No tags found in templ-lucide")
	} else {
		fmt.Printf("  Latest templ-lucide release: %s\n", tags[0])
	}

	// build list of missing tags starting from latest tag
	missingReleases := findMissingTags(tags, releases)
	fmt.Printf("  Found %d missing tags\n", len(missingReleases))
	currRel := missingReleases[len(missingReleases)-1]
	fmt.Println("  Next release to sync:", currRel.TagName)

	// sync the repo into the tmp directory
	fmt.Println("Syncing upstream lucide repo for icons ...")
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

	// switch to tag to sync
	fmt.Printf("  Checking out tag %s\n", currRel.TagName)
	checkoutTag(currRel.TagName, LUCIDE_GIT_DIR)
	svgIcons, err := injestIcons(LUCIDE_GIT_DIR)
	if err != nil {
		fmt.Println("Error reading icons:", err)
		os.Exit(1)
	}

	fmt.Println("--------------------------------------")
	// Delete the templ output folder if it exists
	submoduleTemplPath := filepathPkg.Join(TEMPL_SUBMODULE_PATH, "icons")
	fmt.Println("Cleaning up old templ files...")
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
	fmt.Printf("Generating templ and go files for %d icons ...\n", len(svgIcons))
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

	fmt.Printf("Writing templ and go files ...\n")
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

	fmt.Printf("Done writing files for release %s \n", currRel.TagName)
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

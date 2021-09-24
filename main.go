package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func createCleanPathFromUserSelection(menuSelection *string, userRoot *string) (string, string) {
	reader := bufio.NewReader(os.Stdin)
	var cleanDirPath string
	switch *menuSelection {
	case "1":
		cleanDirPath = *userRoot + "/Desktop"
		return cleanDirPath, "Desktop"
	case "2":
		cleanDirPath = *userRoot + "/Downloads"
		return cleanDirPath, "Downloads"
	case "3":
		cleanDirPath = *userRoot + "/Documents"
		return cleanDirPath, "Documents"
	case "4":
		fmt.Printf("Which directoy:\n")
		userDefinedPath, _ := reader.ReadString('\n')
		userDefinedPath = strings.ReplaceAll(userDefinedPath, "\n", "")
		cleanDirPath = *userRoot + userDefinedPath
		return cleanDirPath, userDefinedPath
	case "0":
		fmt.Printf("Exiting...\n")
		return "", ""
	default:
		fmt.Printf("There was an error, exiting...\n")
		return "", ""
	}
}

func main() {
	var cleanDirPath string
	var screenShotFileNames []string
	var userDirName string
	userRoot := os.Getenv("HOME")
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Screenshot Cleaner\n")
	fmt.Printf("==================\n\n")
	fmt.Printf("Which directory would you like to clean:\n")
	fmt.Printf("1. Desktop\n")
	fmt.Printf("2. Downloads\n")
	fmt.Printf("3. Documents\n")
	fmt.Printf("4. Other\n")
	fmt.Printf("0. Exit\n")

	// Get user root + dir to clean
	menuSelection, _ := reader.ReadString('\n')
	menuSelection = strings.ReplaceAll(menuSelection, "\n", "")

	cleanDirPath, userDirName = createCleanPathFromUserSelection(&menuSelection, &userRoot)
	if cleanDirPath == "" {
		return
	}

	fmt.Printf("Clean all screenshot files in %s - is this correct? [Y/N]\n", cleanDirPath)
	cleanConfirmText, _ := reader.ReadString('\n')
	cleanConfirmText = strings.ToLower(strings.ReplaceAll(cleanConfirmText, "\n", ""))
	if cleanConfirmText == "n" {
		fmt.Printf("Exiting...")
		return
	}

	// create slices of all files in cleanDirPath with token
	files, err := ioutil.ReadDir(cleanDirPath)
	if err != nil {
		fmt.Printf("There was an error, exiting...")
		return
	}
	fmt.Printf("\n[Screenshot files in %s]:\n\n", cleanDirPath)

	for _, file := range files {
		// Get all files with a `Screenshot` substring & .png filetype
		pngFile := strings.Split(file.Name(), ".png")
		if len(pngFile) == 2 {
			screenShotFile := strings.Split(file.Name(), "Screenshot ")
			if len(screenShotFile) == 2 {
				screenShotFileNames = append(screenShotFileNames, file.Name())
				fmt.Printf("- %s\n", file.Name())
			}
		}
	}
	fmt.Printf("The above files will be deleted from the %s folder\n", userDirName)
	fmt.Printf("Do you want to proceed with the deletions: [Y/N]\n")
	cleanConfirm, _ := reader.ReadString('\n')
	cleanConfirm = strings.ToLower(strings.ReplaceAll(cleanConfirm, "\n", ""))
	if cleanConfirm == "n" {
		fmt.Printf("Aborting clean...\n")
		fmt.Printf("Exiting...\n")
		return
	}
	// slice of file names that couldn't be deleted
	// for some reason.
	// TODO it would be better if we checked user permissions
	// before even trying to delete the file.
	var persistingFiles []string
	for _, file := range screenShotFileNames {
		err := os.Remove(cleanDirPath + "/" + file)
		if err != nil {
			persistingFiles = append(persistingFiles, file)
			continue
		}
	}
	if len(persistingFiles) > 0 {
		listPersistingFiles(persistingFiles)
	}
	fmt.Printf("Successfully removed files")
}

func listPersistingFiles(fileNames []string) {
	for _, f := range fileNames {
		fmt.Printf("file %s couldn't be removed\n", f)
	}
}

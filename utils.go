package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func getDownload(link string) string {
	// Instantiate a new Collector
	c := colly.NewCollector()

	// Create a channel to communicate the found link
	linkChan := make(chan string, 1) // Buffered channel with a capacity of 1

	c.OnHTML("img.d_img_holder", func(e *colly.HTMLElement) {
		downloadLink := e.Attr("src")
		// fmt.Println("Checking: " + link + "\n\tGot: \n\t\t" + downloadLink)
		linkChan <- downloadLink // Send the found link back through the channel
	})

	err := c.Visit(link)
	if err != nil {
		fmt.Println("Error visiting website: ", err)
	}

	// Wait to receive the link from the channel
	result := <-linkChan
	close(linkChan) // Close the channel

	return result
}

func removeDuplicates(wallpapers []WallpaperStruct) []WallpaperStruct {
	seen := make(map[string]bool)
	var uniqueWallpapers []WallpaperStruct

	for _, wallpaper := range wallpapers {
		if _, ok := seen[wallpaper.url]; !ok {
			seen[wallpaper.url] = true
			uniqueWallpapers = append(uniqueWallpapers, wallpaper)
		}
	}

	return uniqueWallpapers
}

// Search through string to return display port names...
func findConnectedWords(input string) []string {
	// Define a regular expression to find the word "connected" and capture the word before it
	re := regexp.MustCompile(`(\S+)\s+connected`)

	// Find all matches in the input string
	matches := re.FindAllStringSubmatch(input, -1)

	// Extract the captured words from the matches
	var connectedWords []string
	for _, match := range matches {
		if len(match) >= 2 {
			connectedWords = append(connectedWords, "󰍹 "+match[1])
		}
	}

	return connectedWords
}

// Finds connected displays
func FindDisplays() []string {
	command := "xrandr"
	args := []string{"--query"}
	displays := exec.Command(command, args...)

	// Captures the output of the command
	output, err := displays.CombinedOutput()
	if err != nil {
		// Print error and return empty string
		fmt.Println("Error executing:", err)
		return []string{}
	}
	connectedDisplays := findConnectedWords(string(output))

	return connectedDisplays
}

func ChangeWallpaper(file string, head int) {
	// Separate the command and its arguments
	command := "nitrogen"
	args := []string{"--set-scaled", file, "--head=" + fmt.Sprint(head)}

	// Create the exec.Cmd object
	cmd := exec.Command(command, args...)
	fmt.Println(cmd)

	// Set the command's stdout to the current process's stdout
	cmd.Stdout = os.Stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
}
func listJPGFiles(dir string, prefix string) ([]string, error) {
	var files []string

	// Read the directory
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Extract file names from fileInfos
	for _, fileInfo := range fileInfos {
		// Check if it's a directory
		if fileInfo.IsDir() {
			// Recursively list files in the subdirectory with the updated prefix
			subdir := filepath.Join(dir, fileInfo.Name())
			subfiles, err := listJPGFiles(subdir, filepath.Join(prefix, fileInfo.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, subfiles...)
		} else {
			// Check if the file has a ".jpg" extension
			if strings.HasSuffix(fileInfo.Name(), ".jpg") {
				// Add the file with the updated prefix to the list
				files = append(files, filepath.Join(prefix, fileInfo.Name()))
			}
		}
	}

	return files, nil
}

// downloadImage downloads an image from the specified URL and saves it to a specific path.
func downloadImage(imageURL, savePath string) error {
	// Check if the directory exists, if not, create it
	dir := filepath.Dir(savePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Check if the file already exists
	if _, err := os.Stat(savePath); err == nil {
		return fmt.Errorf("file already exists at %s", savePath)
	}

	// Make a GET request to fetch the image
	response, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check that the GET request was successful
	if response.StatusCode != 200 {
		return fmt.Errorf("received non-200 response code")
	}

	// Create a new file in the specified save path
	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, response.Body)
	return err
}

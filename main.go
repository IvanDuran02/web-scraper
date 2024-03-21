package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
	"unicode/utf8"
)

type WallpaperStruct struct {
	title, url, HD_1080p, HD_1440p, HD_2160p string
}

func getDownload(link string) string {
	// Instantiate a new Collector
	c := colly.NewCollector()

	// Create a channel to communicate the found link
	linkChan := make(chan string, 1) // Buffered channel with a capacity of 1

	c.OnHTML("img.d_img_holder", func(e *colly.HTMLElement) {
		downloadLink := e.Attr("src")
		fmt.Println("Checking: " + link + "\n\tGot: \n\t\t" + downloadLink)
		linkChan <- downloadLink // Send the found link back through the channel
	})

	c.Visit(link)

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

func main() {
	// Instantiate a new Collector
	c := colly.NewCollector()

	// Define a slice to store the wallpapers
	var wallpapers []WallpaperStruct
	urlTemplate := "https://hdqwalls.com/%s"

	// OnHTML callback for 'a' elements
	c.OnHTML("a", func(e *colly.HTMLElement) {
		// Extract the title attribute from the 'a' element
		title := e.Attr("title")

		if trimmedTitle := strings.TrimSpace(title); trimmedTitle != "" {
			// Modify the title for the URL
			modifiedTitle := strings.ToLower(strings.ReplaceAll(title, " ", "-"))

			// Construct the modified URL
			modifiedURL := fmt.Sprintf(urlTemplate, modifiedTitle)

			// Safely remove the last 10 characters if the string is long enough
			if utf8.RuneCountInString(modifiedTitle) > 10 {
				modifiedTitle = string([]rune(modifiedTitle)[:len([]rune(modifiedTitle))-10])
			}

			// Create a WallpaperStruct and append it to the wallpapers slice
			wallpapers = append(wallpapers, WallpaperStruct{
				title:    title,
				url:      modifiedURL,
				HD_1080p: getDownload("https://hdqwalls.com/wallpaper/1920x1080/" + modifiedTitle),
				HD_1440p: getDownload("https://hdqwalls.com/wallpaper/2560x1440/" + modifiedTitle),
				HD_2160p: getDownload("https://hdqwalls.com/wallpaper/3840x2160/" + modifiedTitle),
			})
		}
	})

	// Visit the URL
	err := c.Visit("https://hdqwalls.com/category/anime-wallpapers")
	if err != nil {
		fmt.Println("Error visiting the website:", err)
		return
	}

	// Filter duplicates based on URL
	wallpapers = removeDuplicates(wallpapers)

	// Print the wallpapers
	fmt.Println("Wallpapers:")
	for _, wallpaper := range wallpapers {
		fmt.Printf("Title: %s\nURL: %s\n1080p: %s\n1440p: %s\n4k: %s\n", wallpaper.title, wallpaper.url, wallpaper.HD_1080p, wallpaper.HD_1440p, wallpaper.HD_2160p)
	}
}

package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/gocolly/colly"
	"github.com/manifoldco/promptui"
)

type WallpaperStruct struct {
	title, modifiedTitle, url, HD_1080p, HD_1440p, HD_2160p string
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
				title:         title,
				modifiedTitle: modifiedTitle,
				url:           modifiedURL,
				HD_1080p:      getDownload("https://hdqwalls.com/wallpaper/1920x1080/" + modifiedTitle),
				HD_1440p:      getDownload("https://hdqwalls.com/wallpaper/2560x1440/" + modifiedTitle),
				HD_2160p:      getDownload("https://hdqwalls.com/wallpaper/3840x2160/" + modifiedTitle),
			})
		}
	})
	categories := []string{
		"Popular Wallpapers", "Superheros", "Games", "Artists", "Movies", "Celebrities", "Cars",
		"Nature", "TV Shows", "Girls", "Abstract", "Anime", "Music",
		"Photography", "Computer", "Animals", "Digital Universe", "World",
		"Bikes", "Fantasy Girls", "Flowers", "Love", "Birds", "Sports", "Other",
	}

	categoryPrompt := promptui.Select{
		Label: "Select a Category:",
		Items: categories,
	}

	_, result, err := categoryPrompt.Run()
	if err != nil {
		fmt.Println("Error selecting category: ", err)
		return
	}
	if result != categories[0] {
		result = "category/" + strings.ToLower(strings.ReplaceAll(result, " ", "-")) + "-wallpapers"
	} else {
		result = strings.ToLower(strings.ReplaceAll(categories[0], " ", "-"))
	}

	// Visit the URL
	link := fmt.Sprintf("https://hdqwalls.com/%s", result)
	err = c.Visit(link)
	if err != nil {
		fmt.Println("Error visiting the website:", link, err)
		return
	}

	// Filter duplicates based on URL
	wallpapers = removeDuplicates(wallpapers)

	// Convert []WallpaperStruct to slice for prompt
	titles := make([]string, len(wallpapers))
	for i, wallpaperTitle := range wallpapers {
		titles[i] = wallpaperTitle.title
	}

	Prompt("Select a Wallpaper:", titles, wallpapers)
}

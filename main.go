package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type WallpaperStruct struct {
	title, url, HD_1080p, HD_1440p, HD_2160p string
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
			modifiedTitle := ""
			for _, letter := range title {
				if letter == ' ' {
					modifiedTitle += "-"
				} else {
					modifiedTitle += string(letter)
				}
			}
			modifiedTitle = strings.ToLower(modifiedTitle)

			// Construct the modified URL
			modifiedURL := fmt.Sprintf(urlTemplate, modifiedTitle)

			// Create a WallpaperStruct and append it to the wallpapers slice
			wallpapers = append(wallpapers, WallpaperStruct{
				title: title,
				url:   modifiedURL,
			})
		}
	})

	// Visit the URL
	err := c.Visit("https://hdqwalls.com/category/anime-wallpapers")
	if err != nil {
		fmt.Println("Error visiting the website:", err)
		return
	}

	// Print the wallpapers
	fmt.Println("Wallpapers:")
	for _, wallpaper := range wallpapers {
		fmt.Printf("Title: %s\nURL: %s\n\n", wallpaper.title, wallpaper.url)
	}
}

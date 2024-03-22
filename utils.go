package main

import (
	"fmt"
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

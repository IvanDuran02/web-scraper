package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
)

func (w WallpaperStruct) Find(title string, wallpapers []WallpaperStruct) (WallpaperStruct, bool) {
	for _, wallpaper := range wallpapers {
		if wallpaper.title == title {
			return w, true
		}
	}
	return WallpaperStruct{}, false
}

func Prompt(promptLabel string, options []string, wallpapers []WallpaperStruct) {

	prompt := promptui.Select{
		Label: promptLabel,
		Items: options,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed: ", err)
		return
	}

	fmt.Println("You chose: ", result)

	switch promptLabel {
	case "Select a Wallpaper:":
		var filePath string
		var selectedWallpaper WallpaperStruct
		var found bool
		for _, wallpaper := range wallpapers {
			if wallpaper.title == result {
				selectedWallpaper, found = wallpaper.Find(result, wallpapers)
				break
			}
		}
		if found {
			// Do something with the selected wallpaper
			imgQuality := promptui.Select{
				Label: "Select Image Quality",
				Items: []string{"1080p", "1440p", "2160p"},
			}
			_, result, err := imgQuality.Run()
			if err != nil {
				fmt.Println("Did not find the correct quality:", err)
			}
			switch result {
			case "1080p":
				filePath = fmt.Sprintf("./DownloadedWallpapers/" + selectedWallpaper.modifiedTitle + ".jpg")
				err := downloadImage(selectedWallpaper.HD_1080p, filePath)
				if err != nil {
					fmt.Print("Error occurred during image download: ", err)
				}
			case "1440p":
				filePath = fmt.Sprintf("./DownloadedWallpapers/" + selectedWallpaper.modifiedTitle + ".jpg")
				err := downloadImage(selectedWallpaper.HD_1440p, filePath)
				if err != nil {
					fmt.Print("Error occurred during image download: ", err)
				}
			case "2160p":
				filePath = fmt.Sprintf("./DownloadedWallpapers/" + selectedWallpaper.modifiedTitle + ".jpg")
				err := downloadImage(selectedWallpaper.HD_2160p, filePath)
				if err != nil {
					fmt.Print("Error occurred during image download: ", err)
				}
			}

			displays := FindDisplays()
			displayPrompt := promptui.Select{
				Label: "Select a monitor",
				Items: displays,
			}
			_, selectedDisplay, err := displayPrompt.Run()
			if err != nil {
				fmt.Println("Prompt failed:", err)
				return
			}

			head := 0
			for i, display := range displays {
				if selectedDisplay == display {
					head = i
				}
			}

			ChangeWallpaper(filePath, head)
		} else {
			fmt.Println("Wallpaper not found")
		}
	default:
		// Handle other cases
		fmt.Println("Idk...")
	}
}

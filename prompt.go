package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
)

func Prompt() {
	fmt.Println("Prompting...")

	var arr [3]int8
	for _, index := range arr {
		arr[index] = index
	}

	prompt := promptui.Select{
		Label: "Select Display",
		Items: arr,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Something went wrong with the prompt")
	}

	fmt.Println("You chose: ", result)
}

package main

import (
	"GitSubset/internal"
	"fmt"
)

func main() {
	link, err := internal.Clone("https://github.com/kevo-1/Concepts-of-Programming-Languages-Course")

	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Printf("User: %s\nRepo Name: %s\n", link.User, link.Repo)
	}

	files, err := internal.ListContent(link.Path)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		for i := range files {
			fmt.Printf("    %s\n", files[i])
		}
	}

	chosen := []string{"Functional-Programming/Lecture-4/practiceLec.scala", "Functional-Programming/Lecture-3/higherOrderFunctions.scala"}

	if err := internal.FetchContent(link.Path, chosen); err != nil {
		fmt.Println(err.Error())
	}
}

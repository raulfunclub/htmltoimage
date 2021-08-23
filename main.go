package main

import (
	"fmt"
	"htmltoimage/actions"
)

func main() {

	words := []string{"白日依山尽，", "黄河入海流。"}
	image, err := actions.ConvertWordsToImage("req-123456", words)
	fmt.Printf("image: %v; err: %v", image, err)
}


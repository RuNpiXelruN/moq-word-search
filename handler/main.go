package main

import (
	"fmt"

	moqwordsearch "github.com/RuNpiXelruN/moq-word-search"
)

func main() {
	fmt.Println("hit")
	wsc := moqwordsearch.NewWordSearchClient()
	moqwordsearch.Start(wsc)
}

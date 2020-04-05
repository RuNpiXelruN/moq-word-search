package main

import (
	"fmt"

	moqwordsearch "github.com/RuNpiXelruN/moq-word-search"
)

func main() {
	fmt.Println("hit")
	searchItemClient := moqwordsearch.NewSearchItemClient()
	wsc := moqwordsearch.NewWordSearchClient(searchItemClient)
	moqwordsearch.Start(wsc)
}

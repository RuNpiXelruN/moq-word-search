package main

import (
	"fmt"

	moqwordsearch "github.com/RuNpiXelruN/moq-word-search"
)

func main() {
	fmt.Println("gRPC word search started.. \nREST API Listening on port ':8090'\n\nRun 'make help' to list available commands")
	searchItemClient := moqwordsearch.NewSearchItemClient()
	wsc := moqwordsearch.NewWordSearchClient(searchItemClient)
	moqwordsearch.Start(wsc)
}

package moqwordsearch

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	wsproto "github.com/RuNpiXelruN/moq-word-search/proto"
)

var errNoTerm = errors.New("You must provide a search term")

// IncrementCount func
func (sic *SearchItemClient) IncrementCount(item *wsproto.SearchItem) {
	item.SearchCount++
}

// WordExists func
func (sic *SearchItemClient) WordExists(searchTerm string, items []*wsproto.SearchItem, increment bool) (exists bool) {
	exists = false
	for _, item := range items {
		if item.Term == searchTerm {
			if increment {
				sic.IncrementCount(item)
			}
			exists = true
			break
		}
	}

	return exists
}

// UpdateWordList func
func (wsc *WordSearchClient) UpdateWordList(ctx context.Context, req *wsproto.UpdateWordListRequest) (*wsproto.UpdateWordListResponse, error) {
	termRaw := req.Term
	term := strings.ToLower(termRaw)

	exists := wsc.searchItemService.WordExists(term, searchList, false)
	if exists {
		return &wsproto.UpdateWordListResponse{
			Message:  fmt.Sprintf("Search term '%v', is already on the list.", termRaw),
			WordList: searchList,
		}, nil
	}

	newSearchItem := &wsproto.SearchItem{
		Term:        term,
		SearchCount: 0,
	}

	searchList = append(searchList, newSearchItem)

	return &wsproto.UpdateWordListResponse{
		Message:  fmt.Sprintf("New search term '%v', has been added to the list :)", termRaw),
		WordList: searchList,
	}, nil
}

// TopFiveSearch func
func (wsc *WordSearchClient) TopFiveSearch(ctx context.Context, req *wsproto.TopFiveRequest) (*wsproto.TopFiveResponse, error) {
	sort.Slice(searchList, func(i, j int) bool {
		return searchList[i].SearchCount > searchList[j].SearchCount
	})

	topList := searchList[:5]
	return &wsproto.TopFiveResponse{
		TopFive: topList,
	}, nil
}

// SingleWordSearch func
func (wsc *WordSearchClient) SingleWordSearch(ctx context.Context, req *wsproto.SingleWordRequest) (*wsproto.SingleWordResponse, error) {
	termRaw := req.Term

	if len(termRaw) == 0 {
		return nil, errNoTerm
	}

	term := strings.ToLower(termRaw)

	message := fmt.Sprintf("Sorry, '%v' cannot be found.", termRaw)
	exists := wsc.searchItemService.WordExists(term, searchList, true)

	if exists {
		message = fmt.Sprintf("Yay, '%v' is one of our words.", termRaw)
	}
	return &wsproto.SingleWordResponse{
		Message: message,
	}, nil
}

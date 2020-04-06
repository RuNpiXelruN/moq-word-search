package moqwordsearch

import (
	"context"

	wsproto "github.com/RuNpiXelruN/moq-word-search/proto"
)

//go:generate moq -out mocks/mock_wordsearchservice.go -pkg mocks . WordSearchService

// WordSearchService type
type WordSearchService interface {
	StartGRPC() error
	StartREST() error
	SingleWordSearch(ctx context.Context, req *wsproto.SingleWordRequest) (*wsproto.SingleWordResponse, error)
	TopFiveSearch(ctx context.Context, req *wsproto.TopFiveRequest) (*wsproto.TopFiveResponse, error)
	UpdateWordList(ctx context.Context, req *wsproto.UpdateWordListRequest) (*wsproto.UpdateWordListResponse, error)
}

// WordSearchClient type
type WordSearchClient struct {
	searchItemService SearchItemService
}

// NewWordSearchClient func
func NewWordSearchClient(sis SearchItemService) *WordSearchClient {
	return &WordSearchClient{
		searchItemService: sis,
	}
}

//go:generate moq -out mocks/mock_searchitemservice.go -pkg mocks . SearchItemService

// SearchItemService type
type SearchItemService interface {
	WordExists(searchTerm string, items []*wsproto.SearchItem, increment bool) (exists bool)
	IncrementCount(item *wsproto.SearchItem)
}

// SearchItemClient type
type SearchItemClient struct{}

// NewSearchItemClient func
func NewSearchItemClient() *SearchItemClient {
	return &SearchItemClient{}
}

package moqwordsearch

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	wsproto "github.com/RuNpiXelruN/moq-word-search/proto"
)

var (
	wg         sync.WaitGroup
	searchList []*wsproto.SearchItem
)

//go:generate moq -out mocks/mock_wordsearchservice -pkg mocks . WordSearchService

// WordSearchService type
type WordSearchService interface {
	StartGRPC() error
	StartREST() error
}

//go:generate moq -out mocks/mock_searchservice -pkg mocks . SearchService

// SearchService type
type SearchService interface {
	SingleWordSearch(ctx context.Context, req *wsproto.SingleWordRequest) (*wsproto.SingleWordResponse, error)
	WordExists(searchTerm string, items []*wsproto.SearchItem, increment bool) (exists bool)
	IncrementCount(item *wsproto.SearchItem)
	TopFiveSearch(ctx context.Context, req *wsproto.TopFiveRequest) (*wsproto.TopFiveResponse, error)
	UpdateWordList(ctx context.Context, req *wsproto.UpdateWordListRequest) (*wsproto.UpdateWordListResponse, error)
}

// SearchClient type
type SearchClient struct{}

// NewSearchClient func
func NewSearchClient() *SearchClient {
	return &SearchClient{}
}

// WordSearchClient type
type WordSearchClient struct {
	ss SearchService
}

// NewWordSearchClient func
func NewWordSearchClient(ss SearchService) *WordSearchClient {
	return &WordSearchClient{
		ss: ss,
	}
}

func init() {
	searchList = []*wsproto.SearchItem{
		&wsproto.SearchItem{
			Term:        "hello",
			SearchCount: 0,
		},
		&wsproto.SearchItem{
			Term:        "goodbye",
			SearchCount: 0,
		},
		&wsproto.SearchItem{
			Term:        "simple",
			SearchCount: 0,
		},
		&wsproto.SearchItem{
			Term:        "list",
			SearchCount: 0,
		},
		&wsproto.SearchItem{
			Term:        "search",
			SearchCount: 0,
		},
		&wsproto.SearchItem{
			Term:        "filter",
			SearchCount: 0,
		},
		&wsproto.SearchItem{
			Term:        "yes",
			SearchCount: 0,
		},
		&wsproto.SearchItem{
			Term:        "no",
			SearchCount: 0,
		},
	}
}

// Start func
func Start(wss WordSearchService) {

	wg.Add(1)
	go func() {
		log.Fatal(wss.StartGRPC())
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		log.Fatal(wss.StartREST())
		wg.Done()
	}()

	wg.Wait()
}

// IncrementCount func
func (ss *SearchClient) IncrementCount(item *wsproto.SearchItem) {
	item.SearchCount++
}

// WordExists func
func (ss *SearchClient) WordExists(searchTerm string, items []*wsproto.SearchItem, increment bool) (exists bool) {
	exists = false
	for _, item := range items {
		if item.Term == searchTerm {
			if increment {
				ss.IncrementCount(item)
			}
			exists = true
			break
		}
	}

	return exists
}

// UpdateWordList func
func (ss *SearchClient) UpdateWordList(ctx context.Context, req *wsproto.UpdateWordListRequest) (*wsproto.UpdateWordListResponse, error) {
	termRaw := req.Term
	term := strings.ToLower(termRaw)

	exists := ss.WordExists(term, searchList, false)
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
func (ss *SearchClient) TopFiveSearch(ctx context.Context, req *wsproto.TopFiveRequest) (*wsproto.TopFiveResponse, error) {
	sort.Slice(searchList, func(i, j int) bool {
		return searchList[i].SearchCount > searchList[j].SearchCount
	})

	topList := searchList[:5]
	return &wsproto.TopFiveResponse{
		TopFive: topList,
	}, nil
}

// SingleWordSearch func
func (ss *SearchClient) SingleWordSearch(ctx context.Context, req *wsproto.SingleWordRequest) (*wsproto.SingleWordResponse, error) {
	termRaw := req.Term
	term := strings.ToLower(termRaw)

	message := fmt.Sprintf("Sorry, '%v' cannot be found.", termRaw)
	exists := ss.WordExists(term, searchList, true)

	if exists {
		message = fmt.Sprintf("Yay, '%v' is one of our words.", termRaw)
	}
	return &wsproto.SingleWordResponse{
		Message: message,
	}, nil
}

// StartGRPC func
func (wsc *WordSearchClient) StartGRPC() error {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	wsproto.RegisterWordSearchServer(grpcServer, wsc)

	grpcServer.Serve(lis)
	return nil
}

// StartREST func
func (wsc *WordSearchClient) StartREST() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := wsproto.RegisterWordSearchHandlerFromEndpoint(ctx, mux, ":8080", opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":8090", mux)
}

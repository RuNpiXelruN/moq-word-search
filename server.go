package moqwordsearch

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	wsproto "github.com/RuNpiXelruN/moq-word-search/proto"
)

var wg sync.WaitGroup

// SearchList is the default slice of search terms
var SearchList []*wsproto.SearchItem

func init() {
	SearchList = []*wsproto.SearchItem{
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

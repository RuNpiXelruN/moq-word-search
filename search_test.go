package moqwordsearch

import (
	"context"
	"testing"

	"github.com/RuNpiXelruN/moq-word-search/mocks"
	wsproto "github.com/RuNpiXelruN/moq-word-search/proto"
	"github.com/stretchr/testify/assert"
)

func TestSingleWordSearch(t *testing.T) {

	t.Run("An error should be returned if an empty search term is searched", func(t *testing.T) {

		mockSearchItemClient := &mocks.SearchItemServiceMock{}

		req := &wsproto.SingleWordRequest{
			Term: "",
		}

		wsc := NewWordSearchClient(mockSearchItemClient)
		resp, err := wsc.SingleWordSearch(context.Background(), req)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.IsType(t, errNoTerm, err)
	})

	t.Run("searchItemService.WordExists should be called", func(t *testing.T) {
		mockSearchItemClient := &mocks.SearchItemServiceMock{
			WordExistsFunc: func(searchTerm string, items []*wsproto.SearchItem, increment bool) bool {
				return true
			},
		}

		req := &wsproto.SingleWordRequest{
			Term: "sawyer",
		}

		wsc := NewWordSearchClient(mockSearchItemClient)
		resp, err := wsc.SingleWordSearch(context.Background(), req)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, mockSearchItemClient.WordExistsCalls(), 1)
	})

	t.Run("should return a successful/unsuccessful search message if word found/not found", func(t *testing.T) {
		succMockClient := &mocks.SearchItemServiceMock{
			WordExistsFunc: func(searchTerm string, items []*wsproto.SearchItem, increment bool) bool {
				return true
			},
		}

		failMockClient := &mocks.SearchItemServiceMock{
			WordExistsFunc: func(searchTerm string, items []*wsproto.SearchItem, increment bool) bool {
				return false
			},
		}

		testCases := []struct {
			testName             string
			mockSearchItemClient *mocks.SearchItemServiceMock
			respMessageSnippet   string
		}{
			{
				testName:             "a successful search message should be returned if found",
				mockSearchItemClient: succMockClient,
				respMessageSnippet:   "is one of our words",
			},
			{
				testName:             "an unsuccessful search message should be returned if not found",
				mockSearchItemClient: failMockClient,
				respMessageSnippet:   "cannot be found",
			},
		}

		req := &wsproto.SingleWordRequest{
			Term: "sawyer",
		}

		for _, tc := range testCases {
			t.Run(tc.testName, func(t *testing.T) {
				wsc := NewWordSearchClient(tc.mockSearchItemClient)
				resp, err := wsc.SingleWordSearch(context.Background(), req)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Contains(t, resp.Message, tc.respMessageSnippet)
			})
		}
	})
}

func TestTopFiveSearch(t *testing.T) {

	t.Run("Length of TopFiveSearch response should be 5", func(t *testing.T) {

		mockSearchItemClient := &mocks.SearchItemServiceMock{}

		req := &wsproto.TopFiveRequest{}
		wsc := NewWordSearchClient(mockSearchItemClient)
		resp, err := wsc.TopFiveSearch(context.Background(), req)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.TopFive, 5)
	})
}

func TestUpdateWordList(t *testing.T) {

	t.Run("An error should be returned if an empty term is sent for update", func(t *testing.T) {
		mockSearchItemClient := &mocks.SearchItemServiceMock{}
		req := &wsproto.UpdateWordListRequest{
			Term: "",
		}
		wsc := NewWordSearchClient(mockSearchItemClient)
		resp, err := wsc.UpdateWordList(context.Background(), req)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Equal(t, errNoTerm, err)
	})

	t.Run("searchItemService.WordExists should be called", func(t *testing.T) {
		mockSearchItemClient := &mocks.SearchItemServiceMock{
			WordExistsFunc: func(searchTerm string, items []*wsproto.SearchItem, increment bool) bool {
				return true
			},
		}
		req := &wsproto.UpdateWordListRequest{
			Term: "sawyer",
		}
		wsc := NewWordSearchClient(mockSearchItemClient)
		resp, err := wsc.UpdateWordList(context.Background(), req)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, mockSearchItemClient.WordExistsCalls(), 1)
	})

	t.Run("word should be added/not added if is/is not on list", func(t *testing.T) {
		existsSearchItemClient := &mocks.SearchItemServiceMock{
			WordExistsFunc: func(searchTerm string, items []*wsproto.SearchItem, increment bool) bool {
				return true
			},
		}

		notExistSearchItemClient := &mocks.SearchItemServiceMock{
			WordExistsFunc: func(searchTerm string, items []*wsproto.SearchItem, increment bool) bool {
				return false
			},
		}

		testCases := []struct {
			testName             string
			mockSearchItemClient *mocks.SearchItemServiceMock
			respMessageSnippet   string
		}{
			{
				testName:             "word should be added if not on list",
				mockSearchItemClient: notExistSearchItemClient,
				respMessageSnippet:   "has been added to the list",
			},

			{
				testName:             "word should be not be added if on list",
				mockSearchItemClient: existsSearchItemClient,
				respMessageSnippet:   "is already on the list",
			},
		}

		req := &wsproto.UpdateWordListRequest{
			Term: "sawyer",
		}

		for _, tc := range testCases {
			t.Run(tc.testName, func(t *testing.T) {
				wsc := NewWordSearchClient(tc.mockSearchItemClient)
				resp, err := wsc.UpdateWordList(context.Background(), req)

				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Contains(t, resp.Message, tc.respMessageSnippet)
			})
		}
	})
}

func TestWordExists(t *testing.T) {

	testCases := []struct {
		testName   string
		searchTerm string
		expected   bool
	}{
		{
			testName:   "If word is found should return true",
			searchTerm: "sawyer",
			expected:   true,
		},
		{
			testName:   "If word is not found should return false",
			searchTerm: "brooks",
			expected:   false,
		},
	}
	searchTerms := []*wsproto.SearchItem{
		&wsproto.SearchItem{
			Term:        "sawyer",
			SearchCount: 0,
		},
	}
	searchItemClient := NewSearchItemClient()

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			exists := searchItemClient.WordExists(tc.searchTerm, searchTerms, true)
			assert.Equal(t, tc.expected, exists)
		})
	}
}

func TestIncrement(t *testing.T) {

	t.Run("item.SearchCount should increment", func(t *testing.T) {

		item := &wsproto.SearchItem{
			Term:        "sawyer",
			SearchCount: 0,
		}

		searchItemClient := NewSearchItemClient()
		searchItemClient.IncrementCount(item)

		assert.Equal(t, int64(1), item.SearchCount)

	})
}

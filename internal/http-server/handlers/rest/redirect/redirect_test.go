package redirect_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/lib/consts"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/rest/redirect"
	"url-shortener/internal/http-server/handlers/rest/redirect/mocks"
	"url-shortener/internal/lib/api"
	slogdiscard "url-shortener/internal/lib/logger/slogimpl"
)

type testStruct struct {
	name      string
	alias     string
	url       string
	respError string
	mockError error
}

// TODO: add more tests
func TestSaveHandler(t *testing.T) {
	cases := []*testStruct{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			r := arrange(t, tc)
			ts := httptest.NewServer(r)
			defer ts.Close()

			// act
			redirectedToURL, err := api.GetUrlAfterRedirect(ts.URL + "/" + tc.alias)

			// assert
			require.NoError(t, err)

			assert.Equal(t, tc.url, redirectedToURL)
		})
	}
}

func arrange(t *testing.T, tc *testStruct) *http.ServeMux {
	urlGetterMock := mocks.NewUrlGetter(t)

	if tc.respError == "" || tc.mockError != nil {
		urlGetterMock.On("GetUrl", tc.alias).
			Return(tc.url, tc.mockError).Once()
	}

	r := http.NewServeMux()
	r.HandleFunc(fmt.Sprintf("GET /{%s}", consts.AliasKey), redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))
	return r
}

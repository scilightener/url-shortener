package save_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"url-shortener/internal/http-server/handlers/rest/save"
	"url-shortener/internal/http-server/handlers/rest/save/mocks"
	"url-shortener/internal/lib/bl"
	slogdiscard "url-shortener/internal/lib/logger/slogimpl"
)

type testStruct struct {
	name  string
	alias string
	url   string

	respValidUntil time.Time
	respStatusCode int
	respError      string
	mockError      error
}

// TODO: add more tests
func TestSaveHandler(t *testing.T) {
	cases := []*testStruct{
		{
			name:           "Success",
			alias:          "test_alias",
			url:            "https://google.com",
			respValidUntil: time.Now().Add(bl.ValidDuration).UTC(),
			respStatusCode: http.StatusCreated,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			handler, rr, req := arrange(t, tc)

			// act
			handler.ServeHTTP(rr, req)

			// assert
			assert.Equal(t, rr.Code, tc.respStatusCode)

			body, resp := rr.Body.String(), new(save.Response)
			require.NoError(t, json.Unmarshal([]byte(body), resp))

			require.Equal(t, tc.respError, resp.Error)

			assert.WithinDuration(t, tc.respValidUntil, resp.ValidUntilUTC, time.Second)
		})
	}
}

func arrange(t *testing.T, tc *testStruct) (handler http.HandlerFunc, rr *httptest.ResponseRecorder, req *http.Request) {
	t.Parallel()
	urlSaverMock := mocks.NewUrlRepo(t)
	if tc.respError == "" || tc.mockError != nil {
		urlSaverMock.On("SaveUrl", mock.AnythingOfType("string"), tc.url, mock.AnythingOfType("time.Time")).
			Return(int64(1), tc.mockError).
			Once()
	}

	handler = save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)
	input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)
	rr = httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/api/url", bytes.NewReader([]byte(input)))
	require.NoError(t, err)
	return
}

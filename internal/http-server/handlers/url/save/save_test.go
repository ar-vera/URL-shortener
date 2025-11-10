package save_test

import (
	"URL-shortener/internal/http-server/handlers/url/save"
	"URL-shortener/internal/http-server/handlers/url/save/mocks"
	"URL-shortener/internal/lib/logger/handlers/slogdiscard"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
		setupMock bool
	}{
		{
			name:      "Success",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "",
			setupMock: true,
		},
		{
			name:      "Empty alias",
			alias:     "",
			url:       "https://google.com",
			respError: "",
			setupMock: true,
		},
		{
			name:      "Empty URL",
			url:       "",
			alias:     "some_alias",
			respError: "Field 'URL' is required",
			setupMock: false,
		},
		{
			name:      "Invalid URL",
			url:       "Some invalid URL",
			alias:     "some_alias",
			respError: "Field URL is not a valid URL",
			setupMock: false,
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "Failed to save URL",
			mockError: errors.New("database error"),
			setupMock: true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.setupMock {
				if tc.mockError != nil {
					urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
						Return(int64(0), tc.mockError).Once()
				} else {
					urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
						Return(int64(1), nil).Once()
				}
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			if tc.respError == "" {
				require.Empty(t, resp.Error)
			} else {
				require.Contains(t, resp.Error, tc.respError)
			}
		})
	}
}

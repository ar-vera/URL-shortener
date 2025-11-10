package redirect_test

import (
	"URL-shortener/internal/http-server/handlers/redirect"
	"URL-shortener/internal/http-server/handlers/redirect/mocks"
	"URL-shortener/internal/lib/logger/handlers/slogdiscard"
	"URL-shortener/internal/storage"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name         string
		alias        string
		setupMock    func(*mocks.URLGetter)
		expectedCode int
		checkBody    bool
		expectedErr  string
	}{
		{
			name:  "Success redirect",
			alias: "test_alias",
			setupMock: func(m *mocks.URLGetter) {
				m.On("GetURL", "test_alias").Return("https://google.com", nil).Once()
			},
			expectedCode: http.StatusFound,
			checkBody:    false,
		},
		{
			name:  "URL not found",
			alias: "nonexistent",
			setupMock: func(m *mocks.URLGetter) {
				m.On("GetURL", "nonexistent").Return("", storage.ErrURLNotFound).Once()
			},
			expectedCode: http.StatusOK,
			checkBody:    true,
			expectedErr:  "URL not found",
		},
		{
			name:  "Internal error",
			alias: "test_error",
			setupMock: func(m *mocks.URLGetter) {
				m.On("GetURL", "test_error").Return("", errors.New("database error")).Once()
			},
			expectedCode: http.StatusOK,
			checkBody:    true,
			expectedErr:  "Internal error",
		},
		{
			name:  "Empty alias in URL param",
			alias: " ", // Space to test empty after trim, but chi will treat this as valid param
			setupMock: func(m *mocks.URLGetter) {
				// This will call GetURL with space, which should fail validation or return error
				m.On("GetURL", " ").Return("", storage.ErrURLNotFound).Maybe()
			},
			expectedCode: http.StatusOK,
			checkBody:    true,
			expectedErr:  "URL not found",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlGetterMock := mocks.NewURLGetter(t)
			tc.setupMock(urlGetterMock)

			handler := redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock)

			r := chi.NewRouter()
			r.Get("/{alias}", handler)

			req, err := http.NewRequest(http.MethodGet, "/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code)

			if tc.checkBody {
				var resp struct {
					Status string `json:"status"`
					Error  string `json:"error,omitempty"`
				}

				require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
				require.Contains(t, resp.Error, tc.expectedErr)
			} else {
				// For redirect, check Location header
				location := rr.Header().Get("Location")
				require.NotEmpty(t, location)
			}
		})
	}
}

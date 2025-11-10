package delete_test

import (
	"URL-shortener/internal/http-server/handlers/delete"
	"URL-shortener/internal/http-server/handlers/delete/mocks"
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

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name         string
		alias        string
		setupMock    func(*mocks.URLDeleter)
		expectedCode int
		checkBody    bool
		expectedErr  string
	}{
		{
			name:  "Success delete",
			alias: "test_alias",
			setupMock: func(m *mocks.URLDeleter) {
				m.On("DeleteURL", "test_alias").Return(nil).Once()
			},
			expectedCode: http.StatusOK,
			checkBody:    true,
			expectedErr:  "",
		},
		{
			name:  "URL not found",
			alias: "nonexistent",
			setupMock: func(m *mocks.URLDeleter) {
				m.On("DeleteURL", "nonexistent").Return(storage.ErrURLNotFound).Once()
			},
			expectedCode: http.StatusOK,
			checkBody:    true,
			expectedErr:  "URL not found",
		},
		{
			name:  "Internal error",
			alias: "test_error",
			setupMock: func(m *mocks.URLDeleter) {
				m.On("DeleteURL", "test_error").Return(errors.New("database error")).Once()
			},
			expectedCode: http.StatusOK,
			checkBody:    true,
			expectedErr:  "Failed to delete URL",
		},
		{
			name:  "Empty alias in URL param",
			alias: " ",
			setupMock: func(m *mocks.URLDeleter) {
				m.On("DeleteURL", " ").Return(storage.ErrURLNotFound).Maybe()
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

			urlDeleterMock := mocks.NewURLDeleter(t)
			tc.setupMock(urlDeleterMock)

			handler := delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock)

			r := chi.NewRouter()
			r.Delete("/url/{alias}", handler)

			req, err := http.NewRequest(http.MethodDelete, "/url/"+tc.alias, nil)
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

				if tc.expectedErr == "" {
					require.Equal(t, "OK", resp.Status)
					require.Empty(t, resp.Error)
				} else {
					require.Contains(t, resp.Error, tc.expectedErr)
				}
			}
		})
	}
}

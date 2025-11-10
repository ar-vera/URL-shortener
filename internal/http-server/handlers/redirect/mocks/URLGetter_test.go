package mocks

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLGetter_GetURL(t *testing.T) {
	tests := []struct {
		name        string
		alias       string
		setupMock   func(*URLGetter)
		expectedURL string
		expectedErr error
	}{
		{
			name:  "successful get",
			alias: "test123",
			setupMock: func(m *URLGetter) {
				m.On("GetURL", "test123").Return("https://example.com", nil)
			},
			expectedURL: "https://example.com",
			expectedErr: nil,
		},
		{
			name:  "URL not found",
			alias: "nonexistent",
			setupMock: func(m *URLGetter) {
				m.On("GetURL", "nonexistent").Return("", errors.New("url not found"))
			},
			expectedURL: "",
			expectedErr: errors.New("url not found"),
		},
		{
			name:  "internal error",
			alias: "test456",
			setupMock: func(m *URLGetter) {
				m.On("GetURL", "test456").Return("", errors.New("database error"))
			},
			expectedURL: "",
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockURLGetter := NewURLGetter(t)
			tt.setupMock(mockURLGetter)

			url, err := mockURLGetter.GetURL(tt.alias)

			assert.Equal(t, tt.expectedURL, url)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockURLGetter.AssertExpectations(t)
		})
	}
}

func TestNewURLGetter(t *testing.T) {
	mock := NewURLGetter(t)
	assert.NotNil(t, mock)
}

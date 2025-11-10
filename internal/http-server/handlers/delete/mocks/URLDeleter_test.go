package mocks

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLDeleter_DeleteURL(t *testing.T) {
	tests := []struct {
		name        string
		alias       string
		setupMock   func(*URLDeleter)
		expectedErr error
	}{
		{
			name:  "successful delete",
			alias: "test123",
			setupMock: func(m *URLDeleter) {
				m.On("DeleteURL", "test123").Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:  "URL not found",
			alias: "nonexistent",
			setupMock: func(m *URLDeleter) {
				m.On("DeleteURL", "nonexistent").Return(errors.New("url not found"))
			},
			expectedErr: errors.New("url not found"),
		},
		{
			name:  "internal error",
			alias: "test456",
			setupMock: func(m *URLDeleter) {
				m.On("DeleteURL", "test456").Return(errors.New("database error"))
			},
			expectedErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockURLDeleter := NewURLDeleter(t)
			tt.setupMock(mockURLDeleter)

			err := mockURLDeleter.DeleteURL(tt.alias)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockURLDeleter.AssertExpectations(t)
		})
	}
}

func TestNewURLDeleter(t *testing.T) {
	mock := NewURLDeleter(t)
	assert.NotNil(t, mock)
}

package mocks

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLSaver_SaveURL(t *testing.T) {
	tests := []struct {
		name        string
		urlToSave   string
		alias       string
		setupMock   func(*URLSaver)
		expectedID  int64
		expectedErr error
	}{
		{
			name:      "successful save",
			urlToSave: "https://example.com",
			alias:     "test123",
			setupMock: func(m *URLSaver) {
				m.On("SaveURL", "https://example.com", "test123").Return(int64(1), nil)
			},
			expectedID:  1,
			expectedErr: nil,
		},
		{
			name:      "save with error",
			urlToSave: "https://example.com",
			alias:     "duplicate",
			setupMock: func(m *URLSaver) {
				m.On("SaveURL", "https://example.com", "duplicate").Return(int64(0), errors.New("unique violation"))
			},
			expectedID:  0,
			expectedErr: errors.New("unique violation"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockURLSaver := NewURLSaver(t)
			tt.setupMock(mockURLSaver)

			id, err := mockURLSaver.SaveURL(tt.urlToSave, tt.alias)

			assert.Equal(t, tt.expectedID, id)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockURLSaver.AssertExpectations(t)
		})
	}
}

func TestNewURLSaver(t *testing.T) {
	mock := NewURLSaver(t)
	assert.NotNil(t, mock)
}

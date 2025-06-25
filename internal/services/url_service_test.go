package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/VladimirAzanza/url-shortener/config"
	"github.com/VladimirAzanza/url-shortener/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func getTestConfig() *config.Config {
	return &config.Config{StorageType: "postgres"}
}

func setupTestService(t *testing.T) (*URLService, *mocks.MockIURLRepository, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	// defer ctrl.Finish() // This will be used at every test

	mockRepo := mocks.NewMockIURLRepository(ctrl)
	service := NewURLService(getTestConfig(), mockRepo).(*URLService)
	return service, mockRepo, ctrl
}

func TestShortenURL(t *testing.T) {
	tests := []struct {
		name                string
		originalURL         string
		getRepoReturnsShort string
		getRepoReturnsErr   error
		saveRepoReturnsErr  error
		expectedError       error
		expectRepoSaveCall  bool
		expectedShortIDLen  int
	}{
		{
			name:                "Valid URL - NEW",
			originalURL:         "https://example.com",
			getRepoReturnsShort: "",
			getRepoReturnsErr:   sql.ErrNoRows,
			saveRepoReturnsErr:  nil,
			expectedError:       nil,
			expectRepoSaveCall:  true,
			expectedShortIDLen:  16,
		},
		{
			name:                "valid URL - Exists",
			originalURL:         "https://example.com/existing",
			getRepoReturnsShort: "existingShortID",
			getRepoReturnsErr:   nil,
			saveRepoReturnsErr:  nil,
			expectedError:       nil,
			expectRepoSaveCall:  false,
			expectedShortIDLen:  15,
		},
		{
			name:                "Error at repo",
			originalURL:         "https://example.com/error_get",
			getRepoReturnsShort: "",
			getRepoReturnsErr:   errors.New("db error"),
			saveRepoReturnsErr:  nil,
			expectedError:       errors.New("error checking existing URL: db error"),
			expectRepoSaveCall:  false,
			expectedShortIDLen:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo, ctrl := setupTestService(t)
			defer ctrl.Finish()

			ctx := context.Background()

			mockRepo.EXPECT().
				GetShortIDByOriginalURL(ctx, tt.originalURL).
				Return(tt.getRepoReturnsShort, tt.getRepoReturnsErr).
				Times(1)

			if tt.expectRepoSaveCall {
				mockRepo.EXPECT().
					SaveShortID(ctx, gomock.Any(), tt.originalURL).
					Return(tt.saveRepoReturnsErr).
					Times(1)
			}

			shortID, err := s.ShortenURL(ctx, tt.originalURL)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Empty(t, shortID)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, shortID)
				if tt.expectRepoSaveCall {
					assert.Len(t, shortID, tt.expectedShortIDLen)
				} else {
					assert.Equal(t, tt.getRepoReturnsShort, shortID)
					assert.Len(t, shortID, tt.expectedShortIDLen)
				}
			}
		})
	}
}

// func TestShortenAPIURL(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		originalURL string
// 	}{
// 		{"Valid URL", "https://example.com"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := NewURLService(getTestConfig())
// 			req := &dto.ShortenRequestDTO{URL: tt.originalURL}
// 			shortID, err := s.ShortenAPIURL(context.Background(), req)

// 			assert.NotEmpty(t, shortID)
// 			assert.Len(t, shortID, 16)

// 			originalURL, exists := s.GetOriginalURL(context.Background(), shortID)
// 			assert.True(t, exists)
// 			assert.Equal(t, tt.originalURL, originalURL)
// 		})
// 	}
// }

// func TestGetOriginalURL(t *testing.T) {
// 	s := NewURLService(getTestConfig())
// 	testURL := "https://example.com"
// 	shortID, err := s.ShortenURL(context.Background(), testURL)

// 	tests := []struct {
// 		name     string
// 		shortID  string
// 		expected string
// 		exists   bool
// 	}{
// 		{"Existing URL", shortID, testURL, true},
// 		{"Non Existin URL", "", "", false},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			originalURL, exists := s.GetOriginalURL(context.Background(), tt.shortID)
// 			assert.Equal(t, tt.exists, exists)
// 			if exists {
// 				assert.Equal(t, tt.expected, originalURL)
// 			}
// 		})
// 	}
// }

// func TestGetOriginalAPIURL(t *testing.T) {
// 	s := NewURLService(getTestConfig())
// 	testURL := "https://example.com"
// 	req := &dto.ShortenRequestDTO{URL: testURL}
// 	shortID, err := s.ShortenAPIURL(context.Background(), req)

// 	tests := []struct {
// 		name     string
// 		shortID  string
// 		expected string
// 		exists   bool
// 	}{
// 		{"Existing URL", shortID, testURL, true},
// 		{"Non Existin URL", "", "", false},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			originalURL, exists := s.GetOriginalURL(context.Background(), tt.shortID)
// 			assert.Equal(t, tt.exists, exists)
// 			if exists {
// 				assert.Equal(t, tt.expected, originalURL)
// 			}
// 		})
// 	}
// }

// func TestGenerateUniqueId(t *testing.T) {
// 	fURL := "https://example.com"
// 	sURL := "https://google.com"

// 	id1 := generateUniqueID(fURL)
// 	id2 := generateUniqueID(sURL)

// 	assert.Len(t, id1, 16)
// 	assert.Len(t, id2, 16)
// 	assert.NotEqual(t, id1, id2, "ID's should be unique for different URLS")
// }

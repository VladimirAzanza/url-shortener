package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/VladimirAzanza/url-shortener/config"
	"github.com/VladimirAzanza/url-shortener/internal/dto"
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

func TestShortenAPIURL(t *testing.T) {
	tests := []struct {
		name        string
		originalURL string
	}{
		{"Valid URL", "https://example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo, ctrl := setupTestService(t)
			defer ctrl.Finish()

			ctx := context.Background()
			req := &dto.ShortenRequestDTO{URL: tt.originalURL}

			mockRepo.EXPECT().
				GetShortIDByOriginalURL(ctx, tt.originalURL).
				Return("", sql.ErrNoRows).
				Times(1)

			mockRepo.EXPECT().
				SaveShortID(ctx, gomock.Any(), tt.originalURL).
				Return(nil).
				Times(1)

			shortID, err := s.ShortenAPIURL(ctx, req)

			assert.NoError(t, err)
			assert.NotEmpty(t, shortID)
			assert.Len(t, shortID, 16)
		})
	}
}

func TestGetOriginalURL(t *testing.T) {
	tests := []struct {
		name           string
		shortID        string
		originalURL    string
		exists         bool
		repoReturnsErr error
		expectedError  bool
	}{
		{
			name:           "Existing URL",
			shortID:        "abc123",
			originalURL:    "https://example.com",
			exists:         true,
			repoReturnsErr: nil,
			expectedError:  false,
		},
		{
			name:           "Non-existing URL",
			shortID:        "nonexistent",
			originalURL:    "",
			exists:         false,
			repoReturnsErr: sql.ErrNoRows,
			expectedError:  false,
		},
		{
			name:           "Repo error",
			shortID:        "error",
			originalURL:    "",
			exists:         false,
			repoReturnsErr: errors.New("db error"),
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo, ctrl := setupTestService(t)
			defer ctrl.Finish()

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			mockRepo.EXPECT().
				GetOriginalURL(ctx, tt.shortID).
				Return(tt.originalURL, tt.exists, tt.repoReturnsErr).
				Times(1)

			originalURL, exists := s.GetOriginalURL(ctx, tt.shortID)

			if tt.expectedError {
				assert.False(t, exists)
				assert.Empty(t, originalURL)
			} else {
				assert.Equal(t, tt.exists, exists)
				if exists {
					assert.Equal(t, tt.originalURL, originalURL)
				} else {
					assert.Empty(t, originalURL)
				}
			}
		})
	}
}

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

func TestBatchShortenURL(t *testing.T) {
	tests := []struct {
		name                string
		originalURL         string
		getRepoReturnsShort string
		getRepoReturnsErr   error
		saveRepoReturnsErr  error
		expectedError       error
		expectRepoSaveCall  bool
	}{
		{
			name:                "New URL",
			originalURL:         "https://example.com/new",
			getRepoReturnsShort: "",
			getRepoReturnsErr:   sql.ErrNoRows,
			saveRepoReturnsErr:  nil,
			expectedError:       nil,
			expectRepoSaveCall:  true,
		},
		{
			name:                "Existing URL",
			originalURL:         "https://example.com/existing",
			getRepoReturnsShort: "existing123",
			getRepoReturnsErr:   nil,
			saveRepoReturnsErr:  nil,
			expectedError:       nil,
			expectRepoSaveCall:  false,
		},
		{
			name:                "Repo error",
			originalURL:         "https://example.com/error",
			getRepoReturnsShort: "",
			getRepoReturnsErr:   errors.New("db error"),
			saveRepoReturnsErr:  nil,
			expectedError:       errors.New("error checking existing URL: db error"),
			expectRepoSaveCall:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo, ctrl := setupTestService(t)
			defer ctrl.Finish()

			ctx := context.Background()
			req := dto.BatchRequestDTO{OriginalURL: tt.originalURL}

			mockRepo.EXPECT().
				GetShortIDByOriginalURL(ctx, tt.originalURL).
				Return(tt.getRepoReturnsShort, tt.getRepoReturnsErr).
				Times(1)

			if tt.expectRepoSaveCall {
				mockRepo.EXPECT().
					SaveBatchURL(ctx, gomock.Any(), tt.originalURL).
					Return(tt.saveRepoReturnsErr).
					Times(1)
			}

			shortID, err := s.BatchShortenURL(ctx, req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Empty(t, shortID)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, shortID)
				if !tt.expectRepoSaveCall {
					assert.Equal(t, tt.getRepoReturnsShort, shortID)
				}
			}
		})
	}
}

func TestPingDB(t *testing.T) {
	tests := []struct {
		name          string
		repoReturns   error
		expectedError error
	}{
		{
			name:          "Successful ping",
			repoReturns:   nil,
			expectedError: nil,
		},
		{
			name:          "Failed ping",
			repoReturns:   errors.New("connection failed"),
			expectedError: errors.New("connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo, ctrl := setupTestService(t)
			defer ctrl.Finish()

			ctx := context.Background()

			mockRepo.EXPECT().
				Ping(ctx).
				Return(tt.repoReturns).
				Times(1)

			err := s.PingDB(ctx)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetStorageType(t *testing.T) {
	s, _, ctrl := setupTestService(t)
	defer ctrl.Finish()

	assert.Equal(t, "postgres", s.GetStorageType())
}

func TestBatchDeleteURLs(t *testing.T) {
	tests := []struct {
		name          string
		shortURLs     []string
		repoReturns   error
		expectedError error
	}{
		{
			name:          "Successful delete",
			shortURLs:     []string{"abc123", "def456"},
			repoReturns:   nil,
			expectedError: nil,
		},
		{
			name:          "Failed delete",
			shortURLs:     []string{"abc123", "def456"},
			repoReturns:   errors.New("delete error"),
			expectedError: errors.New("error at deleting urls: delete error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo, ctrl := setupTestService(t)
			defer ctrl.Finish()

			ctx := context.Background()

			mockRepo.EXPECT().
				BatchDeleteURLs(ctx, tt.shortURLs).
				Return(tt.repoReturns).
				Times(1)

			err := s.BatchDeleteURLs(ctx, tt.shortURLs)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConcurrentBatchDelete(t *testing.T) {
	tests := []struct {
		name          string
		shortURLs     []string
		batchReturns  []error
		expectedError error
	}{
		{
			name:          "Successful concurrent delete",
			shortURLs:     []string{"a", "b", "c", "d", "e"},
			batchReturns:  []error{nil, nil, nil},
			expectedError: nil,
		},
		{
			name:          "Failed concurrent delete",
			shortURLs:     []string{"a", "b", "c", "d", "e"},
			batchReturns:  []error{nil, errors.New("batch error"), nil},
			expectedError: errors.New("batch error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, mockRepo, ctrl := setupTestService(t)
			defer ctrl.Finish()

			ctx := context.Background()

			// Expect 3 calls (batch size of 2 for 5 items)
			mockRepo.EXPECT().
				BatchDeleteURLs(ctx, []string{"a", "b"}).
				Return(tt.batchReturns[0]).
				Times(1)

			mockRepo.EXPECT().
				BatchDeleteURLs(ctx, []string{"c", "d"}).
				Return(tt.batchReturns[1]).
				Times(1)

			mockRepo.EXPECT().
				BatchDeleteURLs(ctx, []string{"e"}).
				Return(tt.batchReturns[2]).
				Times(1)

			err := s.ConcurrentBatchDelete(ctx, tt.shortURLs)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

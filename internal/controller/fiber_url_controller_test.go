package controller

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VladimirAzanza/url-shortener/internal/constants"
	"github.com/VladimirAzanza/url-shortener/internal/dto"
	"github.com/VladimirAzanza/url-shortener/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupTestController(t *testing.T) (*FiberURLController, *mocks.MockIURLService, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockIURLService(ctrl)
	controller := NewFiberURLController(mockService)
	return controller, mockService, ctrl
}

func TestGetDBPing(t *testing.T) {
	tests := []struct {
		name           string
		serviceError   error
		storageType    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			serviceError:   nil,
			storageType:    "postgres",
			expectedStatus: fiber.StatusOK,
			expectedBody:   `"status":"success"`,
		},
		{
			name:           "DB Error",
			serviceError:   errors.New("connection failed"),
			storageType:    "",
			expectedStatus: fiber.StatusBadGateway,
			expectedBody:   `"message":"Can not connect to the Database"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, mockService, ctrl := setupTestController(t)
			defer ctrl.Finish()

			app := fiber.New()
			app.Get("/ping", controller.GetDBPing)

			mockService.EXPECT().
				PingDB(gomock.Any()).
				Return(tt.serviceError).
				Times(1)

			if tt.serviceError == nil {
				mockService.EXPECT().
					GetStorageType().
					Return(tt.storageType).
					Times(1)
			}

			req := httptest.NewRequest("GET", "/ping", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body := make([]byte, resp.ContentLength)
			resp.Body.Read(body)
			assert.Contains(t, string(body), tt.expectedBody)
		})
	}
}

func TestHandleAPIDeleteBatch(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		batchSize      int
		serviceError   error
		expectedStatus int
	}{
		{
			name:           "Single item",
			requestBody:    `["abc123"]`,
			batchSize:      1,
			serviceError:   nil,
			expectedStatus: fiber.StatusAccepted,
		},
		{
			name:           "Multiple items",
			requestBody:    `["abc123","def456"]`,
			batchSize:      2,
			serviceError:   nil,
			expectedStatus: fiber.StatusAccepted,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `invalid json`,
			batchSize:      0,
			serviceError:   nil,
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "Empty batch",
			requestBody:    `[]`,
			batchSize:      0,
			serviceError:   nil,
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "Service error",
			requestBody:    `["abc123"]`,
			batchSize:      1,
			serviceError:   errors.New("service error"),
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, mockService, ctrl := setupTestController(t)
			defer ctrl.Finish()

			app := fiber.New()
			app.Post("/api/user/urls", controller.HandleAPIDeleteBatch)

			if tt.batchSize == 1 {
				mockService.EXPECT().
					BatchDeleteURLs(gomock.Any(), []string{"abc123"}).
					Return(tt.serviceError).
					Times(1)
			} else if tt.batchSize == 2 {
				mockService.EXPECT().
					ConcurrentBatchDelete(gomock.Any(), []string{"abc123", "def456"}).
					Return(tt.serviceError).
					Times(1)
			}

			req := httptest.NewRequest("POST", "/api/user/urls", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestHandlePost(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		serviceReturns string
		serviceError   error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			body:           "https://example.com",
			serviceReturns: "abc123",
			serviceError:   nil,
			expectedStatus: fiber.StatusCreated,
			expectedBody:   "http://example.com/abc123",
		},
		{
			name:           "Service error",
			body:           "https://example.com/error",
			serviceReturns: "",
			serviceError:   errors.New("service error"),
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody:   `{"error":"service error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, mockService, ctrl := setupTestController(t)
			defer ctrl.Finish()

			app := fiber.New()
			app.Post("/", controller.HandlePost)

			if tt.body != "" {
				mockService.EXPECT().
					ShortenURL(gomock.Any(), tt.body).
					Return(tt.serviceReturns, tt.serviceError).
					Times(1)
			}

			req := httptest.NewRequest("POST", "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "text/plain")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body := make([]byte, resp.ContentLength)
			resp.Body.Read(body)
			assert.Contains(t, string(body), tt.expectedBody)
		})
	}
}

func TestHandleAPIPost(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		serviceReturns string
		serviceError   error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			requestBody:    `{"url":"https://example.com"}`,
			serviceReturns: "abc123",
			serviceError:   nil,
			expectedStatus: fiber.StatusCreated,
			expectedBody:   `"result":"http://example.com/abc123"`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `invalid json`,
			serviceReturns: "",
			serviceError:   nil,
			expectedStatus: fiber.StatusBadRequest,
			expectedBody:   constants.MsgFailedToParseBody,
		},
		{
			name:           "Service error",
			requestBody:    `{"url":"https://example.com/error"}`,
			serviceReturns: "",
			serviceError:   errors.New("service error"),
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody:   `"error":"service error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, mockService, ctrl := setupTestController(t)
			defer ctrl.Finish()

			app := fiber.New()
			app.Post("/api/shorten", controller.HandleAPIPost)

			if strings.Contains(tt.requestBody, `"url":"https://example.com"`) {
				mockService.EXPECT().
					ShortenAPIURL(gomock.Any(), &dto.ShortenRequestDTO{URL: "https://example.com"}).
					Return(tt.serviceReturns, tt.serviceError).
					Times(1)
			} else if strings.Contains(tt.requestBody, `"url":"https://example.com/error"`) {
				mockService.EXPECT().
					ShortenAPIURL(gomock.Any(), &dto.ShortenRequestDTO{URL: "https://example.com/error"}).
					Return(tt.serviceReturns, tt.serviceError).
					Times(1)
			}

			req := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body := make([]byte, resp.ContentLength)
			resp.Body.Read(body)
			assert.Contains(t, string(body), tt.expectedBody)
		})
	}
}

func TestHandleGet(t *testing.T) {
	tests := []struct {
		name           string
		shortID        string
		originalURL    string
		exists         bool
		ctxError       error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			shortID:        "abc123",
			originalURL:    "https://example.com",
			exists:         true,
			ctxError:       nil,
			expectedStatus: fiber.StatusTemporaryRedirect,
			expectedBody:   "",
		},
		{
			name:           "Not found",
			shortID:        "notfound",
			originalURL:    "",
			exists:         false,
			ctxError:       nil,
			expectedStatus: fiber.StatusNotFound,
			expectedBody:   "URL not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, mockService, ctrl := setupTestController(t)
			defer ctrl.Finish()

			app := fiber.New()
			app.Get("/:id", controller.HandleGet)

			if tt.ctxError == nil {
				mockService.EXPECT().
					GetOriginalURL(gomock.Any(), tt.shortID).
					Return(tt.originalURL, tt.exists).
					Times(1)
			}

			req := httptest.NewRequest("GET", "/"+tt.shortID, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedBody != "" {
				body := make([]byte, resp.ContentLength)
				resp.Body.Read(body)
				assert.Contains(t, string(body), tt.expectedBody)
			}

			if tt.expectedStatus == fiber.StatusTemporaryRedirect {
				assert.Equal(t, tt.originalURL, resp.Header.Get("Location"))
			}
		})
	}
}

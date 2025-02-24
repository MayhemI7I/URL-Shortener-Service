package urlhandler

import (
	"bytes"
	"context"
	"encoding/json"
	
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	
)

// Моки для интерфейсов
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Get(ctx context.Context, shortURL string) (string, error) {
	args := m.Called(ctx, shortURL)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) Save(ctx context.Context, shortURL, longURL string) error {
	args := m.Called(ctx, shortURL, longURL)
	return args.Error(0)
}

type MockURLGenerator struct {
	mock.Mock
}

func (m *MockURLGenerator) GenerateShortURL(longURL string) (string, error) {
	args := m.Called(longURL)
	return args.String(0), args.Error(1)
}

// Тестирование HandleGet
func TestHandleGet(t *testing.T) {
	mockStorage := new(MockStorage)
	mockStorage.On("Get", mock.Anything, "short").Return("http://longurl.com", nil)

	handler := &URLHandler{storage: mockStorage, urlGenerator: nil}
	req := httptest.NewRequest(http.MethodGet, "/short", nil)
	rr := httptest.NewRecorder()

	handler.HandleGet(rr, req)

	assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
	assert.Equal(t, "http://longurl.com", rr.Header().Get("Location"))
	mockStorage.AssertExpectations(t)
}

// Тестирование HandlePost
func TestHandlePost(t *testing.T) {
	mockStorage := new(MockStorage)
	mockURLGenerator := new(MockURLGenerator)
	mockURLGenerator.On("GenerateShortURL", "http://longurl.com").Return("short", nil)
	mockStorage.On("Save", mock.Anything, "short", "http://longurl.com").Return(nil)

	handler := &URLHandler{storage: mockStorage, urlGenerator: mockURLGenerator}
	body := []byte("http://longurl.com")
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.HandlePost(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "short", string(rr.Body.Bytes()))
	mockStorage.AssertExpectations(t)
	mockURLGenerator.AssertExpectations(t)
}

// Тестирование HandJsonPost
func TestHandJsonPost(t *testing.T) {
	mockStorage := new(MockStorage)
	mockURLGenerator := new(MockURLGenerator)
	mockURLGenerator.On("GenerateShortURL", "http://longurl.com").Return("short", nil)
	mockStorage.On("Save", mock.Anything, "short", "http://longurl.com").Return(nil)

	handler := &URLHandler{storage: mockStorage, urlGenerator: mockURLGenerator}
	reqBody := URLRequest{LongURL: "http://longurl.com"}
	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandJsonPost(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "short", response["result"])
	mockStorage.AssertExpectations(t)
	mockURLGenerator.AssertExpectations(t)
}

// Тестирование ошибок при запросах
func TestHandleGet_NotFound(t *testing.T) {
	mockStorage := new(MockStorage)
	mockStorage.On("Get", mock.Anything, "short").Return("", assert.AnError)

	handler := &URLHandler{storage: mockStorage, urlGenerator: nil}
	req := httptest.NewRequest(http.MethodGet, "/short", nil)
	rr := httptest.NewRecorder()

	handler.HandleGet(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockStorage.AssertExpectations(t)
}

func TestHandlePost_InvalidBody(t *testing.T) {
	mockStorage := new(MockStorage)
	mockURLGenerator := new(MockURLGenerator)

	handler := &URLHandler{storage: mockStorage, urlGenerator: mockURLGenerator}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()

	handler.HandlePost(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandJsonPost_InvalidJSON(t *testing.T) {
	handler := &URLHandler{}
	req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewReader([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandJsonPost(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Тестирование на тайм-аут
func TestHandleGet_Timeout(t *testing.T) {
	mockStorage := new(MockStorage)
	mockStorage.On("Get", mock.Anything, "short").Return("", context.DeadlineExceeded)

	handler := &URLHandler{storage: mockStorage, urlGenerator: nil}
	req := httptest.NewRequest(http.MethodGet, "/short", nil)
	rr := httptest.NewRecorder()

	handler.HandleGet(rr, req)

	assert.Equal(t, http.StatusRequestTimeout, rr.Code)
	mockStorage.AssertExpectations(t)
}

func TestHandJsonPost_InternalError(t *testing.T) {
	mockStorage := new(MockStorage)
	mockURLGenerator := new(MockURLGenerator)
	mockURLGenerator.On("GenerateShortURL", "http://longurl.com").Return("", assert.AnError)
	mockStorage.On("Save", mock.Anything, "short", "http://longurl.com").Return(nil)

	handler := &URLHandler{storage: mockStorage, urlGenerator: mockURLGenerator}
	reqBody := URLRequest{LongURL: "http://longurl.com"}
	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/json", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.HandJsonPost(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}


package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"GoApiProject/internal/domain/entity"
)

// mockStructurerService is a manual mock of the input.StructurerService interface.
type mockStructurerService struct {
	response *entity.StructureResponse
	err      error
}

func (m *mockStructurerService) Structure(_ context.Context, _ entity.StructureRequest) (*entity.StructureResponse, error) {
	return m.response, m.err
}

func TestStructureText_Success(t *testing.T) {
	mock := &mockStructurerService{
		response: &entity.StructureResponse{
			Data: map[string]interface{}{"city": "Rome"},
		},
	}
	handler := NewStructureHandler(mock)

	body := strings.NewReader(`{"raw_text":"Travel to Rome"}`)
	req := httptest.NewRequest(http.MethodPost, "/structure", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.StructureText(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var resp entity.StructureResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if resp.Data["city"] != "Rome" {
		t.Errorf("expected city=Rome, got: %v", resp.Data["city"])
	}
}

func TestStructureText_InvalidJSON(t *testing.T) {
	handler := NewStructureHandler(&mockStructurerService{})

	body := strings.NewReader("this is not json")
	req := httptest.NewRequest(http.MethodPost, "/structure", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.StructureText(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}

	var errResp errorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Error != "invalid request body" {
		t.Errorf("expected 'invalid request body', got: %s", errResp.Error)
	}
}

func TestStructureText_EmptyRawText(t *testing.T) {
	handler := NewStructureHandler(&mockStructurerService{})

	body := strings.NewReader(`{"raw_text":""}`)
	req := httptest.NewRequest(http.MethodPost, "/structure", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.StructureText(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}

	var errResp errorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	if errResp.Error != "raw_text is required" {
		t.Errorf("expected 'raw_text is required', got: %s", errResp.Error)
	}
}

func TestStructureText_ServiceError(t *testing.T) {
	mock := &mockStructurerService{
		err: errors.New("all 3 attempts failed"),
	}
	handler := NewStructureHandler(mock)

	body := strings.NewReader(`{"raw_text":"some text"}`)
	req := httptest.NewRequest(http.MethodPost, "/structure", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.StructureText(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", rec.Code)
	}

	var errResp errorResponse
	json.NewDecoder(rec.Body).Decode(&errResp)
	if !strings.Contains(errResp.Error, "all 3 attempts failed") {
		t.Errorf("expected error to contain 'all 3 attempts failed', got: %s", errResp.Error)
	}
}

func TestStructureText_ResponseContentType(t *testing.T) {
	mock := &mockStructurerService{
		response: &entity.StructureResponse{
			Data: map[string]interface{}{"key": "val"},
		},
	}
	handler := NewStructureHandler(mock)

	body := strings.NewReader(`{"raw_text":"text"}`)
	req := httptest.NewRequest(http.MethodPost, "/structure", body)
	rec := httptest.NewRecorder()

	handler.StructureText(rec, req)

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got: %s", ct)
	}
}

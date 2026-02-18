package application

import (
	"context"
	"errors"
	"strings"
	"testing"

	"GoApiProject/internal/domain/entity"
)

// mockLLMClient is a manual mock of the output.LLMClient interface.
type mockLLMClient struct {
	responses []string // sequential responses to return
	errs      []error  // sequential errors to return
	callCount int      // tracks how many times GenerateStructuredJSON was called
}

func (m *mockLLMClient) GenerateStructuredJSON(_ context.Context, _ string) (string, error) {
	idx := m.callCount
	m.callCount++

	var resp string
	var err error

	if idx < len(m.responses) {
		resp = m.responses[idx]
	}
	if idx < len(m.errs) {
		err = m.errs[idx]
	}
	return resp, err
}

// --- Structure method tests ---

func TestStructure_Success_FirstAttempt(t *testing.T) {
	mock := &mockLLMClient{
		responses: []string{`{"name":"John","age":30}`},
		errs:      []error{nil},
	}
	svc := NewStructurerService(mock)

	req := entity.StructureRequest{RawText: "My name is John and I am 30 years old"}
	resp, err := svc.Structure(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.Data["name"] != "John" {
		t.Errorf("expected name=John, got: %v", resp.Data["name"])
	}
	if mock.callCount != 1 {
		t.Errorf("expected 1 LLM call, got %d", mock.callCount)
	}
}

func TestStructure_EmptyRawText(t *testing.T) {
	mock := &mockLLMClient{}
	svc := NewStructurerService(mock)

	req := entity.StructureRequest{RawText: ""}
	_, err := svc.Structure(context.Background(), req)
	if err == nil {
		t.Fatal("expected validation error for empty raw_text")
	}
	if mock.callCount != 0 {
		t.Errorf("expected 0 LLM calls for invalid request, got %d", mock.callCount)
	}
}

func TestStructure_RetryOnLLMError_ThenSuccess(t *testing.T) {
	mock := &mockLLMClient{
		responses: []string{"", `{"status":"ok"}`},
		errs:      []error{errors.New("network timeout"), nil},
	}
	svc := NewStructurerService(mock)

	req := entity.StructureRequest{RawText: "some text"}
	resp, err := svc.Structure(context.Background(), req)
	if err != nil {
		t.Fatalf("expected success on retry, got: %v", err)
	}
	if resp.Data["status"] != "ok" {
		t.Errorf("expected status=ok, got: %v", resp.Data["status"])
	}
	if mock.callCount != 2 {
		t.Errorf("expected 2 LLM calls, got %d", mock.callCount)
	}
}

func TestStructure_AllRetriesFail(t *testing.T) {
	mock := &mockLLMClient{
		responses: []string{"", "", ""},
		errs:      []error{errors.New("err1"), errors.New("err2"), errors.New("err3")},
	}
	svc := NewStructurerService(mock)

	req := entity.StructureRequest{RawText: "some text"}
	_, err := svc.Structure(context.Background(), req)
	if err == nil {
		t.Fatal("expected error after all retries exhausted")
	}
	if !strings.Contains(err.Error(), "all 3 attempts failed") {
		t.Errorf("expected 'all 3 attempts failed' in error, got: %v", err)
	}
	if mock.callCount != 3 {
		t.Errorf("expected 3 LLM calls, got %d", mock.callCount)
	}
}

func TestStructure_LLMReturnsEmptyJSON(t *testing.T) {
	mock := &mockLLMClient{
		responses: []string{`{}`, `{}`, `{}`},
		errs:      []error{nil, nil, nil},
	}
	svc := NewStructurerService(mock)

	req := entity.StructureRequest{RawText: "some text"}
	_, err := svc.Structure(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty JSON response")
	}
}

func TestStructure_RetryOnInvalidJSON_ThenSuccess(t *testing.T) {
	mock := &mockLLMClient{
		responses: []string{"not valid json", `{"key":"value"}`},
		errs:      []error{nil, nil},
	}
	svc := NewStructurerService(mock)

	req := entity.StructureRequest{RawText: "some text"}
	resp, err := svc.Structure(context.Background(), req)
	if err != nil {
		t.Fatalf("expected success on second attempt, got: %v", err)
	}
	if resp.Data["key"] != "value" {
		t.Errorf("expected key=value, got: %v", resp.Data["key"])
	}
	if mock.callCount != 2 {
		t.Errorf("expected 2 LLM calls, got %d", mock.callCount)
	}
}

// --- buildPrompt tests ---

func TestBuildPrompt_ContainsRawText(t *testing.T) {
	rawText := "Flight AZ204 from NYC to Rome"
	prompt := buildPrompt(rawText)
	if !strings.Contains(prompt, rawText) {
		t.Errorf("expected prompt to contain the raw text, got: %s", prompt)
	}
}

func TestBuildPrompt_ContainsInstructions(t *testing.T) {
	prompt := buildPrompt("any text")
	if !strings.Contains(prompt, "JSON") {
		t.Error("expected prompt to mention JSON format")
	}
	if !strings.Contains(prompt, "snake_case") {
		t.Error("expected prompt to mention snake_case convention")
	}
}

// --- parseAndValidate tests ---

func TestParseAndValidate_ValidJSON(t *testing.T) {
	resp, err := parseAndValidate(`{"name":"Alice","age":25}`)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Data["name"] != "Alice" {
		t.Errorf("expected name=Alice, got: %v", resp.Data["name"])
	}
}

func TestParseAndValidate_InvalidJSON(t *testing.T) {
	_, err := parseAndValidate("this is not json")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "invalid JSON") {
		t.Errorf("expected 'invalid JSON' in error, got: %v", err)
	}
}

func TestParseAndValidate_EmptyObject(t *testing.T) {
	_, err := parseAndValidate(`{}`)
	if err == nil {
		t.Fatal("expected validation error for empty JSON object")
	}
}

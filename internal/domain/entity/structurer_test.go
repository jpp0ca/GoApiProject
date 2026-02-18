package entity

import (
	"testing"
)

// --- StructureRequest.Validate tests ---

func TestStructureRequest_Validate_Success(t *testing.T) {
	req := StructureRequest{RawText: "some text to structure"}
	if err := req.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestStructureRequest_Validate_EmptyText(t *testing.T) {
	req := StructureRequest{RawText: ""}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for empty raw_text, got nil")
	}
	if err.Error() != "raw_text is required" {
		t.Errorf("expected 'raw_text is required', got: %v", err)
	}
}

// --- StructureResponse.Validate tests ---

func TestStructureResponse_Validate_Success(t *testing.T) {
	resp := StructureResponse{
		Data: map[string]interface{}{"key": "value"},
	}
	if err := resp.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestStructureResponse_Validate_NilData(t *testing.T) {
	resp := StructureResponse{Data: nil}
	err := resp.Validate()
	if err == nil {
		t.Fatal("expected error for nil data, got nil")
	}
}

func TestStructureResponse_Validate_EmptyMap(t *testing.T) {
	resp := StructureResponse{Data: map[string]interface{}{}}
	err := resp.Validate()
	if err == nil {
		t.Fatal("expected error for empty data map, got nil")
	}
}

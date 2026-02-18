package entity

import "fmt"

// StructureRequest represents the incoming request with raw unstructured text.
type StructureRequest struct {
	RawText string `json:"raw_text"`
}

// Validate checks that the request has a non-empty raw text.
func (r *StructureRequest) Validate() error {
	if r.RawText == "" {
		return fmt.Errorf("raw_text is required")
	}
	return nil
}

// StructureResponse wraps the dynamically inferred JSON data from the LLM.
type StructureResponse struct {
	Data map[string]interface{} `json:"data"`
}

// Validate checks that the response contains at least one extracted field.
func (r *StructureResponse) Validate() error {
	if len(r.Data) == 0 {
		return fmt.Errorf("LLM returned empty or invalid JSON structure")
	}
	return nil
}

package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// doRequest sends a POST (or other method) request to the Calculate handler
// and returns the recorded response, without needing a real HTTP server.
func doRequest(t *testing.T, method, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, "/api/v1/calculate", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	Calculate(rec, req)

	return rec
}

func TestCalculate_Success(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantResult float64
	}{
		{"addition", `{"operation":"add","operands":[4,2.5]}`, 6.5},
		{"division", `{"operation":"divide","operands":[10,2]}`, 5},
		{"square root", `{"operation":"sqrt","operands":[16]}`, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(t, http.MethodPost, tt.body)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}

			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/json")
			}

			var resp calculateResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Result != tt.wantResult {
				t.Errorf("result = %v, want %v", resp.Result, tt.wantResult)
			}
		})
	}
}

func TestCalculate_Errors(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{
			name:       "unknown operation",
			body:       `{"operation":"modulo","operands":[1,2]}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "UNKNOWN_OPERATION",
		},
		{
			name:       "wrong operand count",
			body:       `{"operation":"add","operands":[1]}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_REQUEST",
		},
		{
			name:       "malformed json",
			body:       `{"operation":"add","operands":[1,2]`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_REQUEST",
		},
		{
			name:       "division by zero",
			body:       `{"operation":"divide","operands":[1,0]}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   "DIVISION_BY_ZERO",
		},
		{
			name:       "negative square root",
			body:       `{"operation":"sqrt","operands":[-4]}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   "NEGATIVE_SQRT",
		},
		{
			name:       "invalid overflow result",
			body:       `{"operation":"multiply","operands":[1.7976931348623157e+308,2]}`,
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   "INVALID_RESULT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(t, http.MethodPost, tt.body)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var resp errorResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode error response: %v", err)
			}

			if resp.Error.Code != tt.wantCode {
				t.Errorf("error code = %q, want %q", resp.Error.Code, tt.wantCode)
			}

			if resp.Error.Message == "" {
				t.Errorf("error message should not be empty")
			}
		})
	}
}

func TestCalculate_MethodNotAllowed(t *testing.T) {
	rec := doRequest(t, http.MethodGet, "")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}

	if allow := rec.Header().Get("Allow"); allow != "POST" {
		t.Errorf("Allow header = %q, want %q", allow, "POST")
	}

	var resp errorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if resp.Error.Code != "INVALID_REQUEST" {
		t.Errorf("error code = %q, want %q", resp.Error.Code, "INVALID_REQUEST")
	}
}

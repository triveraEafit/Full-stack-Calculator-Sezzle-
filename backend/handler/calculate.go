// Package handler contains the HTTP layer for the calculator API. It is
// responsible for decoding requests, validating their shape, calling the
// calculator package to do the actual math, and encoding responses. It
// contains no arithmetic itself.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"calculator-backend/calculator"
)

// calculateRequest is the expected JSON body of POST /api/v1/calculate.
type calculateRequest struct {
	Operation string    `json:"operation"`
	Operands  []float64 `json:"operands"`
}

// calculateResponse is returned on success.
type calculateResponse struct {
	Result    float64   `json:"result"`
	Operation string    `json:"operation"`
	Operands  []float64 `json:"operands"`
}

// errorResponse is returned when a request cannot be fulfilled.
type errorResponse struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Calculate handles POST /api/v1/calculate.
func Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		writeError(w, http.StatusMethodNotAllowed, "INVALID_REQUEST", "Only POST is supported")
		return
	}

	var req calculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Request body must be valid JSON")
		return
	}

	expectedOperands, ok := operandCountFor(req.Operation)
	if !ok {
		writeError(w, http.StatusBadRequest, "UNKNOWN_OPERATION", "Unsupported operation: "+req.Operation)
		return
	}

	if len(req.Operands) != expectedOperands {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", operandCountMessage(req.Operation, expectedOperands))
		return
	}

	result, err := compute(req.Operation, req.Operands)
	if err != nil {
		status, code, message := mapCalculatorError(err)
		writeError(w, status, code, message)
		return
	}

	writeJSON(w, http.StatusOK, calculateResponse{
		Result:    result,
		Operation: req.Operation,
		Operands:  req.Operands,
	})
}

// operandCountFor returns how many operands the given operation requires,
// and whether the operation is recognized at all.
func operandCountFor(operation string) (count int, ok bool) {
	switch operation {
	case "add", "subtract", "multiply", "divide", "power":
		return 2, true
	case "sqrt":
		return 1, true
	default:
		return 0, false
	}
}

func operandCountMessage(operation string, expected int) string {
	if expected == 1 {
		return operation + " requires exactly 1 operand"
	}
	return operation + " requires exactly 2 operands"
}

// compute dispatches to the calculator package based on the operation name.
// The operand count has already been validated by the caller, so it is safe
// to index directly into operands here.
func compute(operation string, operands []float64) (float64, error) {
	switch operation {
	case "add":
		return calculator.Add(operands[0], operands[1])
	case "subtract":
		return calculator.Subtract(operands[0], operands[1])
	case "multiply":
		return calculator.Multiply(operands[0], operands[1])
	case "divide":
		return calculator.Divide(operands[0], operands[1])
	case "power":
		return calculator.Power(operands[0], operands[1])
	case "sqrt":
		return calculator.Sqrt(operands[0])
	default:
		// Unreachable: operandCountFor already rejected unknown operations
		// before compute is ever called.
		return 0, calculator.ErrInvalidResult
	}
}

// mapCalculatorError converts an error from the calculator package into an
// HTTP status and API error code/message. Unrecognized errors fall back to
// a generic message rather than leaking Go internals to the client.
func mapCalculatorError(err error) (status int, code string, message string) {
	switch {
	case errors.Is(err, calculator.ErrDivisionByZero):
		return http.StatusUnprocessableEntity, "DIVISION_BY_ZERO", "Cannot divide by zero"
	case errors.Is(err, calculator.ErrNegativeSqrt):
		return http.StatusUnprocessableEntity, "NEGATIVE_SQRT", "Cannot take the square root of a negative number"
	case errors.Is(err, calculator.ErrInvalidResult):
		return http.StatusUnprocessableEntity, "INVALID_RESULT", "Operation resulted in an invalid number"
	default:
		return http.StatusUnprocessableEntity, "INVALID_RESULT", "Operation could not be performed"
	}
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{
		Error: errorDetail{Code: code, Message: message},
	})
}

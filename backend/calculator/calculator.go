// Package calculator implements the core arithmetic operations for the
// calculator service. It has no knowledge of HTTP, JSON, or the API layer —
// it only performs math and reports meaningful errors when an operation is
// invalid.
package calculator

import (
	"errors"
	"math"
)

// Sentinel errors returned by this package. The HTTP handler layer is
// expected to inspect these (with errors.Is) and map them to appropriate
// API error responses and status codes.
var (
	ErrDivisionByZero = errors.New("division by zero")
	ErrNegativeSqrt   = errors.New("cannot take the square root of a negative number")
	ErrInvalidResult  = errors.New("operation produced an invalid result (NaN or Infinity)")
)

// Add returns a + b.
func Add(a, b float64) (float64, error) {
	return validate(a + b)
}

// Subtract returns a - b.
func Subtract(a, b float64) (float64, error) {
	return validate(a - b)
}

// Multiply returns a * b.
func Multiply(a, b float64) (float64, error) {
	return validate(a * b)
}

// Divide returns a / b. Dividing by zero is treated as a user input error,
// not a floating-point Inf/NaN result, so it is checked explicitly.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return validate(a / b)
}

// Power returns base raised to the power of exponent.
func Power(base, exponent float64) (float64, error) {
	return validate(math.Pow(base, exponent))
}

// Sqrt returns the square root of a. Negative input is a user input error,
// not a floating-point NaN result, so it is checked explicitly.
func Sqrt(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSqrt
	}
	return validate(math.Sqrt(a))
}

// validate checks that a computed result is a real, finite number.
// It centralizes the NaN/Infinity check so every operation above reports
// overflow or undefined results the same way, without repeating the check
// three times.
func validate(result float64) (float64, error) {
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrInvalidResult
	}
	return result, nil
}
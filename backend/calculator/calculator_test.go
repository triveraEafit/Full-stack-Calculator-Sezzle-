package calculator

import (
	"errors"
	"math"
	"testing"
)

// epsilon is the tolerance used when comparing float64 results. Direct
// equality (==) is unreliable for floating-point math, so tests compare
// within a small margin instead.
const epsilon = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

// assertResult is a shared helper used by every test below. It checks either
// that the expected error occurred, or that the returned value matches the
// expected result within floating-point tolerance.
func assertResult(t *testing.T, got float64, err error, wantValue float64, wantErr error) {
	t.Helper()

	if wantErr != nil {
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected error %v, got %v", wantErr, err)
		}
		return
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !almostEqual(got, wantValue) {
		t.Fatalf("got %v, want %v", got, wantValue)
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive numbers", 2, 3, 5, nil},
		{"negative numbers", -2, -3, -5, nil},
		{"mixed signs", -2, 3, 1, nil},
		{"zero operands", 0, 0, 0, nil},
		{"overflow to infinity", math.MaxFloat64, math.MaxFloat64, 0, ErrInvalidResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Add(tt.a, tt.b)
			assertResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive numbers", 5, 3, 2, nil},
		{"negative numbers", -5, -3, -2, nil},
		{"result is negative", 3, 5, -2, nil},
		{"zero operands", 0, 0, 0, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Subtract(tt.a, tt.b)
			assertResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive numbers", 4, 5, 20, nil},
		{"negative times positive", -4, 5, -20, nil},
		{"negative times negative", -4, -5, 20, nil},
		{"multiply by zero", 100, 0, 0, nil},
		{"overflow to infinity", math.MaxFloat64, 2, 0, ErrInvalidResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Multiply(tt.a, tt.b)
			assertResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive numbers", 10, 2, 5, nil},
		{"negative dividend", -10, 2, -5, nil},
		{"negative divisor", 10, -2, -5, nil},
		{"zero dividend", 0, 5, 0, nil},
		{"division by zero", 10, 0, 0, ErrDivisionByZero},
		{"zero divided by zero", 0, 0, 0, ErrDivisionByZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			assertResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}

func TestPower(t *testing.T) {
	tests := []struct {
		name           string
		base, exponent float64
		want           float64
		wantErr        error
	}{
		{"positive base and exponent", 2, 3, 8, nil},
		{"exponent of zero", 5, 0, 1, nil},
		{"negative exponent", 2, -1, 0.5, nil},
		{"base of zero", 0, 5, 0, nil},
		{"fractional exponent (square root via power)", 9, 0.5, 3, nil},
		{"overflow to infinity", math.MaxFloat64, 2, 0, ErrInvalidResult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Power(tt.base, tt.exponent)
			assertResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		name    string
		a       float64
		want    float64
		wantErr error
	}{
		{"perfect square", 16, 4, nil},
		{"non-perfect square", 2, math.Sqrt2, nil},
		{"zero", 0, 0, nil},
		{"negative number", -4, 0, ErrNegativeSqrt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sqrt(tt.a)
			assertResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}
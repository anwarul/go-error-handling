package basic

import (
	"testing"
)

func TestDivide_Success(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 10.0, 2.0, 5.0},
		{"negative dividend", -10.0, 2.0, -5.0},
		{"negative divisor", 10.0, -2.0, -5.0},
		{"both negative", -10.0, -2.0, 5.0},
		{"decimal numbers", 7.5, 2.5, 3.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)
			if err != nil {
				t.Errorf("Divide(%v, %v) returned unexpected error: %v", tt.a, tt.b, err)
			}
			if result != tt.expected {
				t.Errorf("Divide(%v, %v) = %v; want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivide_DivisionByZero(t *testing.T) {
	tests := []struct {
		name string
		a    float64
	}{
		{"positive dividend", 10.0},
		{"negative dividend", -10.0},
		{"zero dividend", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, 0)
			if err == nil {
				t.Errorf("Divide(%v, 0) expected error but got none", tt.a)
			}
			if result != 0 {
				t.Errorf("Divide(%v, 0) expected result 0 but got %v", tt.a, result)
			}
			expectedErrMsg := "division by zero"
			if err.Error() != expectedErrMsg {
				t.Errorf("Divide(%v, 0) error = %v; want %v", tt.a, err.Error(), expectedErrMsg)
			}
		})
	}
}

func TestDivide_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"divide by one", 42.0, 1.0, 42.0},
		{"divide zero by number", 0.0, 5.0, 0.0},
		{"very small numbers", 0.001, 0.001, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)
			if err != nil {
				t.Errorf("Divide(%v, %v) returned unexpected error: %v", tt.a, tt.b, err)
			}
			if result != tt.expected {
				t.Errorf("Divide(%v, %v) = %v; want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

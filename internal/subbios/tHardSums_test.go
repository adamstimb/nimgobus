package subbios

import (
	"testing"

	"github.com/adamstimb/nimgobus/internal/subbios/errorcode"
)

func TestFAddTwoReals(t *testing.T) {
	tests := []struct {
		a float64
		b float64
		r float64
		e int
	}{
		{0, 0, 0, errorcode.EOk},
		{5.5, 1.2, 6.7, errorcode.EOk},
	}

	s := Subbios{}
	s.Init()

	for _, tt := range tests {
		result := s.THardSums.FAddTwoReals(tt.a, tt.b)
		if result != tt.r {
			t.Errorf("FAddTwoReals(%f, %f) returned %f, expected %f", tt.a, tt.b, result, tt.r)
		}
	}
}

func TestFSubtractReals(t *testing.T) {
	tests := []struct {
		a float64
		b float64
		r float64
		e int
	}{
		{0, 0, 0, errorcode.EOk},
		{1.2, 5.5, 4.3, errorcode.EOk},
	}

	s := Subbios{}
	s.Init()

	for _, tt := range tests {
		result := s.THardSums.FSubtractReals(tt.a, tt.b)
		if result != tt.r {
			t.Errorf("FSubtractReals(%f, %f) returned %f, expected %f", tt.a, tt.b, result, tt.r)
		}
	}
}

func TestFMultiplyReals(t *testing.T) {
	tests := []struct {
		a float64
		b float64
		r float64
		e int
	}{
		{0, 0, 0, errorcode.EOk},
		{2.0, 4.4, 8.8, errorcode.EOk},
	}

	s := Subbios{}
	s.Init()

	for _, tt := range tests {
		result := s.THardSums.FMultiplyReals(tt.a, tt.b)
		if result != tt.r {
			t.Errorf("FMultiplyReals(%f, %f) returned %f, expected %f", tt.a, tt.b, result, tt.r)
		}
	}
}

func TestFDivideReals(t *testing.T) {
	tests := []struct {
		a float64
		b float64
		r float64
		e int
	}{
		{2.0, 0.0, 0, errorcode.ENotANumber},
		{2.0, 8.8, 4.4, errorcode.EOk},
	}

	s := Subbios{}
	s.Init()

	for _, tt := range tests {
		result := s.THardSums.FDivideReals(tt.a, tt.b)
		if result != tt.r {
			t.Errorf("FDivideReals(%f, %f) returned %f, expected %f", tt.a, tt.b, result, tt.r)
		}
		if s.FunctionError != tt.e {
			t.Errorf("FDivideReals(%f, %f) returned errorcode %d, expected %d", tt.a, tt.b, s.FunctionError, tt.e)
		}
	}
}

func TestFTruncateReal(t *testing.T) {
	tests := []struct {
		a float64
		r float64
		e int
	}{
		{0, 0, errorcode.EOk},
		{3.14, 3.0, errorcode.EOk},
		{3.78, 3.0, errorcode.EOk},
		{-3.78, -3.0, errorcode.EOk},
		{-3.14, -3.0, errorcode.EOk},
	}

	s := Subbios{}
	s.Init()

	for _, tt := range tests {
		result := s.THardSums.FTruncateReal(tt.a)
		if result != tt.r {
			t.Errorf("FTruncateReal(%f) returned %f, expected %f", tt.a, result, tt.r)
		}
	}
}

func TestFRealFromInt(t *testing.T) {
	tests := []struct {
		a int
		r float64
		e int
	}{
		{3, 3.0, errorcode.EOk},
	}

	s := Subbios{}
	s.Init()

	for _, tt := range tests {
		result := s.THardSums.FRealFromInt(tt.a)
		if result != tt.r {
			t.Errorf("FRealFromInt(%d) returned %f, expected %f", tt.a, result, tt.r)
		}
	}
}

func TestFIntLessThanReal(t *testing.T) {
	tests := []struct {
		a float64
		r int
		e int
	}{
		{1.6, 1, errorcode.EOk},
		{-1.6, -2, errorcode.EOk},
		{1.0, 1, errorcode.EOk},
	}

	s := Subbios{}
	s.Init()

	for _, tt := range tests {
		result := s.THardSums.FIntLessThanReal(tt.a)
		if result != tt.r {
			t.Errorf("FIntLessThanReal(%f) returned %d, expected %d", tt.a, result, tt.r)
		}
	}
}

func TestFIntPartOfReal(t *testing.T) {
	tests := []struct {
		a float64
		r int
		e int
	}{
		{0, 0, errorcode.EOk},
		{3.14, 3, errorcode.EOk},
		{3.78, 3, errorcode.EOk},
		{-3.78, -3, errorcode.EOk},
		{-3.14, -3, errorcode.EOk},
	}

	s := Subbios{}
	s.Init()

	for _, tt := range tests {
		result := s.THardSums.FIntPartOfReal(tt.a)
		if result != tt.r {
			t.Errorf("FIntPartOfReal(%f) returned %d, expected %d", tt.a, result, tt.r)
		}
	}
}

func TestFCommonLog(t *testing.T) {
	tests := []struct {
		a float64
		r float64
		e int
	}{
		{0, 0, errorcode.ENotANumber},
		{-3.14, 0, errorcode.ENotANumber},
		{10.0, 1.0, errorcode.EOk},
		{100.0, 2.0, errorcode.EOk},
	}

	s := Subbios{}
	s.Init()

	for _, tt := range tests {
		result := s.THardSums.FCommonLog(tt.a)
		if result != tt.r {
			t.Errorf("FCommonLog(%f) returned %f, expected %f", tt.a, result, tt.r)
		}
		if s.FunctionError != tt.e {
			t.Errorf("FCommonLog(%f) returned errorcode %d, expected %d", tt.a, s.FunctionError, tt.e)
		}
	}
}

func TestFNaturalLog(t *testing.T) {
	tests := []struct {
		a float64
		r float64
		e int
	}{
		{0, 0, errorcode.ENotANumber},
		{-3.14, 0, errorcode.ENotANumber},
	}

	s := Subbios{}
	s.Init()

	for _, tt := range tests {
		result := s.THardSums.FNaturalLog(tt.a)
		if result != tt.r {
			t.Errorf("FNaturalLog(%f) returned %f, expected %f", tt.a, result, tt.r)
		}
		if s.FunctionError != tt.e {
			t.Errorf("FNaturalLog(%f) returned errorcode %d, expected %d", tt.a, s.FunctionError, tt.e)
		}
	}
}

func TestFRealToAscii(t *testing.T) {
	tests := []struct {
		a float64
		r string
		e int
	}{
		{0, "0", errorcode.EOk},
		{-3.14, "-3.14", errorcode.EOk},
		{5.1234e9, "5.1234e+09", errorcode.EOk},
	}

	s := Subbios{}
	s.Init()

	for _, tt := range tests {
		result := s.THardSums.FRealToAscii(tt.a)
		if result != tt.r {
			t.Errorf("FRealToAscii(%f) returned %s, expected %s", tt.a, result, tt.r)
		}
		if s.FunctionError != tt.e {
			t.Errorf("FRealToAscii(%f) returned errorcode %d, expected %d", tt.a, s.FunctionError, tt.e)
		}
	}
}

func TestFAsciiToReal(t *testing.T) {
	tests := []struct {
		a string
		r float64
		e int
	}{
		{"0", 0, errorcode.EOk},
		{"-3.14", -3.14, errorcode.EOk},
		{"5.1234e+09", 5.1234e9, errorcode.EOk},
		{"foobar", 0, errorcode.ENotANumber},
	}

	s := Subbios{}
	s.Init()

	for _, tt := range tests {
		result := s.THardSums.FAsciiToReal(tt.a)
		if result != tt.r {
			t.Errorf("FAsciiToReal(%s) returned %f, expected %f", tt.a, result, tt.r)
		}
		if s.FunctionError != tt.e {
			t.Errorf("FAsciiToReal(%s) returned errorcode %d, expected %d", tt.a, s.FunctionError, tt.e)
		}
	}
}

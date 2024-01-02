package subbios

import (
	"fmt"
	"math"
	"strconv"

	"github.com/adamstimb/nimgobus/internal/subbios/errorcode"
)

// THardSums has all the t_hard_maths functions attached to it.
type THardSums struct {
	s *Subbios
}

// FAddTwoReals adds "a" to "b" and returns the result.
func (t *THardSums) FAddTwoReals(a, b float64) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	return a + b
}

// FSubtractReals subtracts "a" from "b" and returns the result.
func (t *THardSums) FSubtractReals(a, b float64) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	return b - a
}

// FMultiplyReals multiplies "a" by "b" and returns the result.
func (t *THardSums) FMultiplyReals(a, b float64) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	return a * b
}

// FDivideReals divides "b" by "a" and returns the result.
func (t *THardSums) FDivideReals(a, b float64) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	// Returns 0, ENotANumber if a=0.
	if b == 0 {
		t.s.FunctionError = errorcode.ENotANumber
		return 0
	}
	return b / a
}

// FTruncateReal returns the integer part of "a" as a float.
func (t *THardSums) FTruncateReal(a float64) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	if a < 0 {
		return -1.0 * math.Floor(a*-1.0)
	}
	return math.Floor(a)
}

// FRealFromInt converts an integer to floating point.
func (t *THardSums) FRealFromInt(a int) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	return float64(a)
}

// FIntLessThanReal basically floors a real and returns an int.
func (t *THardSums) FIntLessThanReal(a float64) int {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	return int(math.Floor(a))
}

// FIntPartOfReal returns the integer part of "a" as an int.
func (t *THardSums) FIntPartOfReal(a float64) int {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	if a < 0 {
		return int(-1.0 * math.Floor(a*-1.0))
	}
	return int(math.Floor(a))
}

// FCommonLog returns the log to base 10 of "a".
func (t *THardSums) FCommonLog(a float64) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	if a <= 0 {
		t.s.FunctionError = errorcode.ENotANumber
		return 0
	}
	return math.Log10(a)
}

// FNaturalLog returns the natural log of "a".
func (t *THardSums) FNaturalLog(a float64) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	if a <= 0 {
		t.s.FunctionError = errorcode.ENotANumber
		return 0
	}
	return math.Log(a)
}

// FInverseNaturalLog returns the inverse natural log of "a".  Note
// that this function uses 2.7182 as the base, as described in the
// Nimbus SUBBIOS documentation.
func (t *THardSums) FInverseNaturalLog(a float64) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	if a <= 0 {
		t.s.FunctionError = errorcode.ENotANumber
		return 0
	}
	r := math.Pow(2.7182, a)
	if math.IsNaN(r) {
		t.s.FunctionError = errorcode.ENotANumber
		return 0
	}
	return r
}

// FRaiseToPower raises a float to an integer power.
func (t *THardSums) FRaiseToPower(a float64, b int) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	return math.Pow(a, float64(b))
}

// FSquareRoot returns the square root of "a".
func (t *THardSums) FSquareRoot(a float64) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	r := math.Sqrt(a)
	if math.IsNaN(r) {
		t.s.FunctionError = errorcode.ENotANumber
		return 0
	}
	return r
}

// FCosine returns the cosine of "a".
func (t *THardSums) FCosine(a float64) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	r := math.Cos(a)
	if math.IsNaN(r) {
		t.s.FunctionError = errorcode.ENotANumber
		return 0
	}
	return r
}

// FSine returns the sine of "a".
func (t *THardSums) FSine(a float64) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	r := math.Sin(a)
	if math.IsNaN(r) {
		t.s.FunctionError = errorcode.ENotANumber
		return 0
	}
	return r
}

// FTangent returns the tangent of "a".
func (t *THardSums) FTangent(a float64) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	r := math.Tan(a)
	if math.IsNaN(r) {
		t.s.FunctionError = errorcode.ENotANumber
		return 0
	}
	return r
}

// FArctan returns the inverse tangent of "a".
func (t *THardSums) FArctan(a float64) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	r := math.Atan(a)
	if math.IsNaN(r) {
		t.s.FunctionError = errorcode.ENotANumber
		return 0
	}
	return r
}

// FRealToAscii returns the string representation of "a".
func (t *THardSums) FRealToAscii(a float64) string {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	return fmt.Sprintf("%g", a)
}

// FAsciiToReal converts the string representation of a float to float.
func (t *THardSums) FAsciiToReal(a string) float64 {
	t.s.FunctionError = errorcode.EOk
	t.s.FunctionStatus = 0
	if r, err := strconv.ParseFloat(a, 64); err == nil {
		return r
	} else {
		t.s.FunctionError = errorcode.ENotANumber
		return 0
	}
}

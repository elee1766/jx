package jx

import (
	"bytes"
	"strings"

	"github.com/ogen-go/errors"
)

// NumFormat is format of Num.Value.
type NumFormat uint8

// Possible formats of Num.
const (
	NumFormatInvalid  NumFormat = iota // invalid or blank
	NumFormatInt                       // 1234
	NumFormatFloat                     // 1.234
	NumFormatIntStr                    // "1234"
	NumFormatFloatStr                  // "1.234"
)

// Float reports whether format is float.
func (f NumFormat) Float() bool {
	return f == NumFormatFloat || f == NumFormatFloatStr
}

// Invalid reports whether format is invalid.
func (f NumFormat) Invalid() bool {
	return f == NumFormatInvalid || f > NumFormatFloatStr
}

// Int reports whether format is integer.
func (f NumFormat) Int() bool {
	return f == NumFormatInt || f == NumFormatIntStr
}

func (f NumFormat) String() string {
	switch f {
	case NumFormatInt:
		return "integer"
	case NumFormatFloat:
		return "float"
	case NumFormatIntStr:
		return "integer string"
	case NumFormatFloatStr:
		return "float string"
	default:
		return "invalid"
	}
}

// Str reports whether format is string integer or float.
func (f NumFormat) Str() bool {
	return f == NumFormatIntStr || f == NumFormatFloatStr
}

// Num represents number, which can be raw json number or string of number.
//
// Zero value is invalid.
type Num struct {
	// Format is number format for Value.
	Format NumFormat
	// Value is raw json of number, only digits or float characters.
	//
	// If Num is string number, Value does not contain quotes.
	Value []byte
}

// Equal reports whether numbers are strictly equal, including their formats.
func (n Num) Equal(v Num) bool {
	if n.Format != v.Format {
		return false
	}
	return bytes.Equal(n.Value, v.Value)
}

func (n Num) String() string {
	if n.Format.Invalid() {
		return "<invalid>"
	}
	var b strings.Builder
	if n.Format.Str() {
		b.WriteByte('\n')
	}
	_, _ = b.Write(n.Value)
	if n.Format.Str() {
		b.WriteByte('\n')
	}
	return b.String()
}

// Sign reports sign of number.
//
// 0 is zero, 1 is positive, -1 is negative.
func (n Num) Sign() int {
	if n.Format.Invalid() || len(n.Value) == 0 {
		return 0
	}
	switch n.Value[0] {
	case '-':
		return -1
	case '0':
		return 0
	default:
		return 1
	}
}

// Positive reports whether number is positive.
func (n Num) Positive() bool { return n.Sign() > 0 }

// Negative reports whether number is negative.
func (n Num) Negative() bool { return n.Sign() < 0 }

// Zero reports whether number is zero.
func (n Num) Zero() bool {
	if n.Format.Invalid() || len(n.Value) == 0 {
		return false
	}
	if len(n.Value) == 1 {
		return n.Value[0] == '0'
	}
	if n.Format.Int() {
		return false
	}
	for _, c := range n.Value {
		switch c {
		case '.', '0':
			continue
		default:
			return false
		}
	}
	return true
}

// Num encodes number.
func (e *Encoder) Num(v Num) {
	if v.Format.Invalid() {
		e.Null()
		return
	}
	if v.Format.Str() {
		e.byte('"')
	}
	e.RawBytes(v.Value)
	if v.Format.Str() {
		e.byte('"')
	}
}

// Num decodes number.
func (d *Decoder) Num() (Num, error) {
	return d.NumTo(Num{})
}

// NumTo decodes number into Num.
func (d *Decoder) NumTo(v Num) (Num, error) {
	var str bool
	switch d.Next() {
	case String:
		// Consume start of the string.
		d.head++
		str = true
	case Number: // float or integer
	default:
		return v, errors.Errorf("unexpected %s", d.Next())
	}
	if d.reader == nil {
		// Can use underlying buffer directly.
		v.Value = d.number()
	} else {
		buf, err := d.numberAppend(v.Value[:0])
		if err != nil {
			return v, errors.Wrap(err, "decode")
		}
		v.Value = buf
	}

	var dot bool
	for _, c := range v.Value {
		if c != '.' {
			continue
		}
		dot = true
		break
	}
	if dot {
		v.Format = NumFormatFloat
		if str {
			v.Format = NumFormatFloatStr
		}
	} else {
		v.Format = NumFormatInt
		if str {
			v.Format = NumFormatIntStr
		}
	}

	// TODO(ernado): Additional validity checks

	return v, nil
}
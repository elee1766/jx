package jx

import (
	"unicode/utf8"
)

// htmlSafeSet holds the value true if the ASCII character with the given
// array position can be safely represented inside a JSON string, embedded
// inside of HTML <script> tags, without any additional escaping.
//
// All values are true except for the ASCII control characters (0-31), the
// double quote ("), the backslash character ("\"), HTML opening and closing
// tags ("<" and ">"), and the ampersand ("&").
var htmlSafeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      false,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	tDot:     true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      false,
	'=':      true,
	'>':      false,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

// safeSet holds the value true if the ASCII character with the given array
// position can be represented inside a JSON string without any further
// escaping.
//
// All values are true except for the ASCII control characters (0-31), the
// double quote ("), and the backslash character ("\").
var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	tDot:     true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

const hex = "0123456789abcdef"

// StrHTMLEscaped write string to stream with html special characters escaped
func (w *Writer) StrHTMLEscaped(v string) {
	valLen := len(v)
	w.buf = append(w.buf, '"')
	// write string, the fast path, without utf8 and escape support
	i := 0
	for ; i < valLen; i++ {
		c := v[i]
		if c < utf8.RuneSelf && htmlSafeSet[c] {
			w.buf = append(w.buf, c)
		} else {
			break
		}
	}
	if i == valLen {
		w.buf = append(w.buf, '"')
		return
	}
	w.strHTMLEscapedSlow(i, v, valLen)
}

func (w *Writer) strHTMLEscapedSlow(i int, v string, valLen int) {
	start := i
	// for the remaining parts, we process them char by char
	for i < valLen {
		if b := v[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] {
				i++
				continue
			}
			if start < i {
				w.Raw(v[start:i])
			}
			switch b {
			case '\\', '"':
				w.twoBytes('\\', b)
			case '\n':
				w.twoBytes('\\', 'n')
			case '\r':
				w.twoBytes('\\', 'r')
			case '\t':
				w.twoBytes('\\', 't')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				w.Raw(`\u00`)
				w.twoBytes(hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(v[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				w.Raw(v[start:i])
			}
			w.Raw(`\ufffd`)
			i++
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				w.Raw(v[start:i])
			}
			w.Raw(`\u202`)
			w.byte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(v) {
		w.Raw(v[start:])
	}
	w.byte('"')
}

// Str write string to stream without html escape
func (w *Writer) Str(v string) {
	length := len(v)
	w.buf = append(w.buf, '"')
	// write string, the fast path, without utf8 and escape support
	i := 0
	for ; i < length; i++ {
		c := v[i]
		if c > 31 && c != '"' && c != '\\' {
			w.buf = append(w.buf, c)
		} else {
			break
		}
	}
	if i == length {
		w.buf = append(w.buf, '"')
		return
	}
	w.strSlowPath(i, v, length)
}

func (w *Writer) strSlowPath(i int, v string, length int) {
	start := i
	// for the remaining parts, we process them char by char
	for i < length {
		if b := v[i]; b < utf8.RuneSelf {
			if safeSet[b] {
				i++
				continue
			}
			if start < i {
				w.Raw(v[start:i])
			}
			switch b {
			case '\\', '"':
				w.twoBytes('\\', b)
			case '\n':
				w.twoBytes('\\', 'n')
			case '\r':
				w.twoBytes('\\', 'r')
			case '\t':
				w.twoBytes('\\', 't')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				w.Raw(`\u00`)
				w.twoBytes(hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		i++
		continue
	}
	if start < len(v) {
		w.Raw(v[start:])
	}
	w.byte('"')
}
package digest

// Integer type classification for MySQL numeric literals.
// This module determines the appropriate token type (NUM, LONG_NUM, etc.)
// based on the magnitude of integer literals, matching MySQL's int_token()
// function in sql_lex.cc.

// IntegerClassifier determines the token type for integer literals.
// It is stateless and can be reused across multiple classifications.
type IntegerClassifier struct{}

// Boundary constants for integer type classification.
// These match MySQL's limits in sql_lex.cc.
const (
	// Maximum values for each integer type
	maxLong             = "2147483647" // 2^31 - 1
	maxLongLen          = 10
	minSignedLong       = "2147483648"          // 2^31 (absolute value of -2^31)
	maxLongLong         = "9223372036854775807" // 2^63 - 1
	maxLongLongLen      = 19
	minSignedLongLong   = "9223372036854775808"  // 2^63 (absolute value of -2^63)
	maxUnsignedLongLong = "18446744073709551615" // 2^64 - 1
	maxUnsignedLongLen  = 20
)

// NewIntegerClassifier creates a new IntegerClassifier.
func NewIntegerClassifier() *IntegerClassifier {
	return &IntegerClassifier{}
}

// Classify determines the token type for an integer string.
// The string may include an optional leading sign (+/-) and leading zeros.
// Returns NUM, LONG_NUM, ULONGLONG_NUM, or DECIMAL_NUM.
func (c *IntegerClassifier) Classify(s string) int {
	if len(s) == 0 {
		return NUM
	}

	// Quick path for short numbers
	if len(s) < maxLongLen {
		return NUM
	}

	// Parse sign and normalize
	neg := false
	offset := 0

	if s[0] == '+' {
		offset++
	} else if s[0] == '-' {
		offset++
		neg = true
	}

	// Skip leading zeros
	str := s[offset:]
	for len(str) > 0 && str[0] == '0' {
		str = str[1:]
	}

	length := len(str)

	// After normalization, check again
	if length < maxLongLen {
		return NUM
	}

	return c.classifyByMagnitude(str, length, neg)
}

// classifyByMagnitude compares the normalized digit string against type boundaries.
func (c *IntegerClassifier) classifyByMagnitude(str string, length int, neg bool) int {
	var cmp string
	var smaller, bigger int

	if neg {
		// Negative numbers
		if length == maxLongLen {
			cmp = minSignedLong
			smaller = NUM
			bigger = LONG_NUM
		} else if length < maxLongLongLen {
			return LONG_NUM
		} else if length > maxLongLongLen {
			return DECIMAL_NUM
		} else {
			// length == maxLongLongLen
			cmp = minSignedLongLong
			smaller = LONG_NUM
			bigger = DECIMAL_NUM
		}
	} else {
		// Positive numbers
		if length == maxLongLen {
			cmp = maxLong
			smaller = NUM
			bigger = LONG_NUM
		} else if length < maxLongLongLen {
			return LONG_NUM
		} else if length > maxLongLongLen {
			if length > maxUnsignedLongLen {
				return DECIMAL_NUM
			}
			cmp = maxUnsignedLongLong
			smaller = ULONGLONG_NUM
			bigger = DECIMAL_NUM
		} else {
			// length == maxLongLongLen
			cmp = maxLongLong
			smaller = LONG_NUM
			bigger = ULONGLONG_NUM
		}
	}

	// Compare digit by digit
	return c.compareDigits(str, cmp, smaller, bigger)
}

// compareDigits compares two digit strings and returns smaller or bigger token type.
func (c *IntegerClassifier) compareDigits(str, cmp string, smaller, bigger int) int {
	for i := 0; i < len(str) && i < len(cmp); i++ {
		if str[i] < cmp[i] {
			return smaller
		}
		if str[i] > cmp[i] {
			return bigger
		}
	}
	return smaller // Equal means it fits in the smaller type
}

// defaultClassifier is a package-level classifier for convenience.
var defaultClassifier = NewIntegerClassifier()

// ClassifyInteger is a convenience function that uses the default classifier.
func ClassifyInteger(s string) int {
	return defaultClassifier.Classify(s)
}

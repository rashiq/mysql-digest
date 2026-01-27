package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

// Digest represents a computed SQL digest
type Digest struct {
	// Hash is the SHA-256 hash of the binary token array (hex-encoded)
	Hash string
	// Text is the normalized SQL text with literals replaced by placeholders
	Text string
}

// Compute calculates the digest of a SQL statement.
// It normalizes the SQL by:
// - Replacing literal values (strings, numbers) with placeholders (?)
// - Collapsing multiple values in IN(...) to a single placeholder
// - Preserving keywords and identifiers
// - Normalizing whitespace
func Compute(sql string) Digest {
	normalizer := newNormalizer(sql)
	normalizer.normalize()

	// Compute SHA-256 hash on binary token array (like MySQL does)
	hash := sha256.Sum256(normalizer.tokenArray)
	hashStr := hex.EncodeToString(hash[:])

	// Build human-readable digest text
	text := normalizer.buildOutput()

	return Digest{
		Hash: hashStr,
		Text: text,
	}
}

// ComputeWithOptions calculates digest with custom options
func ComputeWithOptions(sql string, opts Options) Digest {
	normalizer := newNormalizer(sql)
	normalizer.opts = opts
	normalizer.normalize()

	hash := sha256.Sum256(normalizer.tokenArray)
	hashStr := hex.EncodeToString(hash[:])

	text := normalizer.buildOutput()

	return Digest{
		Hash: hashStr,
		Text: text,
	}
}

// Options controls digest computation behavior
type Options struct {
	// SQLMode affects lexer behavior (ANSI_QUOTES, NO_BACKSLASH_ESCAPES)
	SQLMode SQLMode
	// MaxLength limits the digest text length (0 = unlimited)
	MaxLength int
}

// storedToken represents a token in the reduction stack (for text output)
type storedToken struct {
	tokType int
	text    string // For identifiers, the identifier text
}

// normalizer handles SQL normalization using MySQL's reduction-based approach
type normalizer struct {
	lexer            *Lexer
	opts             Options
	tokens           []storedToken // Token stack for reductions (used for text output)
	tokenArray       []byte        // Binary token array for hashing (MySQL format)
	lastIdentIndex   int           // Index after last identifier (for peek boundary)
	inOrderOrGroupBy bool          // Inside ORDER BY, GROUP BY, PARTITION BY clause
}

func newNormalizer(sql string) *normalizer {
	return &normalizer{
		lexer:          NewLexer(sql),
		tokens:         make([]storedToken, 0, 256),
		tokenArray:     make([]byte, 0, 1024),
		lastIdentIndex: 0,
	}
}

// storeTokenBinary stores a token in the binary array (2 bytes, little-endian)
func (n *normalizer) storeTokenBinary(token int) {
	n.tokenArray = append(n.tokenArray, byte(token&0xff), byte((token>>8)&0xff))
}

// storeIdentifierBinary stores an identifier token with its text in the binary array
// Format: 2 bytes (token) + 2 bytes (length) + N bytes (identifier text)
func (n *normalizer) storeIdentifierBinary(token int, text string) {
	// Write the token (2 bytes)
	n.tokenArray = append(n.tokenArray, byte(token&0xff), byte((token>>8)&0xff))
	// Write the string length (2 bytes)
	length := len(text)
	n.tokenArray = append(n.tokenArray, byte(length&0xff), byte((length>>8)&0xff))
	// Write the string data
	n.tokenArray = append(n.tokenArray, []byte(text)...)
}

// storeByNumericColumnBinary stores a TOK_BY_NUMERIC_COLUMN with its value
// Format: 2 bytes (token) + 4 bytes (value as little-endian)
func (n *normalizer) storeByNumericColumnBinary(value uint64) {
	token := TOK_BY_NUMERIC_COLUMN
	n.tokenArray = append(n.tokenArray, byte(token&0xff), byte((token>>8)&0xff))
	// Write the value (4 bytes, little-endian)
	n.tokenArray = append(n.tokenArray,
		byte(value&0xff),
		byte((value>>8)&0xff),
		byte((value>>16)&0xff),
		byte((value>>24)&0xff),
	)
}

// removeLastTokenBinary removes the last token (2 bytes) from the binary array
func (n *normalizer) removeLastTokenBinary() {
	if len(n.tokenArray) >= 2 {
		n.tokenArray = n.tokenArray[:len(n.tokenArray)-2]
	}
}

// removeLastTwoTokensBinary removes the last two tokens (4 bytes) from the binary array
func (n *normalizer) removeLastTwoTokensBinary() {
	if len(n.tokenArray) >= 4 {
		n.tokenArray = n.tokenArray[:len(n.tokenArray)-4]
	}
}

func (n *normalizer) normalize() {
	n.lexer.SetSQLMode(n.opts.SQLMode)

	for {
		tok := n.lexer.Lex()

		if tok.Type == END_OF_INPUT {
			// Remove trailing semicolon from both arrays
			if len(n.tokens) > 0 && n.tokens[len(n.tokens)-1].tokType == ';' {
				n.tokens = n.tokens[:len(n.tokens)-1]
				n.removeLastTokenBinary()
			}
			break
		}

		if tok.Type == ABORT_SYM {
			break
		}

		tokType := tok.Type

		// Track ORDER BY / GROUP BY context
		if tokType == ORDER_SYM || tokType == GROUP_SYM || tokType == PARTITION_SYM {
			n.inOrderOrGroupBy = true
		} else if n.inOrderOrGroupBy {
			switch tokType {
			case LIMIT, OFFSET_SYM, FOR_SYM, LOCK_SYM, PROCEDURE_SYM, SELECT_SYM, UPDATE_SYM, DELETE_SYM, INSERT_SYM, UNION_SYM, END_OF_INPUT, ';':
				n.inOrderOrGroupBy = false
			}
		}

		n.addToken(tok)
	}
}

// addToken adds a token to the stack and performs reductions
func (n *normalizer) addToken(tok Token) {
	tokType := tok.Type

	switch tokType {
	case NUM, LONG_NUM, ULONGLONG_NUM, DECIMAL_NUM, FLOAT_NUM, BIN_NUM, HEX_NUM:
		// Handle unary +/- signs
		n.reduceUnarySign()

		// Special handling for ORDER BY / GROUP BY numeric columns
		if n.inOrderOrGroupBy && n.isNumericColumnRef() {
			text := n.lexer.TokenText(tok)
			value, _ := strconv.ParseUint(text, 10, 64)
			n.storeByNumericColumnBinary(value)
			n.storeToken(storedToken{tokType: TOK_BY_NUMERIC_COLUMN, text: text})
			return
		}

		// Reduce to TOK_GENERIC_VALUE
		n.reduceToGenericValueAndList()

	case LEX_HOSTNAME, TEXT_STRING, NCHAR_STRING, PARAM_MARKER:
		// Reduce to TOK_GENERIC_VALUE
		n.reduceToGenericValueAndList()

	case NULL_SYM:
		// NULL is reduced to TOK_GENERIC_VALUE only when used as a literal value,
		// NOT when it follows IS or IS NOT (e.g., "x IS NULL", "x IS NOT NULL").
		// Check if the previous token is IS or NOT (preceded by IS).
		if n.isNullKeywordContext() {
			// Keep NULL as a keyword
			n.storeTokenBinary(NULL_SYM)
			n.storeToken(storedToken{tokType: NULL_SYM})
		} else {
			// Reduce to TOK_GENERIC_VALUE
			n.reduceToGenericValueAndList()
		}

	case ')':
		// On close paren, we store it first, then try to reduce the stack.
		n.storeTokenBinary(')')
		n.storeToken(storedToken{tokType: ')'})
		n.reduceStack()

	case IDENT, IDENT_QUOTED:
		// Store identifier with its text
		identText := n.lexer.TokenText(tok)
		if tokType == IDENT_QUOTED {
			identText = stripIdentifierQuotes(identText)
		}
		// Use TOK_IDENT for both IDENT and IDENT_QUOTED (normalized)
		n.storeIdentifierBinary(TOK_IDENT, identText)
		n.storeToken(storedToken{tokType: TOK_IDENT, text: identText})
		n.lastIdentIndex = len(n.tokens)

	default:
		// Store token as-is
		n.storeTokenBinary(tokType)
		n.storeToken(storedToken{tokType: tokType})
		n.reduceStack() // Try to reduce after other tokens too
	}
}

// isNumericColumnRef checks if we're in a position where a numeric literal is a column reference
func (n *normalizer) isNumericColumnRef() bool {
	if len(n.tokens) == 0 {
		return false
	}
	lastTok := n.tokens[len(n.tokens)-1].tokType
	return lastTok == BY || lastTok == ',' || lastTok == '('
}

// isNullKeywordContext checks if NULL should be kept as a keyword (after IS or IS NOT)
// rather than reduced to a generic value placeholder.
func (n *normalizer) isNullKeywordContext() bool {
	if len(n.tokens) == 0 {
		return false
	}
	lastTok := n.tokens[len(n.tokens)-1].tokType

	// IS NULL
	if lastTok == IS {
		return true
	}

	// IS NOT NULL - check if we have "IS NOT" pattern
	if lastTok == NOT_SYM && len(n.tokens) >= 2 {
		prevTok := n.tokens[len(n.tokens)-2].tokType
		if prevTok == IS {
			return true
		}
	}

	return false
}

// reduceUnarySign handles unary +/- absorption
func (n *normalizer) reduceUnarySign() {
	for {
		last, last2 := n.peekLast2Raw() // Use raw peek to include identifiers
		if (last == '+' || last == '-') && startsExpression(last2) {
			// Absorb the unary sign from both arrays
			n.tokens = n.tokens[:len(n.tokens)-1]
			n.removeLastTokenBinary()
		} else {
			break
		}
	}
}

// reduceToGenericValueAndList reduces literals to TOK_GENERIC_VALUE and handles list building
func (n *normalizer) reduceToGenericValueAndList() {
	last, last2 := n.peekLast2()

	if (last2 == TOK_GENERIC_VALUE || last2 == TOK_GENERIC_VALUE_LIST) && last == ',' {
		// Reduce: TOK_GENERIC_VALUE_LIST := TOK_GENERIC_VALUE/LIST ',' TOK_GENERIC_VALUE
		n.tokens = n.tokens[:len(n.tokens)-2]
		n.removeLastTwoTokensBinary()
		n.storeTokenBinary(TOK_GENERIC_VALUE_LIST)
		n.storeToken(storedToken{tokType: TOK_GENERIC_VALUE_LIST})
	} else {
		n.storeTokenBinary(TOK_GENERIC_VALUE)
		n.storeToken(storedToken{tokType: TOK_GENERIC_VALUE})
	}
	n.reduceStack()
}

// reduceStack applies reduction rules to the token stack until no more apply.
//
// Reduction rules (in order of precedence):
//
//	Rule 1: '(' ? ')'      → (?)         Single value in parens
//	Rule 2: '(' ?, ... ')' → (...)       Multiple values in parens
//	Rule 3: (?) , (?)      → (?), ...    List of single-value rows
//	Rule 4: (...) , (...)  → (...), ...  List of multi-value rows
//	Rule 5: IN (?)         → IN (...)    Collapse IN clause
//	Rule 6: IN (...)       → IN (...)    Collapse IN clause
func (n *normalizer) reduceStack() {
	for {
		if n.tryReduceParenthesizedValue() {
			continue
		}
		if n.tryReduceRowList() {
			continue
		}
		if n.tryReduceInClause() {
			continue
		}
		return // No reduction applied
	}
}

// tryReduceParenthesizedValue handles: '(' VALUE ')' → ROW_VALUE
// Returns true if a reduction was made.
func (n *normalizer) tryReduceParenthesizedValue() bool {
	if len(n.tokens) < 3 {
		return false
	}

	last := n.tokens[len(n.tokens)-1].tokType
	mid := n.tokens[len(n.tokens)-2].tokType
	first := n.tokens[len(n.tokens)-3].tokType

	if last != ')' || first != '(' {
		return false
	}

	switch mid {
	case TOK_GENERIC_VALUE:
		// '(' ? ')' → (?)
		n.tokens = n.tokens[:len(n.tokens)-3]
		// Remove 3 tokens from binary array (6 bytes)
		if len(n.tokenArray) >= 6 {
			n.tokenArray = n.tokenArray[:len(n.tokenArray)-6]
		}
		n.storeTokenBinary(TOK_ROW_SINGLE_VALUE)
		n.storeToken(storedToken{tokType: TOK_ROW_SINGLE_VALUE})
		return true

	case TOK_GENERIC_VALUE_LIST:
		// '(' ?, ... ')' → (...)
		n.tokens = n.tokens[:len(n.tokens)-3]
		if len(n.tokenArray) >= 6 {
			n.tokenArray = n.tokenArray[:len(n.tokenArray)-6]
		}
		n.storeTokenBinary(TOK_ROW_MULTIPLE_VALUE)
		n.storeToken(storedToken{tokType: TOK_ROW_MULTIPLE_VALUE})
		return true
	}

	return false
}

// tryReduceRowList handles: ROW ',' ROW → ROW_LIST
// Returns true if a reduction was made.
func (n *normalizer) tryReduceRowList() bool {
	if len(n.tokens) < 3 {
		return false
	}

	last := n.tokens[len(n.tokens)-1].tokType
	comma := n.tokens[len(n.tokens)-2].tokType
	first := n.tokens[len(n.tokens)-3].tokType

	if comma != ',' {
		return false
	}

	// Single-value rows: (?) , (?) → (?), ...
	if isSingleValueRow(last) && isSingleValueRow(first) {
		n.tokens = n.tokens[:len(n.tokens)-3]
		if len(n.tokenArray) >= 6 {
			n.tokenArray = n.tokenArray[:len(n.tokenArray)-6]
		}
		n.storeTokenBinary(TOK_ROW_SINGLE_VALUE_LIST)
		n.storeToken(storedToken{tokType: TOK_ROW_SINGLE_VALUE_LIST})
		return true
	}

	// Multi-value rows: (...) , (...) → (...), ...
	if isMultiValueRow(last) && isMultiValueRow(first) {
		n.tokens = n.tokens[:len(n.tokens)-3]
		if len(n.tokenArray) >= 6 {
			n.tokenArray = n.tokenArray[:len(n.tokenArray)-6]
		}
		n.storeTokenBinary(TOK_ROW_MULTIPLE_VALUE_LIST)
		n.storeToken(storedToken{tokType: TOK_ROW_MULTIPLE_VALUE_LIST})
		return true
	}

	return false
}

// tryReduceInClause handles: IN ROW → IN (...)
// Returns true if a reduction was made.
func (n *normalizer) tryReduceInClause() bool {
	if len(n.tokens) < 2 {
		return false
	}

	last := n.tokens[len(n.tokens)-1].tokType
	before := n.tokens[len(n.tokens)-2].tokType

	if before != IN_SYM {
		return false
	}

	if last == TOK_ROW_SINGLE_VALUE || last == TOK_ROW_MULTIPLE_VALUE {
		n.tokens = n.tokens[:len(n.tokens)-2]
		n.removeLastTwoTokensBinary()
		n.storeTokenBinary(TOK_IN_GENERIC_VALUE_EXPRESSION)
		n.storeToken(storedToken{tokType: TOK_IN_GENERIC_VALUE_EXPRESSION})
		return true
	}

	return false
}

// Helper: is this a single-value row token?
func isSingleValueRow(tok int) bool {
	return tok == TOK_ROW_SINGLE_VALUE || tok == TOK_ROW_SINGLE_VALUE_LIST
}

// Helper: is this a multi-value row token?
func isMultiValueRow(tok int) bool {
	return tok == TOK_ROW_MULTIPLE_VALUE || tok == TOK_ROW_MULTIPLE_VALUE_LIST
}

// storeToken adds a token to the stack
func (n *normalizer) storeToken(tok storedToken) {
	n.tokens = append(n.tokens, tok)
}

// peekLast2 returns the last two tokens from the stack.
func (n *normalizer) peekLast2() (last, last2 int) {
	last = TOK_UNUSED
	last2 = TOK_UNUSED

	if len(n.tokens) >= 1 {
		last = n.tokens[len(n.tokens)-1].tokType
	}
	if len(n.tokens) >= 2 {
		last2 = n.tokens[len(n.tokens)-2].tokType
	}
	return
}

// peekLast2Raw returns the last two tokens from the stack (including identifiers)
func (n *normalizer) peekLast2Raw() (last, last2 int) {
	last = TOK_UNUSED
	last2 = TOK_UNUSED

	idx := len(n.tokens)
	if idx > 0 {
		idx--
		last = n.tokens[idx].tokType

		if idx > 0 {
			idx--
			last2 = n.tokens[idx].tokType
		}
	}
	return
}

// buildOutput converts the token stack to output string using MySQL's delayed space approach
func (n *normalizer) buildOutput() string {
	var builder strings.Builder
	addSpace := false

	for _, tok := range n.tokens {
		text := n.tokenText(tok)
		if text == "" {
			continue
		}

		// Add delayed space before this token (if previous token requested it)
		if addSpace {
			builder.WriteByte(' ')
		}

		builder.WriteString(text)

		// Check if this token wants a space after it (delayed until next token)
		addSpace = TokenAppendSpace(tok.tokType)
	}

	result := builder.String()

	// Apply max length if set
	if n.opts.MaxLength > 0 && len(result) > n.opts.MaxLength {
		result = result[:n.opts.MaxLength] + "..."
	}

	return result
}

// tokenText returns the output text for a token
func (n *normalizer) tokenText(tok storedToken) string {
	switch tok.tokType {
	case TOK_IDENT:
		return "`" + escapeBackticks(tok.text) + "`"
	case TOK_BY_NUMERIC_COLUMN:
		return tok.text
	default:
		text := TokenString(tok.tokType)
		if text == "(unknown)" {
			return ""
		}
		return text
	}
}

// startsExpression returns true if the token can start an expression
// (meaning a following +/- would be unary, not binary).
func startsExpression(tokType int) bool {
	if tokType == 0 || tokType == TOK_UNUSED {
		return true
	}
	return TokenStartExpr(tokType)
}

// stripIdentifierQuotes removes surrounding quotes from a quoted identifier.
func stripIdentifierQuotes(s string) string {
	if len(s) < 2 {
		return s
	}

	if s[0] == '`' && s[len(s)-1] == '`' {
		inner := s[1 : len(s)-1]
		return strings.ReplaceAll(inner, "``", "`")
	}

	if s[0] == '"' && s[len(s)-1] == '"' {
		inner := s[1 : len(s)-1]
		return strings.ReplaceAll(inner, `""`, `"`)
	}

	return s
}

// escapeBackticks escapes backticks in an identifier by doubling them.
func escapeBackticks(s string) string {
	return strings.ReplaceAll(s, "`", "``")
}

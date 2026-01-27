package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// Digest represents a computed SQL digest
type Digest struct {
	// Hash is the SHA-256 hash of the normalized SQL (hex-encoded)
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
	text := normalizer.normalize()

	// Compute SHA-256 hash
	hash := sha256.Sum256([]byte(text))
	hashStr := hex.EncodeToString(hash[:])

	return Digest{
		Hash: hashStr,
		Text: text,
	}
}

// ComputeWithOptions calculates digest with custom options
func ComputeWithOptions(sql string, opts Options) Digest {
	normalizer := newNormalizer(sql)
	normalizer.opts = opts
	text := normalizer.normalize()

	hash := sha256.Sum256([]byte(text))
	hashStr := hex.EncodeToString(hash[:])

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

// storedToken represents a token in the reduction stack
type storedToken struct {
	tokType int
	text    string // For identifiers, the identifier text
}

// normalizer handles SQL normalization using MySQL's reduction-based approach
type normalizer struct {
	lexer            *Lexer
	opts             Options
	tokens           []storedToken // Token stack for reductions
	lastIdentIndex   int           // Index after last identifier (for peek boundary)
	inOrderOrGroupBy bool          // Inside ORDER BY, GROUP BY, PARTITION BY clause
}

func newNormalizer(sql string) *normalizer {
	return &normalizer{
		lexer:          NewLexer(sql),
		tokens:         make([]storedToken, 0, 256),
		lastIdentIndex: 0,
	}
}

func (n *normalizer) normalize() string {
	n.lexer.SetSQLMode(n.opts.SQLMode)

	for {
		tok := n.lexer.Lex()

		if tok.Type == END_OF_INPUT {
			// Remove trailing semicolon
			if len(n.tokens) > 0 && n.tokens[len(n.tokens)-1].tokType == ';' {
				n.tokens = n.tokens[:len(n.tokens)-1]
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

	return n.buildOutput()
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
			n.storeToken(storedToken{tokType: TOK_BY_NUMERIC_COLUMN, text: n.lexer.TokenText(tok)})
			return
		}

		// Reduce to TOK_GENERIC_VALUE
		n.reduceToGenericValueAndList()

	case LEX_HOSTNAME, TEXT_STRING, NCHAR_STRING, PARAM_MARKER:
		// Reduce to TOK_GENERIC_VALUE
		n.reduceToGenericValueAndList()

	case NULL_SYM:
		// NULL is also reduced to TOK_GENERIC_VALUE
		n.reduceToGenericValueAndList()

	case ')':
		// On close paren, we store it first, then try to reduce the stack.
		n.storeToken(storedToken{tokType: ')'})
		n.reduceStack()

	case IDENT, IDENT_QUOTED:
		// Store identifier with its text
		identText := n.lexer.TokenText(tok)
		if tokType == IDENT_QUOTED {
			identText = stripIdentifierQuotes(identText)
		}
		n.storeToken(storedToken{tokType: TOK_IDENT, text: identText})
		n.lastIdentIndex = len(n.tokens)

	default:
		// Store token as-is
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

// reduceUnarySign handles unary +/- absorption
func (n *normalizer) reduceUnarySign() {
	for {
		last, last2 := n.peekLast2Raw() // Use raw peek to include identifiers
		if (last == '+' || last == '-') && startsExpression(last2) {
			// Absorb the unary sign
			n.tokens = n.tokens[:len(n.tokens)-1]
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
		n.storeToken(storedToken{tokType: TOK_GENERIC_VALUE_LIST})
	} else {
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
		n.storeToken(storedToken{tokType: TOK_ROW_SINGLE_VALUE})
		return true

	case TOK_GENERIC_VALUE_LIST:
		// '(' ?, ... ')' → (...)
		n.tokens = n.tokens[:len(n.tokens)-3]
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
		n.storeToken(storedToken{tokType: TOK_ROW_SINGLE_VALUE_LIST})
		return true
	}

	// Multi-value rows: (...) , (...) → (...), ...
	if isMultiValueRow(last) && isMultiValueRow(first) {
		n.tokens = n.tokens[:len(n.tokens)-3]
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

// buildOutput converts the token stack to output string
func (n *normalizer) buildOutput() string {
	var builder strings.Builder
	lastWritten := 0

	for _, tok := range n.tokens {
		text := n.tokenText(tok)
		if text == "" {
			continue
		}

		// Add space if needed
		if builder.Len() > 0 && needsSpaceBefore(lastWritten, tok.tokType) {
			builder.WriteByte(' ')
		}

		builder.WriteString(text)
		lastWritten = tok.tokType
	}

	result := strings.TrimSpace(builder.String())

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

func needsSpaceBefore(lastWritten, tokType int) bool {
	return TokenAppendSpace(lastWritten) && TokenPrependSpace(tokType)
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

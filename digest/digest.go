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
	inValues         bool          // Saw VALUES keyword
	valuesParenDepth int           // Paren depth when VALUES was seen
	sawFirstRow      bool          // Saw first row in VALUES
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

	parenDepth := 0

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

		// Track paren depth for VALUES context
		if tokType == '(' {
			parenDepth++
		} else if tokType == ')' {
			parenDepth--
		}

		// Track VALUES keyword
		if tokType == VALUES {
			n.inValues = true
			n.valuesParenDepth = parenDepth
			n.sawFirstRow = false
		}

		// Reset VALUES context on certain keywords
		if n.inValues && (tokType == ON_SYM || tokType == WHERE || tokType == SET_SYM) {
			n.inValues = false
		}

		// Track ORDER BY / GROUP BY context
		if tokType == ORDER_SYM || tokType == GROUP_SYM || tokType == PARTITION_SYM {
			n.inOrderOrGroupBy = true
		} else if n.inOrderOrGroupBy {
			switch tokType {
			case LIMIT, OFFSET_SYM, FOR_SYM, LOCK_SYM, PROCEDURE_SYM, SELECT_SYM, UPDATE_SYM, DELETE_SYM, INSERT_SYM, UNION_SYM, END_OF_INPUT, ';':
				n.inOrderOrGroupBy = false
			}
		}

		n.addToken(tok, parenDepth)
	}

	return n.buildOutput()
}

// addToken adds a token to the stack and performs reductions
func (n *normalizer) addToken(tok Token, parenDepth int) {
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
		// Check for row value reductions
		n.reduceCloseParen(parenDepth)

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
}

// reduceCloseParen handles ')' token and row value reductions
func (n *normalizer) reduceCloseParen(parenDepth int) {
	last, last2 := n.peekLast2()

	if last == TOK_GENERIC_VALUE && last2 == '(' {
		// Reduce: TOK_ROW_SINGLE_VALUE := '(' TOK_GENERIC_VALUE ')'
		n.tokens = n.tokens[:len(n.tokens)-2]
		token := TOK_ROW_SINGLE_VALUE

		// Check for IN reduction or row list
		token = n.reduceRowSingleValue(token, parenDepth)
		n.storeToken(storedToken{tokType: token})

	} else if last == TOK_GENERIC_VALUE_LIST && last2 == '(' {
		// Reduce: TOK_ROW_MULTIPLE_VALUE := '(' TOK_GENERIC_VALUE_LIST ')'
		n.tokens = n.tokens[:len(n.tokens)-2]
		token := TOK_ROW_MULTIPLE_VALUE

		// Check for IN reduction or row list
		token = n.reduceRowMultipleValue(token, parenDepth)
		n.storeToken(storedToken{tokType: token})

	} else {
		// No reduction, just store the )
		n.storeToken(storedToken{tokType: ')'})
	}
}

// reduceRowSingleValue checks for further reductions after creating TOK_ROW_SINGLE_VALUE
func (n *normalizer) reduceRowSingleValue(token int, parenDepth int) int {
	last, last2 := n.peekLast2()

	if (last2 == TOK_ROW_SINGLE_VALUE || last2 == TOK_ROW_SINGLE_VALUE_LIST) && last == ',' {
		// Reduce to row list
		n.tokens = n.tokens[:len(n.tokens)-2]
		return TOK_ROW_SINGLE_VALUE_LIST
	} else if last == IN_SYM {
		// Reduce: TOK_IN_GENERIC_VALUE_EXPRESSION := IN_SYM TOK_ROW_SINGLE_VALUE
		n.tokens = n.tokens[:len(n.tokens)-1]
		return TOK_IN_GENERIC_VALUE_EXPRESSION
	}

	// Track first row for VALUES
	if n.inValues && parenDepth == n.valuesParenDepth+1 {
		n.sawFirstRow = true
	}

	return token
}

// reduceRowMultipleValue checks for further reductions after creating TOK_ROW_MULTIPLE_VALUE
func (n *normalizer) reduceRowMultipleValue(token int, parenDepth int) int {
	last, last2 := n.peekLast2()

	if (last2 == TOK_ROW_MULTIPLE_VALUE || last2 == TOK_ROW_MULTIPLE_VALUE_LIST) && last == ',' {
		// Reduce to row list
		n.tokens = n.tokens[:len(n.tokens)-2]
		return TOK_ROW_MULTIPLE_VALUE_LIST
	} else if last == IN_SYM {
		// Reduce: TOK_IN_GENERIC_VALUE_EXPRESSION := IN_SYM TOK_ROW_MULTIPLE_VALUE
		n.tokens = n.tokens[:len(n.tokens)-1]
		return TOK_IN_GENERIC_VALUE_EXPRESSION
	}

	// Track first row for VALUES
	if n.inValues && parenDepth == n.valuesParenDepth+1 {
		n.sawFirstRow = true
	}

	return token
}

// storeToken adds a token to the stack
func (n *normalizer) storeToken(tok storedToken) {
	n.tokens = append(n.tokens, tok)
}

// peekLast2 returns the last two non-identifier tokens from the stack
func (n *normalizer) peekLast2() (last, last2 int) {
	last = TOK_UNUSED
	last2 = TOK_UNUSED

	idx := len(n.tokens)
	if idx > n.lastIdentIndex {
		idx--
		last = n.tokens[idx].tokType

		if idx > n.lastIdentIndex {
			idx--
			last2 = n.tokens[idx].tokType
		}
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

func isLiteral(tokType int) bool {
	switch tokType {
	case NUM, LONG_NUM, ULONGLONG_NUM, DECIMAL_NUM, FLOAT_NUM,
		TEXT_STRING, NCHAR_STRING,
		HEX_NUM, BIN_NUM,
		LEX_HOSTNAME,
		NULL_SYM,
		PARAM_MARKER:
		return true
	}
	return false
}

// isNumericLiteral returns true if the token is a numeric literal.
func isNumericLiteral(tokType int) bool {
	switch tokType {
	case NUM, LONG_NUM, ULONGLONG_NUM, DECIMAL_NUM, FLOAT_NUM, HEX_NUM, BIN_NUM:
		return true
	}
	return false
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

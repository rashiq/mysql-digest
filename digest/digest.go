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

// normalizer handles SQL normalization
type normalizer struct {
	lexer            *Lexer
	builder          strings.Builder
	opts             Options
	lastToken        int
	secondLastToken  int  // Track second-to-last token for context
	inValues         bool // Saw VALUES keyword, looking for row tuples
	inList           bool // Inside IN (...) list or VALUES row
	parenDepth       int  // Parenthesis depth for list tracking
	listStart        int  // Paren depth when list started
	sawFirstLiteral  bool // Saw first literal in list (for collapsing)
	sawFirstRow      bool // Saw first VALUES row (for row collapsing)
	valuesDepth      int  // Paren depth at VALUES context start
	skippingRows     bool // Currently skipping subsequent VALUES rows
	pendingSign      int  // Pending +/- token (0, '+', or '-')
	pendingSignToken int  // The actual token type of pending sign
	absorbedSign     bool // Did we absorb a unary sign this iteration?
	inOrderOrGroupBy bool // Inside ORDER BY, GROUP BY, PARTITION BY clause
}

func newNormalizer(sql string) *normalizer {
	return &normalizer{
		lexer: NewLexer(sql),
	}
}

func (n *normalizer) normalize() string {
	n.lexer.SetSQLMode(n.opts.SQLMode)

	for {
		tok := n.lexer.Lex()

		if tok.Type == END_OF_INPUT {
			break
		}

		if tok.Type == ABORT_SYM {
			// Error in SQL, return what we have
			break
		}

		n.absorbedSign = false
		n.processToken(tok)

		// Update token history, but if we absorbed a unary sign,
		// don't update lastToken (it was already restored)
		if !n.absorbedSign {
			n.secondLastToken = n.lastToken
			n.lastToken = tok.Type
		}
	}

	result := strings.TrimSpace(n.builder.String())

	// Trailing semicolon removal
	if strings.HasSuffix(result, ";") {
		result = strings.TrimSuffix(result, ";")
		result = strings.TrimSpace(result)
	}

	// Apply max length if set
	if n.opts.MaxLength > 0 && len(result) > n.opts.MaxLength {
		result = result[:n.opts.MaxLength] + "..."
	}

	return result
}

func (n *normalizer) processToken(tok Token) {
	tokType := tok.Type

	// Track VALUES keyword for row collapsing
	if tokType == VALUES {
		n.inValues = true
		n.valuesDepth = n.parenDepth
		n.sawFirstRow = false
		n.skippingRows = false
	}

	// If we're skipping subsequent VALUES rows, skip everything until VALUES context ends
	if n.skippingRows {
		// Check for end of VALUES context
		if tokType == ON_SYM || tokType == WHERE || tokType == SET_SYM || tokType == ';' || tokType == END_OF_INPUT {
			n.inValues = false
			n.skippingRows = false
			// Fall through to process this token
		} else {
			// Skip this token
			return
		}
	}

	// Track ORDER BY / GROUP BY context
	if tokType == ORDER_SYM || tokType == GROUP_SYM || tokType == PARTITION_SYM {
		n.inOrderOrGroupBy = true
	} else if n.inOrderOrGroupBy {
		// Reset on clauses that might end ORDER/GROUP BY
		// Note: Most of these actually come *before* ORDER BY, but some
		// like LIMIT, OFFSET, FOR (UPDATE), LOCK (IN SHARE MODE), PROCEDURE
		// can validly follow.
		switch tokType {
		case LIMIT, OFFSET_SYM, FOR_SYM, LOCK_SYM, PROCEDURE_SYM, SELECT_SYM, UPDATE_SYM, DELETE_SYM, INSERT_SYM, UNION_SYM, END_OF_INPUT, ';':
			n.inOrderOrGroupBy = false
		}
	}

	// Handle pending sign (unary +/-)
	// If we have a pending sign and the next token is a numeric literal,
	// and the token before the sign can start an expression, absorb the sign
	if n.pendingSign != 0 {
		if n.isNumericLiteral(tokType) && n.startsExpression(n.secondLastToken) {
			// Unary operator - absorb into the literal
			// Restore lastToken to what it was before the sign
			n.lastToken = n.secondLastToken
			n.absorbedSign = true
			n.pendingSign = 0
			n.pendingSignToken = 0
			// Fall through to process the literal normally
		} else {
			// Binary operator - output the pending sign
			n.appendToken(string(rune(n.pendingSign)), n.pendingSignToken)
			n.pendingSign = 0
			n.pendingSignToken = 0
		}
	}

	// Check if this is a potential unary +/- (defer output)
	if (tokType == '+' || tokType == '-') && n.startsExpression(n.lastToken) {
		n.pendingSign = tokType
		n.pendingSignToken = tokType
		return
	}

	// Handle literals - replace with ?
	if n.isLiteral(tokType) {
		// Special handling for numeric columns in ORDER BY / GROUP BY
		// If we are in an ORDER/GROUP BY clause, and we see a numeric literal,
		// and it is preceded by BY, comma, or open paren (for rollup/cube),
		// we keep the literal as-is instead of replacing with ?.
		if n.inOrderOrGroupBy && n.isNumericLiteral(tokType) {
			if n.lastToken == BY || n.lastToken == ',' || n.lastToken == '(' {
				// Don't replace with ?, keep the number
				// But we still need to append it.
				// Since we are returning early here, we need to handle the append manually
				// or just fall through to the identifier handling logic?
				// Actually, we can just *skip* the "return" here and let it fall through
				// to the end where it appends text.
				// BUT: The default logic at the end wraps identifiers in backticks.
				// Literals shouldn't be wrapped in backticks.
				// So we should append it here and return.

				text := n.lexer.TokenText(tok)
				n.appendToken(text, tokType)
				return
			}
		}

		// Check if we're in a list context (IN clause or VALUES row)
		if n.inList && n.parenDepth >= n.listStart {
			if n.sawFirstLiteral {
				// We already emitted the summary '?' for this list.
				// Skip this token.
				return
			}

			// If this is the first token of interest (literal), emit '?' and mark as seen
			n.sawFirstLiteral = true
			n.appendToken("?", tokType)
			return
		}

		n.appendToken("?", tokType)
		return
	}

	// Track parentheses for list collapsing
	if tokType == '(' {
		n.parenDepth++
		// Check if this starts an IN list
		if n.lastToken == IN_SYM {
			n.inList = true
			n.listStart = n.parenDepth
			n.sawFirstLiteral = false
		}
		// Check if this starts a VALUES row
		if n.inValues && n.parenDepth == n.valuesDepth+1 {
			// This is a row tuple in VALUES
			if n.sawFirstRow {
				// Start skipping subsequent rows
				n.skippingRows = true
				n.skipValuesRow()
				n.parenDepth-- // skipValuesRow consumed the closing paren
				return
			}
			n.inList = true
			n.listStart = n.parenDepth
			n.sawFirstLiteral = false
		}

		// If we are IN a list, and we see a '(', it might be a row constructor IN ((...))
		// We want to collapse the whole list to IN (?)
		if n.inList && n.parenDepth == n.listStart+1 {
			if !n.sawFirstLiteral {
				n.sawFirstLiteral = true
				n.appendToken("?", tokType)
				// We continue, but subsequent tokens will be skipped by the sawFirstLiteral check
				// until we hit the closing ')' of the list.
				// However, we need to be careful NOT to increment parenDepth again?
				// We already did n.parenDepth++ at the top.
				// If we skip subsequent tokens, we need to track parenDepth so we know when the list ends.
				return
			}
			// If we already saw the first literal, this '(' is just skipped content.
			return
		}
	} else if tokType == ')' {
		if n.inList && n.parenDepth == n.listStart {
			n.inList = false
			n.sawFirstLiteral = false
			// If we're in VALUES context and just closed a row tuple
			if n.inValues && n.parenDepth == n.valuesDepth+1 {
				n.sawFirstRow = true
				n.skippingRows = true // Start skipping after first row closes
			}
		}
		n.parenDepth--

		// If we are skipping content in a list, we still need to process the closing paren of that list
		// to exit the list state.
		// The logic above handles exiting the state.
		// But if we are deeper in the list (sawFirstLiteral is true), we should skip this ')' UNLESS it closes the list.
		if n.inList && n.sawFirstLiteral && n.parenDepth >= n.listStart {
			// It was skipped content.
			return
		}
	}

	// Skip commas inside collapsed lists
	if tokType == ',' && n.inList && n.sawFirstLiteral {
		return
	}

	// Reset VALUES context on certain keywords
	if n.inValues && (tokType == ON_SYM || tokType == WHERE || tokType == SET_SYM || tokType == ';') {
		n.inValues = false
	}

	// Get token text
	var text string
	if tokType < 256 {
		// Single character token
		text = string(rune(tokType))
	} else {
		// Use token info for named tokens
		text = TokenString(tokType)
		if text == "" || text == "(unknown)" {
			// Fall back to lexer's token text for identifiers
			text = n.lexer.TokenText(tok)
		}
	}

	// Handle identifiers - wrap in backticks like MySQL does
	if tokType == IDENT || tokType == IDENT_QUOTED {
		identText := n.lexer.TokenText(tok)
		// For IDENT_QUOTED, strip existing delimiters and get raw identifier
		if tokType == IDENT_QUOTED {
			identText = n.stripIdentifierQuotes(identText)
		}
		// Wrap in backticks, escaping any embedded backticks
		text = "`" + n.escapeBackticks(identText) + "`"
	}

	n.appendToken(text, tokType)
}

// skipValuesRow skips tokens until the end of the current VALUES row
func (n *normalizer) skipValuesRow() {
	depth := 1 // We already saw the opening paren
	for depth > 0 {
		tok := n.lexer.Lex()
		if tok.Type == END_OF_INPUT || tok.Type == ABORT_SYM {
			break
		}
		if tok.Type == '(' {
			depth++
		} else if tok.Type == ')' {
			depth--
		}
	}
	// Update our paren depth - we consumed the closing paren
	// No need to update since we didn't increment for the skipped row
}

func (n *normalizer) appendToken(text string, tokType int) {
	// Add space before token if needed
	if n.builder.Len() > 0 && n.needsSpaceBefore(tokType) {
		n.builder.WriteByte(' ')
	}

	n.builder.WriteString(text)
}

func (n *normalizer) needsSpaceBefore(tokType int) bool {
	// No space after opening paren or before closing paren
	if n.lastToken == '(' || tokType == ')' {
		return false
	}
	// No space before comma
	if tokType == ',' {
		return false
	}
	// No space after dot or before dot
	if n.lastToken == '.' || tokType == '.' {
		return false
	}
	// No space between @ symbols
	if n.lastToken == '@' && tokType == '@' {
		return false
	}
	// No space after @
	if n.lastToken == '@' {
		return false
	}
	return true
}

func (n *normalizer) isLiteral(tokType int) bool {
	switch tokType {
	case NUM, LONG_NUM, ULONGLONG_NUM, DECIMAL_NUM, FLOAT_NUM:
		return true
	case TEXT_STRING, NCHAR_STRING:
		return true
	case HEX_NUM, BIN_NUM:
		return true
	case LEX_HOSTNAME:
		return true
	case NULL_SYM:
		// NULL is treated as a literal in digest (becomes ?)
		return true
	case PARAM_MARKER:
		// Parameter markers (?) are already placeholders
		return true
	}
	return false
}

// isNumericLiteral returns true if the token is a numeric literal
func (n *normalizer) isNumericLiteral(tokType int) bool {
	switch tokType {
	case NUM, LONG_NUM, ULONGLONG_NUM, DECIMAL_NUM, FLOAT_NUM, HEX_NUM, BIN_NUM:
		return true
	}
	return false
}

// startsExpression returns true if the token can start an expression
// (meaning a following +/- would be unary, not binary)
func (n *normalizer) startsExpression(tokType int) bool {
	// Start of input
	if tokType == 0 {
		return true
	}

	// Operators and punctuation that start expressions
	switch tokType {
	case '(', ',', '=', '+', '-', '*', '/', '%', '^', '~':
		return true
	case EQ, NE, LT, LE, GT_SYM, GE, EQUAL_SYM: // Comparison operators
		return true
	case AND_AND_SYM, OR_OR_SYM, AND_SYM, OR_SYM, XOR, NOT_SYM: // Logical
		return true
	case BETWEEN_SYM, IN_SYM, LIKE, REGEXP: // Predicates
		return true
	case SELECT_SYM, WHERE, HAVING, SET_SYM, VALUES, CASE_SYM, WHEN_SYM, THEN_SYM, ELSE:
		return true
	case RETURN_SYM, IF, WHILE_SYM, UNTIL_SYM:
		return true
	case BY: // ORDER BY, GROUP BY
		return true
	case LIMIT, OFFSET_SYM:
		return true
	case AS: // For expressions like CAST(x AS type)
		return true
	case SHIFT_LEFT, SHIFT_RIGHT: // Bit shift operators
		return true
	}
	return false
}

// stripIdentifierQuotes removes surrounding quotes from a quoted identifier
// Handles backticks (`ident`) and double quotes ("ident" in ANSI_QUOTES mode)
func (n *normalizer) stripIdentifierQuotes(s string) string {
	if len(s) < 2 {
		return s
	}

	// Check for backtick-quoted identifier
	if s[0] == '`' && s[len(s)-1] == '`' {
		// Remove backticks and unescape doubled backticks
		inner := s[1 : len(s)-1]
		return strings.ReplaceAll(inner, "``", "`")
	}

	// Check for double-quoted identifier (ANSI_QUOTES mode)
	if s[0] == '"' && s[len(s)-1] == '"' {
		// Remove quotes and unescape doubled quotes
		inner := s[1 : len(s)-1]
		return strings.ReplaceAll(inner, `""`, `"`)
	}

	return s
}

// escapeBackticks escapes backticks in an identifier by doubling them
func (n *normalizer) escapeBackticks(s string) string {
	return strings.ReplaceAll(s, "`", "``")
}

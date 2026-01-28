package digest

// reducer applies MySQL's normalization reduction rules to the token stack.
// These rules collapse literals, IN clauses, and VALUES lists into normalized forms.
type reducer struct {
	store *tokenStore
}

// newReducer creates a new reducer operating on the given token store.
func newReducer(store *tokenStore) *reducer {
	return &reducer{store: store}
}

// reduceUnarySign absorbs unary +/- signs before numeric literals.
// A sign is considered unary if the preceding token starts an expression
// (e.g., '(', ',', operators, or start of input).
func (r *reducer) reduceUnarySign() {
	for {
		p := r.store.peek(2)
		last, prev := p[1], p[0]
		if (last == '+' || last == '-') && startsExpression(prev) {
			r.store.pop(1)
		} else {
			break
		}
	}
}

// reduceAfterValue attempts value list reduction after adding a value token.
// Checks for pattern: VALUE/VALUE_LIST ',' VALUE → VALUE_LIST
func (r *reducer) reduceAfterValue() {
	p := r.store.peek(3)
	// p[0]=third-to-last, p[1]=second-to-last, p[2]=last (just added VALUE)
	if p[1] == ',' && isValueOrValueList(p[0]) {
		r.store.pop(3)
		r.store.push(TOK_GENERIC_VALUE_LIST)
	}
	r.reduceAll()
}

// reduceAll applies all applicable reductions until no more apply.
func (r *reducer) reduceAll() {
	for {
		if r.reduceParenthesizedValue() {
			continue
		}
		if r.reduceRowList() {
			continue
		}
		if r.reduceInClause() {
			continue
		}
		break
	}
}

// reduceParenthesizedValue handles: '(' VALUE ')' → ROW_VALUE
// - '(' TOK_GENERIC_VALUE ')' → TOK_ROW_SINGLE_VALUE
// - '(' TOK_GENERIC_VALUE_LIST ')' → TOK_ROW_MULTIPLE_VALUE
func (r *reducer) reduceParenthesizedValue() bool {
	if r.store.len() < 3 {
		return false
	}

	p := r.store.peek(3)
	first, mid, last := p[0], p[1], p[2]

	if first != '(' || last != ')' {
		return false
	}

	switch mid {
	case TOK_GENERIC_VALUE:
		r.store.pop(3)
		r.store.push(TOK_ROW_SINGLE_VALUE)
		return true
	case TOK_GENERIC_VALUE_LIST:
		r.store.pop(3)
		r.store.push(TOK_ROW_MULTIPLE_VALUE)
		return true
	}
	return false
}

// reduceRowList handles: ROW ',' ROW → ROW_LIST
// - Single-value rows: (?) , (?) → (?), ...
// - Multi-value rows: (...) , (...) → (...), ...
func (r *reducer) reduceRowList() bool {
	if r.store.len() < 3 {
		return false
	}

	p := r.store.peek(3)
	first, comma, last := p[0], p[1], p[2]

	if comma != ',' {
		return false
	}

	if isSingleValueRow(first) && isSingleValueRow(last) {
		r.store.pop(3)
		r.store.push(TOK_ROW_SINGLE_VALUE_LIST)
		return true
	}

	if isMultiValueRow(first) && isMultiValueRow(last) {
		r.store.pop(3)
		r.store.push(TOK_ROW_MULTIPLE_VALUE_LIST)
		return true
	}

	return false
}

// reduceInClause handles: IN ROW → IN (...)
// Collapses IN clauses to a single normalized form.
// Note: This is MySQL 8.0+ only. MySQL 5.7 doesn't have TOK_IN_GENERIC_VALUE_EXPRESSION.
func (r *reducer) reduceInClause() bool {
	// Skip this reduction for MySQL 5.7 - it doesn't have TOK_IN_GENERIC_VALUE_EXPRESSION
	if r.store.version == MySQL57 {
		return false
	}

	if r.store.len() < 2 {
		return false
	}

	p := r.store.peek(2)
	before, last := p[0], p[1]

	if before != IN_SYM {
		return false
	}

	if last == TOK_ROW_SINGLE_VALUE || last == TOK_ROW_MULTIPLE_VALUE {
		r.store.pop(2)
		r.store.push(TOK_IN_GENERIC_VALUE_EXPRESSION)
		return true
	}

	return false
}

// isValueOrValueList returns true if the token is a value or value list.
func isValueOrValueList(tok int) bool {
	return tok == TOK_GENERIC_VALUE || tok == TOK_GENERIC_VALUE_LIST
}

// isSingleValueRow returns true if the token represents a single-value row.
func isSingleValueRow(tok int) bool {
	return tok == TOK_ROW_SINGLE_VALUE || tok == TOK_ROW_SINGLE_VALUE_LIST
}

// isMultiValueRow returns true if the token represents a multi-value row.
func isMultiValueRow(tok int) bool {
	return tok == TOK_ROW_MULTIPLE_VALUE || tok == TOK_ROW_MULTIPLE_VALUE_LIST
}

// startsExpression returns true if the token can start an expression,
// meaning a following +/- would be unary rather than binary.
func startsExpression(tokType int) bool {
	if tokType == 0 || tokType == TOK_UNUSED {
		return true
	}
	return TokenStartExpr(tokType)
}

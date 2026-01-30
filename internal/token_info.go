package internal

type TokenInfo struct {
	String      string
	Length      int
	AppendSpace bool
	StartExpr   bool
	IsHintable  bool
}

var TokenInfos [1252]TokenInfo

func init() {
	for i := 0; i < 256; i++ {
		TokenInfos[i] = TokenInfo{
			String:      string(rune(i)),
			Length:      1,
			AppendSpace: true,
			StartExpr:   false,
		}
	}

	for i := 256; i < len(TokenInfos); i++ {
		TokenInfos[i] = TokenInfo{
			AppendSpace: true,
			StartExpr:   false,
		}
	}

	TokenInfos[WITH_ROLLUP_SYM].String = "WITH ROLLUP"
	TokenInfos[NOT2_SYM].String = "!"
	TokenInfos[OR2_SYM].String = "||"
	TokenInfos[PARAM_MARKER].String = "?"
	TokenInfos[SET_VAR].String = ":="
	TokenInfos[UNDERSCORE_CHARSET].String = "(_charset)"
	TokenInfos[JSON_SEPARATOR_SYM].String = "->"
	TokenInfos[JSON_UNQUOTED_SEPARATOR_SYM].String = "->>"
	TokenInfos[BIN_NUM].String = "(bin)"
	TokenInfos[DECIMAL_NUM].String = "(decimal)"
	TokenInfos[FLOAT_NUM].String = "(float)"
	TokenInfos[HEX_NUM].String = "(hex)"
	TokenInfos[LEX_HOSTNAME].String = "(hostname)"
	TokenInfos[LONG_NUM].String = "(long)"
	TokenInfos[NUM].String = "(num)"
	TokenInfos[TEXT_STRING].String = "(text)"
	TokenInfos[NCHAR_STRING].String = "(nchar)"
	TokenInfos[ULONGLONG_NUM].String = "(ulonglong)"
	TokenInfos[IDENT].String = "(id)"
	TokenInfos[IDENT_QUOTED].String = "(id_quoted)"
	TokenInfos[TOK_GENERIC_VALUE].String = "?"
	TokenInfos[TOK_GENERIC_VALUE_LIST].String = "?, ..."
	TokenInfos[TOK_ROW_SINGLE_VALUE].String = "(?)"
	TokenInfos[TOK_ROW_SINGLE_VALUE_LIST].String = "(?) /* , ... */"
	TokenInfos[TOK_ROW_MULTIPLE_VALUE].String = "(...)"
	TokenInfos[TOK_ROW_MULTIPLE_VALUE_LIST].String = "(...) /* , ... */"
	TokenInfos[TOK_IN_GENERIC_VALUE_EXPRESSION].String = "IN (...)"
	TokenInfos[TOK_IDENT].String = "(tok_id)"
	TokenInfos[TOK_IDENT_AT].String = "(tok_id_at)"
	TokenInfos[TOK_BY_NUMERIC_COLUMN].String = "(by_num_col)"
	TokenInfos[TOK_UNUSED].String = "UNUSED"

	// MySQL only sets m_append_space = false for '@' token
	// See gen_lex_token.cc: compiled_token_array[(int)'@'].m_append_space = false;
	TokenInfos['@'].AppendSpace = false

	TokenInfos['('].StartExpr = true
	TokenInfos[','].StartExpr = true
	TokenInfos['='].StartExpr = true
	TokenInfos['~'].StartExpr = true
	TokenInfos['+'].StartExpr = true
	TokenInfos['-'].StartExpr = true
	TokenInfos['*'].StartExpr = true
	TokenInfos['/'].StartExpr = true
	TokenInfos['%'].StartExpr = true
	TokenInfos['^'].StartExpr = true
	TokenInfos['|'].StartExpr = true
	TokenInfos['&'].StartExpr = true
	TokenInfos[EQ].StartExpr = true
	TokenInfos[NE].StartExpr = true
	TokenInfos[LT].StartExpr = true
	TokenInfos[LE].StartExpr = true
	TokenInfos[GT_SYM].StartExpr = true
	TokenInfos[GE].StartExpr = true
	TokenInfos[EQUAL_SYM].StartExpr = true
	TokenInfos[AND_AND_SYM].StartExpr = true
	TokenInfos[OR_OR_SYM].StartExpr = true
	TokenInfos[AND_SYM].StartExpr = true
	TokenInfos[OR_SYM].StartExpr = true
	TokenInfos[OR2_SYM].StartExpr = true
	TokenInfos[XOR].StartExpr = true
	TokenInfos[NOT_SYM].StartExpr = true
	TokenInfos[BETWEEN_SYM].StartExpr = true
	TokenInfos[LIKE].StartExpr = true
	TokenInfos[REGEXP].StartExpr = true
	TokenInfos[SELECT_SYM].StartExpr = true
	TokenInfos[WHERE].StartExpr = true
	TokenInfos[HAVING].StartExpr = true
	TokenInfos[SET_SYM].StartExpr = true
	TokenInfos[VALUES].StartExpr = true
	TokenInfos[CASE_SYM].StartExpr = true
	TokenInfos[WHEN_SYM].StartExpr = true
	TokenInfos[THEN_SYM].StartExpr = true
	TokenInfos[ELSE].StartExpr = true
	TokenInfos[RETURN_SYM].StartExpr = true
	TokenInfos[IF].StartExpr = true
	TokenInfos[ELSEIF_SYM].StartExpr = true
	TokenInfos[WHILE_SYM].StartExpr = true
	TokenInfos[UNTIL_SYM].StartExpr = true
	TokenInfos[BY].StartExpr = true
	TokenInfos[LIMIT].StartExpr = true
	TokenInfos[OFFSET_SYM].StartExpr = true
	TokenInfos[AS].StartExpr = true
	TokenInfos[SHIFT_LEFT].StartExpr = true
	TokenInfos[SHIFT_RIGHT].StartExpr = true
	TokenInfos[INTERVAL_SYM].StartExpr = true
	TokenInfos[DIV_SYM].StartExpr = true
	TokenInfos[MOD_SYM].StartExpr = true
	TokenInfos[EVERY_SYM].StartExpr = true
	TokenInfos[AT_SYM].StartExpr = true
	TokenInfos[STARTS_SYM].StartExpr = true
	TokenInfos[ENDS_SYM].StartExpr = true
	TokenInfos[DEFAULT_SYM].StartExpr = true
	TokenInfos[IN_SYM].StartExpr = true

	// Hint comment tokens
	TokenInfos[TOK_HINT_COMMENT_OPEN].String = "/*+"
	TokenInfos[TOK_HINT_COMMENT_CLOSE].String = "*/"

	// Hintable statement tokens
	TokenInfos[SELECT_SYM].IsHintable = true
	TokenInfos[INSERT_SYM].IsHintable = true
	TokenInfos[UPDATE_SYM].IsHintable = true
	TokenInfos[DELETE_SYM].IsHintable = true
	TokenInfos[REPLACE_SYM].IsHintable = true
}

func TokenString(tok int) string {
	if tok >= 0 && tok < len(TokenInfos) {
		s := TokenInfos[tok].String
		if s != "" {
			return s
		}
	}
	if tok < 256 {
		return string(rune(tok))
	}
	return "(unknown)"
}

func TokenAppendSpace(tok int) bool {
	if tok >= 0 && tok < len(TokenInfos) {
		return TokenInfos[tok].AppendSpace
	}
	return true
}

func TokenStartExpr(tok int) bool {
	if tok >= 0 && tok < len(TokenInfos) {
		return TokenInfos[tok].StartExpr
	}
	return false
}

func TokenIsHintable(tok int) bool {
	if tok >= 0 && tok < len(TokenInfos) {
		return TokenInfos[tok].IsHintable
	}
	return false
}

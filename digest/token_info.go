package digest

// TokenInfo holds metadata about a token
type TokenInfo struct {
	String      string
	Length      int
	AppendSpace bool // If true, add space after this token (delayed until next token)
	StartExpr   bool
	IsHintable  bool
}

// TokenInfos maps token IDs to their metadata
// Size must be large enough to hold all tokens.
// Max token is around 1251 (MY_MAX_TOKEN)
var TokenInfos [1252]TokenInfo

func init() {
	// 1. Initialize all tokens with default values
	// By default: AppendSpace = true, StartExpr = false
	// MySQL's gen_lex_token.cc sets m_append_space = true for all tokens by default
	for i := 0; i < 256; i++ {
		TokenInfos[i] = TokenInfo{
			String:      string(rune(i)),
			Length:      1,
			AppendSpace: true,
			StartExpr:   false,
		}
	}

	// For named tokens (256+), we will rely on overrides or manual setting if needed
	// But gen_lex_token.cc initializes them.
	// We'll set a default "AppendSpace=true" for everything else too, as per `gen_lex_token_string` constructor
	for i := 256; i < len(TokenInfos); i++ {
		TokenInfos[i] = TokenInfo{
			AppendSpace: true,
			StartExpr:   false,
		}
	}

	// 2. Apply generated overrides
	initTokenInfoOverrides()
}

// TokenString returns the string representation of a token ID
func TokenString(tok int) string {
	if tok >= 0 && tok < len(TokenInfos) {
		s := TokenInfos[tok].String
		if s != "" {
			return s
		}
	}
	// Fallback/Default behavior if string is missing?
	// In C++, if not initialized, it uses "(unknown)" or dummy.
	// For single chars, it's the char itself.
	if tok < 256 {
		return string(rune(tok))
	}
	return "(unknown)"
}

// TokenAppendSpace returns whether a space should be appended after this token
func TokenAppendSpace(tok int) bool {
	if tok >= 0 && tok < len(TokenInfos) {
		return TokenInfos[tok].AppendSpace
	}
	return true
}

// TokenStartExpr returns whether this token starts an expression
func TokenStartExpr(tok int) bool {
	if tok >= 0 && tok < len(TokenInfos) {
		return TokenInfos[tok].StartExpr
	}
	return false
}

// TokenIsHintable returns whether this token can be followed by optimizer hints /*+ ... */
func TokenIsHintable(tok int) bool {
	if tok >= 0 && tok < len(TokenInfos) {
		return TokenInfos[tok].IsHintable
	}
	return false
}

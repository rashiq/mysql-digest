package digest

type TokenConfig struct {
	Version      MySQLVersion
	Keywords     map[string]int
	TokenStrings map[int]string
	HashTokens   map[int]int
}

func (c *TokenConfig) LookupKeyword(word string) int {
	return c.Keywords[word]
}

func (c *TokenConfig) GetString(tok int) string {
	if c.TokenStrings != nil {
		if s, ok := c.TokenStrings[tok]; ok {
			return s
		}
	}
	return TokenString(tok)
}

func (c *TokenConfig) TranslateForHash(tok int) int {
	if c.HashTokens == nil {
		return tok
	}
	if tok < 256 {
		return tok
	}
	if translated, ok := c.HashTokens[tok]; ok {
		return translated
	}
	return m57TOK_UNUSED
}

var (
	configMySQL80 *TokenConfig
	configMySQL84 *TokenConfig
	configMySQL57 *TokenConfig
)

func init() {
	configMySQL80 = buildMySQL80Config()
	configMySQL84 = buildMySQL84Config()
	configMySQL57 = buildMySQL57Config()
}

func GetTokenConfig(v MySQLVersion) *TokenConfig {
	switch v {
	case MySQL57:
		return configMySQL57
	case MySQL84:
		return configMySQL84
	default:
		return configMySQL80
	}
}

func buildMySQL80Config() *TokenConfig {
	return &TokenConfig{
		Version:  MySQL80,
		Keywords: TokenKeywords,
	}
}

func buildMySQL84Config() *TokenConfig {
	return &TokenConfig{
		Version:  MySQL84,
		Keywords: TokenKeywords,
	}
}

func buildMySQL57Config() *TokenConfig {
	keywords := make(map[string]int, len(TokenKeywords)+len(mysql57Keywords))
	for k, v := range TokenKeywords {
		if mapped, ok := mysql80To57TokenMap[v]; ok && mapped != m57TOK_UNUSED {
			keywords[k] = v
		}
	}
	for k, v := range mysql57Keywords {
		keywords[k] = v
	}

	tokenStrings := map[int]string{
		OBSOLETE_TOKEN_271: "ANALYSE",
		OBSOLETE_TOKEN_388: "DES_KEY_FILE",
		OBSOLETE_TOKEN_538: "LOCATOR",
		OBSOLETE_TOKEN_550: "MASTER_AUTO_POSITION",
		OBSOLETE_TOKEN_551: "MASTER_BIND",
		OBSOLETE_TOKEN_552: "MASTER_CONNECT_RETRY",
		OBSOLETE_TOKEN_553: "MASTER_DELAY",
		OBSOLETE_TOKEN_554: "MASTER_HOST",
		OBSOLETE_TOKEN_555: "MASTER_LOG_FILE",
		OBSOLETE_TOKEN_556: "MASTER_LOG_POS",
		OBSOLETE_TOKEN_557: "MASTER_PASSWORD",
		OBSOLETE_TOKEN_558: "MASTER_PORT",
		OBSOLETE_TOKEN_559: "MASTER_RETRY_COUNT",
		OBSOLETE_TOKEN_561: "MASTER_SERVER_ID",
		OBSOLETE_TOKEN_562: "MASTER_SSL",
		OBSOLETE_TOKEN_563: "MASTER_SSL_CA",
		OBSOLETE_TOKEN_564: "MASTER_SSL_CAPATH",
		OBSOLETE_TOKEN_565: "MASTER_SSL_CERT",
		OBSOLETE_TOKEN_566: "MASTER_SSL_CIPHER",
		OBSOLETE_TOKEN_567: "MASTER_SSL_CRL",
		OBSOLETE_TOKEN_568: "MASTER_SSL_CRLPATH",
		OBSOLETE_TOKEN_569: "MASTER_SSL_KEY",
		OBSOLETE_TOKEN_570: "MASTER_SSL_VERIFY_SERVER_CERT",
		OBSOLETE_TOKEN_572: "MASTER_TLS_VERSION",
		OBSOLETE_TOKEN_573: "MASTER_USER",
		OBSOLETE_TOKEN_654: "PARSE_GCOL_EXPR",
		OBSOLETE_TOKEN_693: "REDOFILE",
		OBSOLETE_TOKEN_755: "SERVER_OPTIONS",
		OBSOLETE_TOKEN_784: "SQL_CACHE",
		OBSOLETE_TOKEN_820: "TABLE_REF_PRIORITY",
		OBSOLETE_TOKEN_848: "UDF_RETURNS",
		OBSOLETE_TOKEN_893: "WITH_CUBE",
	}

	return &TokenConfig{
		Version:      MySQL57,
		Keywords:     keywords,
		TokenStrings: tokenStrings,
		HashTokens:   mysql80To57TokenMap,
	}
}

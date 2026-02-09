package internal

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
	configMySQL90 *TokenConfig
	configMySQL57 *TokenConfig
)

func init() {
	configMySQL80 = buildMySQL80Config()
	configMySQL84 = buildMySQL84Config()
	configMySQL90 = buildMySQL90Config()
	configMySQL57 = buildMySQL57Config()
}

func GetTokenConfig(v MySQLVersion) *TokenConfig {
	switch v {
	case MySQL57:
		return configMySQL57
	case MySQL84:
		return configMySQL84
	case MySQL90:
		return configMySQL90
	default:
		return configMySQL80
	}
}

// Keywords only in MySQL 8.0 (removed in 8.4+)
var mysql80OnlyKeywords = map[string]bool{
	"GET_MASTER_PUBLIC_KEY":         true,
	"MASTER_AUTO_POSITION":          true,
	"MASTER_BIND":                   true,
	"MASTER_COMPRESSION_ALGORITHMS": true,
	"MASTER_CONNECT_RETRY":          true,
	"MASTER_DELAY":                  true,
	"MASTER_HEARTBEAT_PERIOD":       true,
	"MASTER_HOST":                   true,
	"MASTER_LOG_FILE":               true,
	"MASTER_LOG_POS":                true,
	"MASTER_PASSWORD":               true,
	"MASTER_PORT":                   true,
	"MASTER_PUBLIC_KEY_PATH":        true,
	"MASTER_RETRY_COUNT":            true,
	"MASTER_SSL":                    true,
	"MASTER_SSL_CA":                 true,
	"MASTER_SSL_CAPATH":             true,
	"MASTER_SSL_CERT":               true,
	"MASTER_SSL_CIPHER":             true,
	"MASTER_SSL_CRL":                true,
	"MASTER_SSL_CRLPATH":            true,
	"MASTER_SSL_KEY":                true,
	"MASTER_SSL_VERIFY_SERVER_CERT": true,
	"MASTER_TLS_CIPHERSUITES":       true,
	"MASTER_TLS_VERSION":            true,
	"MASTER_ZSTD_COMPRESSION_LEVEL": true,
}

// Keywords added in MySQL 8.4
var mysql84Keywords = map[string]bool{
	"AUTO":        true,
	"BERNOULLI":   true,
	"GTIDS":       true,
	"LOG":         true,
	"MANUAL":      true,
	"PARALLEL":    true,
	"PARSE_TREE":  true,
	"QUALIFY":     true,
	"S3":          true,
	"TABLESAMPLE": true,
}

// Keywords added in MySQL 9.0
var mysql90Keywords = map[string]bool{
	"ABSENT":                 true,
	"ALLOW_MISSING_FILES":    true,
	"AUTO_REFRESH":           true,
	"AUTO_REFRESH_SOURCE":    true,
	"DUALITY":                true,
	"EXTERNAL":               true,
	"EXTERNAL_FORMAT":        true,
	"FILES":                  true,
	"FILE_FORMAT":            true,
	"FILE_NAME":              true,
	"FILE_PATTERN":           true,
	"FILE_PREFIX":            true,
	"GUIDED":                 true,
	"HEADER":                 true,
	"JSON_DUALITY_OBJECT":    true,
	"LIBRARY":                true,
	"MATERIALIZED":           true,
	"PARAMETERS":             true,
	"RELATIONAL":             true,
	"SETS":                   true,
	"STRICT_LOAD":            true,
	"URI":                    true,
	"VALIDATE":               true,
	"VECTOR":                 true,
	"VERIFY_KEY_CONSTRAINTS": true,
}

func buildKeywordsFor(version MySQLVersion) map[string]int {
	keywords := make(map[string]int, len(TokenKeywords))
	for k, v := range TokenKeywords {
		keywords[k] = v
	}

	switch version {
	case MySQL80:
		for k := range mysql84Keywords {
			delete(keywords, k)
		}
		for k := range mysql90Keywords {
			delete(keywords, k)
		}
	case MySQL84:
		for k := range mysql80OnlyKeywords {
			delete(keywords, k)
		}
		for k := range mysql90Keywords {
			delete(keywords, k)
		}
	case MySQL90:
		for k := range mysql80OnlyKeywords {
			delete(keywords, k)
		}
	}

	return keywords
}

func buildMySQL80Config() *TokenConfig {
	tokenStrings := map[int]string{
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
		OBSOLETE_TOKEN_573: "MASTER_HEARTBEAT_PERIOD",
		OBSOLETE_TOKEN_966: "MASTER_PUBLIC_KEY_PATH",
		OBSOLETE_TOKEN_967: "GET_MASTER_PUBLIC_KEY",
		OBSOLETE_TOKEN_989: "MASTER_COMPRESSION_ALGORITHMS",
		OBSOLETE_TOKEN_990: "MASTER_ZSTD_COMPRESSION_LEVEL",
		OBSOLETE_TOKEN_992: "MASTER_TLS_CIPHERSUITES",
	}
	return &TokenConfig{
		Version:      MySQL80,
		Keywords:     buildKeywordsFor(MySQL80),
		TokenStrings: tokenStrings,
	}
}

func buildMySQL84Config() *TokenConfig {
	return &TokenConfig{
		Version:  MySQL84,
		Keywords: buildKeywordsFor(MySQL84),
	}
}

func buildMySQL90Config() *TokenConfig {
	return &TokenConfig{
		Version:  MySQL90,
		Keywords: buildKeywordsFor(MySQL90),
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

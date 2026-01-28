package digest

// LexerConfig holds configuration options for the lexer.
// This enables dependency injection and customization.
type LexerConfig struct {
	// SQLMode contains SQL mode flags affecting lexer behavior.
	SQLMode SQLMode

	// PrepareMode enables prepared statement mode where '?' is a parameter marker.
	PrepareMode bool

	// MySQLVersion is the target MySQL version (e.g., 80400 for 8.4.0).
	// Used for version comments like /*!80000 ... */.
	MySQLVersion int

	// KeywordResolver resolves identifiers to keyword token types.
	// If nil, DefaultKeywordResolver is used.
	KeywordResolver KeywordResolver

	// StateMapper maps characters to initial lexer states.
	// If nil, DefaultStateMapper is used.
	StateMapper StateMapper
}

// DefaultMySQLVersion is the default MySQL version (8.4.0).
const DefaultMySQLVersion = 80400

// DefaultConfig returns a LexerConfig with default settings.
func DefaultConfig() LexerConfig {
	return LexerConfig{
		SQLMode:         0,
		PrepareMode:     false,
		MySQLVersion:    DefaultMySQLVersion,
		KeywordResolver: DefaultKeywordResolver{},
		StateMapper:     DefaultStateMapper{},
	}
}

// NewLexerWithConfig creates a new lexer with the given configuration.
func NewLexerWithConfig(input string, config LexerConfig) *Lexer {
	// Apply defaults for nil interfaces
	if config.KeywordResolver == nil {
		config.KeywordResolver = DefaultKeywordResolver{}
	}
	if config.StateMapper == nil {
		config.StateMapper = DefaultStateMapper{}
	}
	if config.MySQLVersion == 0 {
		config.MySQLVersion = DefaultMySQLVersion
	}

	return &Lexer{
		input:           input,
		pos:             0,
		tokStart:        0,
		nextState:       MY_LEX_START,
		sqlMode:         config.SQLMode,
		stmtPrepareMode: config.PrepareMode,
		keywordResolver: config.KeywordResolver,
		stateMapper:     config.StateMapper,
		mysqlVersion:    config.MySQLVersion,
	}
}

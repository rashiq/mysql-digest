package internal

// MySQLVersion represents a MySQL version for digest computation.
type MySQLVersion int

const (
	MySQL80 MySQLVersion = iota
	MySQL84
	MySQL57
)

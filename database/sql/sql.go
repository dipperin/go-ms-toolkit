package sql

import "database/sql"

type conn struct {
	*sql.DB
	// todo need add breaker
	conf *Config
}


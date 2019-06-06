package sql

import "time"

type Config struct {
	DSN          string        // write data source name.
	ReadDSN      []string      // read data source name.
	Active       int           // pool
	Idle         int           // pool
	IdleTimeout  time.Duration // connect max life time.
	QueryTimeout time.Duration // query sql timeout
	ExecTimeout  time.Duration // execute sql timeout
	TranTimeout  time.Duration // transaction sql timeout
}



package db

import (
	"database/sql"
	"time"
)

type DBCfg struct {
	Name        string        `json:"Name"`
	Protocol    string        `json:"Protocol"`
	IP          string        `json:"IP"`
	Port        string        `json:"Port"`
	DBName      string        `json:"DBName"`
	Network     string        `json:"Network"`
	UserName    string        `json:"UserName"`
	PassWord    string        `json:"PassWord"`
	NlsLang     string        `json:"NlsLang"`
	LibDir      string        `json:"LibDir"`
	MaxOpen     int           `json:"MaxOpen"`
	MaxIdle     int           `json:"MaxIdle"`
	MaxLifetime time.Duration `json:"MaxLifetime"`
}

type DBConn struct {
	*sql.DB
	*DBCfg
}

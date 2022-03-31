package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
)

// New
func NewDBConn(dbcfg DBCfg) *DBConn {
	var dbc DBConn
	dbc.DBCfg = &dbcfg
	return &dbc
}

func (conn *DBConn) Connect() error {
	cfg := conn.DBCfg
	var db *sql.DB
	var err error

	switch strings.ToUpper(cfg.Protocol) {
	case "MYSQL":
		// username:password@tcp(IP:Port)/dbname
		dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local", cfg.UserName, cfg.PassWord, cfg.Network, cfg.IP, cfg.Port, cfg.DBName)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Open mysql failed,err:%v\n", err)
			return err
		}
	case "POSTGRESQL":
	case "PGSQL":
		dsn := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", cfg.IP, cfg.Port, cfg.UserName, cfg.PassWord, cfg.DBName)
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Open DB connection failed,err:%v\n", err)
			return err
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Open DB connection failed,err:%v\n", err)
			return err
		}
	case "ORACLE":
		_ = os.Setenv("NLS_LANG", cfg.NlsLang)
		// user/name@host:port/sid
		dsn := fmt.Sprintf("%s/%s@%s:%s/%s", cfg.UserName, cfg.PassWord, cfg.IP, cfg.Port, cfg.DBName)
		// dsn := `user="`+cfg.UserName+`" password="`+cfg.PassWord+`" connectString="`+cfg.IP+`:`+cfg.Port+`/`cfg.DBName+`" libDir="`+cfg.Libdir+`"`
		db, err = sql.Open("godror", dsn)
		if err != nil {
			log.Printf("Open DB connection failed,err:%v\n", err)
			return err
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Open DB connection failed,err:%v\n", err)
			return err
		}
	}

	// 设置连接池，20个最大连接，5个闲置连接，超过60秒自动断开
	if conn.DBCfg.MaxOpen == 0 {
		conn.DBCfg.MaxOpen = 20
	}
	if conn.DBCfg.MaxIdle == 0 {
		conn.DBCfg.MaxIdle = 5
	}

	db.SetConnMaxLifetime(conn.DBCfg.MaxLifetime * time.Second) //最大连接周期，超过时间的连接就close
	db.SetMaxOpenConns(conn.DBCfg.MaxOpen)                      //设置最大连接数
	db.SetMaxIdleConns(conn.DBCfg.MaxIdle)                      //设置闲置连接数

	conn.DB = db

	return nil
}

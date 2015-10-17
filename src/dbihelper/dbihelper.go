package dbihelper

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v1"
)

const (
	DEBUG = true

	// Postgresql default data
	poolSize   = 10
	hostDef    = "127.0.0.1"
	portDef    = "5432"
	sslModeDef = "disable"
	loginDef   = "test"
	passDef    = "test"
	dbNameDef  = "test"
)

func GetDBI() *gorp.DbMap {

	//"host=127.0.0.1 port=5432 user=my_pg_login password=my_pg_pass dbname=<db> sslmode=disable"
	dsnLine := pg_dns()

	db, err := _dbi(poolSize, dsnLine)
	if err != nil {
		log.Fatalln(err)
	}

	return &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
}

func pg_dns() string {

	login := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWD")
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	dbName := os.Getenv("DBNAME")
	sslMode := os.Getenv("SSLMODE")

	usr := ""
	if login != "" {
		usr = " user=" + login
		if pass != "" {
			usr += " password=" + pass
		}
	} else {
		usr = " user=" + loginDef + " password=" + passDef
	}

	if host == "" {
		host = hostDef
	}

	if port == "" {
		port = portDef
	}

	if dbName == "" {
		dbName = dbNameDef
	}

	if sslMode == "" {
		sslMode = sslModeDef
	}

	return "host=" + host + " port=" + port + usr + " dbname=" + dbName + " sslmode=" + sslMode
}

func _dbi(poolSize int, dsn string) (*sql.DB, error) {

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(poolSize)

	return db, nil
}

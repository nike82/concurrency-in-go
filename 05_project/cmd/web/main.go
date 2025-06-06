package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/stores/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

func main() {
	// connect to db
	db := initDb()
	db.Ping()
	// create sessions
	session := initSession()

	// create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate| log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate| log.Ltime| log.Lshortfile)
	// create channels

	// create wait group
	wg := sync.WaitGroup{}
	// set up the  pplication config
	app := Config {
		Session: session,
		DB: db,
		Infolog: infoLog,
		ErrorLog: errorLog,
		Wait: &wg,
	}
	// set up mail

	// listen for web connections
	app.serve()

}

func (app * Config) serve() {
	// start http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	app.Infolog.Println("Starting web server...")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func initDb() *sql.DB {
	conn := connectToDB()

	if conn == nil {
		log.Panic("Can't connect to the DB")
	}
	return conn
}

func connectToDB() *sql.DB {
	counts := 0

	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
		} else {
			log.Println("Connected to the DB!")
			return connection
		}

		if counts > 10 {
			return nil
		}

		log.Println("Backing off for 1 sec...")
		time.Sleep(time.Second)
		counts++
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initSession() *scs.SessionManager {
	// setup session
	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session 
}

func initRedis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.ExpandEnv("REDIS"))
		},
	}

	return redisPool
}

package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
	"github.com/mmulyana/sosmed-be/internal/db"
	"github.com/mmulyana/sosmed-be/internal/env"
	"github.com/mmulyana/sosmed-be/internal/store"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Tidak bisa load file .env")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_URL", "root:password@tcp(localhost:3306)/sosmed?parseTime=true"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}

package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Conn struct{}

func (conn *Conn) GetConnection() (context.Context, *sql.DB) {
	godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	ctx := context.Background()
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		log.Println(err)
		return ctx, nil
	}
	return ctx, db
}

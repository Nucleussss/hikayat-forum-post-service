package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(connString string) (*pgxpool.Pool, error) {
	op := "db.InitDB"
	retriesMaxAttemp := 5

	log.Printf("%s Try Init Database with max retry %d times", op, retriesMaxAttemp)

	for i := 0; i < retriesMaxAttemp; i++ {

		// open connection to the database
		pool, err := pgxpool.New(context.Background(), connString)
		if err != nil {
			log.Printf("%s Failed to connect to database (attemp %d): %v", op, i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		// check if the connection is valid by pinging the database
		if err = pool.Ping(context.Background()); err != nil {
			log.Printf("%s Failed to ping database (attemp %d): %v", op, i+1, err)
			pool.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		log.Printf("%s Successfully connected to the database", op)
		return pool, nil
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts", retriesMaxAttemp)
}

func ConnString() string {
	DBHost := os.Getenv("DB_HOST")
	DBUser := os.Getenv("DB_USER")
	DBPass := os.Getenv("DB_PASSWORD")
	DBname := os.Getenv("DB_NAME")
	DBPort := os.Getenv("DB_PORT")
	DBsslMode := os.Getenv("DB_SSL_MODE")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", DBHost, DBUser, DBPass, DBname, DBPort, DBsslMode)
}

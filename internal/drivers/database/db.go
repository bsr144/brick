package database

import (
	"brick/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func NewSQLDatabase(databaseConfig config.Database, log *logrus.Logger) *sql.DB {
	connectionString := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		databaseConfig.Username,
		databaseConfig.Name,
		databaseConfig.Password,
		databaseConfig.Host,
		databaseConfig.Port,
	)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("error while opening db connection")
	}

	err = db.Ping()

	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("error while pinging db")

		db.Close()
	}

	return db

}

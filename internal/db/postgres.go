package db

import (
	"fmt"
	fl "vk/internal/floodcontrol"

	"github.com/jmoiron/sqlx"
)

func NewPosrgresDB(config fl.Config) (*sqlx.DB, error) {
	//Подключение к базе данных PostgreSQL
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		config.DB.Host,
		config.DB.LocalPort,
		config.DB.Username,
		config.DB.Database,
		config.DB.Password,
		config.DB.SSLMode))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

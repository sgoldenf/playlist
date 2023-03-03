package postgresdb

import (
	"fmt"
	ps "github.com/sgoldenf/playlist/api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

func New(config PostgresConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			config.Host, config.Port, config.User, config.DBName, config.Password)),
		&gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(ps.SongInfo{})
	if err != nil {
		return nil, err
	}
	return db, err
}

package server

import (
	ps "github.com/sgoldenf/playlist/api"
	db "github.com/sgoldenf/playlist/db"
	"github.com/sgoldenf/playlist/internal/model/playlist"
	"gorm.io/gorm"
)

type PlaylistService struct {
	ps.UnimplementedPlaylistServiceServer
	DB *gorm.DB
	P  *playlist.Playlist
}

func NewService() (*PlaylistService, error) {
	database, errDB := db.New(db.PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "sgoldenf",
		DBName:   "playlist",
		Password: "sgoldenf",
	})
	if errDB != nil {
		return nil, errDB
	}
	service := &PlaylistService{DB: database}
	service.Init()
	return service, nil
}

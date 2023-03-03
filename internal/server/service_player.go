package server

import (
	"context"
	ps "github.com/sgoldenf/playlist/api"
	"github.com/sgoldenf/playlist/internal/model/playlist"
)

func (s *PlaylistService) Init() {
	res, err := s.GetSongs(context.Background(), &ps.ReadSongsRequest{})
	if err != nil {
		s.P = playlist.NewPlaylist([]*ps.SongInfo{})
	} else {
		s.P = playlist.NewPlaylist(res.Songs)
	}
}

func (s *PlaylistService) Play(_ context.Context, _ *ps.PlayRequest) (*ps.PlayResponse, error) {
	s.P.Play()
	return &ps.PlayResponse{Success: true}, nil
}

func (s *PlaylistService) Pause(_ context.Context, _ *ps.PauseRequest) (*ps.PauseResponse, error) {
	s.P.Pause()
	return &ps.PauseResponse{Success: true}, nil
}

func (s *PlaylistService) Next(_ context.Context, _ *ps.NextSongRequest) (*ps.NextSongResponse, error) {
	s.P.Next()
	return &ps.NextSongResponse{Success: true}, nil
}

func (s *PlaylistService) Prev(_ context.Context, _ *ps.PrevSongRequest) (*ps.PrevSongResponse, error) {
	s.P.Prev()
	return &ps.PrevSongResponse{Success: true}, nil
}

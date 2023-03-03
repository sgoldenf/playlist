package server

import (
	"context"
	"errors"
	"github.com/google/uuid"
	ps "github.com/sgoldenf/playlist/api"
)

func (s *PlaylistService) CreateSong(_ context.Context, req *ps.CreateSongRequest) (*ps.CreateSongResponse, error) {
	info := req.GetSong()
	if info.Title == "" || info.Duration == 0 {
		return nil, errors.New("create song error: empty title/duration==0")
	}
	info.Id = uuid.New().String()
	res := s.DB.Create(info)
	if res.RowsAffected == 0 {
		return nil, errors.New("song creation unsuccessful")
	}
	s.P.AddSong(info)
	return &ps.CreateSongResponse{Song: info}, nil
}

func (s *PlaylistService) GetSong(_ context.Context, req *ps.ReadSongRequest) (*ps.ReadSongResponse, error) {
	var song ps.SongInfo
	res := s.DB.Find(&song, "id = ?", req.GetId())
	if res.RowsAffected == 0 {
		return nil, errors.New("song not found")
	}
	return &ps.ReadSongResponse{Song: &song}, nil
}

func (s *PlaylistService) GetSongs(context.Context, *ps.ReadSongsRequest) (*ps.ReadSongsResponse, error) {
	var songs []*ps.SongInfo
	res := s.DB.Find(&songs)
	if res.RowsAffected == 0 {
		return nil, errors.New("songs not found")
	}
	return &ps.ReadSongsResponse{Songs: songs}, nil
}

func (s *PlaylistService) UpdateSong(_ context.Context, req *ps.UpdateSongRequest) (*ps.UpdateSongResponse, error) {
	var song ps.SongInfo
	reqSong := req.GetSong()
	res := s.DB.Model(&song).Where("id = ?", reqSong.Id).Updates(ps.SongInfo{
		Title:    reqSong.Title,
		Duration: reqSong.Duration,
	})
	if res.RowsAffected == 0 {
		return nil, errors.New("song not found")
	}
	return &ps.UpdateSongResponse{Song: &song}, nil
}

func (s *PlaylistService) DeleteSong(_ context.Context, req *ps.DeleteSongRequest) (*ps.DeleteSongResponse, error) {
	id := req.GetId()
	if s.P.IsPlaying && s.P.Cur.Info.Id == id {
		return &ps.DeleteSongResponse{Success: false}, errors.New("delete error: song is currently playing")
	}
	s.P.DeleteSong(id)
	var song ps.SongInfo
	res := s.DB.Where("id = ?", id).Delete(&song)
	if res.RowsAffected == 0 {
		return nil, errors.New("song not found")
	}
	return &ps.DeleteSongResponse{Success: true}, nil
}

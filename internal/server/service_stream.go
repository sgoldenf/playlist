package server

import (
	ps "github.com/sgoldenf/playlist/api"
	"log"
	"time"
)

func (s *PlaylistService) Player(_ *ps.ConnectRequest, stream ps.PlaylistService_PlayerServer) error {
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-timer.C:
			if s.P.IsPlaying {
				info := s.getPlayerInfo()
				err := stream.Send(info)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}
}

func (s *PlaylistService) getPlayerInfo() *ps.PlayerInfo {
	return &ps.PlayerInfo{
		Title:    s.P.Cur.Info.Title,
		Duration: s.P.Cur.Info.Duration,
		Elapsed:  s.P.Cur.ElapsedTime,
	}
}

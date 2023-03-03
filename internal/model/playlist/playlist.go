package playlist

import (
	ps "github.com/sgoldenf/playlist/api"
	"sync"
	"time"
)

type song struct {
	Info        ps.SongInfo
	ElapsedTime uint64
	prev        *song
	next        *song
}

type Playlist struct {
	len       int
	pause     chan struct{}
	m         sync.Mutex
	IsPlaying bool
	Cur       *song
	head      *song
	tail      *song
}

func NewPlaylist(songs []*ps.SongInfo) *Playlist {
	p := new(Playlist)
	p.pause = make(chan struct{})
	for _, s := range songs {
		p.AddSong(s)
	}
	return p
}

func (p *Playlist) AddSong(info *ps.SongInfo) {
	s := new(song)
	s.Info.Id = info.Id
	s.Info.Title = info.Title
	s.Info.Duration = info.Duration
	s.ElapsedTime = 0
	p.m.Lock()
	if p.head == nil {
		p.head = s
		p.tail = s
	} else if p.head == p.tail {
		p.tail = s
		p.head.next = s
		p.tail.prev = p.head
	} else {
		p.tail.next = s
		s.prev = p.tail
		p.tail = s
	}
	p.len++
	p.m.Unlock()
}

func (p *Playlist) Play() {
	if p.len > 0 && !p.IsPlaying {
		p.IsPlaying = true
		if p.Cur == nil {
			p.Cur = p.head
		}
		go p.playRoutine()
	}
}

func (p *Playlist) playRoutine() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-p.pause:
			return
		case <-ticker.C:
			p.Cur.ElapsedTime++
			if p.Cur.ElapsedTime == p.Cur.Info.Duration {
				p.IsPlaying = false
				if p.Cur == p.tail {
					p.Cur.ElapsedTime = 0
				} else {
					go p.Next()
				}
				return
			}
		}
	}
}

func (p *Playlist) Pause() {
	if p.IsPlaying {
		p.IsPlaying = false
		p.pause <- struct{}{}
	}
}

func (p *Playlist) Next() {
	p.m.Lock()
	if p.Cur != p.tail {
		p.Pause()
		if p.Cur != nil && p.Cur.next != nil {
			p.Cur.ElapsedTime = 0
			p.Cur = p.Cur.next
			p.Play()
		} else if p.Cur == nil && p.len > 2 {
			p.Cur = p.head.next
			p.Play()
		}
	}
	p.m.Unlock()
}

func (p *Playlist) Prev() {
	p.m.Lock()
	if p.Cur != nil && p.Cur != p.head {
		p.Pause()
		if p.Cur != nil && p.Cur.prev != nil {
			p.Cur.ElapsedTime = 0
			p.Cur = p.Cur.prev
			p.Play()
		}
	}
	p.m.Unlock()
}

func (p *Playlist) DeleteSong(id string) {
	p.m.Lock()
	if p.IsPlaying && id == p.Cur.Info.Id {
		p.m.Unlock()
		return
	}
	s := p.head
	for s != nil && s.Info.Id != id {
		s = s.next
	}
	if s != nil {
		if s == p.head {
			p.head = s.next
		}
		if s == p.tail {
			p.tail = s.prev
		}
		if s.prev != nil {
			s.prev.next = s.next
		}
		if s.next != nil {
			s.next.prev = s.prev
		}
		p.len--
	}
	p.m.Unlock()
}

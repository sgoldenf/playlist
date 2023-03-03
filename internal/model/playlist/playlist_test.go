package playlist

import (
	ps "github.com/sgoldenf/playlist/api"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

var song1 = &ps.SongInfo{
	Id:       "uuid1",
	Title:    "artist1 - song1",
	Duration: 4,
}

var song2 = &ps.SongInfo{
	Id:       "uuid2",
	Title:    "artist2 - song2",
	Duration: 3,
}

var song3 = &ps.SongInfo{
	Id:       "uuid3",
	Title:    "artist 3 - song 3",
	Duration: 2,
}

var songs = []*ps.SongInfo{song1, song2, song3}

func (p *Playlist) timeoutForSongChanging(duration time.Duration) bool {
	title := p.Cur.Info.Title
	timer := time.NewTimer(duration)
	for p.Cur != nil && title == p.Cur.Info.Title {
		select {
		case <-timer.C:
			return true
		default:
		}
	}
	return false
}

func TestPlaylist_NewPlaylistEmpty(t *testing.T) {
	p := NewPlaylist([]*ps.SongInfo{})
	if p.len != 0 {
		t.Errorf("len is %d, must be %d", p.len, 0)
	}
	if p.IsPlaying {
		t.Errorf("IsPlaying must be false")
	}
	if p.Cur != nil {
		t.Errorf("Cur must be nil")
	}
	if p.head != nil {
		t.Errorf("head must be nil")
	}
	if p.tail != nil {
		t.Errorf("tail must be nil")
	}
}

func TestPlaylist_AddSong1(t *testing.T) {
	p := NewPlaylist(songs[:1])
	if p.len != 1 {
		t.Errorf("len is %d, expected %d", p.len, 1)
	}
	if p.IsPlaying {
		t.Errorf("expected IsPlaying == false")
	}
	if p.Cur != nil {
		t.Errorf("Cur must be nil")
	}
	if p.head != p.tail {
		t.Errorf("expected p.head == p.tail")
	}
	if p.head.prev != nil {
		t.Errorf("p.head.prev must be nil")
	}
	if p.head.next != nil {
		t.Errorf("p.head.next must be nil")
	}
	if !reflect.DeepEqual(&p.head.Info, song1) {
		t.Errorf("WrongSongInfo\n%v\nexpected\n%v", &p.head.Info, song1)
	}
	if p.head.ElapsedTime != 0 {
		t.Errorf("p.head.ElapsedTime == %d, expected 0", p.Cur.ElapsedTime)
	}
}

func TestPlaylist_AddSong2(t *testing.T) {
	p := NewPlaylist(songs[:2])
	if p.len != 2 {
		t.Errorf("len is %d, expected %d", p.len, 2)
	}
	if p.head == nil {
		t.Errorf("head must not be nil")
	}
	if p.tail == nil {
		t.Errorf("tail must not be nil")
	}
	if p.head.prev != nil {
		t.Errorf("p.head.prev must be nil")
	}
	if p.head.next != p.tail {
		t.Errorf("expected p.head.next == p.tail")
	}
	if p.tail.prev != p.head {
		t.Errorf("expected p.tail.prev == p.head")
	}
	if p.tail.next != nil {
		t.Errorf("p.tail.next must be nil")
	}
	if !reflect.DeepEqual(&p.head.Info, song1) || !reflect.DeepEqual(&p.tail.Info, song2) {
		t.Errorf("WrongSongInfo\nSong1:\n%v\nexpected\n%v\nSong2:\n%v\nexpected\n%v",
			&p.head.Info, song1, &p.tail.Info, song2)
	}
}

func TestPlaylist_AddSong3(t *testing.T) {
	p := NewPlaylist(songs)
	if p.len != 3 {
		t.Errorf("len is %d, expected %d", p.len, 3)
	}
	if p.head == nil {
		t.Errorf("head must not be nil")
	}
	if p.head.prev != nil {
		t.Errorf("p.head.prev must be nil")
	}
	if p.tail == nil {
		t.Errorf("tail must not be nil")
	}
	s2 := p.head.next
	if s2 == nil {
		t.Errorf("p.head.next must not be nil")
	}
	if s2.prev != p.head {
		t.Errorf("s2.prev != p.head, expected s2.prev == p.head")
	}
	if s2.next != p.tail {
		t.Errorf("s2.next != p.tail, expected s2.next == p.tail")
	}
	if p.tail.prev != s2 {
		t.Errorf("expected p.tail.prev == s2")
	}
	if p.tail.next != nil {
		t.Errorf("p.tail.next must be nil")
	}
	if !reflect.DeepEqual(&p.head.Info, song1) ||
		!reflect.DeepEqual(&s2.Info, song2) ||
		!reflect.DeepEqual(&p.tail.Info, song3) {
		t.Errorf("WrongSongInfo\nSong1:\n%v\nexpected\n%v\nSong2:\n%v\nexpected\n%v\nSong3:\n%v\nexpected:\n%v",
			&p.head.Info, song1, &s2.Info, song2, &p.tail.Info, song3)
	}
}

func TestPlaylist_AddSongConcurrent(t *testing.T) {
	p := NewPlaylist([]*ps.SongInfo{})
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s := &ps.SongInfo{Id: string(rune(i))}
			p.AddSong(s)
		}(i)
	}
	wg.Wait()
	if p.len != 1000 {
		t.Errorf("len is %d, expected %d", p.len, 1000)
	}
	addedMap := make(map[string]struct{}, 1000)
	cur := p.head
	for cur != nil {
		addedMap[cur.Info.Id] = struct{}{}
		cur = cur.next
	}
	for i := 0; i < 1000; i++ {
		if _, ok := addedMap[string(rune(i))]; !ok {
			t.Errorf("%d - song haven't been added", i)
		}
	}
}

func TestPlaylist_Play(t *testing.T) {
	p := NewPlaylist(songs[:2])
	p.Play()
	if !p.IsPlaying {
		t.Errorf("expected p.isPlayng == true")
	}
	if p.Cur == nil {
		t.Errorf("expected p.Cur != nil")
	}
	if p.timeoutForSongChanging(5 * time.Second) {
		t.Errorf("timeout while changing song - %s", p.Cur.Info.Id)
	}
	if p.Cur == nil {
		t.Errorf("expected p.Cur != nil")
	} else {
		if !reflect.DeepEqual(&p.Cur.Info, song2) {
			t.Errorf("WrongSongInfo\n%v\nexpected\n%v", &p.Cur.Info, song2)
		}
		if !p.timeoutForSongChanging(3 * time.Second) {
			t.Errorf("timeout while changing song - %s", p.Cur.Info.Id)
		}
		if p.IsPlaying {
			t.Errorf("expected p.isPlayng == false")
		}
		if p.Cur != p.tail {
			t.Errorf("expected p.Cur == p.tail")
		}
	}
}

func TestPlaylist_Pause(t *testing.T) {
	p := NewPlaylist(songs)
	p.Play()
	time.Sleep(1 * time.Second)
	id1 := p.Cur.Info.Id
	elapsed1 := p.Cur.ElapsedTime
	p.Pause()
	id2 := p.Cur.Info.Id
	elapsed2 := p.Cur.ElapsedTime
	if p.IsPlaying {
		t.Errorf("expected IsPlaying == false")
	}
	p.Play()
	id3 := p.Cur.Info.Id
	elapsed3 := p.Cur.ElapsedTime
	if !p.IsPlaying {
		t.Errorf("expected IsPlaying == true")
	}
	if elapsed1 != elapsed2 || elapsed2 != elapsed3 ||
		id1 != id2 || id2 != id3 {
		t.Errorf("pause error, expected %d == %d == %d && %s == %s == %s",
			elapsed1, elapsed2, elapsed3, id1, id2, id3)
	}
}

func TestPlaylist_Next(t *testing.T) {
	p := NewPlaylist(songs)
	p.Next()
	if !p.IsPlaying {
		t.Errorf("expected IsPlaying == true")
	}
	if !reflect.DeepEqual(&p.Cur.Info, song2) {
		t.Errorf("WrongSongInfo\nSong2:\n%v\nexpected\n%v",
			&p.Cur.Info, song2)
	}
	p.Next()
	if !p.IsPlaying {
		t.Errorf("expected IsPlaying == true")
	}
	if !reflect.DeepEqual(&p.Cur.Info, song3) {
		t.Errorf("WrongSongInfo\nSong3:\n%v\nexpected\n%v",
			&p.Cur.Info, song3)
	}
	p.Pause()
	p.Next()
	if !reflect.DeepEqual(&p.Cur.Info, song3) {
		t.Errorf("WrongSongInfo\nSong3:\n%v\nexpected\n%v",
			&p.Cur.Info, song3)
	}
}

func TestPlaylist_NextConcurrent(t *testing.T) {
	p := NewPlaylist(songs)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.Next()
		}()
	}
	wg.Wait()
	if !p.IsPlaying {
		t.Errorf("expected IsPlaying == true")
	}
	if !reflect.DeepEqual(&p.Cur.Info, song3) {
		t.Errorf("WrongSongInfo\nSong3:\n%v\nexpected\n%v",
			&p.Cur.Info, song3)
	}
}

func TestPlaylist_Prev(t *testing.T) {
	p := NewPlaylist(songs)
	p.Prev()
	if p.IsPlaying {
		t.Errorf("expected IsPlaying == false")
	}
	if p.Cur != nil {
		t.Errorf("expected p.Cur == nil")
	}
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.Next()
		}()
	}
	wg.Wait()
	p.Pause()
	p.Prev()
	if !p.IsPlaying {
		t.Errorf("expected IsPlaying == true")
	}
	if !reflect.DeepEqual(&p.Cur.Info, song2) {
		t.Errorf("WrongSongInfo\nSong2:\n%v\nexpected\n%v",
			&p.Cur.Info, song2)
	}
}

func TestPlaylist_PrevNextConcurrent(t *testing.T) {
	p := NewPlaylist(append(songs, songs...))
	var wg sync.WaitGroup
	pos := 0
	for i := 0; i < 10; i++ {
		wg.Add(1)
		rand.Seed(time.Now().UnixNano())
		next := rand.Intn(2) == 0
		go func(next bool) {
			defer wg.Done()
			if next {
				p.Next()
			} else {
				p.Prev()
			}
			if next && pos < 5 {
				pos++
			} else if !next && pos > 0 {
				pos--
			}
		}(next)
	}
	wg.Wait()
	if !reflect.DeepEqual(&p.Cur.Info, songs[pos%3]) {
		t.Errorf("WrongSongInfo\n%v\nexpected\n%v",
			&p.Cur.Info, songs[pos%3])
	}
}

func TestPlaylist_DeleteSong(t *testing.T) {
	p := NewPlaylist(songs)
	p.Play()
	p.DeleteSong("uuid1")
	if !reflect.DeepEqual(&p.head.Info, song1) {
		t.Errorf("deleted playing song")
	}
	p.Pause()
	p.DeleteSong("uuid1")
	if p.len != 2 {
		t.Errorf("len is %d, expected %d", p.len, 2)
	}
	if !reflect.DeepEqual(&p.head.Info, song2) {
		t.Errorf("delete song error")
	}
	p.DeleteSong("uuid3")
	if p.len != 1 {
		t.Errorf("len is %d, expected %d", p.len, 1)
	}
	if p.head != p.tail {
		t.Errorf("expected p.head == p.tail")
	}
	if !reflect.DeepEqual(&p.head.Info, song2) {
		t.Errorf("delete song error")
	}
	p.DeleteSong("uuid2")
	if p.len != 0 {
		t.Errorf("len is %d, expected %d", p.len, 0)
	}
	if p.head != nil || p.tail != nil {
		t.Errorf("expected p.head == nuil && p.tail == nil")
	}
}

func TestPlaylist_DeleteSong2(t *testing.T) {
	p := NewPlaylist(songs)
	p.DeleteSong("uuid2")
	if p.len != 2 {
		t.Errorf("len is %d, expected %d", p.len, 2)
	}
	if p.head.next != p.tail {
		t.Errorf("expected p.head.next == p.tail")
	}
	if !reflect.DeepEqual(&p.tail.Info, song3) {
		t.Errorf("delete song error")
	}
}

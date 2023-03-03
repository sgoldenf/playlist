package server

import (
	"context"
	"errors"
	"fmt"
	ps "github.com/sgoldenf/playlist/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"io"
	"log"
	"net"
	"strconv"
	"testing"
	"time"
)

func runTestServerClientConnection(ctx context.Context) (ps.PlaylistServiceClient, func()) {
	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)

	s := grpc.NewServer()
	service, err := NewService()
	if err != nil {
		log.Fatalf("Failed to serve Database: %v", err)
	}
	ps.RegisterPlaylistServiceServer(s, service)
	go func() {
		if err = s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	conn, errConn := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if errConn != nil {
		log.Fatalf("error connecting to server: %v", errConn)
	}

	closeListener := func() {
		err = lis.Close()
		if err != nil {
			log.Fatalf("error closing listener: %v", err)
		}
		s.Stop()
	}

	client := ps.NewPlaylistServiceClient(conn)

	return client, closeListener
}

func printServerMessage(message *ps.PlayerInfo) {
	if message != nil {
		m := message.Elapsed / 60
		s := message.Elapsed - m*60
		em := strconv.FormatUint(m, 10)
		es := strconv.FormatUint(s, 10)
		if m < 10 {
			em = "0" + em
		}
		if s < 10 {
			es = "0" + es
		}
		m = message.Duration / 60
		s = message.Duration - m*60
		dm := strconv.FormatUint(m, 10)
		ds := strconv.FormatUint(s, 10)
		if m < 10 {
			dm = "0" + dm
		}
		if s < 10 {
			ds = "0" + ds
		}
		fmt.Printf("%s:%s/%s:%s - %s\n", em, es, dm, ds, message.Title)
	}
}

func TestPlaylistService_CreateSong(t *testing.T) {
	ctx := context.Background()
	client, closeListener := runTestServerClientConnection(ctx)
	defer closeListener()

	type expectation struct {
		out *ps.CreateSongResponse
		err error
	}

	song := &ps.SongInfo{
		Title:    "Dave Brubeck Quartet - Take Five",
		Duration: 325,
	}

	tests := map[string]struct {
		in       *ps.CreateSongRequest
		expected expectation
	}{
		"success": {
			in: &ps.CreateSongRequest{Song: song},
			expected: expectation{
				out: &ps.CreateSongResponse{Song: song},
				err: nil,
			}},
		"empty": {
			in: &ps.CreateSongRequest{Song: &ps.SongInfo{}},
			expected: expectation{nil,
				errors.New(
					"rpc error: code = Unknown desc = create song error: empty title/duration==0")},
		},
	}
	for caseName, test := range tests {
		t.Run(caseName, func(t *testing.T) {
			fmt.Printf("Creating song: %v\nExpected: %v err: \"%v\"\n", test.in, test.expected.out, test.expected.err)
			res, err := client.CreateSong(ctx, test.in)
			if err != nil {
				if test.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %v\nGot: %v\n", test.expected.err, err)
				}
			} else {
				if test.expected.out.Song.Duration != res.Song.Duration ||
					test.expected.out.Song.Title != res.Song.Title {
					t.Errorf("Out -> \nWant: %v\nGot : %v", test.expected.out, res)
				}
				if test.expected.out.Song.Id == res.Song.Id {
					t.Errorf(`expected res.Song.Id != ""`)
				}
			}
		})
	}
}

func TestPlaylistService_GetSong(t *testing.T) {
	ctx := context.Background()
	client, closeListener := runTestServerClientConnection(ctx)
	defer closeListener()

	res, err := client.GetSongs(ctx, &ps.ReadSongsRequest{})
	if err != nil {
		t.Errorf("GetSongsError:\nexpected err == nil, got:\n%v", err)
	}

	testSong := res.Songs[0]

	type expectation struct {
		out *ps.ReadSongResponse
		err error
	}

	tests := map[string]struct {
		in       *ps.ReadSongRequest
		expected expectation
	}{
		"success": {
			in: &ps.ReadSongRequest{Id: testSong.GetId()},
			expected: expectation{
				out: &ps.ReadSongResponse{Song: testSong},
				err: nil,
			}},
		"fail": {
			in: &ps.ReadSongRequest{Id: "uuid"},
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = Unknown desc = song not found"),
			},
		},
	}
	for caseName, test := range tests {
		fmt.Printf("Get song: %v\nExpected: %v, err: \"%v\"\n", test.in, test.expected.out, test.expected.err)
		t.Run(caseName, func(t *testing.T) {
			response, errGet := client.GetSong(ctx, test.in)
			if errGet != nil {
				if test.expected.err.Error() != errGet.Error() {
					t.Errorf("Err -> \nWant: %v\nGot: %v\n", test.expected.err, errGet)
				}
			} else {
				if test.expected.out.Song.Duration != response.Song.Duration ||
					test.expected.out.Song.Title != response.Song.Title ||
					test.expected.out.Song.Id != response.Song.Id {
					t.Errorf("Out -> \nWant: %v\nGot : %v", test.expected.out, response)
				}
			}
		})
	}
}

func TestPlaylistService_UpdateSong(t *testing.T) {
	ctx := context.Background()
	client, closeListener := runTestServerClientConnection(ctx)
	defer closeListener()

	res, err := client.GetSongs(ctx, &ps.ReadSongsRequest{})
	if err != nil {
		t.Errorf("GetSongsError:\nexpected err == nil, got:\n%v", err)
	}

	reqSong := res.Songs[len(res.Songs)-1]

	testSong := &ps.SongInfo{
		Title:    "Dave Brubeck Quartet - Blue Rondo A La Turk",
		Duration: 405,
	}

	type expectation struct {
		out *ps.UpdateSongResponse
		err error
	}

	tests := map[string]struct {
		in       *ps.UpdateSongRequest
		expected expectation
	}{
		"success": {
			in: &ps.UpdateSongRequest{Song: &ps.SongInfo{
				Id:       reqSong.Id,
				Title:    testSong.Title,
				Duration: testSong.Duration,
			}},
			expected: expectation{
				out: &ps.UpdateSongResponse{Song: testSong},
				err: nil,
			}},
		"fail": {
			in: &ps.UpdateSongRequest{Song: &ps.SongInfo{Id: "invalid"}},
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = Unknown desc = song not found"),
			},
		},
	}
	for caseName, test := range tests {
		t.Run(caseName, func(t *testing.T) {
			fmt.Printf("Updating song: %v\nExpected: %v, err: %v\n", test.in, test.expected.out, test.expected.err)
			response, errUpdate := client.UpdateSong(ctx, test.in)
			if errUpdate != nil {
				if test.expected.err.Error() != errUpdate.Error() {
					t.Errorf("Err -> \nWant: %v\nGot: %v\n", test.expected.err, errUpdate)
				}
			} else {
				if test.expected.out.Song.Duration != response.Song.Duration ||
					test.expected.out.Song.Title != response.Song.Title ||
					test.expected.out.Song.Id != response.Song.Id {
					t.Errorf("Out -> \nWant: %v\nGot : %v", test.expected.out, response)
				}
			}
		})
	}
}

func TestPlaylistService_DeleteSong(t *testing.T) {
	ctx := context.Background()
	client, closeListener := runTestServerClientConnection(ctx)
	defer closeListener()

	res, err := client.GetSongs(ctx, &ps.ReadSongsRequest{})
	if err != nil {
		t.Errorf("GetSongsError:\nexpected err == nil, got:\n%v", err)
	}

	reqSong := res.Songs[len(res.Songs)-1]
	failSong := res.Songs[0]

	type expectation struct {
		out *ps.DeleteSongResponse
		err error
	}

	tests := map[string]struct {
		in       *ps.DeleteSongRequest
		expected expectation
	}{
		"success": {
			in: &ps.DeleteSongRequest{Id: reqSong.Id},
			expected: expectation{
				out: &ps.DeleteSongResponse{Success: true},
				err: nil,
			}},
		"fail": {
			in: &ps.DeleteSongRequest{Id: failSong.Id},
			expected: expectation{
				out: &ps.DeleteSongResponse{Success: false},
				err: errors.New("rpc error: code = Unknown desc = delete error: song is currently playing"),
			},
		},
		"error": {
			in: &ps.DeleteSongRequest{Id: "invalid"},
			expected: expectation{
				out: nil,
				err: errors.New("rpc error: code = Unknown desc = song not found"),
			},
		},
	}
	_, err = client.Play(ctx, &ps.PlayRequest{})
	fmt.Printf("Now playing: %s\n", failSong)
	if err != nil {
		t.Errorf("play error: %v", err)
	}
	for caseName, test := range tests {
		t.Run(caseName, func(t *testing.T) {
			fmt.Printf("Delete song: %v\nExpected: %v, err: \"%v\"\n", test.in.Id, test.expected.out, test.expected.err)
			response, errDelete := client.DeleteSong(ctx, test.in)
			if errDelete != nil {
				if test.expected.err.Error() != errDelete.Error() {
					t.Errorf("Err -> \nWant: %v\nGot: %v\n", test.expected.err, errDelete)
				}
			} else {
				if test.expected.out.Success != response.Success {
					t.Errorf("Out -> \nWant: %v\nGot : %v", test.expected.out, response)
				}
			}
		})
	}
	_, err = client.Pause(ctx, &ps.PauseRequest{})
	if err != nil {
		t.Errorf("pause error: %v", err)
	}
}

func TestPlaylistService_Player(t *testing.T) {
	ctx := context.Background()
	client, closeListener := runTestServerClientConnection(ctx)
	defer closeListener()

	getSongsResponse, err := client.GetSongs(ctx, &ps.ReadSongsRequest{})
	if err != nil {
		t.Errorf("GetSongsError:\nexpected err == nil, got:\n%v", err)
	}

	if len(getSongsResponse.Songs) < 2 {
		t.Errorf("need at least 2 songs in db to run test")
	} else {
		song1 := getSongsResponse.Songs[0]
		song2 := getSongsResponse.Songs[1]

		testResponses := []*ps.PlayerInfo{
			{
				Title:    song1.Title,
				Duration: song1.Duration,
			},
			{
				Title:    song1.Title,
				Duration: song1.Duration,
			},
			{
				Title:    song2.Title,
				Duration: song2.Duration,
			},
			{
				Title:    song2.Title,
				Duration: song2.Duration,
			},
			{
				Title:    song1.Title,
				Duration: song1.Duration,
			},
		}

		t.Run("player_functions", func(t *testing.T) {
			stream, errPlayer := client.Player(ctx, &ps.ConnectRequest{})
			var messages []*ps.PlayerInfo

			go func() {
				for {
					message, errStream := stream.Recv()
					printServerMessage(message)
					if errors.Is(errStream, io.EOF) {
						break
					}
					messages = append(messages, message)
				}
			}()

			_, errPlayer = client.Play(ctx, &ps.PlayRequest{})
			if errPlayer != nil {
				t.Errorf("play error: %v", errPlayer)
			}
			fmt.Printf("Playing...\n")
			time.Sleep(2*time.Second + 1*time.Millisecond)
			_, errPlayer = client.Next(ctx, &ps.NextSongRequest{})
			if errPlayer != nil {
				t.Errorf("next error: %v", errPlayer)
			}
			fmt.Printf("Next...\n")
			time.Sleep(2*time.Second + 1*time.Millisecond)
			_, errPlayer = client.Prev(ctx, &ps.PrevSongRequest{})
			if errPlayer != nil {
				t.Errorf("prev error: %v", errPlayer)
			}
			fmt.Printf("Prev...\n")
			time.Sleep(1*time.Second + 1*time.Millisecond)
			_, errPlayer = client.Pause(ctx, &ps.PauseRequest{})
			if errPlayer != nil {
				t.Errorf("pause error: %v", errPlayer)
			}
			fmt.Printf("Paused.\n")
			ctx.Done()

			if len(messages) != len(testResponses) {
				t.Errorf("Out -> \nWant: %v\nGot : %v", testResponses, messages)
			} else {
				for i, message := range messages {
					if message.Title != testResponses[i].Title ||
						message.Duration != testResponses[i].Duration {
						t.Errorf("Out -> \nWant: %v\nGot : %v", testResponses[i], message)
					}
					if i > 0 && message.Elapsed == messages[i-1].Elapsed && message.Title == messages[i-1].Title {
						t.Errorf("Dublicate Message %d: %v", i, message)
					}
				}
			}
		})
	}
}

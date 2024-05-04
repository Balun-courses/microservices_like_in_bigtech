package main

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"

	pb "github.com/Balun-courses/microservices_like_in_bigtech/grpc_stream/pkg/api/notes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var idSerial uint64

// server is used to implement pb.NotesServiceServer.
type server struct {
	// UnimplementedNotesServiceServer must be embedded to have forward compatible implementations.
	pb.UnimplementedNotesServiceServer

	mx    sync.RWMutex
	notes map[uint64]*pb.NoteInfo
}

func NewServer() *server {
	return &server{notes: make(map[uint64]*pb.NoteInfo)}
}

func (s *server) SaveNotesStream(stream pb.NotesService_SaveNotesStreamServer) error {
	log.Println("SaveNotesStream: received")

	var (
		rpcError error
		wg       sync.WaitGroup
	)
	for {
		note, err := stream.Recv()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				rpcError = err
			}
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			id := atomic.AddUint64(&idSerial, 1)

			s.mx.Lock()
			s.notes[id] = note
			s.mx.Unlock()

			log.Printf("save note %d", id)

			if err := stream.Send(&pb.Note{
				Id:   id,
				Info: note,
			}); err != nil {
				log.Println(err)
			}
		}()
	}
	wg.Wait()

	return rpcError
}

func (s *server) ListNotesStream(req *pb.ListNotesStreamRequest, stream pb.NotesService_ListNotesStreamServer) error {
	log.Println("ListNotesStream: received")

	ch := make(chan *pb.Note, 100)
	go func() {
		s.mx.RLock()
		for id, note := range s.notes {
			ch <- &pb.Note{
				Id:   id,
				Info: note,
			}
		}
		s.mx.RUnlock()
		close(ch)
	}()

	for note := range ch {
		log.Printf("ListNotesStream: send note %d", note.Id)
		if err := stream.Send(note); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterNotesServiceServer(s, NewServer())

	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

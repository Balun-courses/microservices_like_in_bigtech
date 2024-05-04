package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	pb "github.com/Balun-courses/microservices_like_in_bigtech/grpc_openapi/pkg/api/notes"
	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var idSerial uint64

// server is used to implement pb.NotesServiceServer.
type server struct {
	// UnimplementedNotesServiceServer must be embedded to have forward compatible implementations.
	pb.UnimplementedNotesServiceServer

	mx        sync.RWMutex
	notes     map[uint64]*pb.NoteInfo
	validator *protovalidate.Validator
}

func NewServer() (*server, error) {
	srv := &server{
		notes: make(map[uint64]*pb.NoteInfo),
	}

	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(true),
		protovalidate.WithMessages(
			// Добавляем сюда все запросы наши
			&pb.SaveNoteRequest{},
			&pb.ListNotesRequest{},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	srv.validator = validator
	return srv, nil
}

func protovalidateVialationsToGoogleViolations(vs []*validate.Violation) []*errdetails.BadRequest_FieldViolation {
	res := make([]*errdetails.BadRequest_FieldViolation, len(vs))
	for i, v := range vs {
		res[i] = &errdetails.BadRequest_FieldViolation{
			Field:       v.FieldPath,
			Description: v.Message,
		}
	}
	return res
}

func convertProtovalidateValidationErrorToErrdetailsBadRequest(valErr *protovalidate.ValidationError) *errdetails.BadRequest {
	return &errdetails.BadRequest{
		FieldViolations: protovalidateVialationsToGoogleViolations(valErr.Violations),
	}
}

func rpcValidationError(err error) error {
	if err == nil {
		return nil
	}

	var valErr *protovalidate.ValidationError
	if ok := errors.As(err, &valErr); ok {
		st, err := status.New(codes.InvalidArgument, codes.InvalidArgument.String()).
			WithDetails(convertProtovalidateValidationErrorToErrdetailsBadRequest(valErr))
		if err == nil {
			return st.Err()
		}
	}

	return status.Error(codes.Internal, err.Error())
}

// SaveNote implements pb.NotesServiceServer
func (s *server) SaveNote(_ context.Context, req *pb.SaveNoteRequest) (*pb.SaveNoteResponse, error) {
	info := req.GetInfo()
	log.Printf("SaveNote: received: %s", info.GetTitle())

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	id := atomic.AddUint64(&idSerial, 1)

	s.mx.Lock()
	s.notes[id] = info
	s.mx.Unlock()

	return &pb.SaveNoteResponse{
		Id: id,
	}, nil
}

// ListNotes implements pb.NotesServiceServer
func (s *server) ListNotes(_ context.Context, req *pb.ListNotesRequest) (*pb.ListNotesResponse, error) {
	log.Println("ListNotes: received")

	s.mx.RLock()
	defer s.mx.RUnlock()

	notes := make([]*pb.Note, 0, len(s.notes))
	for id, note := range s.notes {
		notes = append(notes, &pb.Note{
			Id:   id,
			Info: note,
		})
	}

	return &pb.ListNotesResponse{
		Notes: notes,
	}, nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server, err := NewServer()
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		grpcServer := grpc.NewServer()
		pb.RegisterNotesServiceServer(grpcServer, server)

		reflection.Register(grpcServer)

		lis, err := net.Listen("tcp", ":8082")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Printf("server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

		// SaveNote:
		// grpc_cli call --json_input --json_output localhost:8082 NotesService/SaveNote '{"info":{"title":"my note","content":"my note content"}}'
		// ListNotes:
		// grpc_cli call --json_input --json_output localhost:8082 NotesService/ListNotes '{}'
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Register gRPC server endpoint
		// Note: Make sure the gRPC server is running properly and accessible
		mux := runtime.NewServeMux()
		if err = pb.RegisterNotesServiceHandlerServer(ctx, mux, server); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		httpServer := &http.Server{Handler: mux}

		lis, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		// Start HTTP server (and proxy calls to gRPC server endpoint)
		log.Printf("server listening at %v", lis.Addr())
		if err := httpServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	wg.Wait()

	// SaveNote:
	// curl --location 'localhost:8080/api/v1/notes' --header 'Content-Type: application/json' --data '{ "title": "1","content": "1"}'
	// ListNotes:
	// curl --location 'localhost:8080/api/v1/notes'
}

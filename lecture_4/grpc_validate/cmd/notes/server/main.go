package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	pb "github.com/Balun-courses/microservices_like_in_bigtech/grpc_validate/pkg/api/notes"
	"github.com/bufbuild/protovalidate-go"
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
	server, err := NewServer()
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotesServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// SaveNote:
	// grpc_cli call --json_input --json_output localhost:8082 NotesService/SaveNote '{"info":{"title":"my note","content":"my note content"}}'
	// ListNotes:
	// grpc_cli call --json_input --json_output localhost:8082 NotesService/ListNotes '{}'
}

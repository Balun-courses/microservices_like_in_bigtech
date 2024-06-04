// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: api/notes/service.proto

package notes

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	NotesService_SaveNote_FullMethodName  = "/github.com.moguchev.microservices.grpc_gateway.NotesService/SaveNote"
	NotesService_ListNotes_FullMethodName = "/github.com.moguchev.microservices.grpc_gateway.NotesService/ListNotes"
)

// NotesServiceClient is the client API for NotesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotesServiceClient interface {
	// SaveNote - save note
	SaveNote(ctx context.Context, in *SaveNoteRequest, opts ...grpc.CallOption) (*SaveNoteResponse, error)
	// ListNotes - list all notes
	ListNotes(ctx context.Context, in *ListNotesRequest, opts ...grpc.CallOption) (*ListNotesResponse, error)
}

type notesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNotesServiceClient(cc grpc.ClientConnInterface) NotesServiceClient {
	return &notesServiceClient{cc}
}

func (c *notesServiceClient) SaveNote(ctx context.Context, in *SaveNoteRequest, opts ...grpc.CallOption) (*SaveNoteResponse, error) {
	out := new(SaveNoteResponse)
	err := c.cc.Invoke(ctx, NotesService_SaveNote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notesServiceClient) ListNotes(ctx context.Context, in *ListNotesRequest, opts ...grpc.CallOption) (*ListNotesResponse, error) {
	out := new(ListNotesResponse)
	err := c.cc.Invoke(ctx, NotesService_ListNotes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotesServiceServer is the server API for NotesService service.
// All implementations must embed UnimplementedNotesServiceServer
// for forward compatibility
type NotesServiceServer interface {
	// SaveNote - save note
	SaveNote(context.Context, *SaveNoteRequest) (*SaveNoteResponse, error)
	// ListNotes - list all notes
	ListNotes(context.Context, *ListNotesRequest) (*ListNotesResponse, error)
	mustEmbedUnimplementedNotesServiceServer()
}

// UnimplementedNotesServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNotesServiceServer struct {
}

func (UnimplementedNotesServiceServer) SaveNote(context.Context, *SaveNoteRequest) (*SaveNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveNote not implemented")
}
func (UnimplementedNotesServiceServer) ListNotes(context.Context, *ListNotesRequest) (*ListNotesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNotes not implemented")
}
func (UnimplementedNotesServiceServer) mustEmbedUnimplementedNotesServiceServer() {}

// UnsafeNotesServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotesServiceServer will
// result in compilation errors.
type UnsafeNotesServiceServer interface {
	mustEmbedUnimplementedNotesServiceServer()
}

func RegisterNotesServiceServer(s grpc.ServiceRegistrar, srv NotesServiceServer) {
	s.RegisterService(&NotesService_ServiceDesc, srv)
}

func _NotesService_SaveNote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveNoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotesServiceServer).SaveNote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotesService_SaveNote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotesServiceServer).SaveNote(ctx, req.(*SaveNoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NotesService_ListNotes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListNotesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotesServiceServer).ListNotes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NotesService_ListNotes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotesServiceServer).ListNotes(ctx, req.(*ListNotesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NotesService_ServiceDesc is the grpc.ServiceDesc for NotesService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NotesService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "github.com.moguchev.microservices.grpc_gateway.NotesService",
	HandlerType: (*NotesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SaveNote",
			Handler:    _NotesService_SaveNote_Handler,
		},
		{
			MethodName: "ListNotes",
			Handler:    _NotesService_ListNotes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/notes/service.proto",
}

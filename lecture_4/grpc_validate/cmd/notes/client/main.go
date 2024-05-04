package main

import (
	"context"
	"log"

	pb "github.com/Balun-courses/microservices_like_in_bigtech/grpc_validate/pkg/api/notes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	conn, err := grpc.Dial(":8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	cli := pb.NewNotesServiceClient(conn)

	// /SaveNote
	{
		resp, err := cli.SaveNote(context.Background(), &pb.SaveNoteRequest{
			Info: &pb.NoteInfo{
				Title:   "client note",
				Content: "content",
			},
		})
		if err != nil {
			log.Fatalf("SaveNote error: %v", err)
		} else {
			log.Printf("note id is %d", resp.GetId())
		}
	}

	// /ListNotes
	{
		resp, err := cli.ListNotes(context.Background(), &pb.ListNotesRequest{})
		if err != nil {
			log.Fatalf("ListNotes error: %v", err)
		} else {
			// для Marshal proto сообщений в JSON необходимо использовать пакет protojson
			notes, err := protojson.Marshal(resp)
			if err != nil {
				log.Fatalf(" protojson.Marshal error: %v", err)
			} else {
				log.Printf("notes: %s", string(notes))
			}
		}
	}
}

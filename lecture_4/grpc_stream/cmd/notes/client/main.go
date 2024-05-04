package main

import (
	"context"
	"errors"
	"io"
	"log"
	"strconv"
	"sync"

	pb "github.com/Balun-courses/microservices_like_in_bigtech/grpc_stream/pkg/api/notes"
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
		stream, err := cli.SaveNotesStream(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		const n = 5
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer stream.CloseSend() // ОБЯЗАТЕЛЬНО!

			for i := 0; i < n; i++ {
				note := &pb.NoteInfo{
					Title:   "title " + strconv.Itoa(i),
					Content: "content",
				}

				log.Printf("send note: %s", note.Title)
				if err := stream.Send(note); err != nil {
					log.Println(err)
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				resp, err := stream.Recv()
				if err != nil {
					if !errors.Is(err, io.EOF) {
						log.Println(err)
					}
					break
				}
				// для Marshal proto сообщений в JSON необходимо использовать пакет protojson
				note, err := protojson.Marshal(resp)
				if err != nil {
					log.Fatalf(" protojson.Marshal error: %v", err)
				} else {
					log.Printf("recieve note: %s", string(note))
				}
			}
		}()

		wg.Wait()
	}

	// /ListNotes
	{
		stream, err := cli.ListNotesStream(context.Background(), &pb.ListNotesStreamRequest{})
		if err != nil {
			log.Fatal(err)
		}

		for {
			resp, err := stream.Recv()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					log.Println(err)
				}
				return
			}
			// для Marshal proto сообщений в JSON необходимо использовать пакет protojson
			note, err := protojson.Marshal(resp)
			if err != nil {
				log.Fatalf(" protojson.Marshal error: %v", err)
			} else {
				log.Printf("recieve note: %s", string(note))
			}
		}
	}

}

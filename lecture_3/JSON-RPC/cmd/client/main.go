package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/moguchev/microservices_courcse/jsonrpc/model"
)

type JSONRPC2Request struct {
	Method string `json:"method"`
	Params [1]any `json:"params"`
	ID     string `json:"id"`
}

type JSONRPC2Response[T any] struct {
	ID     string  `json:"id"`
	Error  *string `json:"error"`
	Result *T      `json:"result"`
}

func NewJSONRPC2Response[T any]() JSONRPC2Response[T] {
	return JSONRPC2Response[T]{}
}

func main() {
	rpcReq := JSONRPC2Request{
		Method: "Service.Multiply",
		Params: [1]any{
			model.MultiplyRequest{
				A: 1,
				B: 10,
			}},
		ID: "1",
	}

	request, err := json.Marshal(rpcReq)
	if err != nil {
		log.Fatalln(err)
	}

	response, err := http.Post("http://localhost:8080/rpc", "application/json", bytes.NewBuffer(request))
	if err != nil {
		log.Println(err)
	} else {
		// var reply map[string]any
		var reply = NewJSONRPC2Response[model.MultiplyResponse]()
		if err = json.NewDecoder(response.Body).Decode(&reply); err != nil {
			log.Println(err)
		} else {
			log.Print("rpc id: ", reply.ID)
			if reply.Error != nil {
				log.Println("rpc error: ", err)
			} else if reply.Result != nil {
				log.Printf("rpc response: %#v\n", reply.Result)
			}
		}
	}
}

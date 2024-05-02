package model

// MultiplyRequest represents the arguments passed in the JSON-RPC Multiply request.
type MultiplyRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

// MultiplyRequest represents the reply in the JSON-RPC Multiply method.
type MultiplyResponse struct {
	Value int `json:"value"`
}

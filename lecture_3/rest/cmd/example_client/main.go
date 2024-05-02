package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rest/server/models"
	"strings"
)

const host = "http://localhost:8080"

func createUser(httpClient *http.Client, request *models.CreateUserRequest) (*models.CreateUserResponse, error) {
	const uri = "/api/v1/users"

	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(request); err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(http.MethodPost, strings.Join([]string{host, uri}, ""), body)
	if err != nil {
		return nil, err
	}

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusCreated {
		var errorMsg models.ErrorMessage
		if err := json.NewDecoder(httpResponse.Body).Decode(&errorMsg); err != nil {
			return nil, fmt.Errorf("user not created: status code %d", httpResponse.StatusCode)
		}

		return nil, fmt.Errorf("user not created: status code %d, error: %s", httpResponse.StatusCode, errorMsg.Message)
	}

	response := new(models.CreateUserResponse)
	if err := json.NewDecoder(httpResponse.Body).Decode(response); err != nil {
		return nil, fmt.Errorf("can't unmarshal reponse: %s", err.Error())
	}

	return response, nil
}

func newString(s string) *string {
	return &s
}
func main() {
	httpClient := &http.Client{}

	resp, err := createUser(httpClient, &models.CreateUserRequest{
		Email:    newString("tet@mail.ru"),
		Name:     newString("ddde"),
		Password: newString("1234"),
	})
	if err != nil {
		log.Printf("createUser error: %v", err)
	} else {
		log.Printf("createUser: %#v", resp)
	}
}

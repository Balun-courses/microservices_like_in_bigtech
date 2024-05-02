package main

import (
	"encoding/json"
	"net/http"
	"rest/server/models"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
)

var empty struct{}

func httpErrorMsg(err error) *models.ErrorMessage {
	if err == nil {
		return nil
	}
	return &models.ErrorMessage{
		Message: err.Error(),
	}
}

func createUser(c echo.Context) error {
	var request models.CreateUserRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	if err := request.Validate(strfmt.Default); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	// ...

	response := models.CreateUserResponse{ID: 1}
	return c.JSON(http.StatusCreated, response)
}

func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	// ...

	response := models.GetUserResponse{
		ID:    int64(userID),
		Email: "test@mail.ru",
		Name:  "test",
	}
	return c.JSON(http.StatusOK, response)
}

func updateUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	var request models.UpdateUserRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&request); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	if err := request.Validate(strfmt.Default); err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	// ...

	response := models.UpdateUserResponse{
		ID:    int64(userID),
		Email: "test@mail.ru",
		Name:  "test",
	}
	return c.JSON(http.StatusOK, response)
}

func deleteUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpErrorMsg(err))
	}

	_ = userID

	return c.JSON(http.StatusOK, empty)
}

func main() {
	e := echo.New()

	e.POST("/api/v1/users", createUser)
	e.GET("/api/v1/users/:id", getUser)
	e.PUT("/api/v1/users/:id", updateUser)
	e.DELETE("/api/v1/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":8080"))
}

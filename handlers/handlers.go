package handlers

import (
	"github.com/JayJamieson/sekret/util"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Sekret struct {
	store map[string]SecretData
}

type SecretData struct {
	Data      string `json:"data"`
	CreatedAt int64  `json:"createdAt"`
	Owner     string `json:"owner"`
}

type createdResponse struct {
	Name string `json:"name"`
}

func NewSekret() *Sekret {
	return &Sekret{
		store: make(map[string]SecretData),
	}
}

func (s *Sekret) CreateSecret(c echo.Context) error {

	createdAt := time.Now().Add(time.Hour * time.Duration(24)).UnixNano()

	secret := &SecretData{
		CreatedAt: createdAt,
	}

	if err := c.Bind(secret); err != nil {
		return err
	}

	var key string = util.GetRandomName(0)

	if _, ok := s.store[key]; ok {
		key = util.GetRandomName(1)
	}

	s.store[key] = *secret

	return c.JSON(http.StatusAccepted, createdResponse{key})
}

func (s *Sekret) GetSecret(c echo.Context) error {
	key := c.Param("key")

	secret, ok := s.store[key]

	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	delete(s.store, key)

	if secret.CreatedAt < time.Now().UnixNano() {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, secret)
}

package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type SecretData struct {
	Secret    string `json:"secret" form:"secret"`
	CreatedAt int64
	Owner     string
}

type createdResponse struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type SecretStore struct {
	store map[string]SecretData
}

func New() *SecretStore {
	return &SecretStore{
		store: make(map[string]SecretData),
	}
}

func (s *SecretStore) Health(c echo.Context) error {
	return c.String(http.StatusOK, "Ok")
}

func (s *SecretStore) CreateSecret(c echo.Context) error {

	createdAt := time.Now().Add(time.Hour * time.Duration(24)).UnixNano()

	secret := &SecretData{
		CreatedAt: createdAt,
	}

	if err := c.Bind(secret); err != nil {
		return err
	}

	var key string = GetRandomName(0)

	if _, ok := s.store[key]; ok {
		key = GetRandomName(1)
	}

	s.store[key] = *secret

	// TODO: Handle in content type specific handlers
	contentType := c.Request().Header.Get(echo.HeaderContentType)

	if contentType == echo.MIMEApplicationJSON {
		return c.JSON(http.StatusOK, createdResponse{Name: key, Link: fmt.Sprintf("%s://%s/secret/%s", c.Scheme(), c.Request().Host, key)})
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/private/%s", key))
}

func (s *SecretStore) GetSecret(c echo.Context) error {
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

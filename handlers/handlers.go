package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type SecretStore struct {
	store map[string]SecretData
}

type SecretData struct {
	Secret    string `json:"secret"`
	CreatedAt int64  `json:"createdAt"`
	Owner     string `json:"owner"`
}

type createdResponse struct {
	Name string `json:"name"`
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

	return c.JSON(http.StatusAccepted, createdResponse{key})
}

func (s *SecretStore) ViewSecret(c echo.Context) error {
	key := c.Param("key")

	secret := s.store[key]
	delete(s.store, key)

	return c.Render(http.StatusOK, "view", map[string]interface{}{
		"secret": secret.Secret,
	})
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

	if ct, ok := c.Request().Header[http.CanonicalHeaderKey("Content-type")]; ok && ct[0] == "application/json" {
		return c.JSON(http.StatusOK, secret)
	}

	return c.Redirect(http.StatusFound, "/")
}

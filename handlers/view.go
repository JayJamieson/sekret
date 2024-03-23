package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (s *SecretStore) ViewIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func (s *SecretStore) ViewSecret(c echo.Context) error {
	key := c.Param("key")
	data, _ := c.FormParams()
	secret, ok := s.store[key]

	if !ok {
		return c.Render(http.StatusOK, "view", map[string]interface{}{
			"key":    key,
			"show":   data.Has("show"),
			"used":   !ok,
			"secret": "",
			"data":   data,
		})
	}

	if !data.Has("show") {
		return c.Render(http.StatusOK, "view", map[string]interface{}{
			"key":    key,
			"show":   data.Has("show"),
			"secret": "",
			"data":   data,
		})
	}

	delete(s.store, key)

	if secret.CreatedAt < time.Now().UnixNano() {
		return c.Render(http.StatusOK, "view", map[string]interface{}{
			"key":    key,
			"show":   data.Has("show"),
			"used":   true,
			"secret": "",
			"data":   data,
		})
	}

	return c.Render(http.StatusOK, "view", map[string]interface{}{
		"key":    key,
		"show":   data.Has("show"),
		"secret": secret.Secret,
		"data":   data,
	})
}

func (s *SecretStore) ViewPrivate(c echo.Context) error {
	key := c.Param("key")

	_, ok := s.store[key]

	if !ok {
		return c.Redirect(http.StatusFound, "/")
	}

	return c.Render(http.StatusOK, "private", map[string]interface{}{
		"link": fmt.Sprintf("%s://%s/secret/%s", c.Scheme(), c.Request().Host, key),
	})
}

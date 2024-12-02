package handlers

import (
	"database/sql"
	"errors"
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

	secret, err := s.get(key)

	errNoRows := err != nil && errors.Is(err, sql.ErrNoRows)

	if errNoRows {
		return c.Render(http.StatusOK, "view", map[string]interface{}{
			"key":    key,
			"show":   data.Has("show"),
			"used":   errNoRows,
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

	s.remove(key)

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
		"key":         key,
		"show":        data.Has("show"),
		"secret":      secret.Secret,
		"iv":          secret.IV,
		"salt":        secret.Salt,
		"hasPassword": secret.HasPassword,
		"data":        data,
	})
}

func (s *SecretStore) ViewPrivate(c echo.Context) error {
	key := c.Param("key")

	_, err := s.get(key)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Redirect(http.StatusFound, "/")
		}

		return err
	}

	return c.Render(http.StatusOK, "private", map[string]interface{}{
		"link": fmt.Sprintf("%s://%s/secret/%s", c.Scheme(), c.Request().Host, key),
	})
}

package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"database/sql"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
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
	db    *sql.DB
}

func New(url string) (*SecretStore, error) {
	db, err := sql.Open("libsql", url)

	if err != nil {
		return nil, err
	}

	// stops STREAM_EXPIRED issues with connections timing out but not released
	db.SetConnMaxIdleTime(9)

	return &SecretStore{
		store: make(map[string]SecretData),
		db:    db,
	}, nil
}

func (s *SecretStore) set(data *SecretData) (string, error) {
	tx, err := s.db.Begin()

	if err != nil {
		return "", err
	}

	var name string = GetRandomName(0)

	createSql := "INSERT INTO secret(name, value, created_at, owner) VALUES(?, ?, ?, ?)"
	_, err = tx.Exec(createSql, name, data.Secret, data.CreatedAt, data.Owner)

	if err != nil {
		txErr := tx.Rollback()
		return "", errors.Join(err, txErr)
	}

	tx.Commit()

	return name, nil
}

func (s *SecretStore) get(name string) (*SecretData, error) {
	fetchSql := "SELECT value, created_at, owner FROM secret WHERE name = ?"
	row := s.db.QueryRow(fetchSql, name)

	if err := row.Err(); err != nil {
		return nil, err
	}

	var secret SecretData
	err := row.Scan(&secret.Secret, &secret.CreatedAt, &secret.Owner)

	if err != nil {
		return nil, err
	}

	return &secret, nil
}

func (s *SecretStore) remove(name string) error {
	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	removeSql := "DELETE FROM secret WHERE name = ?"

	_, err = tx.Exec(removeSql, name)

	if err != nil {
		txErr := tx.Rollback()
		return errors.Join(err, txErr)
	}

	return tx.Commit()
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

	key, err := s.set(secret)

	if err != nil {
		return err
	}

	// TODO: Handle in content type specific handlers
	contentType := c.Request().Header.Get(echo.HeaderContentType)

	if contentType == echo.MIMEApplicationJSON {
		return c.JSON(http.StatusOK, createdResponse{Name: key, Link: fmt.Sprintf("%s://%s/secret/%s", c.Scheme(), c.Request().Host, key)})
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/private/%s", key))
}

func (s *SecretStore) GetSecret(c echo.Context) error {
	key := c.Param("key")

	secret, err := s.get(key)
	errNoRows := err != nil && errors.Is(err, sql.ErrNoRows)

	if errNoRows {
		return c.NoContent(http.StatusNotFound)
	}

	s.remove(key)

	if secret.CreatedAt < time.Now().UnixNano() {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, secret)
}

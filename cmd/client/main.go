package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/JayJamieson/sekret/handlers"
	"github.com/spf13/viper"
)

type Client struct {
	endpoint   *url.URL
	httpClient *http.Client
}

// fetch secret using nice name
func (client *Client) fetch(key string) error {

	u := client.endpoint.ResolveReference(&url.URL{Path: fmt.Sprintf("/secret/%v", key)})

	resp, err := client.httpClient.Get(u.String())

	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New("secret not found")
	}

	defer resp.Body.Close()

	data := handlers.SecretData{}

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return err
	}

	fmt.Printf("Secret: %v\n\n", data.Data)

	return nil
}

// create secret, output location and password
func (client *Client) create(password string) error {
	u := client.endpoint.ResolveReference(&url.URL{Path: "/secret"})

	reqData, err := json.Marshal(map[string]string{"data": password, "owner": "sekret cli"})

	if err != nil {
		return err
	}

	resp, err := client.httpClient.Post(
		u.String(),
		"application/json",
		bytes.NewBuffer(reqData))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var respData map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&respData)

	if err != nil {
		return err
	}

	fmt.Printf("%v/%v\n\n", u.String(), respData["name"])

	return nil
}

func main() {
	viper.SetEnvPrefix("sk")
	err := viper.BindEnv("endpoint")

	if err != nil {
		fmt.Printf("Error binding to endpoint %v\n", err.Error())
		os.Exit(1)
	}

	u, err := url.Parse(viper.GetString("endpoint"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &Client{
		endpoint: u,
		httpClient: &http.Client{
			Timeout: time.Duration(10) * time.Second,
		},
	}

	fs := flag.NewFlagSet("Sekret CLI", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)

	fs.Func("fetch", "`Secret name` to fetch", client.fetch)
	fs.Func("create", "Generate secret", client.create)

	err = fs.Parse(os.Args[1:])

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

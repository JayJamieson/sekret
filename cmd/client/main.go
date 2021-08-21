package main

import (
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

// generate secret, output location and password
func (client *Client) gen(s string) error {
	fmt.Println("generating")
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
	// fs.Func("gen", "Generate secret", client.gen)
	gen := fs.Bool("gen", false, "Generate secret")

	err = fs.Parse(os.Args[1:])

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%v", *gen)

	//
	//fileName := filepath.Join(wd, secretName)
	//fmt.Printf("Writing secret to %v\n", fileName)
	//
	//secretFile, err := os.Create(fileName)
	//
	//fmt.Printf("Fetching secret by name - %v\n", secretName)
	//
	//// TODO http request to fetchin secret
	//if err != nil {
	//	panic(err)
	//}
	//
	//defer secretFile.Close()
	//
	//secretFile.WriteString("test")
	//secretFile.Sync()
}

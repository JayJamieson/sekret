package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	secretName := os.Args[1]

	fmt.Printf("Fetching secret by name - %v\n", secretName)

	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	fileName := filepath.Join(wd, secretName)
	fmt.Printf("Writing secret to %v\n", fileName)

	secretFile, err := os.Create(fileName)

	// TODO http request to fetchin secret
	if err != nil {
		panic(err)
	}

	defer secretFile.Close()

	secretFile.WriteString("test")
	secretFile.Sync()
}

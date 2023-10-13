package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/devries/envcrypt"
)

var inputFile string
var outputFile string

func init() {
	flag.StringVar(&inputFile, "i", "-", "input filename or '-' for standard input")
	flag.StringVar(&outputFile, "o", "-", "output filename or '-' for standard output")
}

func main() {
	flag.Parse()

	var fin *os.File
	var err error
	if inputFile == "-" {
		fin = os.Stdin
	} else {
		fin, err = os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open file: %s\n", inputFile)
			os.Exit(1)
		}
	}
	defer fin.Close()

	keyspec := os.Getenv("KMS_KEYSPEC")
	if keyspec == "" {
		fmt.Fprintf(os.Stderr, "The environment variable KMS_KEYSPEC must be set to a Google\n"+
			"Cloud KMS key in the format: projects/{project_id}/locations/{location}/keyRings/{keyring}/cryptoKeys/{key}\n")
		os.Exit(1)
	}

	var fout *os.File
	if outputFile == "-" {
		fout = os.Stdout
	} else {
		fout, err = os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create file: %s\n", outputFile)
			os.Exit(1)
		}
	}
	defer fout.Close()

	ctx := context.Background()
	message, err := envcrypt.EncryptMessage(ctx, keyspec, fin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encrypting message: %s\n", err)
		os.Exit(2)
	}

	err = json.NewEncoder(fout).Encode(message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing message: %s\n", err)
		os.Exit(2)
	}
}

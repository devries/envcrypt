package main

import (
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

	var message envcrypt.EncodedMessage
	err = json.NewDecoder(fin).Decode(&message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading message: %s\n", err)
		os.Exit(2)
	}

	err = envcrypt.DecryptMessage(keyspec, &message, fout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decrypting message: %s\n", err)
		os.Exit(2)
	}
}

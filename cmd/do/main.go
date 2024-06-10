package main

import (
	"io"
	"log"
	"os"

	"github.com/jibaru/do/internal/parser"
	"github.com/jibaru/do/internal/request"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("use: do <filename.do>")
		return
	}

	filename := os.Args[1]
	doFile, err := parser.Filename(filename)
	if err != nil {
		log.Printf("error parsing: %v\n", err)
		return
	}

	log.Println(doFile)

	resp, err := request.Do(*doFile)
	if err != nil {
		log.Printf("error in request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response: %v\n", err)
		return
	}

	log.Printf("response: %s\n", body)
}

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("timeout in get exchange operation")
		}
		panic(err)
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var bid string
	err = json.Unmarshal(response, &bid)
	if err != nil {
		panic(err)
	}

	fmt.Println("Bid value:", bid)

	f, fileCreateErr := os.Create("./cotacao.txt")
	if fileCreateErr != nil {
		panic(fileCreateErr)
	}
	defer f.Close()

	contentToBeWrite := fmt.Sprintf("Dolar: %s", bid)
	_, fileWriteErr := f.WriteString(contentToBeWrite)
	if fileWriteErr != nil {
		panic(fileWriteErr)
	}
}

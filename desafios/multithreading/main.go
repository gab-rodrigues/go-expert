package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type BrasilAPICep struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Erro        string `json:"erro"`
}

type CepError struct {
	StatusCode int
	Message    string
}

func (e CepError) Error() string {
	return fmt.Sprintf("CepError: %d - %s", e.StatusCode, e.Message)
}

func main() {
	brasilApiChannel := make(chan BrasilAPICep)
	viaCepChannel := make(chan ViaCep)
	errChannel := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Digite o cep desejado: ")
	terminalInput, _ := reader.ReadString('\n')
	cep := strings.TrimSpace(terminalInput)

	var finishedByTimeout bool

	// Poderia ter uma validação de formato cep

	go clientBrasilAPI(ctx, cep, brasilApiChannel, errChannel)

	go clientViaCep(ctx, cep, viaCepChannel, errChannel)

	select {
	case brasilAPICep := <-brasilApiChannel:
		fmt.Println("Dados do BrasilAPI:")
		fmt.Println(brasilAPICep)
	case viaCep := <-viaCepChannel:
		fmt.Println("Dados do ViaCep API:")
		fmt.Println(viaCep)
	case <-time.After(1 * time.Second):
		finishedByTimeout = true
		fmt.Println("Timeout")
		cancel()
	}

	// Mostra os erros das APIs em caso de erro
	if finishedByTimeout {
		for err := range errChannel {
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func clientBrasilAPI(ctx context.Context, cep string, ch chan<- BrasilAPICep, errCh chan<- error) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep + "cep"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		errCh <- CepError{StatusCode: http.StatusInternalServerError, Message: "Error creating request"}
		return
	}

	resp, errHttp := http.DefaultClient.Do(req)
	if errHttp != nil {
		errCh <- CepError{StatusCode: http.StatusInternalServerError, Message: "Error doing request to BrasilAPI API"}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseJson, _ := io.ReadAll(resp.Body)
		errCh <- CepError{StatusCode: resp.StatusCode, Message: "Error returned by request to BrasilAPI API:" + string(responseJson)}
		return
	}

	var brasilAPICep BrasilAPICep
	responseJson, err := io.ReadAll(resp.Body)
	if err != nil {
		errCh <- CepError{StatusCode: http.StatusInternalServerError, Message: "Error reading response from BrasilAPI API: " + err.Error()}
		return
	}
	err = json.Unmarshal(responseJson, &brasilAPICep)
	if err != nil {
		errCh <- CepError{StatusCode: http.StatusInternalServerError, Message: "Error unmarshalling response from BrasilAPI API: " + err.Error()}
		return
	}

	//time.Sleep(1 * time.Second)

	select {
	case ch <- brasilAPICep:
		return
	case <-ctx.Done():
		errCh <- CepError{StatusCode: http.StatusRequestTimeout, Message: "Timeout"}
		return
	}
}

func clientViaCep(ctx context.Context, cep string, ch chan<- ViaCep, errCh chan<- error) {
	url := "http://viacep.com.br/ws/" + cep + "/json"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		errCh <- CepError{StatusCode: http.StatusInternalServerError, Message: "Error creating request"}
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errCh <- CepError{StatusCode: http.StatusInternalServerError, Message: "Error doing request to ViaCep API"}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseErrJson, _ := io.ReadAll(resp.Body)
		errCh <- CepError{StatusCode: resp.StatusCode, Message: "Error returned by request to ViaCep API: " + string(responseErrJson)}
		return
	}

	var viaCep ViaCep
	responseJson, err := io.ReadAll(resp.Body)
	if err != nil {
		errCh <- CepError{StatusCode: http.StatusInternalServerError, Message: "Error reading response from ViaCep API: " + err.Error()}
		return
	}

	err = json.Unmarshal(responseJson, &viaCep)
	if err != nil {
		errCh <- CepError{StatusCode: http.StatusInternalServerError, Message: "Error unmarshalling response from ViaCep API: " + err.Error()}
		return
	}

	if viaCep.Erro != "" {
		errCh <- CepError{StatusCode: http.StatusInternalServerError, Message: "Error returned by request to ViaCep API: " + viaCep.Erro}
		return
	}

	//time.Sleep(1 * time.Second)

	select {
	case ch <- viaCep:
		return
	case <-ctx.Done():
		errCh <- CepError{StatusCode: http.StatusRequestTimeout, Message: "Timeout"}
		return
	}
}

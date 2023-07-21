/**
 * @author Enéas Almeida <eneas.eng@yahoo.com>
 * @description Busca por CEP
 */

package src

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

const URL = "https://viacep.com.br/ws/{cep}/json/"

func GetCEP() {
	url := func(cep string) string {
		return strings.Replace(URL, "{cep}", cep, 1)
	}

	checkCepResult := func(res []byte) bool {
		mapRes := map[string]bool{"erro": false}

		json.Unmarshal(res, &mapRes)

		value, exists := mapRes["erro"]

		if !exists || value {
			return false
		}

		return true
	}

	convertToViaCEP := func(body []byte) (*ViaCEP, error) {
		viaCEP := &ViaCEP{}

		err := json.Unmarshal(body, viaCEP)

		if err != nil {
			return nil, err
		}

		return viaCEP, nil
	}

	for _, cep := range os.Args[1:] {
		req, err := http.Get(url(cep))

		if err != nil || req.StatusCode != 200 {
			fmt.Fprintf(os.Stderr, "Erro ao buscar CEP: %v\n", err)
			continue
		}

		defer req.Body.Close()

		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler resposta do serviço: %v\n", err)
			continue
		}

		isValidResult := checkCepResult(body)

		if !isValidResult {
			fmt.Fprintf(os.Stderr, "CEP inválido: %v\n", cep)
			continue
		}

		viaCEP, err := convertToViaCEP(body)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao converter para o objeto ViaCEP: %v\n", err)
			continue
		}

		fmt.Println(string(body))
		fmt.Println(viaCEP)
	}
}

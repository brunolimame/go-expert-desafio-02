package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	cepNumber      = "49504-337"
	endpointApiCep = "https://cdn.apicep.com/file/apicep/%v.json"
	endpointViaCep = "http://viacep.com.br/ws/%v/json/"
)

type ViaCep struct {
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

type ApiCep struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

type Cep interface {
	ViaCep | ApiCep
}

func main() {
	chApiCep := make(chan string)
	chViaCep := make(chan string)
	var apiCepEntity ApiCep
	var viaCepEntity ViaCep

	go requestApi("viacep.com.br", fmt.Sprintf(endpointViaCep, cepNumber), &viaCepEntity, chViaCep)
	go requestApi("apicep.com", fmt.Sprintf(endpointApiCep, cepNumber), &apiCepEntity, chApiCep)

	select {
	case resApiCep := <-chApiCep:
		fmt.Println(resApiCep, apiCepEntity)
	case resViaCep := <-chViaCep:
		fmt.Println(resViaCep, viaCepEntity)
	case <-time.After(time.Second):
		fmt.Println("timeout")
	}

}

func requestApi[C Cep](name string, urlRequest string, entity *C, ch chan<- string) {

	req, err := http.NewRequest("GET", urlRequest, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&entity)
	if err != nil {
		panic(err)
	}
	ch <- name
}

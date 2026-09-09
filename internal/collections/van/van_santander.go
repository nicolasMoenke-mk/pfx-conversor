package van

import (
	"encoding/json"
	"errors"

	c "io.github.com/conv-pfx/internal/collections"
	b "io.github.com/conv-pfx/internal/collections/body"
)

const (
	santanderVanHost = "trust-open.api.santander.com.br"
	santanderVanBase = "https://" + santanderVanHost
)

type SantanderVanConfig struct {
	ClientId              string
	ClientSecret          string
	SantanderAppKey       string
	SantanderWkId         string
	SantanderClientNumber string
	SantanderBankNumb     string
	SantanderConventNumb  string
}

func AuthVanSantander(cfg SantanderVanConfig) c.PostmanItem {

	return c.PostmanItem{
		Name: "Autenticar Van Santander",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "x-www-form-urlencoded"},
			},
			Body: &c.PostmanBody{
				Mode: "urlencoded",
				Urlencoded: []c.PostmanParam{
					{Key: "grant_type", Value: "client_credentials", Type: "string"},
					{Key: "client_id", Value: cfg.ClientId, Type: "string"},
					{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  santanderVanBase + "/auth/oauth/v2/token",
				Host: []string{santanderVanHost},
				Path: []string{"auth", "oauth", "v2", "token"},
			},
		},
	}

}

func CreateSantanderVan(cfg SantanderVanConfig) (c.PostmanItem, error) {

	createBody := b.DadosBoletoSantander{}

	body, err := json.MarshalIndent(createBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar o body da requisição!")
	}

	return c.PostmanItem{
		Name: "Criar Boleto com QR Code",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
				{Key: "X-Application-Key", Value: cfg.SantanderAppKey, Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(body),
			},
			Url: c.PostmanURL{
				Raw:  santanderVanBase + "/collection_bill_management/v2/workspaces/" + cfg.SantanderWkId + "/bank_slips",
				Host: []string{santanderVanHost},
				Path: []string{"collection_bill_management", "v2", "workspaces", cfg.SantanderWkId, "bank_slips"},
			},
		},
	}, nil
}

func RemoveSantanderVan(cfg SantanderVanConfig) (c.PostmanItem, error) {

	return c.PostmanItem{
		Name: "Remover Boleto",
		Request: &c.PostmanRequest{
			Method: "PATCH",
		},
	}, nil
}

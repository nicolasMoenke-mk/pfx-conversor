package van

import (
	"encoding/json"
	"errors"
	"fmt"

	c "io.github.com/conv-pfx/internal/collections"
)

const (
	sicrediVanHost = "api-parceiro.sicredi.com.br"
	sicrediVanBase = "https://" + sicrediVanHost
	sicrediScope   = "cobranca"
)

type SicrediVanConfig struct {
	SicrediUsername string
	SicrediPassword string
	SicrediApiKey   string
	SicrediCoop     string
	SicrediPosto    string
	SicrediCodBenef string
	SicrediNN       string
	Scope           string
}

func AuthVanSicredi() (c.PostmanItem, error) {

	return c.PostmanItem{
		Name: "Autenticar Van Sicredi",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/x-www-form-urlencoded", Type: "string"},
				{Key: "x-api-key", Value: "{{sicredi_api_key}}", Type: "string"},
				{Key: "context", Value: "COBRANCA", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "urlencoded",
				Urlencoded: []c.PostmanParam{
					{Key: "username", Value: "{{username}}", Type: "string"},
					{Key: "password", Value: "{{password}}", Type: "string"},
					{Key: "grant_type", Value: "password", Type: "string"},
					{Key: "scope", Value: "{{scope}}", Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  sicrediVanBase + "/auth/openapi/token",
				Host: []string{sicrediVanHost},
				Path: []string{"auth", "openapi", "token"},
			},
		},
	}, nil

}

func CreateSicrediVan() (c.PostmanItem, error) {

	createBody := map[string]any{
		"dados": map[string]any{
			"codigoBeneficiario": "{{sicredi_codBenef}}",
			"dataVencimento":     "2027-02-02",
			"seuNumero":          "251454",
			"especieDocumento":   "DUPLICATA_MERCANTIL_INDICACAO",
			"pagador": map[string]any{
				"documento":  "96050176876",
				"nome":       "CLIENTE TESTE MK",
				"tipoPessoa": "PESSOA_FISICA",
				"endereco":   "Rua Amazonas",
				"cidade":     "Votuporanga",
				"cep":        "15500004",
				"uf":         "SP",
			},
			"tipoCobranca": "HIBRIDO",
			"valor":        "10",
			"tipoJuros":    "PERCENTUAL",
			"juros":        "0.33",
			"multa":        "1",
		},
	}

	body, err := json.MarshalIndent(createBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar body da requisição!")
	}

	return c.PostmanItem{
		Name: "Criar boleto Sicredi",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
				{Key: "Cooperativa", Value: "{{sicredi_coop}}", Type: "string"},
				{Key: "Posto", Value: "{{sicredi_posto}}", Type: "string"},
				{Key: "X-Api-Key", Value: "{{sicredi_api_key}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(body),
			},
			Url: c.PostmanURL{
				Raw:  sicrediVanBase + "/cobranca/boleto/v1/boletos",
				Host: []string{sicrediVanHost},
				Path: []string{"cobranca", "boleto", "v1", "boletos"},
			},
		},
	}, nil
}

func RemoveSicrediVan() (c.PostmanItem, error) {

	return c.PostmanItem{
		Name: "Remover boleto Sicredi",
		Request: &c.PostmanRequest{
			Method: "PATCH",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
				{Key: "Cooperativa", Value: "{{sicredi_coop}}", Type: "string"},
				{Key: "Posto", Value: "{{sicredi_posto}}", Type: "string"},
				{Key: "X-Api-Key", Value: "{{sicredi_api_key}}", Type: "string"},
				{Key: "CodigoBeneficiario", Value: "{{sicredi_codBenef}}", Type: "string"},
			},
			Url: c.PostmanURL{
				Raw:  sicrediVanBase + "/cobranca/boleto/v1/boletos/{{sicredi_nossoNumero}}/baixa",
				Host: []string{sicrediVanHost},
				Path: []string{"cobranca", "boleto", "v1", "boletos", "{{sicredi_nossoNumero}}", "baixa"},
			},
		},
	}, nil

}

func SicrediVanCollection(cfg *SicrediVanConfig) ([]byte, error) {

	if cfg == nil {
		return nil, fmt.Errorf("Configuração não pode ser nula")
	}

	scope := cfg.Scope
	if scope == "" {
		scope = sicrediScope
	}

	auth, err := AuthVanSicredi()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Auth Van Sicredi'")
	}

	create, err := CreateSicrediVan()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Criar Boleto Sicredi': %w", err)
	}

	remove, err := RemoveSicrediVan()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Remover Boleto Sicredi': %w", err)
	}

	collection := c.PostmanCollection{
		Info: c.PostmanInfo{
			Name:        "Sicredi Van",
			Description: "Collection Completa Sicredi Van",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []c.PostmanVariable{
			{Key: "bearer_token", Value: "", Type: "string"},
			{Key: "username", Value: cfg.SicrediUsername, Type: "string"},
			{Key: "password", Value: cfg.SicrediPassword, Type: "string"},
			{Key: "x-api-key", Value: cfg.SicrediApiKey, Type: "string"},
			{Key: "cooperativa", Value: cfg.SicrediCoop, Type: "string"},
			{Key: "posto", Value: cfg.SicrediPosto, Type: "string"},
			{Key: "sicredi_codBenef", Value: cfg.SicrediCodBenef, Type: "string"},
			{Key: "sicredi_nossoNumero", Value: cfg.SicrediNN, Type: "string"},
			{Key: "scope", Value: scope, Type: "string"},
		},
		Item: []c.PostmanItem{
			auth,
			create,
			remove,
		},
	}
	return json.MarshalIndent(collection, "", "  ")
}

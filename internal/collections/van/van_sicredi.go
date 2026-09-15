package van

import (
	"encoding/json"
	"errors"
	"fmt"

	c "io.github.com/conv-pfx/internal/collections"
	b "io.github.com/conv-pfx/internal/collections/body"
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

func AuthVanSicredi() c.PostmanItem {

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
	}

}

func CreateSicrediVan() (c.PostmanItem, error) {

	createBody := b.DadosBoletoSicredi{
		CodBeneficiario: "{{sicredi_cod_benef}}",
		DtVencimento:    "2027-02-02",
		SeuNumero:       "251454",
		EspecieDoc:      "DUPLICATA_MERCANTIL_INDICACAO",
		Pagador: b.SicrediPagador{
			Documento:  "96050176876",
			Nome:       "CLIENTE TESTE MK",
			TipoPessoa: "PESSOA_FISICA",
			Endereco:   "Rua Amazonas",
			Cidade:     "Votuporanga",
			Cep:        "15500004",
			Uf:         "SP",
		},
		TipoConbranca: "HIBRIDO",
		Valor:         "10",
		TipoJuros:     "PERCENTUAL",
		Juros:         "0.33",
		Multa:         "1",
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

func RemoveSicrediVan() c.PostmanItem {

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
				{Key: "CodigoBeneficiario", Value: "{{sicredi_cod_benef}}", Type: "string"},
			},
			Url: c.PostmanURL{
				Raw:  sicrediVanBase + "/cobranca/boleto/v1/boletos/{{sicredi_nn}}/baixa",
				Host: []string{sicrediVanHost},
				Path: []string{"cobranca", "boleto", "v1", "boletos", "{{sicredi_nn}}", "baixa"},
			},
		},
	}

}

func SicrediVanCollection(cfg *SicrediVanConfig) ([]byte, error) {

	if cfg == nil {
		return nil, fmt.Errorf("Configuração não pode ser nula")
	}

	scope := cfg.Scope
	if scope == "" {
		scope = sicrediScope
	}

	auth := AuthVanSicredi()

	create, err := CreateSicrediVan()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Criar Boleto': %w", err)
	}

	remove := RemoveSicrediVan()

	collection := c.PostmanCollection{
		Info: c.PostmanInfo{
			Name:        "Sicredi Van",
			Description: "Collection Completa Sicredi Van",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []c.PostmanVariable{
			{Key: "bearer_token", Value: "", Type: "string"},
			{Key: "sicredi_username", Value: cfg.SicrediUsername, Type: "string"},
			{Key: "sicredi_password", Value: cfg.SicrediPassword, Type: "string"},
			{Key: "x-api-key", Value: cfg.SicrediApiKey, Type: "string"},
			{Key: "cooperativa", Value: cfg.SicrediCoop, Type: "string"},
			{Key: "posto", Value: cfg.SicrediPosto, Type: "string"},
			{Key: "sicredi_cod_benef", Value: cfg.SicrediCodBenef, Type: "string"},
			{Key: "sicredi_nn", Value: cfg.SicrediNN, Type: "string"},
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

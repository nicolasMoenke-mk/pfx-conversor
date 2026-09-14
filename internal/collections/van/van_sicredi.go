package van

import (
	"encoding/json"
	"errors"

	c "io.github.com/conv-pfx/internal/collections"
	b "io.github.com/conv-pfx/internal/collections/body"
)

const (
	sicrediVanHost = "api-parceiro.sicredi.com.br"
	sicrediVanBase = "https://" + sicrediVanHost
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

func AuthVanSicredi(cfg SicrediVanConfig) c.PostmanItem {

	cfg.Scope = "cobranca"

	return c.PostmanItem{
		Name: "Autenticar Van Sicredi",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/x-www-form-urlencoded", Type: "string"},
				{Key: "x-api-key", Value: cfg.SicrediApiKey, Type: "string"},
				{Key: "context", Value: "COBRANCA", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "urlencoded",
				Urlencoded: []c.PostmanParam{
					{Key: "username", Value: cfg.SicrediUsername, Type: "string"},
					{Key: "password", Value: cfg.SicrediPassword, Type: "string"},
					{Key: "grant_type", Value: "password", Type: "string"},
					{Key: "scope", Value: cfg.Scope, Type: "string"},
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

func CreateSicrediVan(cfg SicrediVanConfig) (c.PostmanItem, error) {

	createBody := b.DadosBoletoSicredi{
		CodBeneficiario: cfg.SicrediCodBenef,
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
				{Key: "Cooperativa", Value: cfg.SicrediCoop, Type: "string"},
				{Key: "Posto", Value: cfg.SicrediPosto, Type: "string"},
				{Key: "X-Api-Key", Value: cfg.SicrediApiKey, Type: "string"},
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

func RemoveSicrediVan(cfg SicrediVanConfig) c.PostmanItem {

	return c.PostmanItem{
		Name: "Remover boleto Sicredi",
		Request: &c.PostmanRequest{
			Method: "PATCH",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
				{Key: "Cooperativa", Value: cfg.SicrediCoop, Type: "string"},
				{Key: "Posto", Value: cfg.SicrediPosto, Type: "string"},
				{Key: "X-Api-Key", Value: cfg.SicrediApiKey, Type: "string"},
				{Key: "CodigoBeneficiario", Value: cfg.SicrediCodBenef, Type: "string"},
			},
			Url: c.PostmanURL{
				Raw:  sicrediVanBase + "/cobranca/boleto/v1/boletos/" + cfg.SicrediNN + "/baixa",
				Host: []string{sicrediVanHost},
				Path: []string{"cobranca", "boleto", "v1", "boletos", cfg.SicrediNN, "baixa"},
			},
		},
	}

}

func SicrediVanCollection(cfg SicrediVanConfig) ([]byte, error) {

	auth := AuthVanSicredi(cfg)

	create, err := CreateSicrediVan(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	remove := RemoveSicrediVan(cfg)

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
			{Key: "scope", Value: cfg.Scope, Type: "string"},
		},
		Item: []c.PostmanItem{
			auth,
			create,
			remove,
		},
	}
	return json.MarshalIndent(collection, "", "  ")
}

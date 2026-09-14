package pix

import (
	"encoding/json"
	"errors"

	c "io.github.com/conv-pfx/internal/collections"
	b "io.github.com/conv-pfx/internal/collections/body"
)

const (
	santanderPixHost = "trust-pix.santander.com.br"
	santanderPixBase = "https://" + santanderPixHost
	santanderTXID    = "AHUDAHIUDHAKDUHA"
)

type SantanderPixConfig struct {
	ClientId     string
	ClientSecret string
	ChavePix     string
	Scope        string
}

func AuthPixSantander(cfg SantanderPixConfig) c.PostmanItem {

	cfg.Scope = "pix.read pix.write cobv.read cobv.write webhook.write webhook.read"

	return c.PostmanItem{
		Name: "Autenticar Pix Santander",
		Request: &c.PostmanRequest{
			Auth: &c.PostmanAuth{
				Type: "basic",
				Basic: []c.PostmanAuthParam{
					{Key: "client_id", Value: cfg.ClientId, Type: "string"},
					{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
				},
			},
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/x-www-form-urlencoded", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "urlencoded",
				Urlencoded: []c.PostmanParam{
					{Key: "grant_type", Value: "client_credentials", Type: "string"},
					{Key: "client_id", Value: cfg.ClientId, Type: "string"},
					{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
					{Key: "scope", Value: cfg.Scope, Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  santanderPixBase + "/oauth/token?grant_type=client_credentials",
				Host: []string{santanderPixHost},
				Path: []string{"oauth", "token?grant_type=client_credentials"},
			},
		},
	}

}

func CreatePixSantander(cfg SantanderPixConfig) (c.PostmanItem, error) {

	createPix := b.SantanderCriarPix{
		Calendario: b.SantanderCalendario{
			DtVencimento:          "2027-02-02",
			ValidadePosVencimento: 1,
		},
		Chave: b.SantanderChavePix{
			Chave: cfg.ChavePix,
		},
		Devedor: b.SantanderPixDevedor{
			Cpf:  "96050176876",
			Nome: "MK - Teste Comunicação-Call",
		},
		Valor: b.SantanderPixValor{
			Valor: "0.25",
		},
	}

	body, err := json.MarshalIndent(createPix, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar body da requisição!")
	}

	return c.PostmanItem{
		Name: "Criar um Pix",
		Request: &c.PostmanRequest{
			Method: "PUT",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(body),
			},
			Url: c.PostmanURL{
				Raw:  santanderPixBase + "/api/v1/cobv/" + santanderTXID,
				Host: []string{santanderPixHost},
				Path: []string{"api", "v1", "cobv", santanderTXID},
			},
		},
	}, nil

}

func RemoverPixSantander(cfg SantanderPixConfig) (c.PostmanItem, error) {

	removeBody := b.SantanderStatusPix{
		Status: "REMOVIDA_PELO_USUARIO_RECEBEDOR",
	}

	rm, err := json.MarshalIndent(removeBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar o body da requisição!")
	}

	return c.PostmanItem{
		Name: "Remover Pix Ativo",
		Request: &c.PostmanRequest{
			Method: "PATCH",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(rm),
			},
			Url: c.PostmanURL{
				Raw:  santanderPixBase + "/api/v1/cobv/" + santanderTXID,
				Host: []string{santanderPixHost},
				Path: []string{"api", "v1", "cobv", santanderTXID},
			},
		},
	}, nil
}

func SantanderPixCollection(cfg SantanderPixConfig) ([]byte, error) {

	auth := AuthPixSantander(cfg)

	create, err := CreatePixSantander(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	remove, err := RemoverPixSantander(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	collection := c.PostmanCollection{
		Info: c.PostmanInfo{
			Name:        "Collection Pix Santander",
			Description: "Collection completa Santander",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []c.PostmanVariable{
			{Key: "bearer_token", Value: "", Type: "string"},
			{Key: "client_id", Value: cfg.ClientId, Type: "string"},
			{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
			{Key: "chave_pix", Value: cfg.ChavePix, Type: "string"},
			{Key: "scope", Value: cfg.Scope, Type: "string"},
			{Key: "pix_txid", Value: santanderTXID, Type: "string"},
		},
		Item: []c.PostmanItem{
			auth,
			create,
			remove,
		},
	}
	return json.MarshalIndent(collection, "", "  ")
}

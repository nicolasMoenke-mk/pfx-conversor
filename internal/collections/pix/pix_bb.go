package pix

import (
	"encoding/json"
	"errors"

	c "io.github.com/conv-pfx/internal/collections"
	b "io.github.com/conv-pfx/internal/collections/body"
)

const (
	bbPixAuthHost = "api.bb.com.br"
	bbPixAuthBase = "https://" + bbPixAuthHost
	bbPixHost     = "api-pix.bb.com.br"
	bbPixBase     = "https://" + bbPixHost
	bbTXID        = "AJIDJALISJSHDKUSHD"
)

type BancoDoBrasilPixConfig struct {
	ClientId     string
	ClientSecret string
	BBChavePix   string
	BBAppKey     string
	Scope        string
}

func AuthPixBB(cfg BancoDoBrasilPixConfig) c.PostmanItem {

	cfg.Scope = "cobv.write cobv.read pix.read"

	return c.PostmanItem{
		Name: "Autenticar Pix BB",
		Request: &c.PostmanRequest{
			Auth: &c.PostmanAuth{
				Type: "basic",
				Basic: []c.PostmanAuthParam{
					{Key: "username", Value: cfg.ClientId, Type: "string"},
					{Key: "password", Value: cfg.ClientSecret, Type: "string"},
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
					{Key: "scope", Value: cfg.Scope, Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  bbPixAuthBase + "/oauth/token",
				Host: []string{bbPixAuthHost},
				Path: []string{"oauth", "token"},
			},
		},
	}
}

func CreatePixBB(cfg BancoDoBrasilPixConfig) (c.PostmanItem, error) {

	createBody := b.BBCriarPix{
		Calendario: b.BBCalendario{
			DtVencimento:          "2027-02-02",
			ValidadePosVencimento: 1,
		},
		ChavePix: b.BBChavePix{
			Chave: cfg.BBChavePix,
		},
		Devedor: b.BBDevedor{
			Cpf:  "96050176876",
			Nome: "MK - Teste Comunicação-Call",
		},
		Valor: b.BBValor{
			Valor: "0.25",
		},
	}

	body, err := json.MarshalIndent(createBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar body da requisição!")
	}

	return c.PostmanItem{
		Name: "Criar Pix BB",
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
				Raw:  bbPixBase + "/pix/v2/cobv/" + bbTXID + "?gw-dev-app-key=" + cfg.BBAppKey,
				Host: []string{bbPixHost},
				Path: []string{"pix", "v2", "cobv", bbTXID, "?gw-dev-app-key=", cfg.BBAppKey},
			},
		},
	}, nil

}

func RemovePixBB(cfg BancoDoBrasilPixConfig) (c.PostmanItem, error) {

	removeBody := b.BBStatusPix{
		Status: "REMOVIDA_PELO_USUARIO_RECEBEDOR",
	}

	rm, err := json.MarshalIndent(removeBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar o body da requisição!")
	}

	return c.PostmanItem{
		Name: "Remover Pix BB",
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
				Raw:  bbPixBase + "/pix/v2/cobv/" + bbTXID + "?gw-dev-app-key=" + cfg.BBAppKey,
				Host: []string{bbPixHost},
				Path: []string{"pix", "v2", "cobv", bbTXID, "?gw-dev-app-key=", cfg.BBAppKey},
			},
		},
	}, nil
}

func BancoDoBrasilPixCollection(cfg BancoDoBrasilPixConfig) ([]byte, error) {

	auth := AuthPixBB(cfg)

	create, err := CreatePixBB(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	remove, err := RemovePixBB(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	collection := c.PostmanCollection{
		Info: c.PostmanInfo{
			Name:        "Pix BB",
			Description: "Collection completa BB",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []c.PostmanVariable{
			{Key: "client_id", Value: cfg.ClientId, Type: "string"},
			{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
			{Key: "bb_dev_app_key", Value: cfg.BBAppKey, Type: "string"},
			{Key: "chave_pix", Value: cfg.BBChavePix, Type: "string"},
			{Key: "pix_txid", Value: bbTXID, Type: "string"},
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

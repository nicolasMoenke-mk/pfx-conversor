package pix

import (
	"encoding/json"
	"errors"
	"fmt"

	c "io.github.com/conv-pfx/internal/collections"
)

const (
	bbPixAuthHost = "api.bb.com.br"
	bbPixAuthBase = "https://" + bbPixAuthHost
	bbPixHost     = "api-pix.bb.com.br"
	bbPixBase     = "https://" + bbPixHost
	bbTXID        = "AJIDJALISJSHDKUSHD"
	bbScope       = "cobv.write cobv.read pix.read"
)

type BancoDoBrasilPixConfig struct {
	ClientId     string
	ClientSecret string
	BBChavePix   string
	BBAppKey     string
	Scope        string
}

func AuthPixBB() (c.PostmanItem, error) {

	return c.PostmanItem{
		Name:  "Autenticar Pix BB",
		Event: c.AutoSetBearerTokenEvent(),
		Request: &c.PostmanRequest{
			Auth: &c.PostmanAuth{
				Type: "basic",
				Basic: []c.PostmanAuthParam{
					{Key: "username", Value: "{{client_id}}", Type: "string"},
					{Key: "password", Value: "{{client_secret}}", Type: "string"},
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
					{Key: "client_id", Value: "{{client_id}}", Type: "string"},
					{Key: "scope", Value: "{{scope}}", Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  bbPixAuthBase + "/oauth/token",
				Host: []string{bbPixAuthHost},
				Path: []string{"oauth", "token"},
			},
		},
	}, nil
}

func CreatePixBB() (c.PostmanItem, error) {

	createBody := map[string]any{
		"calendario": map[string]any{
			"dataDeVencimento":       "2027-02-02",
			"validadeAposVencimento": 1,
		},
		"chave": map[string]any{
			"chave_pix": "{{chave_pix}}",
		},
		"devedor": map[string]any{
			"cpf":  "96050176876",
			"nome": "MK - Teste Comunicação-Call",
		},
		"valor": map[string]any{
			"original": "0.25",
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
			Body: c.RawJSONBody(string(body)),
			Url: c.PostmanURL{
				Raw:  bbPixBase + "/pix/v2/cobv/{{pix_txid}}?gw-dev-app-key={{bb_dev_app_key}}",
				Host: []string{bbPixHost},
				Path: []string{"pix", "v2", "cobv", "{{pix_txid}}", "?gw-dev-app-key=", "{{bb_dev_app_key}}"},
			},
		},
	}, nil

}

func RemovePixBB() (c.PostmanItem, error) {

	removeBody := map[string]any{
		"status": "REMOVIDA_PELO_USUARIO_RECEBEDOR",
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
			Body: c.RawJSONBody(string(rm)),
			Url: c.PostmanURL{
				Raw:  bbPixBase + "/pix/v2/cobv/{{pix_txid}}?gw-dev-app-key={{bb_dev_app_key}}",
				Host: []string{bbPixHost},
				Path: []string{"pix", "v2", "cobv", "{{pix_txid}}", "?gw-dev-app-key=", "{{bb_dev_app_key}}"},
			},
		},
	}, nil
}

func BancoDoBrasilPixCollection(cfg *BancoDoBrasilPixConfig) ([]byte, error) {

	auth, err := AuthPixBB()
	if err != nil {
		return nil, fmt.Errorf("Configuração não pode ser nula")
	}

	create, err := CreatePixBB()
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	remove, err := RemovePixBB()
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	scope := cfg.Scope
	if scope == "" {
		scope = bbScope
	}

	collection := c.PostmanCollection{
		Info: c.PostmanInfo{
			Name:        "Pix BB",
			Description: "Collection completa BB",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []c.PostmanVariable{
			{Key: "bearer_token", Value: "", Type: "string"},
			{Key: "client_id", Value: cfg.ClientId, Type: "string"},
			{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
			{Key: "bb_dev_app_key", Value: cfg.BBAppKey, Type: "string"},
			{Key: "chave_pix", Value: cfg.BBChavePix, Type: "string"},
			{Key: "pix_txid", Value: bbTXID, Type: "string"},
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

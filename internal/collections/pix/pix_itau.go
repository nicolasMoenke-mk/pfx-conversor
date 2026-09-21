package pix

import (
	"encoding/json"
	"fmt"

	c "io.github.com/conv-pfx/internal/collections"
)

const (
	itauHostAuth    = "sts.itau.com.br"
	itauBaseUrlAuth = "https://" + itauHostAuth
	itauHostApi     = "api.itau.com.br"
	itauBaseUrlApi  = "https://" + itauHostApi
	itauTXID        = "HAJSHDSLJMFIDHLU"
)

type ItauPixConfig struct {
	ClientId     string
	ClientSecret string
	ItauAppKey   string
	ChavePix     string
}

func authPixItau() (c.PostmanItem, error) {

	return c.PostmanItem{
		Name: "Login Pix Itau",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "x-www-form-urlencoded", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "urlencoded",
				Urlencoded: []c.PostmanParam{
					{Key: "grant_type", Value: "client_credentials", Type: "string"},
					{Key: "client_id", Value: "{{client_id}}", Type: "string"},
					{Key: "client_secret", Value: "{{client_secret}}", Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  itauBaseUrlAuth + "/api/oauth/token",
				Host: []string{itauHostAuth},
				Path: []string{"api", "oauth", "token"},
			},
		},
	}, nil
}

func createPixItau() (c.PostmanItem, error) {

	createBody := map[string]any{
		"calendario": map[string]any{
			"dataDeVencimento":       "2027-02-02",
			"validadeAposVencimento": 1,
		},
		"chave": "{{chave_pix}}",
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
		return c.PostmanItem{}, fmt.Errorf("Erro ao serializar o body da requisição")
	}

	return c.PostmanItem{
		Name: "Criar Pix Itau",
		Request: &c.PostmanRequest{
			Method: "PUT",
			Header: []c.PostmanHeader{
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "x-itau-apikey", Value: "{{itau_app_key}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(body),
			},
			Url: c.PostmanURL{
				Raw:  itauBaseUrlApi + "/pix/v2/cob/{{pix_txid}}",
				Host: []string{itauHostApi},
				Path: []string{"pix", "v2", "cob", "{{pix_txid}}"},
			},
		},
	}, nil
}

func removePixItau() (c.PostmanItem, error) {

	removeBody := map[string]any{
		"status": "REMOVIDA_PELO_USUARIO_RECEBEDOR",
	}

	rm, err := json.MarshalIndent(removeBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, fmt.Errorf("Erro ao serializar body da requisição")
	}

	return c.PostmanItem{
		Name: "Remover Pix Itau",
		Request: &c.PostmanRequest{
			Method: "PATCH",
			Header: []c.PostmanHeader{
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "x-itau-apikey", Value: "{{itau_app_key}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(rm),
			},
			Url: c.PostmanURL{
				Raw:  itauBaseUrlApi + "/pix/v2/cob/{{pix_txid}}",
				Host: []string{itauHostApi},
				Path: []string{"pix", "v2", "cob", "{{pix_txid}}"},
			},
		},
	}, nil
}

func ItauPixCollection(cfg *ItauPixConfig) ([]byte, error) {

	if cfg == nil {
		return nil, fmt.Errorf("Configuração não pode ser nula!")
	}

	auth, err := authPixItau()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Auth Itau Pix': %w", err)
	}

	create, err := createPixItau()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Criar Itau Pix': %w", err)
	}

	remove, err := removePixItau()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Criar Itau Pix': %w", err)
	}

	collection := c.PostmanCollection{
		Info: c.PostmanInfo{
			Name:        "Pix Itau",
			Description: "Collection Pix Completa Itau",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []c.PostmanVariable{
			{Key: "bearer_token", Value: "", Type: "string"},
			{Key: "client_id", Value: cfg.ClientId, Type: "string"},
			{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
			{Key: "x-itau-apikey", Value: cfg.ItauAppKey, Type: "string"},
			{Key: "chave_pix", Value: cfg.ChavePix, Type: "string"},
			{Key: "pix_txid", Value: itauTXID, Type: "string"},
		},
		Item: []c.PostmanItem{
			auth,
			create,
			remove,
		},
	}

	return json.MarshalIndent(collection, "", "  ")
}

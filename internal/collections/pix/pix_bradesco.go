package pix

import (
	"encoding/json"
	"errors"
	"fmt"

	c "io.github.com/conv-pfx/internal/collections"
)

const (
	bradescoPixHost = "qrpix.bradesco.com.br"
	bradescoPixBase = "https://" + bradescoPixHost
	bradescoTXID    = "HUAJDKJSIOHHUAHSUGD"
	bradescoScope   = "pix.read pix.write cobv.read cobv.write webhook.write webhook.read"
)

type BradescoPixConfig struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ChavePix     string `json:"chave_pix"`
	Scope        string `json:"scope"`
}

func AuthBradescoPix() (c.PostmanItem, error) {

	return c.PostmanItem{
		Name: "Auth Pix Bradesco",
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
				{Key: "Content-Type", Value: "application/x-www-form-urlencoded"},
			},
			Body: &c.PostmanBody{
				Mode: "urlencoded",
				Urlencoded: []c.PostmanParam{
					{Key: "grant_type", Value: "client_credentials", Type: "string"},
					{Key: "client_id", Value: "{{client_id}}", Type: "string"},
					{Key: "client_secret", Value: "{{client_secret}}", Type: "string"},
					{Key: "scope", Value: "{{scope}}", Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  bradescoPixBase + "/oauth/token",
				Host: []string{bradescoPixHost},
				Path: []string{"ouath", "token"},
			},
		},
	}, nil

}

func CreatePixBradesco() (c.PostmanItem, error) {

	createBody := map[string]any{
		"calendario": map[string]any{
			"dataVencimento":         "2027-02-02",
			"validadeAposVencimento": 1,
		},
		"chave": map[string]any{
			"chave": "{{chave_pix}}",
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
		return c.PostmanItem{}, errors.New("Erro ao serializar o body da requisição")
	}

	return c.PostmanItem{
		Name: "Criar Pix Bradesco",
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
				Raw:  bradescoPixBase + "/v2/cobv/{{pix_txid}}",
				Host: []string{bradescoPixHost},
				Path: []string{"v2", "cobv", "{{pix_txid}}"},
			},
		},
	}, nil

}

func RemoverPixBradesco() (c.PostmanItem, error) {

	removeBody := map[string]any{
		"status": "REMOVIDA_PELO_USUARIO_RECEBEDOR",
	}

	rm, err := json.MarshalIndent(removeBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar o body da requisição")
	}

	return c.PostmanItem{
		Name: "Remover Pix Bradesco",
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
				Raw:  bradescoPixBase + "/v2/cobv/{{pix_txid}}",
				Host: []string{bradescoPixHost},
				Path: []string{"v2", "cobv", "{{pix_txid}}"},
			},
		},
	}, nil
}

func BradescoPixCollection(cfg *BradescoPixConfig) ([]byte, error) {

	if cfg == nil {
		return nil, fmt.Errorf("Configuração não pode ser nula")
	}

	scope := bradescoScope
	if scope == "" {
		scope = bradescoScope
	}

	auth, err := AuthBradescoPix()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Auth Pix Bradesco': %w", err)
	}

	create, err := CreatePixBradesco()
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	remove, err := RemoverPixBradesco()
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	collection := c.PostmanCollection{
		Info: c.PostmanInfo{
			Name:        "Bradesco Pix Avulso",
			Description: "Collection completa Bradesco",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []c.PostmanVariable{
			{Key: "bearer_token", Value: "", Type: "string"},
			{Key: "client_id", Value: cfg.ClientId, Type: "string"},
			{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
			{Key: "chave_pix", Value: cfg.ChavePix, Type: "string"},
			{Key: "scope", Value: scope, Type: "string"},
			{Key: "pix_txid", Value: bradescoTXID, Type: "string"},
		},
		Item: []c.PostmanItem{
			auth,
			create,
			remove,
		},
	}

	return json.MarshalIndent(collection, "", "  ")
}

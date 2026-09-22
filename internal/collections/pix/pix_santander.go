package pix

import (
	"encoding/json"
	"fmt"

	c "io.github.com/conv-pfx/internal/collections"
)

const (
	santanderPixHost = "trust-pix.santander.com.br"
	santanderPixBase = "https://" + santanderPixHost
	santanderTXID    = "AHUDAHIUDHAKDUHA"
	santanderScope   = "pix.read pix.write cobv.read cobv.write webhook.write webhook.read"
)

type SantanderPixConfig struct {
	ClientId     string
	ClientSecret string
	ChavePix     string
}

func AuthPixSantander() (c.PostmanItem, error) {

	return c.PostmanItem{
		Name:  "Autenticar Pix Santander",
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
					{Key: "scope", Value: "{{scope}}", Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  santanderPixBase + "/oauth/token",
				Host: []string{santanderPixHost},
				Path: []string{"oauth", "token"},
			},
		},
	}, nil
}

func CreatePixSantander() (c.PostmanItem, error) {
	createPix := map[string]any{
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

	body, err := json.MarshalIndent(createPix, "", "  ")
	if err != nil {
		return c.PostmanItem{}, fmt.Errorf("erro ao serializar body de criação Pix Santander: %w", err)
	}

	return c.PostmanItem{
		Name: "Criar um Pix",
		Request: &c.PostmanRequest{
			Method: "PUT",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
			},
			Body: c.RawJSONBody(string(body)),
			Url: c.PostmanURL{
				Raw:  santanderPixBase + "/api/v1/cobv/{{pix_txid}}",
				Host: []string{santanderPixHost},
				Path: []string{"api", "v1", "cobv", "{{pix_txid}}"},
			},
		},
	}, nil
}

func RemoverPixSantander() (c.PostmanItem, error) {
	removeBody := map[string]any{
		"status": "REMOVIDA_PELO_USUARIO_RECEBEDOR",
	}

	rm, err := json.MarshalIndent(removeBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, fmt.Errorf("erro ao serializar o body de remoção Pix Santander: %w", err)
	}

	return c.PostmanItem{
		Name: "Remover Pix Ativo",
		Request: &c.PostmanRequest{
			Method: "PATCH",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
			},
			Body: c.RawJSONBody(string(rm)),
			Url: c.PostmanURL{
				Raw:  santanderPixBase + "/api/v1/cobv/{{pix_txid}}",
				Host: []string{santanderPixHost},
				Path: []string{"api", "v1", "cobv", "{{pix_txid}}"},
			},
		},
	}, nil
}

func SantanderPixCollection(cfg *SantanderPixConfig) ([]byte, error) {

	if cfg == nil {
		return nil, fmt.Errorf("Configuração não pode ser nula")
	}

	scope := santanderScope
	if scope == "" {
		scope = santanderScope
	}

	auth, err := AuthPixSantander()
	if err != nil {
		return nil, fmt.Errorf("Falha ao gerar item 'Auth Pix Santander'")
	}

	create, err := CreatePixSantander()
	if err != nil {
		return nil, fmt.Errorf("Falha ao gerar item 'Criar Pix Santander': %w", err)
	}

	remove, err := RemoverPixSantander()
	if err != nil {
		return nil, fmt.Errorf("Falha ao gerar item 'Remover Pix Santander': %w", err)
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
			{Key: "scope", Value: scope, Type: "string"},
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

package van

import (
	"encoding/json"
	"fmt"

	c "io.github.com/conv-pfx/internal/collections"
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
	SantanderType         string
	SantanderDictKey      string
}

func AuthVanSantander() (c.PostmanItem, error) {
	return c.PostmanItem{
		Name:  "Autenticar Van Santander",
		Event: c.AutoSetBearerTokenEvent(),
		Request: &c.PostmanRequest{
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
				},
			},
			Url: c.PostmanURL{
				Raw:  santanderVanBase + "/auth/oauth/v2/token",
				Host: []string{santanderVanHost},
				Path: []string{"auth", "oauth", "v2", "token"},
			},
		},
	}, nil
}

func CreateSantanderVan() (c.PostmanItem, error) {
	createBody := map[string]any{
		"clientNumber": "{{santander_client_number}}",
		"nsuCode":      "{{santander_bank_number}}",
		"covenantCode": "{{santander_convent_code}}",
		"bankNumber":   "{{santander_bank_number}}",
		"dueDate":      "2027-02-02",
		"nsuDate":      "2026-09-15",
		"issueDate":    "2026-09-15",
		"key": map[string]any{
			"type":    "{{santander_dict_type}}",
			"dictKey": "{{santander_dict_key}}",
		},
		"environment":  "PRODUCAO",
		"nominalValue": "10.00",
		"payer": map[string]any{
			"name":           "CLIENTE TESTE MK",
			"documentType":   "CPF",
			"documentNumber": "96050176876",
			"address":        "Rua Amazonas",
			"neighborhood":   "Patrimonio Novo",
			"city":           "Votuporanga",
			"state":          "SP",
			"zipCode":        "15500-004",
		},
		"documentKind":       "DUPLICATA_MERCANTIL",
		"finePercentage":     "2.00",
		"fineQuantityDays":   "0",
		"interestPercentage": "0.99",
		"protestType":        "SEM_PROTESTO",
		"paymentType":        "REGISTRO",
	}

	body, err := json.MarshalIndent(createBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, fmt.Errorf("erro ao serializar o body do boleto Santander: %w", err)
	}

	return c.PostmanItem{
		Name: "Criar Boleto com QR Code",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
				{Key: "X-Application-Key", Value: "{{santander_app_key}}", Type: "string"},
			},
			Body: c.RawJSONBody(string(body)),
			Url: c.PostmanURL{
				Raw:  santanderVanBase + "/collection_bill_management/v2/workspaces/{{santander_workspace_id}}/bank_slips",
				Host: []string{santanderVanHost},
				Path: []string{"collection_bill_management", "v2", "workspaces", "{{santander_workspace_id}}", "bank_slips"},
			},
		},
	}, nil
}

func RemoveSantanderVan() (c.PostmanItem, error) {
	removeBody := map[string]any{
		"convenantCode": "{{santander_convent_code}}",
		"bankNumber":    "{{santander_bank_number}}",
		"operation":     "BAIXA",
	}

	rm, err := json.MarshalIndent(removeBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, fmt.Errorf("erro ao serializar o body de baixa Santander: %w", err)
	}

	return c.PostmanItem{
		Name: "Remover Boleto",
		Request: &c.PostmanRequest{
			Method: "PATCH",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
				{Key: "X-Application-Key", Value: "{{santander_app_key}}", Type: "string"},
			},
			Body: c.RawJSONBody(string(rm)),
			Url: c.PostmanURL{
				Raw:  santanderVanBase + "/collection_bill_management/v2/workspaces/{{santander_workspace_id}}/bank_slips",
				Host: []string{santanderVanHost},
				Path: []string{"collection_bill_management", "v2", "workspaces", "{{santander_workspace_id}}", "bank_slips"},
			},
		},
	}, nil
}

func SantanderVanCollection(cfg *SantanderVanConfig) ([]byte, error) {

	if cfg == nil {
		return nil, fmt.Errorf("Configuração não pode ser nula")
	}

	auth, err := AuthVanSantander()
	if err != nil {
		return nil, fmt.Errorf("Falha ao gerar item 'Auth Van Santander': %w", err)
	}

	create, err := CreateSantanderVan()
	if err != nil {
		return nil, fmt.Errorf("Falha ao gerar item 'Criar Van Santander': %w", err)
	}

	remove, err := RemoveSantanderVan()
	if err != nil {
		return nil, fmt.Errorf("Falha ao gerar item 'Remover Van Santander': %w", err)
	}

	collection := c.PostmanCollection{
		Info: c.PostmanInfo{
			Name:        "Collection Van Santander",
			Description: "Criar boletos com QR Code Dinâmico por API",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []c.PostmanVariable{
			{Key: "bearer_token", Value: "", Type: "string"},
			{Key: "client_id", Value: cfg.ClientId, Type: "string"},
			{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
			{Key: "santander_app_key", Value: cfg.SantanderAppKey, Type: "string"},
			{Key: "santander_workspace_id", Value: cfg.SantanderWkId, Type: "string"},
			{Key: "santander_client_number", Value: cfg.SantanderClientNumber, Type: "string"},
			{Key: "santander_bank_number", Value: cfg.SantanderBankNumb, Type: "string"},
			{Key: "santander_convent_code", Value: cfg.SantanderConventNumb, Type: "string"},
			{Key: "santander_dict_type", Value: cfg.SantanderType, Type: "string"},
			{Key: "santander_dict_key", Value: cfg.SantanderDictKey, Type: "string"},
		},
		Item: []c.PostmanItem{
			auth,
			create,
			remove,
		},
	}

	return json.MarshalIndent(collection, "", "  ")
}

package van

import (
	"encoding/json"
	"errors"
	"fmt"

	c "io.github.com/conv-pfx/internal/collections"
)

const (
	sicoobHostAuthVan = "auth.sicoob.com.br"
	sicoobHostApiVan  = "api.sicoob.com.br"
	sicoobBaseUrlAuth = "https://" + sicoobHostAuthVan
	sicoobBaseUrlApi  = "https://" + sicoobHostApiVan
	sicoobScope       = "boletos_inclusao boletos_consulta boletos_alteracao webhooks_alteracao webhooks_consulta webhooks_inclusao"
)

type SicoobVanConfig struct {
	ClientId         string `json:"client_id"`
	ClientSecret     string `json:"client_secret"`
	Scope            string `json:"scope"`
	NumeroCliente    int    `json:"numero_cliente"`
	CodigoModalidade int    `json:"codigo_modalidade"`
	ContaCorrente    int    `json:"conta_corrente"`
	DataVencimento   string `json:"data_vencimento"`
	DataEmissao      string `json:"data_emissao"`
	SeuNumero        int    `json:"seu_numero"`
	NossoNumero      int    `json:"nosso_numero"`
	IdBoletoEmpresa  int    `json:"id_boleto_empresa"`
}

func AuthVanSicoob() (c.PostmanItem, error) {
	return c.PostmanItem{
		Name: "Autenticação Van Gerar Bearer Token",
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
					{Key: "scope", Value: "{{scope}}", Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  sicoobBaseUrlAuth + "/auth/realms/cooperado/protocol/openid-connect/token",
				Host: []string{sicoobHostAuthVan},
				Path: []string{"auth", "realms", "cooperado", "protocol", "openid-connect", "token"},
			},
		},
	}, nil
}

func CreateSicoobVan() (c.PostmanItem, error) {

	createBody := map[string]any{
		"dados": map[string]any{
			"numeroCliente":                   "{{numeroCliente}}",
			"codigoModalidade":                "{{codigoModalidade}}",
			"numeroContaCorrente":             "{{numeroContaCorrente}}",
			"dataVencimento":                  "{{dataVencimento}}",
			"dataEmissao":                     "{{dataEmissao}}",
			"seuNumero":                       "{{seuNumero}}",
			"nossoNumero":                     "{{nossoNumero}}",
			"identificacaoBoletoEmpresa":      "{{identificacaoBoletoEmpresa}}",
			"codigoEspecieDocumento":          "0",
			"identificacaoEmissaoBoleto":      0,
			"identificacaoDistribuicaoBoleto": 0,
			"valor":                           25.0,
		},
		"pagador": map[string]any{
			"numeroCpfCnpj": "1452875",
			"nome":          "MK Solutions",
			"endereco":      "Rua Amazonas",
			"bairro":        "Patrimonio Novo",
			"cidade":        "Votuporanga",
			"cep":           "15500004",
			"uf":            "SP",
			"email":         "",
		},
		"pix": map[string]any{
			"codigoCadastrarPIX": 1,
		},
	}

	bolBody, err := json.MarshalIndent(createBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar o body da requisição!")
	}

	return c.PostmanItem{
		Name: "Criar Boleto",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
				{Key: "client_id", Value: "{{client_id}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(bolBody),
			},
			Url: c.PostmanURL{
				Raw:  sicoobBaseUrlApi + "/cobranca-bancaria/v3/boletos",
				Host: []string{sicoobHostApiVan},
				Path: []string{"cobranca-bancaria", "v3", "boletos"},
			},
		},
	}, nil
}

func RemoveSicoobVan() (c.PostmanItem, error) {

	rm := map[string]any{
		"boleto": map[string]any{
			"numeroCliente":    "{{numeroCliente}}",
			"codigoModalidade": "{{codigoModalidade}}",
		},
	}

	rmBody, err := json.MarshalIndent(rm, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar o body da requisição!")
	}

	return c.PostmanItem{
		Name: "Remover Boleto",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
				{Key: "client_id", Value: "{{client_id}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(rmBody),
			},
			Url: c.PostmanURL{
				Raw:  sicoobBaseUrlApi + "/cobranca-bancaria/v3/boletos/{{sicoob_nossoNumero}}/baixar",
				Host: []string{sicoobHostApiVan},
				Path: []string{"cobranca-bancaria", "v3", "boletos", "{{sicoob_nossoNumero}}", "baixar"},
			},
		},
	}, nil
}

func SicoobVanCollection(cfg *SicoobVanConfig) ([]byte, error) {

	if cfg == nil {
		return nil, fmt.Errorf("Configuração não pode ser nula!")
	}

	scope := sicoobScope
	if scope == "" {
		scope = sicoobScope
	}

	auth, err := AuthVanSicoob()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Auth Van Sicoob': %w", err)
	}

	create, err := CreateSicoobVan()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Criar Boleto Sicoob': %w", err)
	}

	remove, err := RemoveSicoobVan()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Remover Boleto Sicoob': %w", err)
	}

	collection := c.PostmanCollection{
		Info: c.PostmanInfo{
			Name:        "Sicoob - Criação de Boletos com QR Code",
			Description: "Criar boletos com QR Code Dinâmico por API",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []c.PostmanVariable{
			{Key: "bearer_token", Value: "", Type: "string"},
			{Key: "client_id", Value: cfg.ClientId, Type: "string"},
			{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
			{Key: "numeroCliente", Value: fmt.Sprintf("%d", cfg.NumeroCliente), Type: "string"},
			{Key: "scope", Value: cfg.Scope, Type: "string"},
			{Key: "codigoModalidade", Value: fmt.Sprintf("%d", cfg.CodigoModalidade), Type: "string"},
			{Key: "numeroContaCorrente", Value: fmt.Sprintf("%d", cfg.ContaCorrente), Type: "string"},
			{Key: "dataVencimento", Value: cfg.DataVencimento, Type: "string"},
			{Key: "dataEmissao", Value: cfg.DataEmissao, Type: "string"},
			{Key: "seuNumero", Value: fmt.Sprintf("%d", cfg.SeuNumero), Type: "string"},
			{Key: "nossoNumero", Value: fmt.Sprintf("%d", cfg.NossoNumero), Type: "string"},
			{Key: "identificacaoBoletoEmpresa", Value: fmt.Sprintf("%d", cfg.IdBoletoEmpresa), Type: "string"},
		},
		Item: []c.PostmanItem{
			auth,
			create,
			remove,
		},
	}

	return json.MarshalIndent(collection, "", "  ")
}

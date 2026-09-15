package pix

import (
	"encoding/json"
	"errors"
	"fmt"

	c "io.github.com/conv-pfx/internal/collections"
	b "io.github.com/conv-pfx/internal/collections/body"
)

const (
	sicoobHostAuth    = "auth.sicoob.com.br"
	sicoobHostApi     = "api.sicoob.com.br"
	sicoobBaseUrlAuth = "https://" + sicoobHostAuth
	sicoobBaseUrlApi  = "https://" + sicoobHostApi
	sicoobTXID        = "JHAUDHFJANXUGVFUABND"
	sicoobScope       = "cobv.write cobv.read pix.read"
)

type SicoobPixConfig struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ChavePix     string `json:"chave_pix"`
	Scope        string `json:"scope"`
}

func AuthPixSicoob() c.PostmanItem {

	return c.PostmanItem{
		Name: "Autenticação Gerar Bearer Token",
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
					{Key: "grant_type", Value: "client_credentials"},
					{Key: "client_id", Value: "{{client_id}}", Type: "string"},
					{Key: "scope", Value: "{{scope}}", Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  sicoobBaseUrlAuth + "/auth/realms/cooperado/protocol/openid-connect/token",
				Host: []string{sicoobHostAuth},
				Path: []string{"auth", "realms", "cooperado", "protocol", "openid-connect", "token"},
			},
		},
	}
}

func CreatePixSicoob() (c.PostmanItem, error) {

	createBody := b.SicoobCriarPix{
		Calendario: b.SicoobCalendario{
			DtVencimento: "2027-02-02",
			Validade:     1,
		},
		Chave: b.SicoobChavePix{
			Chave: "{{chave_pix}}",
		},
		Devedor: b.SicoobDevedor{
			Cpf:  "96050176876",
			Nome: "MK - Teste Comunicação-Call",
		},
		Valor: b.SicoobValor{
			Valor: "0.25",
		},
	}

	cBody, err := json.MarshalIndent(createBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar o body da requisição!")
	}

	return c.PostmanItem{
		Name: "Criação de Cobrança via Pix",
		Request: &c.PostmanRequest{
			Method: "PUT",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(cBody),
			},
			Url: c.PostmanURL{
				Raw:  sicoobBaseUrlApi + "/pix/api/v2/cobv/{{pix_txid}}",
				Host: []string{sicoobHostApi},
				Path: []string{"pix", "api", "v2", "cobv", "{{pix_txid}}"},
			},
		},
	}, nil

}

func RemoverPixSicoob() (c.PostmanItem, error) {

	rm := b.SicoobStatusPix{
		Status: "REMOVIDA_PELO_USUARIO_RECEBEDOR",
	}

	remove, err := json.MarshalIndent(rm, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar o body da requisição!")
	}

	return c.PostmanItem{
		Name: "Remover Cobrança via Pix",
		Request: &c.PostmanRequest{
			Method: "PATCH",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(remove),
			},
			Url: c.PostmanURL{
				Raw:  sicoobBaseUrlApi + "/pix/api/v2/cobv/{{pix_txid}}",
				Host: []string{sicoobHostApi},
				Path: []string{"pix", "api", "v2", "cobv", "{{pix_txid}}"},
			},
		},
	}, nil

}

func SicoobPixCollection(cfg *SicoobPixConfig) ([]byte, error) {

	if cfg == nil {
		return nil, fmt.Errorf("Configuração não pode ser nula!")
	}

	scope := cfg.Scope
	if scope == "" {
		scope = sicoobScope
	}

	auth := AuthPixSicoob()

	create, err := CreatePixSicoob()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Criar Pix': %w", err)
	}

	remove, err := RemoverPixSicoob()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Remover Pix': %w", err)
	}

	collection := c.PostmanCollection{
		Info: c.PostmanInfo{
			Name:        "Sicoob - Pix Avulso",
			Description: "Criação de Pix QR Code dinâmico",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []c.PostmanVariable{
			{Key: "bearer_token", Value: "", Type: "string"},
			{Key: "client_id", Value: cfg.ClientId, Type: "string"},
			{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
			{Key: "pix_txid", Value: sicoobTXID, Type: "string"},
			{Key: "chave_pix", Value: cfg.ChavePix, Type: "string"},
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

package pix

import (
	"encoding/json"
	"errors"
	"fmt"

	c "io.github.com/conv-pfx/internal/collections"
	b "io.github.com/conv-pfx/internal/collections/body"
)

const (
	sicrediPixHost = "api-parceiro.sicredi.com.br"
	sicrediPixBase = "https://" + sicrediPixHost
	sicrediTXID    = "HAUSKJLDSIJUKDHSK"
	sicrediScope   = "cobv.write cobv.read pix.read"
)

type SicrediPixConfig struct {
	SicrediApiKey      string
	SicrediPixUsername string
	SicrediPiXPassword string
	SicrediChavePix    string
	SicrediPixCoop     string
	SicrediPixPosto    string
	Scope              string
}

func AuthSicrediPix() c.PostmanItem {

	return c.PostmanItem{
		Name: "Autenticar Pix Sicredi",
		Request: &c.PostmanRequest{
			Auth: &c.PostmanAuth{
				Type: "basic",
				Basic: []c.PostmanAuthParam{
					{Key: "username", Value: "{{username}}", Type: "string"},
					{Key: "password", Value: "{{password}}", Type: "string"},
				},
			},
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/x-www-form-urlencoded", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "urlencoded",
				Urlencoded: []c.PostmanParam{
					{Key: "grant_type", Value: "client_credentials", Type: "string"},
					{Key: "client_id", Value: "{{username}}", Type: "string"},
					{Key: "scope", Value: "{{scope}}", Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  sicrediPixBase + "/oauth/token",
				Host: []string{sicrediPixHost},
				Path: []string{"oauth", "token"},
			},
		},
	}
}

func CreatePixSicredi() (c.PostmanItem, error) {

	createBody := b.SicrediCriarPix{
		Calendario: b.SicrediCalendario{
			DtVencimento:          "2027-02-02",
			ValidadePosVencimento: 1,
		},
		ChavePix: b.SicrediChavePix{
			Chave: "{{chave_pix}}",
		},
		Devedor: b.SicrediDevedor{
			Cpf:  "96050176876",
			Nome: "MK - Teste Comunicação-Call",
		},
		Valor: b.SicrediValor{
			Valor: "0.25",
		},
	}

	body, err := json.MarshalIndent(createBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar body da requisição")
	}

	return c.PostmanItem{
		Name: "Criar Pix Sicredi",
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
				Raw:  sicrediPixBase + "/api/v2/cobv/{{pix_txid}}",
				Host: []string{sicrediPixHost},
				Path: []string{"api", "v2", "cobv", "{{pix_txid}}"},
			},
		},
	}, nil
}

func RemoverPixSicredi() (c.PostmanItem, error) {

	removeBody := b.SicrediStatusPix{
		Status: "REMOVIDA_PELO_USUARIO_RECEBEDOR",
	}

	rm, err := json.MarshalIndent(removeBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar body da requisição")
	}

	return c.PostmanItem{
		Name: "Remover Cobrança Pix",
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
				Raw:  sicrediPixBase + "/api/v2/cobv/{{pix_txid}}",
				Host: []string{sicrediPixHost},
				Path: []string{"api", "v2", "cobv", "{{pix_txid}}"},
			},
		},
	}, nil

}

func SicrediPixCollection(cfg *SicrediPixConfig) ([]byte, error) {

	if cfg == nil {
		return nil, fmt.Errorf("Configuração não pode ser nula")
	}

	scope := cfg.Scope
	if scope == "" {
		scope = sicrediScope
	}

	auth := AuthSicrediPix()

	create, err := CreatePixSicredi()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Criar Pix': %w", err)
	}

	remove, err := RemoverPixSicredi()
	if err != nil {
		return nil, fmt.Errorf("Falha ao criar item 'Remover Pix': %w", err)
	}

	collection := c.PostmanCollection{
		Info: c.PostmanInfo{
			Name:        "Pix Sicredi",
			Description: "Collection Completa Sicredi",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []c.PostmanVariable{
			{Key: "bearer_token", Value: "", Type: "string"},
			{Key: "sicredi_api_key", Value: cfg.SicrediApiKey, Type: "string"},
			{Key: "username", Value: cfg.SicrediPixUsername, Type: "string"},
			{Key: "password", Value: cfg.SicrediPiXPassword, Type: "string"},
			{Key: "cooperativa", Value: cfg.SicrediPixCoop, Type: "string"},
			{Key: "posto", Value: cfg.SicrediPixPosto, Type: "string"},
			{Key: "chave_pix", Value: cfg.SicrediChavePix, Type: "string"},
			{Key: "pix_txid", Value: sicrediTXID, Type: "string"},
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

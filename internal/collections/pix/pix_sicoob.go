package pix

import (
	"encoding/json"
	"errors"

	c "io.github.com/conv-pfx/internal/collections"
	b "io.github.com/conv-pfx/internal/collections/body"
)

const (
	sicoobHostAuth    = "auth.sicoob.com.br"
	sicoobHostApi     = "api.sicoob.com.br"
	sicoobBaseUrlAuth = "https://" + sicoobHostAuth
	sicoobBaseUrlApi  = "https://" + sicoobHostApi
	sicoobTXID        = "JHAUDHFJANXUGVFUABND"
)

type SicoobPixConfig struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ChavePix     string `json:"chave_pix"`
	Scope        string `json:"scope"`
}

func AuthPixSicoob(cfg SicoobPixConfig) c.PostmanItem {

	cfg.Scope = "cobv.write cobv.read pix.read"

	return c.PostmanItem{
		Name: "Autenticação Gerar Bearer Token",
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
				{Key: "content-type", Value: "application/x-www-form-urlencoded"},
			},
			Body: &c.PostmanBody{
				Mode: "urlencoded",
				Urlencoded: []c.PostmanParam{
					{Key: "grant_type", Value: "client_credentials"},
					{Key: "client_id", Value: cfg.ClientId, Type: "string"},
					{Key: "scope", Value: cfg.Scope},
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

func CreatePixSicoob(cfg SicoobPixConfig) (c.PostmanItem, error) {

	createBody := b.SicoobCriarPix{
		Calendario: b.SicoobCalendario{
			DtVencimento: "2027-02-02",
			Validade:     1,
		},
		Chave: b.SicoobChavePix{
			Chave: cfg.ChavePix,
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
				{Key: "content-type", Value: "application/json", Type: "string"},
				{Key: "authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(cBody),
			},
			Url: c.PostmanURL{
				Raw:  sicoobBaseUrlApi + "/pix/api/v2/cobv/" + sicoobTXID,
				Host: []string{sicoobHostApi},
				Path: []string{"pix", "api", "v2", "cobv", sicoobTXID},
			},
		},
	}, nil

}

func RemoverPixSicoob(cfg SicoobPixConfig) (c.PostmanItem, error) {

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
				{Key: "content-type", Value: "application/json", Type: "string"},
				{Key: "authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(remove),
			},
			Url: c.PostmanURL{
				Raw:  sicoobBaseUrlApi + "/pix/api/v2/cobv/" + sicoobTXID,
				Host: []string{sicoobHostApi},
				Path: []string{"pix", "api", "v2", "cobv", sicoobTXID},
			},
		},
	}, nil

}

func SicoobPixCollection(cfg SicoobPixConfig) ([]byte, error) {

	auth := AuthPixSicoob(cfg)

	create, err := CreatePixSicoob(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	remove, err := RemoverPixSicoob(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	txid := sicoobTXID

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
			{Key: "pix_txid", Value: txid, Type: "string"},
			{Key: "chave_pix_conta", Value: cfg.ChavePix, Type: "string"},
		},
		Item: []c.PostmanItem{
			auth,
			create,
			remove,
		},
	}

	return json.MarshalIndent(collection, "", "  ")
}

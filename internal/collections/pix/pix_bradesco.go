package pix

import (
	"encoding/json"
	"errors"

	c "io.github.com/conv-pfx/internal/collections"
	b "io.github.com/conv-pfx/internal/collections/body"
)

const (
	bradescoPixHost = "qrpix.bradesco.com.br"
	bradescoPixBase = "https://" + bradescoPixHost
	bradescoTXID    = "HUAJDKJSIOHHUAHSUGD"
)

type BradescoPixConfig struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ChavePix     string `json:"chave_pix"`
	Scope        string `json:"scope"`
}

func AuthBradescoPix(cfg BradescoPixConfig) c.PostmanItem {

	cfg.Scope = "pix.read pix.write cobv.read cobv.write webhook.write webhook.read"

	return c.PostmanItem{
		Name: "Auth Pix Bradesco",
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
				{Key: "Content-Type", Value: "application/x-www-form-urlencoded"},
			},
			Body: &c.PostmanBody{
				Mode: "urlencoded",
				Urlencoded: []c.PostmanParam{
					{Key: "grant_type", Value: "client_credentials", Type: "string"},
					{Key: "client_id", Value: cfg.ClientId, Type: "string"},
					{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
					{Key: "scope", Value: cfg.Scope, Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  bradescoPixBase + "/oauth/token",
				Host: []string{bradescoPixHost},
				Path: []string{"ouath", "token"},
			},
		},
	}

}

func CreatePixBradesco(cfg BradescoPixConfig) (c.PostmanItem, error) {

	createBody := b.BradescoCriarPix{
		Calendario: b.BradescoCalendario{
			DtVencimento:          "2027-02-02",
			ValidadePosVencimento: 1,
		},
		Chave: b.BradescoChavePix{
			ChavePix: cfg.ChavePix,
		},
		Devedor: b.BradescoPixDevedor{
			Cpf:  "96050176876",
			Nome: "MK - Teste Comunicação-Call",
		},
		Valor: b.BradescoPixValor{
			Valor: "0.25",
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
				{Key: "authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(body),
			},
			Url: c.PostmanURL{
				Raw:  bradescoPixBase + "/v2/cobv/" + bradescoTXID,
				Host: []string{bradescoPixHost},
				Path: []string{"v2", "cobv", bradescoTXID},
			},
		},
	}, nil

}

func RemoverPixBradesco(cfg BradescoPixConfig) (c.PostmanItem, error) {

	removeBody := b.BradescoStatusPix{
		Status: "REMOVIDA_PELO_USUARIO_RECEBEDOR",
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
				{Key: "authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(rm),
			},
			Url: c.PostmanURL{
				Raw:  bradescoPixBase + "/v2/cobv/" + bradescoTXID,
				Host: []string{bradescoPixHost},
				Path: []string{"v2", "cobv", bradescoTXID},
			},
		},
	}, nil
}

func BradescoPixCollection(cfg BradescoPixConfig) ([]byte, error) {

	auth := AuthBradescoPix(cfg)

	create, err := CreatePixBradesco(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	remove, err := RemoverPixBradesco(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	txid := bradescoTXID

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
			{Key: "pix_txid", Value: txid, Type: "string"},
		},
		Item: []c.PostmanItem{
			auth,
			create,
			remove,
		},
	}

	return json.MarshalIndent(collection, "", "  ")
}

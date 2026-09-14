package van

import (
	"encoding/json"
	"errors"

	c "io.github.com/conv-pfx/internal/collections"
	b "io.github.com/conv-pfx/internal/collections/body"
)

const (
	bbVanAuthHost = "oauth.bb.com.br"
	bbVanAuthBase = "https://" + bbVanAuthHost
	bbVanHost     = "api.bb.com.br"
	bbVanBase     = "https://" + bbVanHost
)

type BancoDoBrasilVanConfig struct {
	BbBasicToken   string
	BbAppKey       string
	BbConvenio     string
	BbNumCarteira  string
	BbVariacaoCart string
	BbModalidade   string
	BbTituloBenef  string
	BbTituloClient string
	Scope          string
}

func AuthVanBB(cfg BancoDoBrasilVanConfig) c.PostmanItem {

	cfg.Scope = "cobrancas.boletos-requisicao cobrancas.boletos-info"

	return c.PostmanItem{
		Name: "Autenticar Van BB",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/x-www-form-urlencoded", Type: "string"},
				{Key: "Authorization", Value: cfg.BbBasicToken, Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "urlencoded",
				Urlencoded: []c.PostmanParam{
					{Key: "grant_type", Value: "client_credentials", Type: "string"},
					{Key: "scope", Value: cfg.Scope, Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  bbVanAuthBase + "/oauth/token",
				Host: []string{bbVanAuthHost},
				Path: []string{"oauth", "token"},
			},
		},
	}
}

func CreateBBVan(cfg BancoDoBrasilVanConfig) (c.PostmanItem, error) {

	createBody := b.DadosBoletoBB{
		NumConvenio:            cfg.BbConvenio,
		NumCarteira:            cfg.BbNumCarteira,
		NumVariacaoCart:        cfg.BbVariacaoCart,
		CodigoModalidade:       cfg.BbModalidade,
		DtEmissao:              "2026-09-15",
		DtVencimento:           "2027-02-02",
		VlOriginal:             "10",
		IndicadorTitVencido:    "S",
		NumDiasLimiteReceb:     0,
		CdTipoTitulo:           2,
		IndPermissaoRecParcial: "N",
		NumTituloBenef:         cfg.BbTituloBenef,
		NumTituloClient:        cfg.BbTituloClient,
		Desconto: b.BoletoDescontoBB{
			Tipo:        1,
			DtExpiracao: "2027-02-02",
			Valor:       "20",
		},
		JurosMora: b.BoletoJurosBB{
			Tipo:        2,
			Porcentagem: "0.99",
		},
		Multa: b.BoletoMultaBB{
			Tipo:         2,
			DtVencimento: "2027-02-02",
			Porcentagem:  "2",
		},
		Pagador: b.BoletoPagadorBB{
			Tipo:         1,
			NumInscricao: 96050176876,
		},
		IndicadorPix: "S",
	}

	body, err := json.MarshalIndent(createBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar body da requisição!")
	}

	return c.PostmanItem{
		Name: "Criar boleto BB",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(body),
			},
			Url: c.PostmanURL{
				Raw:  bbVanBase + "/cobrancas/v2/boletos?gw-app-key=" + cfg.BbAppKey,
				Host: []string{bbVanHost},
				Path: []string{"cobrancas", "v2", "boletos?gw-app-key=", cfg.BbAppKey},
			},
		},
	}, nil
}

func RemoveBBVan(cfg BancoDoBrasilVanConfig) (c.PostmanItem, error) {

	removeBody := b.BaixaBoletoBB{
		NumeroConv: cfg.BbConvenio,
	}

	rm, err := json.MarshalIndent(removeBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar body da requisição!")
	}

	return c.PostmanItem{
		Name: "Baixar Boleto BB",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(rm),
			},
			Url: c.PostmanURL{
				Raw:  bbVanBase + "/cobrancas/v2/boletos/numTitulo/baixar?gw-app-key=" + cfg.BbAppKey,
				Host: []string{bbVanHost},
				Path: []string{"cobrancas", "v2", "boletos", "numTitulo", "baixar?gw-app-key=", cfg.BbAppKey},
			},
		},
	}, nil
}

func BancoDoBrasilVanCollection(cfg BancoDoBrasilVanConfig) ([]byte, error) {

	auth := AuthVanBB(cfg)

	create, err := CreateBBVan(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	remove, err := RemoveBBVan(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado")
	}

	collection := c.PostmanCollection{
		Info: c.PostmanInfo{
			Name:        "BB",
			Description: "Collection Completa BB",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []c.PostmanVariable{
			{Key: "bearer_token", Value: "", Type: "string"},
			{Key: "bb_basic_token", Value: cfg.BbBasicToken, Type: "string"},
			{Key: "bb_app_key", Value: cfg.BbAppKey, Type: "string"},
			{Key: "bb_convenio", Value: cfg.BbConvenio, Type: "string"},
			{Key: "bb_num_carteira", Value: cfg.BbNumCarteira, Type: "string"},
			{Key: "bb_variacao_carteira", Value: cfg.BbVariacaoCart, Type: "string"},
			{Key: "bb_modalidade", Value: cfg.BbModalidade, Type: "string"},
			{Key: "bb_titulo_beneficiario", Value: cfg.BbTituloBenef, Type: "string"},
			{Key: "bb_titulo_cliente", Value: cfg.BbTituloClient, Type: "string"},
			{Key: "scope", Value: cfg.Scope, Type: "string"},
		},
		Item: []c.PostmanItem{
			auth,
			create,
			remove,
		},
	}
	return json.MarshalIndent(collection, "", "  ")
}

package van

import (
	"encoding/json"
	"errors"
	"strconv"

	c "io.github.com/conv-pfx/internal/collections"
	b "io.github.com/conv-pfx/internal/collections/body"
)

const (
	sicoobHostAuthVan = "auth.sicoob.com.br"
	sicoobHostApiVan  = "api.sicoob.com.br"
	sicoobBaseUrlAuth = "https://" + sicoobHostAuthVan
	sicoobBaseUrlApi  = "https://" + sicoobHostApiVan
)

type SicoobVanConfig struct {
	ClientId     string              `json:"client_id"`
	ClientSecret string              `json:"client_secret"`
	Scope        string              `json:"scope"`
	Dados        b.SicoobDadosBoleto `json:"dados_sicoob"`
}

func AuthVanSicoob(cfg SicoobVanConfig) c.PostmanItem {

	cfg.Scope = "boletos_inclusao boletos_consulta boletos_alteracao webhooks_alteracao webhooks_consulta webhooks_inclusao"

	return c.PostmanItem{
		Name: "Autenticação Van Gerar Bearer Token",
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
				{Key: "content-type", Value: "application/x-www-urlencoded"},
			},
			Body: &c.PostmanBody{
				Mode: "urlencoded",
				Urlencoded: []c.PostmanParam{
					{Key: "grant_type", Value: "client_credentials", Type: "string"},
					{Key: "client_id", Value: cfg.ClientId, Type: "string"},
					{Key: "scope", Value: cfg.Scope, Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  sicoobBaseUrlAuth + "/auth/realms/cooperado/protocol/openid-connect/token",
				Host: []string{sicoobHostAuthVan},
				Path: []string{"auth", "realms", "cooperado", "protocol", "openid-connect", "token"},
			},
		},
	}

}

func CreateSicoobVan(cfg SicoobVanConfig) (c.PostmanItem, error) {

	createBody := b.SicoobCriarBoleto{
		Dados: b.SicoobDadosBoleto{
			NumeroCliente:          cfg.Dados.NumeroCliente,
			CodigoModalidade:       cfg.Dados.CodigoModalidade,
			ContaCorrente:          cfg.Dados.ContaCorrente,
			DataVencimento:         cfg.Dados.DataVencimento,
			DataEmissao:            cfg.Dados.DataEmissao,
			SeuNumero:              cfg.Dados.SeuNumero,
			NossoNumero:            cfg.Dados.NossoNumero,
			IdBoletoEmpresa:        cfg.Dados.IdBoletoEmpresa,
			CodigoEspecieDocumento: "0",
			IdEmissaoBoleto:        0,
			IdDistribuicaoBoleto:   0,
			Valor:                  25.0,
		},
		Pagador: b.SicoobPagadorBoleto{
			Documento: "1452875",
			Nome:      "MK Solutions",
			Endereco:  "Rua Amazonas",
			Bairro:    "Patrimonio Novo",
			Cidade:    "Votuporanga",
			Cep:       "15500004",
			Uf:        "SP",
			Email:     "",
		},
		Pix: b.SicoobPixBoleto{
			CdCadastroPix: 1,
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
				{Key: "content-type", Value: "application/json", Type: "string"},
				{Key: "authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
				{Key: "client_id", Value: cfg.ClientId, Type: "string"},
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

func RemoveSicoobVan(cfg SicoobVanConfig) (c.PostmanItem, error) {

	nossoNum := strconv.Itoa(cfg.Dados.NossoNumero)

	rm := b.SicoobRemoverBoleto{
		Boleto: b.SicoobDadosBoleto{
			NumeroCliente:    cfg.Dados.NumeroCliente,
			CodigoModalidade: cfg.Dados.CodigoModalidade,
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
				{Key: "content-type", Value: "application/json", Type: "string"},
				{Key: "authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
				{Key: "client_id", Value: cfg.ClientId, Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(rmBody),
			},
			Url: c.PostmanURL{
				Raw:  sicoobBaseUrlApi + "/cobranca-bancaria/v3/boletos/" + nossoNum + "/baixar",
				Host: []string{sicoobHostApiVan},
				Path: []string{"cobranca-bancaria", "v3", "boletos", nossoNum, "baixar"},
			},
		},
	}, nil
}

func SicoobVanCollection(cfg SicoobVanConfig) ([]byte, error) {

	auth := AuthVanSicoob(cfg)

	create, err := CreateSicoobVan(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado!")
	}

	remove, err := RemoveSicoobVan(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado!")
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
			{Key: "numeroCliente", Value: strconv.Itoa(cfg.Dados.NumeroCliente), Type: "string"},
			{Key: "codigoModalidade", Value: strconv.Itoa(cfg.Dados.CodigoModalidade), Type: "string"},
			{Key: "numeroContaCorrente", Value: strconv.Itoa(cfg.Dados.ContaCorrente), Type: "string"},
			{Key: "dataVencimento", Value: cfg.Dados.DataVencimento, Type: "string"},
			{Key: "dataEmissao", Value: cfg.Dados.DataEmissao, Type: "string"},
			{Key: "seuNumero", Value: strconv.Itoa(cfg.Dados.SeuNumero), Type: "string"},
			{Key: "nossoNumero", Value: strconv.Itoa(cfg.Dados.NossoNumero), Type: "string"},
			{Key: "identificacaoBoletoEmpresa", Value: strconv.Itoa(cfg.Dados.IdBoletoEmpresa), Type: "string"},
		},
		Item: []c.PostmanItem{
			auth,
			create,
			remove,
		},
	}

	return json.MarshalIndent(collection, "", "  ")

}

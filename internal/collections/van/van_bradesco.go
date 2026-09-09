package van

import (
	"encoding/json"
	"errors"
	"strconv"

	c "io.github.com/conv-pfx/internal/collections"
	b "io.github.com/conv-pfx/internal/collections/body"
)

const (
	bradescoHost    = "openapi.bradesco.com.br"
	bradescoBaseUrl = "https://" + bradescoHost
)

type BradescoVanConfig struct {
	ClientId     string                `json:"client_id"`
	ClientSecret string                `json:"client_secret"`
	Dados        b.DadosBoletoBradesco `json:"dados_bradesco"`
}

func AuthVanBradesco(cfg BradescoVanConfig) c.PostmanItem {

	return c.PostmanItem{
		Name: "Autenticação Van Gerar Bearer Token ",
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
					{Key: "grant-type", Value: "client_credentials", Type: "string"},
					{Key: "client_id", Value: cfg.ClientId, Type: "string"},
					{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  bradescoBaseUrl + "/auth/server-mtls/v2/token",
				Host: []string{bradescoHost},
				Path: []string{"auth", "server-mtls", "v2", "token"},
			},
		},
	}

}

func CreateBradescoVan(cfg BradescoVanConfig) (c.PostmanItem, error) {

	createBody := b.BoletoBradesco{
		BoletoBradesco: b.DadosBoletoBradesco{
			NuCPFCNPJ:                            cfg.Dados.NuCPFCNPJ,
			FilialCPFCNPJ:                        cfg.Dados.FilialCPFCNPJ,
			CtrlCPFCNPJ:                          cfg.Dados.CtrlCPFCNPJ,
			IdProduto:                            cfg.Dados.IdProduto,
			NuNegociacao:                         cfg.Dados.NuNegociacao,
			NuTitulo:                             cfg.Dados.NuTitulo,
			DtEmissaoTitulo:                      cfg.Dados.DtEmissaoTitulo,
			DtVencimentoTitulo:                   cfg.Dados.DtVencimentoTitulo,
			CdPagamentoParcial:                   "N",
			QtdePagamentoParcial:                 0,
			PercentualJuros:                      1,
			VlJuros:                              0.00,
			QtdeDiasJuros:                        1,
			PercentualMulta:                      2,
			VlMulta:                              0.00,
			QtdeDiasMulta:                        1,
			PercentualDesconto1:                  0,
			VlDesconto1:                          0.00,
			DataLimiteDesconto1:                  "",
			PercentualDesconto2:                  0,
			VlDesconto2:                          0.00,
			DataLimiteDesconto2:                  "",
			PercentualDesconto3:                  0,
			VlDesconto3:                          0.00,
			DataLimiteDesconto3:                  "",
			PrazoBonificacao:                     0,
			PercentualBonificacao:                0,
			VlBonificacao:                        0.00,
			DtLimiteBonificacao:                  "",
			VlAbatimento:                         0.00,
			VlIOF:                                0.00,
			NomePagador:                          "JESSICASANTOS",
			LogradouroPagador:                    "RuaGercyPdaSilva",
			NuLogradouroPagador:                  "0",
			ComplementoLogradouroPagador:         "LT10QD17",
			CepPagador:                           99050,
			ComplementoCepPagador:                100,
			BairroPagador:                        "Petropolis",
			MunicipioPagador:                     "PassoFundo",
			UfPagador:                            "RS",
			CdIndCpfcnpjPagador:                  1,
			NuCpfcnpjPagador:                     89397645013,
			EndEletronicoPagador:                 "",
			NomeSacadorAvalista:                  "",
			LogradouroSacadorAvalista:            "",
			NuLogradouroSacadorAvalista:          "",
			ComplementoLogradouroSacadorAvalista: "",
			CepSacadorAvalista:                   0,
			ComplementoCepSacadorAvalista:        0,
			BairroSacadorAvalista:                "",
			MunicipioSacadorAvalista:             "",
			UfSacadorAvalista:                    "",
			CdIndCpfcnpjSacadorAvalista:          0,
			NuCpfcnpjSacadorAvalista:             0,
			EnderecoSacadorAvalista:              "",
		},
	}

	body, err := json.MarshalIndent(createBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar body de requisição!")
	}

	return c.PostmanItem{
		Name: "Criar Boleto Van",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
				{Key: "content-type", Value: "application/json", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(body),
			},
			Url: c.PostmanURL{
				Raw:  bradescoBaseUrl + "/boleto/cobranca-registro/v1/cobranca",
				Host: []string{bradescoHost},
				Path: []string{"boleto", "cobranca-registro", "v1", "cobranca"},
			},
		},
	}, nil

}

func RemoveBradescoVan(cfg BradescoVanConfig) (c.PostmanItem, error) {

	rm := b.BaixaBoletoBradesco{
		CpfCnpj: b.DocPagadorBradesco{
			CpfCnpj:  cfg.Dados.NuCPFCNPJ,
			Filial:   cfg.Dados.FilialCPFCNPJ,
			Controle: cfg.Dados.CtrlCPFCNPJ,
		},
		Produto:     cfg.Dados.IdProduto,
		Negociacao:  cfg.Dados.NuNegociacao,
		NossoNumero: cfg.Dados.NuTitulo,
		Sequencia:   0,
		CodigoBaixa: 57,
	}

	rmBody, err := json.MarshalIndent(rm, "", "  ")
	if err != nil {
		return c.PostmanItem{}, errors.New("Erro ao serializar o body da requisição!")
	}

	return c.PostmanItem{
		Name: "Baixar Boleto Bradesco",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
				{Key: "content-type", Value: "application/json", Type: "string"},
			},
			Body: &c.PostmanBody{
				Mode: "raw",
				Raw:  string(rmBody),
			},
			Url: c.PostmanURL{
				Raw:  bradescoBaseUrl + "/boleto/cobranca-baixa/v1/baixar",
				Host: []string{bradescoHost},
				Path: []string{"boleto", "cobranca-baixa", "v1", "baixar"},
			},
		},
	}, nil
}

func BradescoVanCollection(cfg BradescoVanConfig) ([]byte, error) {

	auth := AuthVanBradesco(cfg)

	create, err := CreateBradescoVan(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado!")
	}

	remove, err := RemoveBradescoVan(cfg)
	if err != nil {
		return nil, errors.New("Erro inesperado!")
	}

	collection := c.PostmanCollection{
		Info: c.PostmanInfo{
			Name:        "Bradesco - Criação de Boletos com QR Code dinâmico",
			Description: "Criar boletos com QR Code Dinâmico por API",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []c.PostmanVariable{
			{Key: "bearer_token", Value: "Bearer {{bearer_token}}", Type: "string"},
			{Key: "client_id", Value: cfg.ClientId, Type: "string"},
			{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
			{Key: "nuCPFCNPJ", Value: strconv.Itoa(cfg.Dados.NuCPFCNPJ), Type: "string"},
			{Key: "filialCPFCNPJ", Value: strconv.Itoa(cfg.Dados.FilialCPFCNPJ), Type: "string"},
			{Key: "ctrlCPFCNPJ", Value: strconv.Itoa(cfg.Dados.CtrlCPFCNPJ), Type: "string"},
			{Key: "idProduto", Value: strconv.Itoa(cfg.Dados.IdProduto), Type: "string"},
			{Key: "nuNegociacao", Value: strconv.Itoa(int(cfg.Dados.NuNegociacao)), Type: "string"},
			{Key: "nuTitulo", Value: strconv.Itoa(cfg.Dados.NuTitulo), Type: "string"},
			{Key: "nuCliente", Value: strconv.Itoa(cfg.Dados.NuTitulo), Type: "string"},
			{Key: "nuTitulo", Value: cfg.Dados.NuCliente, Type: "string"},
			{Key: "dtEmissaoTitulo", Value: cfg.Dados.DtEmissaoTitulo, Type: "string"},
			{Key: "dtVencimentoTitulo", Value: cfg.Dados.DtVencimentoTitulo, Type: "string"},
		},
		Item: []c.PostmanItem{
			auth,
			create,
			remove,
		},
	}
	return json.MarshalIndent(collection, "", "  ")
}

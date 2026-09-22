package van

import (
	"encoding/json"
	"fmt"

	c "io.github.com/conv-pfx/internal/collections"
)

const (
	bradescoHost    = "openapi.bradesco.com.br"
	bradescoBaseUrl = "https://" + bradescoHost
)

type BradescoVanConfig struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func AuthVanBradesco(cfg BradescoVanConfig) (c.PostmanItem, error) {
	return c.PostmanItem{
		Name:  "Autenticação Van Gerar Bearer Token ",
		Event: c.AutoSetBearerTokenEvent(),
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/x-www-form-urlencoded"},
			},
			Body: &c.PostmanBody{
				Mode: "urlencoded",
				Urlencoded: []c.PostmanParam{
					{Key: "grant-type", Value: "client_credentials", Type: "string"},
					{Key: "client_id", Value: "{{client_id}}", Type: "string"},
					{Key: "client_secret", Value: "{{client_secret}}", Type: "string"},
				},
			},
			Url: c.PostmanURL{
				Raw:  bradescoBaseUrl + "/auth/server-mtls/v2/token",
				Host: []string{bradescoHost},
				Path: []string{"auth", "server-mtls", "v2", "token"},
			},
		},
	}, nil
}

func CreateBradescoVan(cfg BradescoVanConfig) (c.PostmanItem, error) {
	createBody := map[string]any{
		"boleto_bradesco": map[string]any{
			"nuCPFCNPJ":                            "",
			"filialCPFCNPJ":                        "",
			"ctrlCPFCNPJ":                          "",
			"idProduto":                            "",
			"nuNegociacao":                         "",
			"nuTitulo":                             "",
			"dtEmissaoTitulo":                      "",
			"dtVencimentoTitulo":                   "",
			"cdPagamentoParcial":                   "N",
			"qtdePagamentoParcial":                 0,
			"percentualJuros":                      1,
			"vlJuros":                              0.00,
			"qtdeDiasJuros":                        1,
			"percentualMulta":                      2,
			"vlMulta":                              0.00,
			"qtdeDiasMulta":                        1,
			"percentualDesconto1":                  0,
			"vlDesconto1":                          0.00,
			"dataLimiteDesconto1":                  "",
			"percentualDesconto2":                  0,
			"vlDesconto2":                          0.00,
			"dataLimiteDesconto2":                  "",
			"percentualDesconto3":                  0,
			"vlDesconto3":                          0.00,
			"dataLimiteDesconto3":                  "",
			"prazoBonificacao":                     0,
			"percentualBonificacao":                0,
			"vlBonificacao":                        0.00,
			"dtLimiteBonificacao":                  "",
			"vlAbatimento":                         0.00,
			"vlIOF":                                0.00,
			"nomePagador":                          "JESSICASANTOS",
			"logradouroPagador":                    "RuaGercyPdaSilva",
			"nuLogradouroPagador":                  "0",
			"complementoLogradouroPagador":         "LT10QD17",
			"cepPagador":                           99050,
			"complementoCepPagador":                100,
			"bairroPagador":                        "Petropolis",
			"municipioPagador":                     "PassoFundo",
			"ufPagador":                            "RS",
			"cdIndCpfcnpjPagador":                  1,
			"nuCpfcnpjPagador":                     89397645013,
			"endEletronicoPagador":                 "",
			"nomeSacadorAvalista":                  "",
			"logradouroSacadorAvalista":            "",
			"nuLogradouroSacadorAvalista":          "",
			"complementoLogradouroSacadorAvalista": "",
			"cepSacadorAvalista":                   0,
			"complementoCepSacadorAvalista":        0,
			"bairroSacadorAvalista":                "",
			"municipioSacadorAvalista":             "",
			"ufSacadorAvalista":                    "",
			"cdIndCpfcnpjSacadorAvalista":          0,
			"nuCpfcnpjSacadorAvalista":             0,
			"enderecoSacadorAvalista":              "",
		},
	}

	body, err := json.MarshalIndent(createBody, "", "  ")
	if err != nil {
		return c.PostmanItem{}, fmt.Errorf("erro ao serializar body de requisição: %w", err)
	}

	return c.PostmanItem{
		Name: "Criar Boleto Van",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
			},
			Body: c.RawJSONBody(string(body)),
			Url: c.PostmanURL{
				Raw:  bradescoBaseUrl + "/boleto/cobranca-registro/v1/cobranca",
				Host: []string{bradescoHost},
				Path: []string{"boleto", "cobranca-registro", "v1", "cobranca"},
			},
		},
	}, nil
}

func RemoveBradescoVan(cfg BradescoVanConfig) (c.PostmanItem, error) {
	rm := map[string]any{
		"cpfCnpj": map[string]any{
			"cpfCnpj":  "",
			"filial":   "",
			"controle": "",
		},
		"produto":     "",
		"negociacao":  "",
		"nossoNumero": "",
		"sequencia":   0,
		"codigoBaixa": 57,
	}

	rmBody, err := json.MarshalIndent(rm, "", "  ")
	if err != nil {
		return c.PostmanItem{}, fmt.Errorf("erro ao serializar o body da requisição: %w", err)
	}

	return c.PostmanItem{
		Name: "Baixar Boleto Bradesco",
		Request: &c.PostmanRequest{
			Method: "POST",
			Header: []c.PostmanHeader{
				{Key: "Content-Type", Value: "application/json", Type: "string"},
				{Key: "Authorization", Value: "Bearer {{bearer_token}}", Type: "string"},
			},
			Body: c.RawJSONBody(string(rmBody)),
			Url: c.PostmanURL{
				Raw:  bradescoBaseUrl + "/boleto/cobranca-baixa/v1/baixar",
				Host: []string{bradescoHost},
				Path: []string{"boleto", "cobranca-baixa", "v1", "baixar"},
			},
		},
	}, nil
}

func BradescoVanCollection(cfg *BradescoVanConfig) ([]byte, error) {
	auth, err := AuthVanBradesco(*cfg)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar item 'Auth Van Bradesco': %w", err)
	}

	create, err := CreateBradescoVan(*cfg)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar item de criação VAN Bradesco: %w", err)
	}

	remove, err := RemoveBradescoVan(*cfg)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar item de remoção VAN Bradesco: %w", err)
	}

	collection := c.PostmanCollection{
		Info: c.PostmanInfo{
			Name:        "Bradesco - Criação de Boletos com QR Code dinâmico",
			Description: "Criar boletos com QR Code Dinâmico por API",
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Variable: []c.PostmanVariable{
			{Key: "bearer_token", Value: "", Type: "string"},
			{Key: "client_id", Value: cfg.ClientId, Type: "string"},
			{Key: "client_secret", Value: cfg.ClientSecret, Type: "string"},
			{Key: "nuCPFCNPJ", Value: fmt.Sprint(), Type: "string"},
			{Key: "filialCPFCNPJ", Value: fmt.Sprint(), Type: "string"},
			{Key: "ctrlCPFCNPJ", Value: fmt.Sprint(), Type: "string"},
			{Key: "idProduto", Value: fmt.Sprint(), Type: "string"},
			{Key: "nuNegociacao", Value: fmt.Sprint(), Type: "string"},
			{Key: "nuTitulo", Value: fmt.Sprint(), Type: "string"},
			{Key: "nuCliente", Value: fmt.Sprint(), Type: "string"},
			{Key: "dtEmissaoTitulo", Value: fmt.Sprint(), Type: "string"},
			{Key: "dtVencimentoTitulo", Value: fmt.Sprint(), Type: "string"},
		},
		Item: []c.PostmanItem{
			auth,
			create,
			remove,
		},
	}

	return json.MarshalIndent(collection, "", "  ")
}

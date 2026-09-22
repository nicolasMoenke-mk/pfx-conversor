package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"io.github.com/conv-pfx/internal/assets"
	p "io.github.com/conv-pfx/internal/collections/pix"
	v "io.github.com/conv-pfx/internal/collections/van"
	"io.github.com/conv-pfx/internal/convert"
	"io.github.com/conv-pfx/utils"
)

type Page struct {
	Success   bool   `json:"sucesso"`
	Error     string `json:"erro,omitempty"`
	NomeCert  string `json:"nomeCert,omitempty"`
	CertPEM   string `json:"certPem,omitempty"`
	NomeChave string `json:"nomeChave,omitempty"`
	KeyPEM    string `json:"keyPem,omitempty"`
}

type FaturaResponse struct {
	Success bool     `json:"sucesso"`
	Error   string   `json:"erro,omitempty"`
	Valores []string `json:"valores,omitempty"`
}

type RemessaResponse struct {
	Success         bool           `json:"sucesso"`
	Error           string         `json:"erro,omitempty"`
	TipoCnab        utils.TipoCnab `json:"tipoCnab,omitempty"`
	TotalLinhas     int            `json:"totalLinhas,omitempty"`
	QuantidadeLotes int            `json:"quantidadeLotes,omitempty"`
	Erros           []utils.Erro   `json:"erros,omitempty"`
	Valido          bool           `json:"valido,omitempty"`
}

func formatarNomeArquivo(nomeInput, prefixo, extensao string) string {

	nomeLimpo := strings.TrimSpace(nomeInput)

	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".pem")
	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".crt")
	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".key")
	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".pfx")

	if nomeLimpo == "" {
		return prefixo + extensao
	}

	return fmt.Sprintf("%s%s%s", prefixo, nomeLimpo, extensao)
}

func abrirNavegador(url string) {
	time.Sleep(300 * time.Millisecond)
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "chrome", url)
	case "darwin":
		cmd = exec.Command("open", "-a", "Google Chrome", url)
	default:
		cmd = exec.Command("google-chrome", url)
	}

	_ = cmd.Start()

}

func responseJSON(w http.ResponseWriter, res any) {
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

// atoiCampo converte um campo numérico do formulário. Campo vazio é
// tratado como zero (válido); um valor preenchido que não seja numérico
// retorna erro em vez de virar 0 silenciosamente.
func atoiCampo(values map[string]string, campo string) (int, error) {
	raw := strings.TrimSpace(values[campo])
	if raw == "" {
		return 0, nil
	}

	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("o campo %q precisa ser um número válido (valor recebido: %q)", campo, raw)
	}

	return n, nil
}

// createCollections monta a collection Postman de cada banco/produto a partir
// dos valores preenchidos pelo usuário no front-end. Chamar com values == nil
// gera a collection "vazia", usada apenas para descobrir os campos existentes.
var createCollections = map[string]func(values map[string]string) ([]byte, error){
	"sicoob:van": func(values map[string]string) ([]byte, error) {
		numeroCliente, err := atoiCampo(values, "numeroCliente")
		if err != nil {
			return nil, err
		}

		codigoModalidade, err := atoiCampo(values, "codigoModalidade")
		if err != nil {
			return nil, err
		}

		numeroContaCorrente, err := atoiCampo(values, "numeroContaCorrente")
		if err != nil {
			return nil, err
		}

		seuNumero, err := atoiCampo(values, "seuNumero")
		if err != nil {
			return nil, err
		}

		nossoNumero, err := atoiCampo(values, "nossoNumero")
		if err != nil {
			return nil, err
		}

		idBoletoEmpresa, err := atoiCampo(values, "identificacaoBoletoEmpresa")
		if err != nil {
			return nil, err
		}

		return v.SicoobVanCollection(&v.SicoobVanConfig{
			ClientId:         values["client_id"],
			ClientSecret:     values["client_secret"],
			Scope:            values["scope"],
			NumeroCliente:    numeroCliente,
			CodigoModalidade: codigoModalidade,
			ContaCorrente:    numeroContaCorrente,
			DataVencimento:   values["dataVencimento"],
			DataEmissao:      values["dataEmissao"],
			SeuNumero:        seuNumero,
			NossoNumero:      nossoNumero,
			IdBoletoEmpresa:  idBoletoEmpresa,
		})
	},
	"sicoob:pix": func(values map[string]string) ([]byte, error) {
		return p.SicoobPixCollection(&p.SicoobPixConfig{
			ClientId:     values["client_id"],
			ClientSecret: values["client_secret"],
			ChavePix:     values["chave_pix"],
			Scope:        values["scope"],
		})
	},
	"santander:van": func(values map[string]string) ([]byte, error) {
		return v.SantanderVanCollection(&v.SantanderVanConfig{
			ClientId:              values["client_id"],
			ClientSecret:          values["client_secret"],
			SantanderAppKey:       values["santander_app_key"],
			SantanderWkId:         values["santander_workspace_id"],
			SantanderClientNumber: values["santander_client_number"],
			SantanderBankNumb:     values["santander_bank_number"],
			SantanderConventNumb:  values["santander_convent_code"],
			SantanderType:         values["santander_dict_type"],
			SantanderDictKey:      values["santander_dict_key"],
		})
	},
	"santander:pix": func(values map[string]string) ([]byte, error) {
		return p.SantanderPixCollection(&p.SantanderPixConfig{
			ClientId:     values["client_id"],
			ClientSecret: values["client_secret"],
			ChavePix:     values["chave_pix"],
		})
	},
	"bradesco:van": func(values map[string]string) ([]byte, error) {
		return v.BradescoVanCollection(&v.BradescoVanConfig{
			ClientId:     values["client_id"],
			ClientSecret: values["client_secret"],
		})
	},
	"bradesco:pix": func(values map[string]string) ([]byte, error) {
		return p.BradescoPixCollection(&p.BradescoPixConfig{
			ClientId:     values["client_id"],
			ClientSecret: values["client_secret"],
			ChavePix:     values["chave_pix"],
			Scope:        values["scope"],
		})
	},
	"sicredi:van": func(values map[string]string) ([]byte, error) {
		return v.SicrediVanCollection(&v.SicrediVanConfig{
			SicrediUsername: values["username"],
			SicrediPassword: values["password"],
			SicrediApiKey:   values["x-api-key"],
			SicrediCoop:     values["cooperativa"],
			SicrediPosto:    values["posto"],
			SicrediCodBenef: values["sicredi_codBenef"],
			SicrediNN:       values["sicredi_nossoNumero"],
			Scope:           values["scope"],
		})
	},
	"sicredi:pix": func(values map[string]string) ([]byte, error) {
		return p.SicrediPixCollection(&p.SicrediPixConfig{
			SicrediApiKey:      values["sicredi_api_key"],
			SicrediPixUsername: values["username"],
			SicrediPiXPassword: values["password"],
			SicrediChavePix:    values["chave_pix"],
			SicrediPixCoop:     values["cooperativa"],
			SicrediPixPosto:    values["posto"],
			Scope:              values["scope"],
		})
	},
	"bb:van": func(values map[string]string) ([]byte, error) {
		return v.BancoDoBrasilVanCollection(v.BancoDoBrasilVanConfig{
			BbBasicToken:   values["bb_basic_token"],
			BbAppKey:       values["bb_app_key"],
			BbConvenio:     values["bb_convenio"],
			BbNumCarteira:  values["bb_num_carteira"],
			BbVariacaoCart: values["bb_variacao_carteira"],
			BbModalidade:   values["bb_modalidade"],
			BbTituloBenef:  values["bb_titulo_beneficiario"],
			BbTituloClient: values["bb_titulo_cliente"],
			Scope:          values["scope"],
		})
	},
	"bb:pix": func(values map[string]string) ([]byte, error) {
		return p.BancoDoBrasilPixCollection(&p.BancoDoBrasilPixConfig{
			ClientId:     values["client_id"],
			ClientSecret: values["client_secret"],
			BBChavePix:   values["chave_pix"],
			BBAppKey:     values["bb_dev_app_key"],
			Scope:        values["scope"],
		})
	},
}

func pathConstruct(r *http.Request) string {
	banco := r.PathValue("banco")
	produto := r.URL.Query().Get("produto")
	if produto == "" {
		return banco
	}

	return banco + ":" + produto
}

type campoCollection struct {
	Chave string `json:"chave"`
}

// extractField lista as chaves das variáveis de uma collection gerada,
// para que o front-end monte o formulário de campos a preencher.
func extractField(data []byte) ([]campoCollection, error) {
	var part struct {
		Variable []struct {
			Key string `json:"key"`
		} `json:"variable"`
	}

	if err := json.Unmarshal(data, &part); err != nil {
		return nil, errors.New("Erro inesperado")
	}

	fields := make([]campoCollection, 0, len(part.Variable))
	for _, v := range part.Variable {
		if v.Key == "bearer_token" {
			continue
		}
		fields = append(fields, campoCollection{Chave: v.Key})
	}
	return fields, nil
}

func handlerFieldsCollection(w http.ResponseWriter, r *http.Request) {
	gen, ok := createCollections[pathConstruct(r)]

	if !ok {
		http.NotFound(w, r)
		return
	}

	data, err := gen(nil)
	if err != nil {
		responseJSON(w, Page{Success: false, Error: err.Error()})
		return
	}

	fields, err := extractField(data)
	if err != nil {
		responseJSON(w, Page{Success: false, Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(fields)
}

func handlerGenCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	key := pathConstruct(r)

	gen, ok := createCollections[key]
	if !ok {
		http.NotFound(w, r)
		return
	}

	var values map[string]string

	err := json.NewDecoder(r.Body).Decode(&values)
	if err != nil {
		responseJSON(w, Page{Success: false, Error: err.Error()})
		return
	}

	data, err := gen(values)
	if err != nil {
		responseJSON(w, Page{Success: false, Error: err.Error()})
		return
	}

	fileName := strings.ReplaceAll(key, ":", "-") + ".json"

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", `attachment; filename="`+fileName+`"`)
	w.Write(data)
}

func main() {
	http.Handle("/", http.FileServer(assets.GetCoreFile()))

	http.HandleFunc("/convert-pem", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			responseJSON(w, Page{Success: false, Error: "O arquivo excede o limite permitido 10MB"})
			return
		}

		senha := r.FormValue("password")

		file, _, err := r.FormFile("pfxfile")
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Falha ao carregar o arquivo pfx"})
			return
		}

		defer file.Close()

		pfxByte, err := io.ReadAll(file)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Falha ao ler o conteúdo do pfx"})
			return
		}

		keyPEM, certPEM, err := convert.ConvertPFX(pfxByte, senha)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: err.Error()})
			return
		}

		nomePub := formatarNomeArquivo(r.FormValue("nomeCert"), "ChavePublica", ".pem")
		nomePriv := formatarNomeArquivo(r.FormValue("nomeChave"), "ChavePrivada", ".pem")

		responseJSON(w, Page{
			Success:   true,
			NomeCert:  nomePub,
			CertPEM:   string(certPEM),
			NomeChave: nomePriv,
			KeyPEM:    string(keyPEM),
		})

	})

	http.HandleFunc("/convert-next", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			responseJSON(w, Page{Success: false, Error: "O arquivo excede o limite permitido 10MB"})
			return
		}

		senha := r.FormValue("password")
		novaSenha := r.FormValue("novaSenha")

		if novaSenha == "" {
			novaSenha = senha
		}

		file, _, err := r.FormFile("pfxfile")
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Falha ao carregar o arquivo pfx"})
			return
		}

		defer file.Close()

		certPFX, err := io.ReadAll(file)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Falha ao ler o conteúdo o pfx"})
			return
		}

		newPfx, err := convert.ConvertPfxNext(certPFX, senha, novaSenha)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: err.Error()})
			return
		}

		nomeSaida := formatarNomeArquivo(r.FormValue("nomeSaida"), "certificadonovo", ".pfx")

		w.Header().Set("Content-Type", "application/x-pkcs12")
		w.Header().Set("Content-Disposition", `attachment; filename="`+nomeSaida+`"`)
		w.Write(newPfx)

	})

	http.HandleFunc("/convert-csr", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var dados convert.ItauCSR

		if err := json.NewDecoder(r.Body).Decode(&dados); err != nil {
			responseJSON(w, Page{Success: false, Error: "Falha ao ler os dados enviados"})
			return
		}

		csrPEM, keyPEM, err := convert.GenNewCSR(dados)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: err.Error()})
			return
		}

		responseJSON(w, Page{
			Success:   true,
			NomeCert:  formatarNomeArquivo(dados.CommonName, "CSR_", ".csr"),
			CertPEM:   string(csrPEM),
			NomeChave: formatarNomeArquivo(dados.CommonName, "ChavePrivada_", ".key"),
			KeyPEM:    string(keyPEM),
		})

	})

	http.HandleFunc("/planilha/faturas", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			responseJSON(w, FaturaResponse{Success: false, Error: "O arquivo excede o limite permitido 10MB"})
			return
		}

		file, header, err := r.FormFile("planilha")
		if err != nil {
			responseJSON(w, FaturaResponse{Success: false, Error: "Falha ao carregar a planilha"})
			return
		}

		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			responseJSON(w, FaturaResponse{Success: false, Error: "Falha ao ler o conteúdo da planilha"})
			return
		}

		valores, err := utils.ExtractValues(data, header.Filename)
		if err != nil {
			responseJSON(w, FaturaResponse{Success: false, Error: err.Error()})
			return
		}

		responseJSON(w, FaturaResponse{Success: true, Valores: valores})

	})

	http.HandleFunc("/remessa/validar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := r.ParseMultipartForm(20 << 20); err != nil {
			responseJSON(w, RemessaResponse{Success: false, Error: "O arquivo excede o limite permitido 20MB"})
			return
		}

		file, _, err := r.FormFile("remessa")
		if err != nil {
			responseJSON(w, RemessaResponse{Success: false, Error: "Falha ao carregar o arquivo de remessa"})
			return
		}

		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			responseJSON(w, RemessaResponse{Success: false, Error: "Falha ao ler o conteúdo do arquivo"})
			return
		}

		resultado, err := utils.ValidarRemessa(data)
		if err != nil {
			responseJSON(w, RemessaResponse{Success: false, Error: err.Error()})
			return
		}

		responseJSON(w, RemessaResponse{
			Success:         true,
			TipoCnab:        resultado.TipoCnab,
			TotalLinhas:     resultado.TotalLinhas,
			QuantidadeLotes: resultado.QuantidadeLotes,
			Erros:           resultado.Erros,
			Valido:          resultado.Valido,
		})

	})

	http.HandleFunc("GET /collections/{banco}/campos", handlerFieldsCollection)
	http.HandleFunc("POST /collections/{banco}", handlerGenCollection)

	go abrirNavegador("http://localhost:8080")

	fmt.Println("Servidor inicializado em http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Erro ao iniciar o servidor")
	}
}

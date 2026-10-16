package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"io.github.com/conv-pfx/internal/assets"
	"io.github.com/conv-pfx/internal/convert"
	"io.github.com/conv-pfx/utils"
)

// Page é a resposta JSON comum a todos os endpoints de conversão de
// certificado — só os campos relevantes pra cada banco vêm preenchidos,
// o resto fica omitido (omitempty) no JSON de saída.
type Page struct {
	Success       bool   `json:"sucesso"`
	Error         string `json:"erro,omitempty"`
	NomeCert      string `json:"nomeCert,omitempty"`
	CertPEM       string `json:"certPem,omitempty"`
	NomeChave     string `json:"nomeChave,omitempty"`
	KeyPEM        string `json:"keyPem,omitempty"`
	NomeRaiz      string `json:"nomeRaiz,omitempty"`
	RaizPEM       string `json:"raizPem,omitempty"`
	NomeIntermed1 string `json:"nomeIntermed1,omitempty"`
	Intermed1PEM  string `json:"intermed1Pem,omitempty"`
	NomeIntermed2 string `json:"nomeIntermed2,omitempty"`
	Intermed2PEM  string `json:"intermed2Pem,omitempty"`
}

// FaturaResponse é a resposta da extração de códigos de fatura de uma
// planilha — Valores vem na ordem em que os códigos aparecem na planilha.
type FaturaResponse struct {
	Success bool     `json:"sucesso"`
	Error   string   `json:"erro,omitempty"`
	Valores []string `json:"valores,omitempty"`
}

// RemessaResponse é a resposta da validação de um arquivo de remessa
// bancária — Erros lista cada inconsistência estrutural encontrada, já
// com a linha exata do arquivo original.
type RemessaResponse struct {
	Success         bool           `json:"sucesso"`
	Error           string         `json:"erro,omitempty"`
	TipoCnab        utils.TipoCnab `json:"tipoCnab,omitempty"`
	TotalLinhas     int            `json:"totalLinhas,omitempty"`
	QuantidadeLotes int            `json:"quantidadeLotes,omitempty"`
	Erros           []utils.Erro   `json:"erros,omitempty"`
	Valido          bool           `json:"valido,omitempty"`
}

// formatarNomeArquivo monta o nome de download a partir do apelido digitado
// pelo usuário, removendo qualquer extensão que ele já tenha incluído sem
// querer e aplicando a extensão correta. Se o campo vier vazio, usa só o prefixo
func formatarNomeArquivo(nomeInput, prefixo, extensao string) string {

	nomeLimpo := strings.TrimSpace(nomeInput)

	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".pem")
	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".crt")
	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".key")
	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".pfx")
	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".cer")

	if nomeLimpo == "" {
		return prefixo + extensao
	}

	return fmt.Sprintf("%s%s%s", prefixo, nomeLimpo, extensao)
}

// abrirNavegador espera o servidor subir e abre o Chrome apontando pra ele
// — o comando varia por sistema operacional porque não existe um jeito
// único de abrir o navegador padrão com argumentos pelo shell em
// Windows/macOS/Linux. Falha ao abrir é ignorada de propósito: o servidor
// continua funcionando normalmente, o usuário só precisaria abrir o
// endereço manualmente.
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

// responseJSON serializa qualquer struct de resposta (Page, FaturaResponse,
// RemessaResponse) como JSON — centraliza o content-type pra não repetir em
// cada handler.
func responseJSON(w http.ResponseWriter, res any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

// lastHeartbeat guarda o horário do último sinal de vida recebido do
// navegador. Começa em time.Now() pra dar um tempo de folga entre o
// servidor subir e a página terminar de carregar e mandar o primeiro sinal.
var (
	lastHeartbeatMu sync.Mutex
	lastHeartbeat   = time.Now()
)

// heartbeat é chamado pelo frontend em intervalos curtos enquanto a aba
// estiver aberta — só atualiza o horário do último sinal recebido.
func heartbeat(w http.ResponseWriter, r *http.Request) {
	lastHeartbeatMu.Lock()
	lastHeartbeat = time.Now()
	lastHeartbeatMu.Unlock()
}

// vigiarAba desliga o programa sozinho quando os sinais de vida param de
// chegar — ou seja, quando o usuário fecha a aba ou o navegador inteiro.
// Sem isso, o processo fica esquecido rodando em segundo plano ocupando a
// porta 8080 indefinidamente, exigindo encerrar manualmente pelo Gerenciador
// de Tarefas antes de conseguir abrir o programa de novo.
func vigiarAba() {
	for {
		time.Sleep(2 * time.Second)

		lastHeartbeatMu.Lock()
		ultimo := lastHeartbeat
		lastHeartbeatMu.Unlock()

		if time.Since(ultimo) > 10*time.Second {
			fmt.Println("Nenhum sinal do navegador recebido, encerrando o servidor.")
			os.Exit(0)
		}
	}
}

func main() {
	http.Handle("/", http.FileServer(assets.GetCoreFile()))

	// /convert-pem: converte um .pfx em chave pública + privada PEM. Pro
	// Banco do Brasil (banco == "bb") usa ConvertPFXBB, que também extrai a
	// cadeia de certificação (raiz + intermediários) exigida pelo portal
	// Developer do BB; pros demais bancos usa o conversor genérico.
	http.HandleFunc("/convert-pem", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			responseJSON(w, Page{Success: false, Error: "The file exceeds the 10MB limit"})
			return
		}

		senha := r.FormValue("password")

		file, _, err := r.FormFile("pfxfile")
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Failed to load certificate"})
			return
		}

		defer file.Close()

		pfxByte, err := io.ReadAll(file)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Failed to read the contents of the pfx file"})
			return
		}

		nomePub := formatarNomeArquivo(r.FormValue("nomeCert"), "ChavePublica", ".pem")
		nomePriv := formatarNomeArquivo(r.FormValue("nomeChave"), "ChavePrivada", ".pem")

		if r.FormValue("banco") == "bb" {
			keyPEM, certPEM, cadeia, err := convert.ConvertPFXBB(pfxByte, senha)
			if err != nil {
				responseJSON(w, Page{Success: false, Error: err.Error()})
				return
			}

			responseJSON(w, Page{
				Success:       true,
				NomeCert:      nomePub,
				CertPEM:       string(certPEM),
				NomeChave:     nomePriv,
				KeyPEM:        string(keyPEM),
				NomeRaiz:      formatarNomeArquivo(r.FormValue("nomeCert"), "CertificadoRaiz", ".cer"),
				RaizPEM:       cadeia.Raiz,
				NomeIntermed1: formatarNomeArquivo(r.FormValue("nomeCert"), "CertificadoIntermediario1", ".cer"),
				Intermed1PEM:  cadeia.Intermed1,
				NomeIntermed2: formatarNomeArquivo(r.FormValue("nomeCert"), "CertificadoIntermediario2", ".cer"),
				Intermed2PEM:  cadeia.Intermed2,
			})
			return
		}

		keyPEM, certPEM, err := convert.ConvertPFX(pfxByte, senha)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: err.Error()})
			return
		}

		responseJSON(w, Page{
			Success:   true,
			NomeCert:  nomePub,
			CertPEM:   string(certPEM),
			NomeChave: nomePriv,
			KeyPEM:    string(keyPEM),
		})

	})

	// /convert-sicredi-pem: normaliza o .cer + .key que o Sicredi entrega
	// prontos (banco não usa .pfx) para o formato PEM.
	http.HandleFunc("/convert-sicredi-pem", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			responseJSON(w, Page{Success: false, Error: "The file exceeds the 10MB limit"})
			return
		}

		cerFile, _, err := r.FormFile("cerfile")
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Failed to load certificate .cer"})
			return
		}
		defer cerFile.Close()

		cerByte, err := io.ReadAll(cerFile)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Failed to read the contents of the .cer file"})
			return
		}

		keyFile, _, err := r.FormFile("keyfile")
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Failed to load private key .key"})
			return
		}
		defer keyFile.Close()

		keyByte, err := io.ReadAll(keyFile)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Failed to read the contents of the .key file"})
			return
		}

		certPEM, err := convert.ConvertCERtoPEM(cerByte)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: err.Error()})
			return
		}

		keyPEM, err := convert.ConvertKEYtoPEM(keyByte)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: err.Error()})
			return
		}

		responseJSON(w, Page{
			Success:   true,
			NomeCert:  formatarNomeArquivo(r.FormValue("nomeCert"), "CertificadoPublico", ".pem"),
			CertPEM:   string(certPEM),
			NomeChave: formatarNomeArquivo(r.FormValue("nomeChave"), "ChavePrivada", ".pem"),
			KeyPEM:    string(keyPEM),
		})

	})

	// /convert-next: re-encoda um .pfx existente com uma nova senha (e
	// devolve o arquivo .pfx pronto pra download, não um JSON). Se o campo
	// de nova senha vier vazio, mantém a mesma senha original.
	http.HandleFunc("/convert-next", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			responseJSON(w, Page{Success: false, Error: "The file exceeds the 10MB limit"})
			return
		}

		senha := r.FormValue("password")
		novaSenha := r.FormValue("novaSenha")

		if novaSenha == "" {
			novaSenha = senha
		}

		file, _, err := r.FormFile("pfxfile")
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Failed to load pfx file"})
			return
		}

		defer file.Close()

		certPFX, err := io.ReadAll(file)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: "Failed to read the contents of the pfx"})
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

	// /convert-csr: gera um CSR + chave privada novos para o credenciamento
	// mTLS do Itaú — esse banco não envolve conversão de .pfx, o par é
	// criado do zero a partir dos dados do formulário.
	http.HandleFunc("/convert-csr", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		var dados convert.ItauCSR

		if err := json.NewDecoder(r.Body).Decode(&dados); err != nil {
			responseJSON(w, Page{Success: false, Error: "Failed to read data"})
			return
		}

		csrPEM, keyPEM, err := convert.GenNewCSR(dados)
		if err != nil {
			responseJSON(w, Page{Success: false, Error: err.Error()})
			return
		}

		responseJSON(w, Page{
			Success:   true,
			NomeCert:  "csr_itau_certificado.csr",
			CertPEM:   string(csrPEM),
			NomeChave: "itau_chave_privada.key",
			KeyPEM:    string(keyPEM),
		})

	})

	// /planilha/faturas: lê uma planilha .csv/.xlsx e devolve só os
	// códigos de fatura encontrados, prontos pra copiar (ferramenta
	// "Substituição de Códigos" da aba Utilidades).
	http.HandleFunc("/planilha/faturas", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			responseJSON(w, FaturaResponse{Success: false, Error: "The file exceeds the 10MB limit"})
			return
		}

		file, header, err := r.FormFile("planilha")
		if err != nil {
			responseJSON(w, FaturaResponse{Success: false, Error: "Failed to load xlsx"})
			return
		}

		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			responseJSON(w, FaturaResponse{Success: false, Error: "Failed to read the contents of the xlsx"})
			return
		}

		valores, err := utils.ExtractValues(data, header.Filename)
		if err != nil {
			responseJSON(w, FaturaResponse{Success: false, Error: err.Error()})
			return
		}

		responseJSON(w, FaturaResponse{Success: true, Valores: valores})

	})

	// /remessa/validar: confere a integridade estrutural de um arquivo de
	// remessa bancária CNAB 240/400 (layout detectado automaticamente) —
	// não valida o conteúdo de negócio dos campos, só a estrutura do
	// arquivo em si.
	http.HandleFunc("/remessa/validar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := r.ParseMultipartForm(20 << 20); err != nil {
			responseJSON(w, RemessaResponse{Success: false, Error: "The file exceeds the 20MB limit"})
			return
		}

		file, _, err := r.FormFile("remessa")
		if err != nil {
			responseJSON(w, RemessaResponse{Success: false, Error: "Failed to load file"})
			return
		}

		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			responseJSON(w, RemessaResponse{Success: false, Error: "Failed to read the contents of the file"})
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

	// /heartbeat: recebe o sinal de vida do frontend; sem ele por 10s
	// seguidos, o programa se encerra sozinho (ver vigiarAba).
	http.HandleFunc("/heartbeat", heartbeat)

	go abrirNavegador("http://localhost:8080")
	go vigiarAba()

	fmt.Println("Server initializing locally 🚀")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start the server")
	}
}

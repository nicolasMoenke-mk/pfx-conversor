package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"io.github.com/conv-pfx/internal/assets"
	"io.github.com/conv-pfx/internal/convert"
	"io.github.com/conv-pfx/utils"
)

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
	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".cer")

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

func main() {
	http.Handle("/", http.FileServer(assets.GetCoreFile()))

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

	go abrirNavegador("http://localhost:8080")

	fmt.Println("Servidor inicializando localmente 🚀")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Erro ao iniciar o servidor")
	}
}

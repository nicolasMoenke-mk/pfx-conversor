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
)

type Page struct {
	Success   bool   `json:"sucesso"`
	Error     string `json:"erro,omitempty"`
	NomeCert  string `json:"nomeCert,omitempty"`
	CertPEM   string `json:"certPem,omitempty"`
	NomeChave string `json:"nomeChave,omitempty"`
	KeyPEM    string `json:"keyPem,omitempty"`
}

func formatarNomeArquivo(nomeInput, prefixo string) string {

	nomeLimpo := strings.TrimSpace(nomeInput)

	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".pem")
	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".crt")
	nomeLimpo = strings.TrimSuffix(nomeLimpo, ".key")

	if nomeLimpo == "" {
		return prefixo + ".pem"
	}

	return fmt.Sprintf("%s%s.pem", prefixo, nomeLimpo)
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

func responseJSON(w http.ResponseWriter, res Page) {
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

		nomePub := formatarNomeArquivo(r.FormValue("nomeCert"), "ChavePublica")
		nomePriv := formatarNomeArquivo(r.FormValue("nomeChave"), "ChavePrivada")

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

		nomeSaida := formatarNomeArquivo(r.FormValue("nomeSaida"), "certificadonovo")
		nomeSaida = strings.TrimSuffix(nomeSaida, ".pfx")

		w.Header().Set("Content-Type", "application/pkcs-12")
		w.Header().Set("Content-Disposition", `attachment, filename="`+nomeSaida+`"`)
		w.Write(newPfx)

	})

	http.HandleFunc("/collections", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

	})

	go abrirNavegador("http://localhost:8080")

	fmt.Println("Servidor inicializado em http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Erro ao iniciar o servidor")
	}
}
